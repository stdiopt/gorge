package gorlet

import (
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type (
	renderabler interface {
		Renderable() *gorge.RenderableComponent
	}
)

type (
	// EntityFunc will be used in a container and will define clients rect.
	EntityFunc func(w *Entity) // OnAdd in the Entity
	// ObserverFunc is the type of the function function that will be
	// called when the named property is set.
	ObserverFunc = func(any)
)

// Debug prints reference debug
var Debug = false

// debugCounter debug release stuff
var debugCounter int64

// DebugCounter returns the number of gorlets instantiated.
func DebugCounter() int64 {
	return atomic.LoadInt64(&debugCounter)
}

type gcref struct{ n int }

// Entity is a gui component
type Entity struct {
	gcref *gcref
	gorgeui.ElementComponent
	gorgeui.RectComponent

	container gorge.Container
	children  []*Entity

	name string // debug purposes

	// TODO: {lpf} give it a proper name.
	onAdd      EntityFunc
	clientArea *Entity
	layouter   Layouter

	observers map[string][]ObserverFunc
}

// Create creates builds and prepares a guilet
func Create(fn Func) *Entity {
	fnVal := reflect.ValueOf(fn)
	fi := runtime.FuncForPC(fnVal.Pointer())
	ename := fi.Name()
	ename = ename[strings.LastIndex(ename, "/")+1:]

	g := &gcref{1}

	// debug object release tracker
	runtime.SetFinalizer(g, func(v any) {
		if Debug {
			log.Println("Finalizing Gorlet OBJ:", ename, v)
		}
		atomic.AddInt64(&debugCounter, -1)
	})
	atomic.AddInt64(&debugCounter, 1)
	if Debug {
		log.Println("Creating gorlet:", ename)
	}

	// fi := runtime.FuncForPC((uintptr)unsafe.Pointer(&fn))
	defEntity := &Entity{
		gcref:         g,
		name:          ename,
		RectComponent: *gorgeui.NewRectComponent(),
	}
	// defEntity.SetLayout(gorgeui.AutoHeight(1))
	// defEntity.SetAnchor(0)
	// defEntity.SetRect(0, 0, 30, 5)
	defEntity.SetPivot(0)
	b := Builder{
		root: &curEntity{entity: defEntity},
	}
	fn(&b)

	entityUpdate(b.root.entity)
	return b.root.entity
}

func entityUpdate(ent *Entity) {
	ent.HandleEvent(gorgeui.EventUpdate(0))
	// huh how to solve this?
	event.Trigger(ent, gorgeui.EventUpdate(0))
	// ent.Trigger(gorgeui.EventUpdate(0))
	for _, c := range ent.Children() {
		entityUpdate(c)
	}
}

// XXX:
// we could use layouter func here instead of gorgeui
// this way we could just run through *Entity since this is
// an higher level entity.
//
// Other way would be use the eventBus on the Create func.
// Since then we wouldn't have HandleEvent exposed

// HandleEvent handles events.
func (e *Entity) HandleEvent(evt event.Event) {
	switch evt.(type) {
	case gorgeui.EventUpdate:
		if e.layouter != nil {
			e.layouter.Layout(e)
		}
	default:
	}
}

func (e *Entity) String() string {
	return e.name
}

// Client returns the client area of the entity, the client area
// is an Entity where the children will be added using Add method.
func (e *Entity) Client() *Entity {
	if e.clientArea == nil {
		return e
	}
	return e.clientArea
}

// SetClientArea is the child Entity where Add will put Entities.
func (e *Entity) SetClientArea(c *Entity) {
	e.clientArea = c
}

// OnAdd triggers when the entitywhen the entity is added to the parent.
// missing consistency
func (e *Entity) OnAdd(fn func(e *Entity)) {
	if e.clientArea != nil {
		e.clientArea.OnAdd(fn)
		return
	}
	e.onAdd = fn
}

// SetLayout Will set the layouter thing on client Entity.
func (e *Entity) SetLayout(l Layouter) {
	if e.clientArea != nil {
		e.clientArea.SetLayout(l)
		return
	}
	e.layouter = l
}

// GetEntities implement gorge.Container.
func (e *Entity) GetEntities() []gorge.Entity {
	return e.container
}

// Children returns this entity children.
func (e *Entity) Children() []*Entity {
	// Should it be fromClientArea?!
	return e.children
}

// Add adds a children to entity.
func (e *Entity) Add(child *Entity) {
	if e.clientArea != nil {
		e.clientArea.Add(child)
		if e.layouter != nil {
			e.layouter.Layout(e)
		}
		return
	}

	if e.indexOf(child) != -1 {
		return
	}
	if e.onAdd != nil {
		e.onAdd(child)
	}
	e.children = append(e.children, child)
	e.add(child)

	// TODO: {lpf} This could be only on gorlet and do a tree call
	// Relayout this entity as well as children will also relayout
	// themselves on add.
	if e.layouter != nil {
		e.layouter.Layout(e)
	}
}

// Remove removes a child from the entity.
func (e *Entity) Remove(child *Entity) {
	if e.clientArea != nil {
		e.clientArea.Remove(child)
		if e.layouter != nil {
			e.layouter.Layout(e)
		}
		return
	}

	n := e.indexOf(child)
	if n == -1 {
		return
	}
	// remove from container if exists
	e.remove(child)

	t := e.children
	e.children = append(e.children[:n], e.children[n+1:]...)
	t[len(t)-1] = nil
}

func (e *Entity) exists(c *Entity) bool {
	for _, cc := range e.children {
		if c == cc {
			return true
		}
	}
	return false
}

// AddElement adds an element to the entity
func (e *Entity) AddElement(ents ...gorge.Entity) {
	e.add(ents...)
}

// RemoveElement removes an element.
// TODO: {lpf} should this check for children here, since this can remove a child.
func (e *Entity) RemoveElement(ents ...gorge.Entity) {
	e.remove(ents...)
}

// Attached implements the Attacher interface
func (e *Entity) Attached(ent gorgeui.Entity) {
	// Sounds bad
	for _, c := range e.GetEntities() {
		if r, ok := c.(renderabler); ok {
			if ui := gorgeui.RootUI(e); ui != nil {
				r.Renderable().CullMask = ui.CullMask
			}
		}
	}
	// TODO: This will endup calling update multiple times since
	// Attach is called in every added entity
	// but should be ok since some updates depends on layouting
	// which requires settling
	entityUpdate(e)
}

// Set invoke any observer attached to the named propery.
func (e *Entity) Set(name string, value any) {
	if e.observers == nil {
		return
	}
	if fns, ok := e.observers[name]; ok {
		for _, fn := range fns {
			fn(value)
		}
	}
}

// PropSetter returns a func that will set the named property when called.
func (e *Entity) PropSetter(name string) func(v any) {
	return func(v any) { e.Set(name, v) }
}

// Observe adds a named observer setting nil will delete all observers.
func (e *Entity) Observe(k string, fn ObserverFunc) {
	if e.observers == nil {
		e.observers = map[string][]ObserverFunc{}
	}
	if fn == nil {
		delete(e.observers, k)
		return
	}
	e.observers[k] = append(e.observers[k], fn)
}

// FillParent will reset anchor to 0,0 1,1 and Rect to 0,0,0,0.
func (e *Entity) FillParent(n float32) {
	e.SetAnchor(0, 0, 1, 1)
	e.SetRect(n)
}

// SetDimRect sets anchor and pivot to 0 and use rect from v... params.
func (e *Entity) SetDimRect(v ...float32) {
	e.SetAnchor(0)
	e.SetPivot(0)
	e.SetRect(v...)
}

// SetRelRect sets anchor to fill parent and use rect from v... params.
func (e *Entity) SetRelRect(v ...float32) {
	e.SetAnchor(0, 0, 1, 1)
	e.SetRect(v...)
}

func (e *Entity) indexOf(c *Entity) int {
	for i, c2 := range e.children {
		if c == c2 {
			return i
		}
	}
	return -1
}

// add it will add to container if the entity is attached it will add to
// underlying gorge instance.
func (e *Entity) add(children ...gorge.Entity) {
	for _, c := range children {
		if t, ok := c.(gorge.ParentSetter); ok {
			t.SetParent(e)
		}
		e.container.Add(c)
	}
	if e.ElementComponent.Attached {
		ui := gorgeui.RootUI(e)
		ui.Add(children...)
	}
}

// add it will remove from container if the entity is attached it will add to
// underlying gorge instance.
func (e *Entity) remove(children ...gorge.Entity) {
	for _, c := range children {
		if t, ok := c.(gorge.ParentSetter); ok {
			t.SetParent(nil)
		}
		e.container.Remove(c)
	}
	if e.ElementComponent.Attached {
		ui := gorgeui.RootUI(e)
		ui.Remove(children...)
	}
}