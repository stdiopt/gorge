package particle

import (
	"github.com/stdiopt/gorge"
)

type particle interface {
	Transform() *gorge.TransformComponent
	Colorable() *gorge.ColorableComponent
	Particle() *Component
}

// Used in container.
type generator interface {
	init(g *gorge.Context, em emitter)
	update(g *gorge.Context, em emitter, dt float32)
	destroy(g *gorge.Context)
}

// Emitter will emit particles based on the given parameters.
type EmitterComponent struct { // Component
	Camera gorge.Transformer
	*gorge.Mesh
	*gorge.Material

	Enabled bool
	Local   bool
	Count   int
	Step    float32
	Rate    float32 // Number of particles per frame

	// tracked particles
	Generator generator
}

// NewEmitterComponent creates a new emitter component with defaults.
func NewEmitterComponent[T any]() *EmitterComponent {
	return &EmitterComponent{
		Enabled:   true,
		Count:     1000,
		Rate:      100,
		Step:      0.016,
		Generator: &Generator[T]{},
	}
}

func (c *EmitterComponent) init(g *gorge.Context, em emitter) {
	if c.Generator == nil {
		// This breaks stuff
		c.Generator = &Generator[Entity]{}
	}
	c.Generator.init(g, em)
}

func (c *EmitterComponent) update(g *gorge.Context, em emitter, dt float32) {
	if c.Step <= 0 {
		c.Generator.update(g, em, dt)
		return
	}
	c.Generator.update(g, em, c.Step)
}

func (c *EmitterComponent) destroy(g *gorge.Context) {
	c.Generator.destroy(g)
}

// Emitter implements emitter component.
func (c *EmitterComponent) Emitter() *EmitterComponent { return c }

// SetCamera particles will turn towards the camera.
func (c *EmitterComponent) SetCamera(t gorge.Transformer) {
	c.Camera = t
}

// SetMesh sets the particle mesh, defaults to triangle.
func (c *EmitterComponent) SetMesh(m *gorge.Mesh) {
	c.Mesh = m
}

// SetMaterial sets the particle material defaults to unlit additive.
func (c *EmitterComponent) SetMaterial(m *gorge.Material) {
	c.Material = m
}

// SetEnabled enables or disables the emitter.
func (c *EmitterComponent) SetEnabled(b bool) {
	c.Enabled = b
}

// SetLocal will generate particles relative to local emitter transform.
func (c *EmitterComponent) SetLocal(b bool) {
	c.Local = b
}

// SetCount sets the particle count, changing this will reset the emitter.
func (c *EmitterComponent) SetCount(i int) {
	c.Count = i
}

// SetRate sets the number of particles generated per frame.
func (c *EmitterComponent) SetRate(f float32) {
	c.Rate = f
}
