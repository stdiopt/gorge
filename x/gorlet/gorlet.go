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

// Entity is a gui component
type Entity struct {
	gorgeui.ElementComponent
	gorgeui.RectComponent

	// Maybe an array of funcs for multiple observers.
	props map[string][]func(interface{})
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
	if e.props == nil {
		return
	}
	if fns, ok := e.props[name]; ok {
		for _, fn := range fns {
			fn(value)
		}
	}
}

func (e *Entity) observe(k string, fn interface{}) {
	if e.props == nil {
		e.props = map[string][]func(interface{}){}
	}
	e.props[k] = append(e.props[k], makePropFunc(fn))
}

// Add adds a children to entity.
func (e *Entity) Add(children ...*Entity) {
	for _, c := range children {
		gorgeui.AddChildrenTo(e, c)
	}
	// Relayout
	if e.LayoutFunc != nil {
		e.LayoutFunc(e)
	}
}

// AddElement adds an UI element to entity.
func (e *Entity) AddElement(children ...*Entity) {
	for _, c := range children {
		gorgeui.AddElementTo(e, c)
	}
}
