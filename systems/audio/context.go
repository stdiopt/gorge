package audio

import "github.com/stdiopt/gorge"

type audio = Audio

// Context to be used on func binders
type Context struct {
	*audio
}

// FromContext returns an audio context from a gorge.Context.
func FromContext(g *gorge.Context) *Context {
	var ret *Context
	if err := g.BindProps(func(c *Context) { ret = c }); err != nil {
		g.Error(err)
	}
	return ret
}
