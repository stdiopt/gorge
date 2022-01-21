package particle

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
)

type emitter interface {
	Mat4() gm.Mat4
	Transform() *gorge.TransformComponent
	Emitter() *EmitterComponent
}

type Context struct{}

func System(g *gorge.Context) error {
	if _, ok := gorge.GetContext[*Context](g); ok {
		return nil
	}
	gorge.AddContext(g, &Context{})

	emitters := []emitter{}

	event.Handle(g, func(e gorge.EventAddEntity) {
		em, ok := e.Entity.(emitter)
		if !ok {
			return
		}
		em.Emitter().init(g, em)
		emitters = append(emitters, em)
	})
	event.Handle(g, func(e gorge.EventRemoveEntity) {
		em, ok := e.Entity.(emitter)
		if !ok {
			return
		}
		em.Emitter().destroy(g)
		for i, eem := range emitters {
			if eem == em {
				t := emitters
				emitters = append(emitters[:i], emitters[i+1:]...)
				t[len(t)-1] = nil
				break
			}
		}
	})
	event.Handle(g, func(e gorge.EventUpdate) {
		for _, em := range emitters {
			em.Emitter().update(g, em, e.DeltaTime())
			// Fixed rate
			// em.Emitter().update(g, em, 0.016)
		}
	})

	return nil
}
