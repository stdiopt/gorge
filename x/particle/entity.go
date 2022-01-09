package particle

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/render"
)

type ParticleComponent struct {
	enabled bool
	age     float32
	life    float32
}

func (p *ParticleComponent) Particle() *ParticleComponent { return p }

func (p *ParticleComponent) RenderDisable() bool {
	return !p.enabled
}

// Single particle
type Entity struct {
	ParticleComponent
	gorge.TransformComponent
	gorge.ColorableComponent
	*gorge.RenderableComponent
}

var _ render.Renderable = &Entity{}
