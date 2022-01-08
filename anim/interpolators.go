package anim

import "github.com/stdiopt/gorge/m32"

func Float32(a, b float32, dt float32) float32 {
	return m32.Lerp(a, b, dt)
}

func Vec3(a, b m32.Vec3, dt float32) m32.Vec3 {
	return a.Lerp(b, dt)
}

// Vec4 interpolate a vec4 pointer
func Vec4(a, b m32.Vec4, dt float32) m32.Vec4 {
	return a.Lerp(b, dt)
}

// Quat spherical interpolates a Quat
func Quat(a, b m32.Quat, dt float32) m32.Quat {
	return a.Slerp(b, dt).Normalize()
}
