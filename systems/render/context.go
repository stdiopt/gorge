package render

import "github.com/stdiopt/gorge"

type render = Render

// Context rendering context to be used on gorge systems.
type Context struct {
	*render
}

// FromContext returns a Context from gorge Context
// triggers an error and returns nil if doesn't exists.
func FromContext(g *gorge.Context) *Context {
	var ret *Context
	if err := g.BindProps(func(c *Context) { ret = c }); err != nil {
		g.Error(err)
	}
	return ret
}
