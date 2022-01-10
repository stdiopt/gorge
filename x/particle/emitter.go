package particle

import (
	"math/rand"

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
	update(em emitter, dt float32)
	destroy(g *gorge.Context)
}

// Emitter will emit particles based on the given parameters.
type EmitterComponent struct { // Component
	Camera     gorge.Transformer
	Renderable *gorge.RenderableComponent
	Rand       rand.Rand

	Enabled bool
	Local   bool
	Count   int
	Rate    float32 // Number of particles per second

	// tracked particles
	Particles  generator
	count      int
	lastEmited int
}

// bad?!
func (c *EmitterComponent) Emitter() *EmitterComponent { return c }

// Emitter entity will emit particles based on the given parameters when added
// to gorge, particle.System must be in gorge initialization list.
type Emitter struct {
	gorge.TransformComponent
	EmitterComponent
}

/*
func (c *EmitterComponent) destroy(g *gorge.Context) {
	for i := range c.particles {
		g.Remove(&c.particles[i])
	}
	c.particles = c.particles[:0]
}

func (c *EmitterComponent) init(g *gorge.Context) {
	count := c.Count
	if len(c.particles) == count {
		return
	}

	// Reset all
	c.particles = make([]T, count)
	for i := range c.particles {
		c.InitFunc(&c.particles[i])
		pc := any(&c.particles[i]).(particle).Particle()
		pc.RenderableComponent = c.Renderable
		g.Add(&c.particles[i])
	}
}

func (c *EmitterComponent[T]) initParticle(em emitter, p particle) {
	// ec := em.Emitter()

	pc := p.Particle()
	pc.Age = 0
	pc.Life = 1

	t := p.Transform()
	if c.Local {
		// p.TransformComponent = *gorge.NewTransformComponent()
		t.SetParent(em)
	} else {
		t.SetParent(nil)
		t.Position = em.Transform().WorldPosition()
		t.Rotation = em.Transform().WorldRotation()
		t.Scale = em.Transform().Scale
	}
}

func (c *EmitterComponent[T]) update(em emitter, dt float32) {
	// New initializations
	newPerFrame := c.Rate * dt
	numNewParticles := int(newPerFrame)
	if rand.Float32() < newPerFrame-float32(numNewParticles) {
		numNewParticles++
	}

	created := 0

	for i := range c.particles {
		// var p particle = (&c.particles[i])
		p := any(&c.particles[i]).(particle)
		pc := p.Particle()
		if !pc.enabled && created < numNewParticles && c.Enabled {
			c.initParticle(em, p)
			if c.InitFunc != nil {
				c.InitFunc(&c.particles[i])
			}

			pc.enabled = true
			created++
		}

		pc.Age += dt
		// log.Println("p.Life:", p.life)
		if pc.Age >= pc.Life {
			pc.enabled = false
			continue
		}
		if c.UpdateFunc != nil {
			c.UpdateFunc(&c.particles[i])
		}

		// This might be something to handle
		// Rotate to camera
		if c.Camera != nil {
			t := p.Transform()
			camT := c.Camera.Transform()
			// forward := camT.Forward().Normalize()
			// axisAngle := m32.QAxisAngle(forward, p.rot)
			// p.SetRotation(axisAngle.Mul(camT.Mat4().Quat()))
			t.SetRotation(camT.Mat4().Quat())
		}
	}
}
*/
