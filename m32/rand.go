package m32

import (
	"math"
	"math/rand"
	"sync"
)

type lockedSource struct {
	mu sync.Mutex
	rand.Source
}

func (s *lockedSource) Int63() (n int64) {
	s.mu.Lock()
	n = s.Source.Int63()
	s.mu.Unlock()
	return
}

func (s *lockedSource) Seed(seed int64) {
	s.mu.Lock()
	s.Source.Seed(seed)
	s.mu.Unlock()
}

var globalRand = rand.New(&lockedSource{Source: rand.NewSource(1)})

func NewRand(seed int64) *Rand {
	src := &lockedSource{
		Source: rand.NewSource(seed),
	}
	return &Rand{Rand: rand.New(src)}
}

// Rand helper function that handlers vec3, vec4 etc..
type Rand struct {
	*rand.Rand
}

func (r Rand) rand() *rand.Rand {
	if r.Rand == nil {
		return globalRand
	}
	return r.Rand
}

// Float32 returns a float32 with random values [0, 1] using inner rand
// generator. if the inner source is nil it will use the global one.
func (r Rand) Float32() float32 {
	return r.rand().Float32()
}

// UnitSphere generate random points in a sphere
// From: https://karthikkaranth.me/blog/generating-random-points-in-a-sphere/
// Still needs improvements as it doesn't seem right yet
func (r Rand) UnitSphere() Vec3 {
	u := r.Float32()
	v := r.Float32()

	theta := u * 2.0 * math.Pi

	phi := Acos(2.0*v - 1.0)
	rr := Cbrt(r.Float32())

	sinTheta := Sin(theta)
	cosTheta := Cos(theta)
	sinPhi := Sin(phi)
	cosPhi := Cos(phi)
	x := rr * sinPhi * cosTheta
	y := rr * sinPhi * sinTheta
	z := rr * cosPhi
	return Vec3{x, y, z}
}

// SphereSurface returns a random surface points on a sphere.
func (r Rand) SphereSurface() Vec3 {
	u := r.Float32()
	v := r.Float32()

	theta := u * 2.0 * math.Pi

	phi := Acos(2.0*v - 1.0)
	rr := float32(1) // Cbrt(rand.Float32())

	sinTheta := Sin(theta)
	cosTheta := Cos(theta)
	sinPhi := Sin(phi)
	cosPhi := Cos(phi)
	x := rr * sinPhi * cosTheta
	y := rr * sinPhi * sinTheta
	z := rr * cosPhi
	return Vec3{x, y, z}
}

// Vec3 returns a vec3 with random values [0, 1]
func (r Rand) Vec3() Vec3 {
	return Vec3{
		r.Float32(),
		r.Float32(),
		r.Float32(),
	}
}

// Vec4 returns a vec4 with random values [0, 1]
func (r Rand) Vec4() Vec4 {
	return Vec4{
		r.Float32(),
		r.Float32(),
		r.Float32(),
		r.Float32(),
	}
}

/*
function getPoint() {
    var u = Math.random();
    var v = Math.random();
    var theta = u * 2.0 * Math.PI;
    var phi = Math.acos(2.0 * v - 1.0);
    var r = Math.cbrt(Math.random());
    var sinTheta = Math.sin(theta);
    var cosTheta = Math.cos(theta);
    var sinPhi = Math.sin(phi);
    var cosPhi = Math.cos(phi);
    var x = r * sinPhi * cosTheta;
    var y = r * sinPhi * sinTheta;
    var z = r * cosPhi;
    return {x: x, y: y, z: z};
}
*/
