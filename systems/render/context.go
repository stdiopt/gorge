package render

import (
	"log"
	"time"

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
	if ctx, ok := gorge.GetContext[*Context](g); ok {
		return ctx
	}

	log.Println("initializing system")
	r := newRenderer(g)
	ctx := &Context{render: r}
	gorge.AddContext(g, ctx)

	var (
		totalTime      float32
		statTimeCount  float32 = 3
		renderDuration time.Duration
	)
	gorge.HandleFunc(g, func(gorge.EventStart) {
		r.Init()
	})
	gorge.HandleFunc(g, func(e gorge.EventAddEntity) {
		// Not a switch since the entity can be both a Camera and a Light.
		if v, ok := e.Entity.(Camera); ok {
			r.AddCamera(v)
		}
		if v, ok := e.Entity.(Light); ok {
			r.AddLight(v)
		}
		if v, ok := e.Entity.(Renderable); ok {
			r.AddRenderable(v)
		}
	})
	gorge.HandleFunc(g, func(e gorge.EventRemoveEntity) {
		if v, ok := e.Entity.(Camera); ok {
			r.RemoveCamera(v)
		}
		if v, ok := e.Entity.(Light); ok {
			r.RemoveLight(v)
		}
		if v, ok := e.Entity.(Renderable); ok {
			r.RemoveRenderable(v)
		}
	})
	gorge.HandleFunc(g, func(e gorge.EventRender) {
		dt := float32(e)
		statTimeCount -= dt
		if statTimeCount < 0 {
			gorge.Trigger(g, EventStat{
				Textures:       r.textures.count,
				VBOs:           r.vbos.count,
				Shaders:        r.shaders.count,
				Instances:      len(r.Renderables),
				Buffers:        r.buffers.count,
				RenderDuration: renderDuration,
				DrawCalls:      r.DrawCalls,
			})
			statTimeCount = 1
		}

		totalTime += dt

		if r.DisableRender {
			return
		}
		mark := time.Now()
		r.Render()
		renderDuration = time.Since(mark)
	})
	gorge.HandleFunc(g, func(e gorge.EventResourceUpdate) {
		switch rr := e.Resource.(type) {
		case *gorge.TextureData:
			r.textures.Update(rr)
		case *gorge.MeshData:
			r.vbos.Update(rr)
		}
	})
	return ctx
}
