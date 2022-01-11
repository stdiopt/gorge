package particle

import (
	"github.com/stdiopt/gorge"
)

type Component struct {
	*gorge.RenderableComponent
	Age      float32
	Life     float32
	Rotation float32

	enabled bool
}

func (p *Component) Particle() *Component { return p }

func (p *Component) RenderDisable() bool {
	return !p.enabled
}
