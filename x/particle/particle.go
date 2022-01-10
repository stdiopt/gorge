package particle

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/render"
)

type Component struct {
	*gorge.RenderableComponent
	enabled  bool
	Age      float32
	Life     float32
	Rotation float32
}

func (p *Component) Particle() *Component { return p }

func (p *Component) RenderDisable() bool {
	return !p.enabled
}

// Single particle
type Entity struct {
	Component
	gorge.TransformComponent
	gorge.ColorableComponent
}

var _ render.Renderable = &Entity{}
