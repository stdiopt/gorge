package m32

import "math"

// Vec2 vector2
type Vec2 [2]float32

// Len returns the len of v.
func (v Vec2) Len() float32 {
	return float32(math.Hypot(float64(v[0]), float64(v[1])))
}

// Add sums the elements with v2 and returns a new vec2.
func (v Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v[0] + v2[0], v[1] + v2[1]}
}

// Sub returns a new vec2 the subtraction with v2.
func (v Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{v[0] - v2[0], v[1] - v2[1]}
}

// Mul returns a new vec2 with the multiplication of each element with c.
func (v Vec2) Mul(c float32) Vec2 {
	return Vec2{v[0] * c, v[1] * c}
}

// Vec3 returns a vec3 with the extra value z.
func (v Vec2) Vec3(z float32) Vec3 {
	return Vec3{v[0], v[1], z}
}

// Vec4 returns a vec4 with the extra values z and w.
func (v Vec2) Vec4(z, w float32) Vec4 {
	return Vec4{v[0], v[1], z, w}
}

// V2Abs returns the vec2 with abs values for each element.
func V2Abs(a Vec2) Vec2 {
	return Vec2{
		Abs(a[0]),
		Abs(a[1]),
	}
}

// V2Lerp Linear interpolation between 2 vecs2.
func V2Lerp(a, b Vec2, t float32) Vec2 {
	return Vec2{
		a[0] + t*(b[0]-a[0]),
		a[1] + t*(b[1]-a[1]),
	}
}

// V2Clamp clamps the vec2 elements to specific min and max floats.
func V2Clamp(a Vec2, min, max float32) Vec2 {
	return Vec2{
		Clamp(a[0], min, max),
		Clamp(a[1], min, max),
	}
}
