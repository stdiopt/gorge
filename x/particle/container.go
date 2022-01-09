package particle

import (
	"math/rand"

	"github.com/stdiopt/gorge"
)

// Particle container

type particle interface {
	Transform() *gorge.TransformComponent
	Colorable() *gorge.ColorableComponent
	Renderable() *gorge.RenderableComponent
	Particle() *ParticleComponent
}

type typedParticle interface {
	init(g *gorge.Context, count int)
	update(em emitter, dt float32)
	destroy(g *gorge.Context)
}

type Container[T any] struct {
	particles []T
}

func (c *Container[T]) destroy(g *gorge.Context) {
	for _, p := range c.particles {
		g.Remove(&p)
	}
}

func (c *Container[T]) init(g *gorge.Context, count int) {
	if len(c.particles) == count {
		return
	}

	// Reset all
	c.particles = make([]T, count)
	for i := range c.particles {
		pc := c.particles[i].Particle()
		pc.age = 0
		pc.life = 1
	}
}

func (c *Container[T]) initParticle(em emitter, i int) {
	const origin = 0.2
	ec := em.Emitter()

	pc := c.particles[i].Particle()
	pc.age = 0
	pc.life = 1

	if ec.LifeFunc != nil {
		pc.life = ec.LifeFunc()
	}

	// p.lifeScale = c.LifeScale[0] + rand.Float32()*(c.LifeScale[1]-c.LifeScale[0])
	// p.rotFactor = (-1 + rand.Float32()*2)

	// On init only
	// p.RenderableComponent = em.Renderable
	c.particles[i].Colorable().SetColor(
		.5+rand.Float32()*0.5,
		.5+rand.Float32()*0.5,
		.5+rand.Float32()*0.5,
		.2,
	)
	t := c.particles[i].Transform()
	if ec.Local {
		// p.TransformComponent = *gorge.NewTransformComponent()
		t.SetParent(em)
	} else {
		t.SetParent(nil)
		t.Position = em.Transform().WorldPosition()
		t.Rotation = em.Transform().WorldRotation()
		t.Scale = em.Transform().Scale
	}
	t.Translate(
		2*(rand.Float32()-.5)*origin,
		2*(rand.Float32()-.5)*origin,
		2*(rand.Float32()-.5)*origin,
	)
	t.SetScale(rand.Float32())
}

func (c *Container[T]) update(em emitter, dt float32) {
	ec := em.Emitter()

	// New initializations
	newPerFrame := ec.Rate * dt
	numNewParticles := int(newPerFrame)
	if rand.Float32() < newPerFrame-float32(numNewParticles) {
		numNewParticles++
	}

	created := 0

	pp := c.particles
	for i := range pp {
		pc := pp[i].Particle()
		if !pc.enabled && created < numNewParticles {
			c.initParticle(em, i)
			pc.enabled = true
			created++
		}

		pc.age += dt
		// log.Println("p.Life:", p.life)
		if pc.age >= pc.life {
			pc.enabled = false
			continue
		}
		lifeStage := pc.age / pc.life
		t := pp[i].Transform()
		if ec.TranslateFunc != nil {
			t.Translatev(ec.TranslateFunc(lifeStage))
		}
		if ec.ScaleFunc != nil {
			t.SetScalev(ec.ScaleFunc(lifeStage))
		}

		if ec.ColorFunc != nil {
			pp[i].Colorable().SetColorv(ec.ColorFunc(lifeStage))
		}

		// p.rot += p.age * p.rotFactor * p.life * 10

		// This might be something to handle
		if ec.Camera != nil {
			camT := ec.Camera.Transform()
			// forward := camT.Forward().Normalize()
			// axisAngle := m32.QAxisAngle(forward, p.rot)
			// p.SetRotation(axisAngle.Mul(camT.Mat4().Quat()))
			t.SetRotation(camT.Mat4().Quat())
		}
	}
}
