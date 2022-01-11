package particle

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/render"
)

// Controller is a configurable particle controller
type Controller[T any] struct {
	CreateFunc func(*T)
	InitFunc   func(*T)
	UpdateFunc func(*T, float32)
}

// Create calls the CreateFunc if set.
func (c Controller[T]) Create(p *T) {
	if c.CreateFunc != nil {
		c.CreateFunc(p)
	}
}

// Init calls the init func if set.
func (c Controller[T]) Init(p *T) {
	if c.InitFunc != nil {
		c.InitFunc(p)
	}
}

// Update calls the update func if set.
func (c Controller[T]) Update(p *T, dt float32) {
	if c.UpdateFunc != nil {
		c.UpdateFunc(p, dt)
	}
}

// Emitter entity will emit particles based on the given parameters when added
// to gorge, particle.System must be in gorge initialization list.
type Emitter[T any] struct {
	gorge.TransformComponent
	EmitterComponent
}

// NewEmitter creates a new default emitter.
func NewEmitter[T any]() *Emitter[T] {
	return &Emitter[T]{
		TransformComponent: *gorge.NewTransformComponent(),
		EmitterComponent:   *NewEmitterComponent[T](),
	}
}

// SetController sets the controller for the particles.
func (e *Emitter[T]) SetController(c controller[T]) {
	e.Generator.(*Generator[T]).Controller = c
}

// Entity default particle entity
type Entity struct {
	gorge.TransformComponent
	gorge.ColorableComponent
	Component
}

var _ render.Renderable = &Entity{}
