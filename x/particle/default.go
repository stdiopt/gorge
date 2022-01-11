package particle

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/render"
)

// Emitter entity will emit particles based on the given parameters when added
// to gorge, particle.System must be in gorge initialization list.
type Emitter[T any] struct {
	gorge.TransformComponent
	EmitterComponent
}

func NewEmitter[T any]() *Emitter[T] {
	return &Emitter[T]{
		TransformComponent: *gorge.NewTransformComponent(),
		EmitterComponent:   *NewEmitterComponent[T](),
	}
}

func (e *Emitter[T]) SetInitFunc(fn func(*T)) {
	e.Generator.(*Generator[T]).InitFunc = fn
}

func (e *Emitter[T]) SetUpdateFunc(fn func(*T)) {
	e.Generator.(*Generator[T]).UpdateFunc = fn
}

// Single particle entity
type Entity struct {
	gorge.TransformComponent
	gorge.ColorableComponent
	Component
}

var _ render.Renderable = &Entity{}

/*
type DefaultEntity struct {
	gorge.TransformComponent
	gorge.ColorableComponent
	*gorge.RenderableComponent
	Component

	Dir   m32.Vec3
	Scale m32.Vec3
}

type DefaultGenerator struct {
	*Generator[DefaultEntity]
}

func (c *DefaultGenerator) lazyInit() {
	if c.Generator != nil {
		return
	}
	c.Generator = Generator[DefaultEntity]{
		InitFunc: func(e *DefaultEntity) {
		},
		UpdateFunc: func(e *DefaultEntity, dt float32) {
		},
	}
}

func (c *DefaultGenerator) destroy(g *gorge.Context) {
	c.lazyInit()
	c.Generator.destroy(g)
}

func (c *DefaultGenerator) init(g *gorge.Context, em emitter) {
	c.lazyInit()
	c.Generator.init(g, em)
}

func (c *DefaultGenerator) update(em emitter, dt float32) {
	c.lazyInit()
	c.Generator.update(em, dt)
}
*/
