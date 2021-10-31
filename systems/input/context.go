package input

import "github.com/stdiopt/gorge"

type input = Input

// Context to be used in gorge systems
type Context struct {
	*input
}

// FromContext returns a input.Context from a gorge.Context
func FromContext(g *gorge.Context) *Context {
	var ret *Context
	if err := g.BindProps(func(c *Context) { ret = c }); err != nil {
		g.Error(err)
	}
	return ret
}
