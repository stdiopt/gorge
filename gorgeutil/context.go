package gorgeutil

import "github.com/stdiopt/gorge"

// Context extends gorge context with default entity instantiators.
type Context struct {
	*gorge.Context
}

// FromContext returns a gorgeutil context
func FromContext(g *gorge.Context) *Context {
	return &Context{g}
}
