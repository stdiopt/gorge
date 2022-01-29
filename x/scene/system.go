package scene

import (
	"log"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
)

type Context struct {
	gorge *gorge.Context
}

func System(g *gorge.Context) error {
	FromContext(g)
	return nil
}

type sceneGetter interface {
	GetScene() *Scene
}

func FromContext(g *gorge.Context) *Context {
	if ctx, ok := gorge.GetContext[*Context](g); ok {
		return ctx
	}
	log.Println("Initializing system")

	event.Handle(g, func(e gorge.EventAddEntity) {
		if sg, ok := e.Entity.(sceneGetter); ok {
			s := sg.GetScene()
			s.initScene(g)
			event.Trigger(s, EventAttached{g})
		}
	})
	event.Handle(g, func(e gorge.EventRemoveEntity) {
		if sg, ok := e.Entity.(sceneGetter); ok {
			s := sg.GetScene()
			s.destroyScene(g)
		}
	})
	ctx := &Context{
		gorge: g,
	}

	// Handle scene management
	return gorge.SetContext(g, ctx)
}
