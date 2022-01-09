package particle

import (
	"math/rand"

	"github.com/stdiopt/gorge/m32"
)

func Fixed[T any](v T) func() T {
	return func() T {
		return v
	}
}

func RangeF32(min, max float32) func() float32 {
	return func() float32 {
		return min + (max-min)*rand.Float32()
	}
}

func RangeVec3(min, max m32.Vec3) func() m32.Vec3 {
	return func() m32.Vec3 {
		return m32.Vec3{
			min[0] + (max[0]-min[0])*rand.Float32(),
			min[1] + (max[1]-min[1])*rand.Float32(),
			min[2] + (max[2]-min[2])*rand.Float32(),
		}
	}
}
