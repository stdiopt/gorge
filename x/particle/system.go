package particle

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
)

type emitter interface {
	Mat4() m32.Mat4
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

	gorge.HandleFunc(g, func(e gorge.EventAddEntity) {
		em, ok := e.Entity.(emitter)
		if !ok {
			return
		}
		em.Emitter().init(g, em)
		emitters = append(emitters, em)
	})
	gorge.HandleFunc(g, func(e gorge.EventRemoveEntity) {
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
	gorge.HandleFunc(g, func(e gorge.EventUpdate) {
		for _, em := range emitters {
			em.Emitter().update(g, em, e.DeltaTime())
		}
	})

	return nil
}
