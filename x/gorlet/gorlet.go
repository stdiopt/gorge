package gorlet

import (
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

	// TODO: {lpf} give it a proper name.
	onAdd PlacementFunc

	observers map[string][]func(interface{})
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

func (e *Entity) observe(k string, fn interface{}) {
	if e.observers == nil {
		e.observers = map[string][]func(interface{}){}
	}
	e.observers[k] = append(e.observers[k], makePropFunc(k, fn))
}

// Add adds a children to entity.
func (e *Entity) Add(children ...*Entity) {
	for _, c := range children {
		if e.onAdd != nil {
			e.onAdd(c)
		}
		// This adds it to children class.
		gorgeui.AddChildrenTo(e, c)
	}
	// Relayout
	if e.Layouter != nil {
		e.Layouter.Layout(e)
	}
}

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
