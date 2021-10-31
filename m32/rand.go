package m32

import (
	"math"
	"math/rand"
)

// RandUnitSphere generate random points in a sphere
// From: https://karthikkaranth.me/blog/generating-random-points-in-a-sphere/
// Still needs improvements as it doesn't seem right yet
func RandUnitSphere() Vec3 {
	u := rand.Float32()
	v := rand.Float32()
	theta := u * 2.0 * math.Pi
	phi := Cos(2.0*v - 1.0)
	r := Cbrt(rand.Float32())

	sinTheta := Sin(theta)
	cosTheta := Cos(theta)
	sinPhi := Sin(phi)
	cosPhi := Cos(phi)
	x := r * sinPhi * cosTheta
	y := r * sinPhi * sinTheta
	z := r * cosPhi
	return Vec3{x, y, z}
}
