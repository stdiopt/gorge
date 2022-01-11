package particle

import (
	"math/rand"
	"runtime"
	"sync"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/primitive"
	"github.com/stdiopt/gorge/static"
)

type controller[T any] interface {
	Init(*T)
	Update(*T, float32)
}

// Particle container

// Type is any but it should be a type where pointer implements particle
type Generator[T any] struct {
	Controller controller[T]
	// CreateFunc func(*T)
	// InitFunc   func(*T)
	// UpdateFunc func(*T, float32)

	particles  []T
	renderable gorge.RenderableComponent
	totTime    float32
}

func (g *Generator[T]) destroy(gg *gorge.Context) {
	for i := range g.particles {
		gg.Remove(&g.particles[i])
	}
	g.particles = nil
}

func (g *Generator[T]) init(gg *gorge.Context, em emitter) {
	count := em.Emitter().Count

	mat := em.Emitter().Material
	if mat == nil {
		mat = gorge.NewShaderMaterial(static.Shaders.UnlitAdditive)
		mat.Queue = 1000
		mat.DoubleSided = true
		mat.Depth = gorge.DepthNone
		mat.Blend = gorge.BlendOneOne
		mat.DisableShadow = true
	}
	mesh := em.Emitter().Mesh
	if mesh == nil {
		mesh = primitive.NewPoly(3)
	}
	g.renderable.SetMaterial(mat)
	g.renderable.SetMesh(mesh)

	// Reset all
	g.particles = make([]T, count)
	for i := range g.particles {
		p := &g.particles[i]
		pc := any(p).(particle).Particle()
		pc.RenderableComponent = &g.renderable

		if creator, ok := g.Controller.(interface{ Create(*T) }); ok {
			creator.Create(p)
		}
		if g.Controller != nil {
			g.Controller.Init(p)
		}
		gg.Add(p)
	}
}

func (g *Generator[T]) update(gg *gorge.Context, em emitter, dt float32) {
	g.totTime += dt
	ec := em.Emitter()
	if ec.Count != len(g.particles) {
		g.destroy(gg)
		g.init(gg, em)
	}

	if ec.Material != nil {
		g.renderable.SetMaterial(ec.Material)
	}
	if ec.Mesh != nil {
		g.renderable.SetMesh(ec.Mesh)
	}

	// New initializations
	newPerFrame := ec.Rate * dt
	numNewParticles := uint32(newPerFrame)
	if rand.Float32() < newPerFrame-float32(numNewParticles) {
		numNewParticles++
	}

	created := uint32(0)
	lifeParticles := uint32(0)
	// Generate particles
	for i := range g.particles {
		p := &g.particles[i]
		pp := any(p).(particle)
		pc := pp.Particle()
		if pc.enabled {
			lifeParticles++
		}
		if !pc.enabled && created < numNewParticles && ec.Enabled {
			g.initParticle(em, pp)
			if g.Controller != nil {
				g.Controller.Init(p)
			}
			pc.enabled = true
			created++
		}
	}

	camT := ec.Camera.Transform()
	forward := camT.Forward().Normalize()
	camQuat := camT.Mat4().Quat()

	NSplit := runtime.NumCPU()
	sz := len(g.particles) / NSplit
	wg := sync.WaitGroup{}
	wg.Add(NSplit)
	for i := 0; i < NSplit; i++ {
		go func(i int) {
			off := i * sz
			defer wg.Done()
			particles := g.particles[off : off+sz]
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
				if g.Controller != nil {
					g.Controller.Update(p, dt)
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

// initParticle initializes the particle
func (c *Generator[T]) initParticle(em emitter, p particle) {
	const origin = 0.2
	ec := em.Emitter()

	pc := p.Particle()
	pc.Age = 0
	pc.Life = 1
	t := p.Transform()
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
