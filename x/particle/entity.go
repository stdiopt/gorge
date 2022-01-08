package particle

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/render"
)

// Single particle
type Entity struct {
	gorge.TransformComponent
	gorge.ColorableComponent
	*gorge.RenderableComponent

	// ParticleComponent
	enabled   bool
	animTick  float32
	life      float32
	curTex    int
	lifeScale float32
	rot       float32
	rotFactor float32
}

func (p *Entity) RenderDisable() bool {
	return !p.enabled
}

var _ render.Renderable = &Entity{}
