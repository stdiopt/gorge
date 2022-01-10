package m32

import (
	"math"
	"math/rand"
)

type Rand struct {
	rand.Rand
}

// UnitSphere generate random points in a sphere
// From: https://karthikkaranth.me/blog/generating-random-points-in-a-sphere/
// Still needs improvements as it doesn't seem right yet
func (r Rand) UnitSphere() Vec3 {
	u := rand.Float32()
	v := rand.Float32()

	theta := u * 2.0 * math.Pi

	phi := Acos(2.0*v - 1.0)
	rr := Cbrt(rand.Float32())

	sinTheta := Sin(theta)
	cosTheta := Cos(theta)
	sinPhi := Sin(phi)
	cosPhi := Cos(phi)
	x := rr * sinPhi * cosTheta
	y := rr * sinPhi * sinTheta
	z := rr * cosPhi
	return Vec3{x, y, z}
}

func (r Rand) SphereSurface() Vec3 {
	u := rand.Float32()
	v := rand.Float32()

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

func (r Rand) Vec3() Vec3 {
	return Vec3{
		r.Float32(),
		r.Float32(),
		r.Float32(),
	}
}

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
