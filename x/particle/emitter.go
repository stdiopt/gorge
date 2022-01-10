package particle

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
)

func Range(min, max float32) [2]float32 {
	return [2]float32{min, max}
}

// Emitter will emit particles based on the given parameters.
type EmitterComponent struct { // Component
	Camera     gorge.Transformer
	Renderable *gorge.RenderableComponent

	Enabled bool
	Local   bool
	Count   int
	Rate    float32 // Number of particles per second

	LifeFunc      func() float32
	TranslateFunc func(float32) m32.Vec3
	ColorFunc     func(float32) m32.Vec4
	ScaleFunc     func(float32) m32.Vec3

	// tracked particles
	Particles typedParticle
	// particles  []Entity
	count      int
	lastEmited int
}

// bad?!
func (c *EmitterComponent) Emitter() *EmitterComponent { return c }

/*
func update(g *gorge.Context, em emitter, dt float32) {
	c := em.Emitter()
	if c.count < c.Count {
		c.particles = make([]Entity, c.Count)
		for i := range c.particles {
			initParticle(em, &c.particles[i])
			c.particles[i].enabled = false
			g.Add(&c.particles[i])
		}
		c.count = c.Count
	}

	newPerFrame := c.Rate * dt
	numNewParticles := int(newPerFrame)
	if rand.Float32() < newPerFrame-float32(numNewParticles) {
		numNewParticles++
	}

	created := 0
	// do go routines
	for i := range c.particles {
		p := &c.particles[i]
		if !p.enabled && created < numNewParticles {
			initParticle(em, p)
			p.enabled = true
			created++
		}

		p.age += dt
		// log.Println("p.Life:", p.life)
		if p.age >= p.life {
			p.enabled = false
			continue
		}
		lifeStage := p.age / p.life
		if c.TranslateFunc != nil {
			p.Translatev(c.TranslateFunc(lifeStage))
		}
		if c.ScaleFunc != nil {
			p.SetScalev(c.ScaleFunc(lifeStage))
		}

		if c.ColorFunc != nil {
			p.SetColorv(c.ColorFunc(lifeStage))
		}

		// p.rot += p.age * p.rotFactor * p.life * 10

		// This might be something to handle
		if c.Camera != nil {
			camT := c.Camera.Transform()
			// forward := camT.Forward().Normalize()
			// axisAngle := m32.QAxisAngle(forward, p.rot)
			// p.SetRotation(axisAngle.Mul(camT.Mat4().Quat()))
			p.SetRotation(camT.Mat4().Quat())
		}

	}
	// emit particles
}
func initParticle(em emitter, p *Entity) {
	const origin = 0.2
	c := em.Emitter()
	p.age = 0
	p.life = 1
	if c.LifeFunc != nil {
		p.life = c.LifeFunc()
	}

	// p.lifeScale = c.LifeScale[0] + rand.Float32()*(c.LifeScale[1]-c.LifeScale[0])
	// p.rotFactor = (-1 + rand.Float32()*2)

	p.RenderableComponent = c.Renderable
	p.ColorableComponent = *gorge.NewColorableComponent(
		.5+rand.Float32()*0.5,
		.5+rand.Float32()*0.5,
		.5+rand.Float32()*0.5,
		.2,
	)
	if c.Local {
		p.TransformComponent = *gorge.NewTransformComponent()
		p.SetParent(em)
	} else {
		p.Position = em.Transform().WorldPosition()
		p.Rotation = em.Transform().WorldRotation()
		p.Scale = em.Transform().Scale
	}
	p.Translate(
		2*(rand.Float32()-.5)*origin,
		2*(rand.Float32()-.5)*origin,
		2*(rand.Float32()-.5)*origin,
	)
	p.SetScale(rand.Float32())
}
*/
