package gorlet

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type (
	renderabler interface {
		Renderable() *gorge.RenderableComponent
	}
)

// PlacementFunc will be used in a container and will define clients rect.
type PlacementFunc func(w *Entity) // OnAdd in the Entity

// Entity is a gui component
type Entity struct {
	gorgeui.ElementComponent
	gorgeui.RectComponent

	name string // debug purposes

	// TODO: {lpf} give it a proper name.
	onAdd      PlacementFunc
	clientArea *Entity
	layouter   Layouter

	observers map[string][]func(interface{})
}

// XXX:
// we could use layouter func here instead of gorgeui
// this way we could just run through *Entity since this is
// an higher level entity.
//
// Other way would be use the eventBus on the Create func.
// Since then we wouldn't have HandleEvent exposed

/*func (e *Entity) HandleEvent(evt event.Event) {
	switch e := evt.(type) {
	case gorgeui.EventUpdate:
		if e.layouter != nil {
			e.layouter.Layout(e)
		}
	}
}*/

func (e *Entity) String() string {
	return e.name
}

// Create creates builds and prepares a guilet
func Create(fn BuildFunc) *Entity {
	fnVal := reflect.ValueOf(fn)
	fi := runtime.FuncForPC(fnVal.Pointer())
	ename := fi.Name()
	ename = ename[strings.LastIndex(ename, "/")+1:]

	// fi := runtime.FuncForPC((uintptr)unsafe.Pointer(&fn))
	defEntity := &Entity{
		name:          ename,
		RectComponent: *gorgeui.NewRectComponent(),
	}
	// root.SetLayouter(gorgeui.AutoHeight(1))
	// root.SetAnchor(0)
	// root.SetRect(0, 0, 30, 5)
	defEntity.SetPivot(0)
	b := Builder{
		root: &curEntity{entity: defEntity},
	}
	fn(&b)
	return b.root.entity
}

// Client returns the client area of the entity, the client area
// is an Entity where the children will be added using Add method.
func (e *Entity) Client() *Entity {
	if e.clientArea == nil {
		return e
	}
	return e.clientArea
}

// TODO: with clientArea we might not need Element since Listing childs
// would be directly in the client area.?

// SetClientArea is the child Entity where Add will put Entities.
func (e *Entity) SetClientArea(c *Entity) {
	e.clientArea = c
}

// OnAdd triggers when the entitywhen the entity is added to the parent.
func (e *Entity) OnAdd(fn func(e *Entity)) {
	if e.clientArea != nil {
		e.clientArea.OnAdd(fn)
		return
	}
	e.onAdd = fn
}

// SetLayout Will set the layouter thing on client Entity.
func (e *Entity) SetLayout(l gorgeui.Layouter) {
	if e.clientArea != nil {
		e.clientArea.SetLayout(l)
		return
	}
	e.Layouter = l
}

// Add adds a children to entity.
func (e *Entity) Add(children ...*Entity) {
	if e.clientArea != nil {
		e.clientArea.Add(children...)
		if e.Layouter != nil {
			e.Layouter.Layout(e)
		}
		return
	}

	for _, c := range children {
		if e.onAdd != nil {
			e.onAdd(c)
		}
		// This adds it to children class.
		gorgeui.AddChildrenTo(e, c)
	}

	// TODO: {lpf} This could be only on gorlet and do a tree call

	// Relayout this entity as well as children will also relayout
	// themselves on add.
	if e.Layouter != nil {
		e.Layouter.Layout(e)
	}
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
}

// Set invoke any observer attached to the named propery.
func (e *Entity) Set(name string, value interface{}) {
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
func (e *Entity) PropSetter(name string) func(v interface{}) {
	return func(v interface{}) { e.Set(name, v) }
}

// Observe adds a named observer setting nil will delete all observers.
func (e *Entity) Observe(k string, fn interface{}) {
	if e.observers == nil {
		e.observers = map[string][]func(interface{}){}
	}
	if fn == nil {
		delete(e.observers, k)
		return
	}
	e.observers[k] = append(e.observers[k], makePropFunc(k, fn))
}

/*
func (e *Entity) Add(children ...*Entity) {
	if e.clientArea != nil {
		e.clientArea.Add(children...)
		if e.Layouter != nil {
			e.Layouter.Layout(e)
		}
		return
	}
	for _, c := range children {
		if e.onAdd != nil {
			e.onAdd(c)
		}
		// This adds it to children class.
		gorgeui.AddChildrenTo(e, c)
	}

	// TODO: {lpf} This could be only on gorlet and do a tree call
	// Relayout
	if e.Layouter != nil {
		e.Layouter.Layout(e)
	}
}
*/

// AddElement adds an UI element to entity.
func (e *Entity) AddElement(els ...gorge.Entity) {
	for _, c := range els {
		gorgeui.AddElementTo(e, c)
	}
}

// RemoveElement removes element from entity and resets elements parent,
// if the element is attached it will trigger gorge remove event.
func (e *Entity) RemoveElement(els ...gorge.Entity) {
	for _, c := range els {
		gorgeui.RemoveElementFrom(e, c)
	}
}

// FillParent will reset anchor to 0,0 1,1 and Rect to 0,0,0,0.
func (e *Entity) FillParent(n float32) {
	e.SetAnchor(0, 0, 1, 1)
	e.SetRect(n)
}
