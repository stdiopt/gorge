package particle

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/render"
)

// Emitter entity will emit particles based on the given parameters when added
// to gorge, particle.System must be in gorge initialization list.
type Emitter[T any, Tp particler[T]] struct {
	gorge.TransformComponent
	EmitterComponent

	CreateFunc func(*T)
	InitFunc   func(*T)
	UpdateFunc func(*T, float32)
}

// NewEmitter creates a new default emitter.
func NewEmitter[T any, Tp particler[T]]() *Emitter[T, Tp] {
	return &Emitter[T, Tp]{
		TransformComponent: *gorge.NewTransformComponent(),
		EmitterComponent:   *NewEmitterComponent[T, Tp](),
	}
}

func (e *Emitter[T, Tp]) SetCreateFunc(f func(*T)) {
	e.CreateFunc = f
}

func (e *Emitter[T, Tp]) SetInitFunc(f func(*T)) {
	e.InitFunc = f
}

func (e *Emitter[T, Tp]) SetUpdateFunc(f func(*T, float32)) {
	e.UpdateFunc = f
}

// CreateParticle implements particle creator method.
func (e *Emitter[T, Tp]) CreateParticle(p *T) {
	if e.CreateFunc != nil {
		e.CreateFunc(p)
	}
}

// InitParticle implements particle initializer method.
func (e *Emitter[T, Tp]) InitParticle(p *T) {
	if e.InitFunc != nil {
		e.InitFunc(p)
	}
}

// UpdateParticle implements particle updater method.
func (e *Emitter[T, Tp]) UpdateParticle(p *T, dt float32) {
	if e.UpdateFunc != nil {
		e.UpdateFunc(p, dt)
	}
}

// Entity default particle entity
type Entity struct {
	gorge.TransformComponent
	gorge.ColorableComponent
	Component
}

var _ render.Renderable = &Entity{}
