package render

import (
	"time"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
)

// System gorge system initializer
/*func System(g *gorge.Context, glw *gl.Wrapper) error {
	log.Println("initializing system")
	gl.Init(glw)

	r := newRenderer(g)
	g.PutProp(func() *Context {
		return &Context{r}
	})

	g.Handle(&system{
		gorge:         g,
		renderer:      r,
		statTimeCount: 3,
	})
	return nil
}*/

type system struct {
	gorge    *gorge.Context
	renderer *Render

	totalTime      float32
	statTimeCount  float32
	renderDuration time.Duration
}

func (s *system) HandleEvent(v event.Event) {
	switch e := v.(type) {
	case gorge.EventAddEntity:
		// Not a switch since the entity can be both a Camera and a Light.
		if v, ok := e.Entity.(Camera); ok {
			s.renderer.AddCamera(v)
		}
		if v, ok := e.Entity.(Light); ok {
			s.renderer.AddLight(v)
		}
		if v, ok := e.Entity.(Renderable); ok {
			s.renderer.AddRenderable(v)
		}
	case gorge.EventRemoveEntity:
		if v, ok := e.Entity.(Camera); ok {
			s.renderer.RemoveCamera(v)
		}
		if v, ok := e.Entity.(Light); ok {
			s.renderer.RemoveLight(v)
		}
		if v, ok := e.Entity.(Renderable); ok {
			s.renderer.RemoveRenderable(v)
		}
	case gorge.EventRender:
		s.handleRender(e)
	case gorge.EventStart:
		s.renderer.Init()
	case gorge.EventResourceUpdate:
		switch r := e.Resource.(type) {
		case *gorge.TextureData:
			// log.Println("Update received")
			// s.renderer.textures.GetByRef(r)
			s.renderer.textures.Update(r)
		case *gorge.MeshData:
			// s.renderer.vbos.GetByRef(r)
			s.renderer.vbos.Update(r)
		}
	}
}

func (s *system) handleRender(e gorge.EventRender) {
	dt := float32(e)
	s.statTimeCount -= dt
	if s.statTimeCount < 0 {
		s.TriggerStat()
		s.statTimeCount = 1
	}

	s.totalTime += dt

	if s.renderer.DisableRender {
		return
	}
	mark := time.Now()
	s.renderer.Render()
	s.renderDuration = time.Since(mark)
}

// EventStat track gpu resources for debugging
type EventStat struct {
	VBOs           int
	Textures       int
	Shaders        int
	Instances      int
	Buffers        int
	RenderDuration time.Duration
	DrawCalls      int
}

func (s *system) TriggerStat() {
	stat := EventStat{
		Textures:       s.renderer.textures.count,
		VBOs:           s.renderer.vbos.count,
		Shaders:        s.renderer.shaders.count,
		Instances:      len(s.renderer.Renderables),
		Buffers:        s.renderer.buffers.count,
		RenderDuration: s.renderDuration,
		DrawCalls:      s.renderer.DrawCalls,
	}
	s.gorge.Trigger(stat) // nolint: errcheck
}
