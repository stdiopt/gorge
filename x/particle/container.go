package particle

import (
	"math/rand"
	"runtime"
	"sync"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
)

// Particle container

// Type is any but it should be a type where pointer implements particle
type Particles[T any] struct {
	particles  []T
	CreateFunc func(*T)
	InitFunc   func(*T)
	UpdateFunc func(*T)
}

func (c *Particles[T]) destroy(g *gorge.Context) {
	for i := range c.particles {
		g.Remove(&c.particles[i])
	}
	c.particles = nil
}

func (c *Particles[T]) init(g *gorge.Context, em emitter) {
	count := em.Emitter().Count
	if len(c.particles) == count {
		return
	}

	// Reset all
	c.particles = make([]T, count)
	for i := range c.particles {
		p := &c.particles[i]
		pp := any(p).(particle)
		if c.CreateFunc != nil {
			c.CreateFunc(p)
		}
		if c.InitFunc != nil {
			c.InitFunc(p)
		}
		// p := any(&c.particles[i]).(particle)
		// For now because we want to use the same particle type
		pc := pp.Particle()
		pc.Life = 1
		pc.RenderableComponent = em.Emitter().Renderable
		g.Add(p)
	}
}

func (c *Particles[T]) initParticle(em emitter, pu particle) {
	const origin = 0.2
	ec := em.Emitter()

	pc := pu.Particle()
	pc.Age = 0
	pc.Life = 1
	t := pu.Transform()
	if ec.Local {
		// p.TransformComponent = *gorge.NewTransformComponent()
		t.SetParent(em)
	} else {
		t.SetParent(nil)
		t.Position = em.Transform().WorldPosition()
		t.Rotation = em.Transform().WorldRotation()
		t.Scale = em.Transform().Scale
	}
}

func (c *Particles[T]) update(em emitter, dt float32) {
	ec := em.Emitter()

	// New initializations
	newPerFrame := ec.Rate * dt
	numNewParticles := uint32(newPerFrame)
	if rand.Float32() < newPerFrame-float32(numNewParticles) {
		numNewParticles++
	}

	created := uint32(0)
	// Generate particles
	for i := range c.particles {
		p := &c.particles[i]
		pp := any(p).(particle)
		pc := pp.Particle()
		if !pc.enabled && created < numNewParticles && ec.Enabled {
			c.initParticle(em, pp)
			if c.InitFunc != nil {
				c.InitFunc(p)
			}
			pc.enabled = true
			created++
		}
	}

	camT := ec.Camera.Transform()
	forward := camT.Forward().Normalize()
	camQuat := camT.Mat4().Quat()

	NSplit := runtime.NumCPU()
	sz := len(c.particles) / NSplit
	wg := sync.WaitGroup{}
	wg.Add(NSplit)
	for i := 0; i < NSplit; i++ {
		go func(i int) {
			off := i * sz
			defer wg.Done()
			particles := c.particles[off : off+sz]
			for i := range particles {
				p := &particles[i]
				pp := any(p).(particle)
				pc := pp.Particle()
				if !pc.enabled {
					continue
				}

				pc.Age += dt
				if pc.Age >= pc.Life {
					pc.enabled = false
					continue
				}
				if c.UpdateFunc != nil {
					c.UpdateFunc(p)
				}
				t := pp.Transform()
				// This might be something to handle
				if ec.Camera != nil {
					axisAngle := m32.QAxisAngle(forward, pc.Rotation)
					t.SetRotation(axisAngle.Mul(camQuat))
					// t.SetRotation(camT.Mat4().Quat())
				}
			}
		}(i)
	}
	wg.Wait()

	/*
		for i := range c.particles {
			p := any(&c.particles[i]).(particle)
			pc := p.Particle()
			if !pc.enabled && created < numNewParticles && ec.Enabled {
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
			t := p.Transform()
			// This might be something to handle
			if ec.Camera != nil {
				axisAngle := m32.QAxisAngle(forward, pc.Rotation)
				t.SetRotation(axisAngle.Mul(camQuat))
				// t.SetRotation(camT.Mat4().Quat())
			}
		}
	*/
}
