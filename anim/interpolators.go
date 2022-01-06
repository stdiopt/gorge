package anim

import "github.com/stdiopt/gorge/m32"

// Funcf32 returns an interpolator which triggers the specific function param.
func Funcf32(fn func(v float32)) InterpolatorFunc[float32] {
	return func(a, b float32, dt float32) {
		fn(m32.Lerp(a, b, dt))
	}
}

// Float32 returns an interpolator that will change the value pointed by the param.
func Float32(f *float32) InterpolatorFunc[float32] {
	return func(a, b float32, dt float32) {
		*f = m32.Lerp(a, b, dt)
	}
}

// Vec3 interpolate a vec3 pointer
func Vec3(v *m32.Vec3) InterpolatorFunc[m32.Vec3] {
	return func(a, b m32.Vec3, dt float32) {
		*v = a.Lerp(b, dt)
	}
}

// Vec3 interpolate a vec3 pointer
func FuncVec3(fn func(m32.Vec3)) InterpolatorFunc[m32.Vec3] {
	return func(a, b m32.Vec3, dt float32) {
		fn(a.Lerp(b, dt))
	}
}

// Vec4 interpolate a vec4 pointer
func Vec4(v *m32.Vec4) InterpolatorFunc[m32.Vec4] {
	return func(a, b m32.Vec4, dt float32) {
		*v = a.Lerp(b, dt)
	}
}

// Quat spherical interpolates a Quat
func Quat(v *m32.Quat) InterpolatorFunc[m32.Quat] {
	return func(a, b m32.Quat, dt float32) {
		*v = a.Slerp(b, dt).Normalize()
	}
}
