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

// Element is a gui component
type Element struct {
	gorgeui.ElementComponent
	gorgeui.RectComponent

	// Maybe an array of funcs for multiple observers.
	props map[string][]func(interface{})
}

// Attached implements the Attacher interface
func (g *Element) Attached(e gorgeui.Entity) {
	// Sounds bad
	for _, c := range g.GetEntities() {
		if r, ok := c.(renderabler); ok {
			if ui := gorgeui.RootUI(g); ui != nil {
				r.Renderable().CullMask = ui.CullMask
			}
		}
	}
}

// Set sets a property on the guilet.
func (g *Element) Set(name string, value interface{}) {
	if g.props == nil {
		return
	}
	if fns, ok := g.props[name]; ok {
		for _, fn := range fns {
			fn(value)
		}
	}
}

func (g *Element) observe(k string, fn interface{}) {
	if g.props == nil {
		g.props = map[string][]func(interface{}){}
	}
	g.props[k] = append(g.props[k], makePropFunc(fn))
}
