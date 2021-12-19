package audio

import (
	"log"

	"github.com/stdiopt/gorge"
)

var ctxKey = struct{ string }{"audio"}

type audio = Audio

// Context to be used on func binders
type Context struct {
	*audio
}

// FromContext returns an audio context from a gorge.Context.
func FromContext(g *gorge.Context) *Context {
	if ctx, ok := gorge.GetSystem(g, ctxKey).(*Context); ok {
		return ctx
	}

	log.Println("Initializing system")
	audio := &Audio{
		sources: map[*gorge.AudioSource]*Processor{},
	}
	ctx := &Context{audio}
	gorge.AddSystem(g, ctxKey, ctx)
	// g.PutProp(&Context{audio})
	g.Handle(audio)

	return ctx
}
