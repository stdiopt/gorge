package gorlet

import (
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/math/ray"
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
	observers
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

	// observers map[string][]ObserverFunc

	// Temp solution for the thing
	Masked bool
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
	if b.clientArea != nil && b.clientArea != b.root.entity {
		b.root.entity.SetClientArea(b.clientArea)
	}

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

func (e *Entity) IsMasked() bool {
	return e.Masked
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
	if e.clientArea != nil {
		return e.clientArea.Children()
	}
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
// exceptional case:
// _maskDepth:
//   - if entity has observers it will call them
//   - if not it will be sent to children of the entity
func (e *Entity) Set(name string, value any) {
	if ok := e.set(name, value); ok {
		return
	}

	// Extra case we pass this to all children.
	if name == "_maskDepth" {
		for _, c := range e.children {
			c.Set(name, value)
		}
	}
}

// PropSetter returns a func that will set the named property when called.
func (e *Entity) PropSetter(name string) func(v any) {
	return func(v any) { e.Set(name, v) }
}

// Observe adds a named observer setting nil will delete all observers.
func (e *Entity) Observe(k string, ifn any) {
	e.observe(k, ifn)
}

// Link observs k1 in entity and sets k2 on e2 entity
// if there ar eno listeners on e2 it won't be added
// TODO: maybe it should be added for the late observers.
// or disable observing in entity completely (allow only in builder)
func (e *Entity) ObserveTo(k1 string, e2 *Entity, k2 string) {
	o := e2.observer(k2)
	if o == nil {
		return
	}
	e.observeWithType(k1, o.Type, o.Call)
}

func (e *Entity) Observer(k string) *Observer {
	return e.observer(k)
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

// CalcBounds calculates children bounds and positions and return min max
// size
func (e *Entity) CalcBounds() gm.Vec4 {
	var ret gm.Vec4
	sz := e.CalcSize()
	ret[2] = sz[0] + e.Margin[0] + e.Margin[2]
	ret[3] = sz[1] + e.Margin[1] + e.Margin[3]
	for _, e := range e.children {
		b := e.CalcBounds()
		ret[0] = gm.Min(ret[0], e.Position[0]+b[0])
		ret[1] = gm.Min(ret[1], e.Position[1]+b[1])
		ret[2] = gm.Max(ret[2], e.Position[0]+b[2])
		ret[3] = gm.Max(ret[3], e.Position[1]+b[3])
	}
	return ret
}

// IntersectFromScreen intersects the entity rect from screen coordinates.
func (e *Entity) IntersectFromScreen(pos gm.Vec2) ray.Result {
	sz := e.CalcSize()
	m := e.Mat4()
	v0 := m.MulV4(gm.Vec4{0, 0, 0, 1}).Vec3()     // 0
	v1 := m.MulV4(gm.Vec4{sz[0], 0, 0, 1}).Vec3() // right
	v2 := m.MulV4(gm.Vec4{0, sz[1], 0, 1}).Vec3() // up)

	ui := gorgeui.RootUI(e)
	r := ray.FromScreen(ui.ScreenSize(), ui.Camera, pos)
	return ray.IntersectRect(r, v0, v1, v2)
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
		// Wrong
		ui := gorgeui.RootUI(e)
		ui.GAdd(children...)
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
		ui.GRemove(children...)
	}
}
