package gorgeui

import "github.com/stdiopt/gorge"

// Context to be used in gorge systems.
type Context struct {
	*system
}

// FromContext retrieve gorgeui context from gorge
func FromContext(g *gorge.Context) *Context {
	var ret *Context
	if err := g.BindProps(func(c *Context) { ret = c }); err != nil {
		g.Error(err)
	}
	return ret
}

// New returns a new UI
func (c Context) New(cam cameraEntity) *UI {
	ui := New(c.gorge)
	ui.SetCamera(cam)
	return ui
}
