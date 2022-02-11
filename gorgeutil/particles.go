package gorgeutil

import (
	"time"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/anim"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/primitive"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/x/particle"
)

// ParticleEntity is a particle entity.
type ParticleEntity struct {
	gorge.TransformComponent
	gorge.ColorableComponent
	particle.Component

	speed    gm.Vec3
	rotSpeed float32

	ColorAnim *anim.Channel[gm.Vec4]
	ScaleAnim *anim.Channel[gm.Vec3]
}

// ParticleEmitter is a particle emitter and a controller.
type ParticleEmitter struct {
	gorge.TransformComponent
	particle.EmitterComponent

	rr       *gm.Rand
	Life     float32
	Friction float32
	Gravity  gm.Vec3
	Spread   float32
	SpeedMul float32
	Rotation float32
	Color    *anim.Channel[gm.Vec4]
	Scale    *anim.Channel[gm.Vec3]
}

func NewParticleEmitter() *ParticleEmitter {
	return &ParticleEmitter{
		TransformComponent: *gorge.NewTransformComponent(),
		EmitterComponent: particle.EmitterComponent{
			Mesh: primitive.NewPlane(primitive.PlaneDirZ),
			Material: func() *gorge.Material {
				mat := gorge.NewShaderMaterial(static.Shaders.UnlitAdditive)
				mat.DoubleSided = false
				mat.Blend = gorge.BlendOneOne
				mat.Depth = gorge.DepthRead
				mat.DisableShadow = true
				mat.Queue = 4000
				return mat
			}(),
			Enabled:   true,
			Count:     1000,
			Step:      0,
			Rate:      100,
			Particles: &particle.Particles[ParticleEntity]{},
		},
		rr:       gm.NewRand(time.Now().UnixNano()),
		Life:     2,
		Friction: 0.99,
		Spread:   0.2,
		SpeedMul: 2,
		Rotation: 0.1,
	}
}

func (c *ParticleEmitter) InitParticle(p *ParticleEntity) {
	wantDir := gm.Up()
	p.Translatev(c.rr.SphereSurface())
	p.speed = c.rr.Cone(wantDir, c.Spread).Mul(c.SpeedMul)
	p.rotSpeed = (c.rr.Float32() * c.Rotation) * 0.01
	p.Life = c.Life
	if c.Color != nil {
		p.ColorAnim = c.Color.Clone()
	} else {
		p.ColorAnim = anim.NewChannelWithKeys(anim.Vec4, map[float32]gm.Vec4{
			0: {.4, 0, 0, 1},
			1: {.2, .2, 0, 0},
		})
	}
	if c.Scale != nil {
		p.ScaleAnim = c.Scale.Clone()
	} else {
		p.ScaleAnim = anim.NewChannelWithKeys(anim.Vec3, map[float32]gm.Vec3{
			0: {.2, .2, .2},
			1: {0, 0, 0},
		})
	}
}

func (c *ParticleEmitter) UpdateParticle(p *ParticleEntity, dt float32) {
	l := p.Age / p.Life

	p.speed = p.speed.Add(c.Gravity.Mul(dt))
	p.Translatev(p.speed.Mul(dt))
	// We can use color directly from controller but per particle gives a fine fx
	// when changing parameters
	if p.ColorAnim != nil {
		p.Color = p.ColorAnim.Get(l)
	}
	if p.ScaleAnim != nil {
		p.Scale = p.ScaleAnim.Get(l)
	}
	// ground collision
	/*if p.Position[1] < 0 {
		p.Position[1] = 0
		p.speed[1] *= -.5
		p.speed = p.speed.Mul(1 - c.Friction)
	}*/
	p.Rotation += p.rotSpeed * (1 - l)
}

func AddParticleEmitter(a Contexter) *ParticleEmitter {
	em := NewParticleEmitter()
	a.Add(em)
	return em
}
