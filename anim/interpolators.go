package anim

import "github.com/stdiopt/gorge/m32"

// Funcf32 returns an interpolator which triggers the specific function param.
func Funcf32(fn func(v float32)) InterpolatorFunc {
	return func(a, b interface{}, dt float32) {
		fa := a.(float32)
		fb := b.(float32)

		fn(m32.Lerp(fa, fb, dt))
	}
}

// Float32 returns an interpolator that will change the value pointed by the param.
func Float32(f *float32) InterpolatorFunc {
	return func(a, b interface{}, dt float32) {
		fa := a.(float32)
		fb := b.(float32)

		*f = m32.Lerp(fa, fb, dt)
	}
}

// Vec3 interpolate a vec3 pointer
func Vec3(v *m32.Vec3) InterpolatorFunc {
	return func(a, b interface{}, dt float32) {
		va := a.(m32.Vec3)
		vb := b.(m32.Vec3)

		*v = va.Lerp(vb, dt)
	}
}

// Vec3 interpolate a vec3 pointer
func FuncVec3(fn func(m32.Vec3)) InterpolatorFunc {
	return func(a, b interface{}, dt float32) {
		va := a.(m32.Vec3)
		vb := b.(m32.Vec3)

		fn(va.Lerp(vb, dt))
	}
}

// Vec4 interpolate a vec4 pointer
func Vec4(v *m32.Vec4) InterpolatorFunc {
	return func(a, b interface{}, dt float32) {
		va := a.(m32.Vec4)
		vb := b.(m32.Vec4)

		*v = va.Lerp(vb, dt)
	}
}

// Quat spherical interpolates a Quat
func Quat(v *m32.Quat) InterpolatorFunc {
	return func(a, b interface{}, dt float32) {
		qa := a.(m32.Quat)
		qb := b.(m32.Quat)

		*v = qa.Slerp(qb, dt).Normalize()
	}
}
