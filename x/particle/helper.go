package particle

import (
	"math/rand"

	"github.com/stdiopt/gorge/math/gm"
)

func Fixed[T any](v T) func() T {
	return func() T {
		return v
	}
}

func RangeF32(min, max float32) float32 {
	return min + (max-min)*float32(rand.Float64())
}

func RangeVec3(min, max gm.Vec3) gm.Vec3 {
	return gm.Vec3{
		min[0] + (max[0]-min[0])*float32(rand.Float64()),
		min[1] + (max[1]-min[1])*float32(rand.Float64()),
		min[2] + (max[2]-min[2])*float32(rand.Float64()),
	}
}
