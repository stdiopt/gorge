package render

import (
	"log"

	"github.com/stdiopt/gorge"
)

var ctxKey = struct{ string }{"render"}

type render = Render

// Context rendering context to be used on gorge systems.
type Context struct {
	*render
}

// FromContext returns a Context from gorge Context
// triggers an error and returns nil if doesn't exists.
func FromContext(g *gorge.Context) *Context {
	if ctx, ok := gorge.GetSystem(g, ctxKey).(*Context); ok {
		return ctx
	}
	log.Println("initializing system")

	r := newRenderer(g)
	ctx := &Context{render: r}
	gorge.AddSystem(g, ctxKey, ctx)
	g.Handle(&system{
		gorge:         g,
		renderer:      r,
		statTimeCount: 3,
	})
	return ctx
}
