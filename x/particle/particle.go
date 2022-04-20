package particle

import (
	"math/rand"
	"runtime"
	"sync"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/primitive"
	"github.com/stdiopt/gorge/static"
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

// particler used to have a T constraint to allow the pointers on methods.
type particler interface {
	gorge.Transformer
	Particle() *Component
}

type (
	pCreator[T any]     interface{ CreateParticle(*T) }
	pInitializer[T any] interface{ InitParticle(*T) }
	pUpdater[T any]     interface{ UpdateParticle(*T, float32) }
)

// Particle container

// Type is any but it should be a type where pointer implements particle
type Particles[T any] struct {
	// Controller controller[T]

	particles  []T
	renderable gorge.RenderableComponent
	totTime    float32

	createFunc func(*T)
	initFunc   func(*T)
	updateFunc func(*T, float32)
}

func (g *Particles[T]) destroy(gg *gorge.Context) {
	for i := range g.particles {
		gg.Remove(&g.particles[i])
	}
	g.particles = nil
}

func (g *Particles[T]) init(gg *gorge.Context, em emitter) {
	if c, ok := em.(pCreator[T]); ok {
		g.createFunc = c.CreateParticle
	}
	if c, ok := em.(pInitializer[T]); ok {
		g.initFunc = c.InitParticle
	}
	if c, ok := em.(pUpdater[T]); ok {
		g.updateFunc = c.UpdateParticle
	}

	count := em.Emitter().Count

	mat := em.Emitter().Material
	if mat == nil {
		mat = gorge.NewShaderMaterial(static.Shaders.UnlitAdditive)
		mat.Queue = 1000
		mat.DoubleSided = true
		mat.Depth = gorge.DepthRead
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
		// TODO: {lpf} Extra conversion for particle until we find solution
		// it was Tp(p).Particle() before
		pc := any(p).(particler).Particle()
		pc.RenderableComponent = &g.renderable

		if g.createFunc != nil {
			g.createFunc(p)
		}
		if g.initFunc != nil {
			g.initFunc(p)
		}

		gg.Add(p)
	}
}

func (g *Particles[T]) update(gg *gorge.Context, em emitter, dt float32) {
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
	if float32(rand.Float64()) < newPerFrame-float32(numNewParticles) {
		numNewParticles++
	}

	created := uint32(0)
	lifeParticles := uint32(0)
	// Generate particles
	for i := range g.particles {
		p := &g.particles[i]
		// TODO: Extra conversion
		pp := any(p).(particler)
		pc := pp.Particle()
		if pc.enabled {
			lifeParticles++
		}
		if !pc.enabled && created < numNewParticles && ec.Enabled {
			g.initParticle(em, p)
			if g.initFunc != nil {
				g.initFunc(p)
			}
			pc.enabled = true
			created++
		}
	}
	var camT *gorge.TransformComponent
	if ec.Camera != nil {
		camT = ec.Camera.Transform()
	} else {
		camT = gorge.NewTransformComponent()
	}
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
				// Extra conversion
				pp := any(p).(particler)
				pc := pp.Particle()
				if !pc.enabled {
					continue
				}

				pc.Age += dt
				if pc.Age >= pc.Life {
					pc.enabled = false
					continue
				}
				if g.updateFunc != nil {
					g.updateFunc(p, dt)
				}
				t := pp.Transform()
				// This might be something to handle
				if ec.Camera != nil {
					axisAngle := gm.QAxisAngle(forward, pc.Rotation)
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
				axisAngle := gm.QAxisAngle(forward, pc.Rotation)
				t.SetRotation(axisAngle.Mul(camQuat))
				// t.SetRotation(camT.Mat4().Quat())
			}
		}
	*/
}

// initParticle initializes the particle
func (c *Particles[T]) initParticle(em emitter, p *T) {
	const origin = 0.2
	ec := em.Emitter()
	// TODO: Extra conversion
	pp := any(p).(particler)
	pc := pp.Particle()
	pc.Age = 0
	pc.Life = 1
	t := pp.Transform()
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
