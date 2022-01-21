package anim

import (
	"github.com/stdiopt/gorge/math/gm"
)

func Float32(a, b float32, dt float32) float32 {
	return gm.Lerp(a, b, dt)
}

func Vec3(a, b gm.Vec3, dt float32) gm.Vec3 {
	return a.Lerp(b, dt)
}

// Vec4 interpolate a vec4 pointer
func Vec4(a, b gm.Vec4, dt float32) gm.Vec4 {
	return a.Lerp(b, dt)
}

// Quat spherical interpolates a Quat
func Quat(a, b gm.Quat, dt float32) gm.Quat {
	return a.Slerp(b, dt).Normalize()
}
