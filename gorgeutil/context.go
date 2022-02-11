package gorgeutil

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/input"
	"github.com/stdiopt/gorge/systems/resource"
)

type (
	gorgeContext    = gorge.Context
	inputContext    = input.Context
	resourceContext = resource.Context
)

// Context extends gorge context with default entity instantiators.
type Context struct {
	*gorgeContext
	*inputContext
	*resourceContext
}

// FromContext returns a gorgeutil context
func FromContext(g *gorge.Context) *Context {
	return &Context{
		gorgeContext:    g,
		inputContext:    input.FromContext(g),
		resourceContext: resource.FromContext(g),
	}
}

func Wrapper(fns ...func(*Context)) gorge.InitFunc {
	return func(g *gorge.Context) error {
		ctx := FromContext(g)
		for _, fn := range fns {
			fn(ctx)
		}
		return nil
	}
}
