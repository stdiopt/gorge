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
	if ctx, ok := gorge.GetContext[*Context](g); ok {
		return ctx
	}

	log.Println("Initializing system")
	audio := &Audio{
		sources: map[*gorge.AudioSource]*Processor{},
	}
	ctx := &Context{audio}
	gorge.AddContext(g, ctx)
	// g.PutProp(&Context{audio})
	g.AddHandler(audio)

	return ctx
}
