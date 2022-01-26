package gm

import (
	"math"
)

type (
	// Vec4 Float array of size 4 with vector operation methods.
	Vec4 [4]Float
	// Color alias of vec4
)

func V4(v ...Float) Vec4 {
	switch len(v) {
	case 0:
		return Vec4{}
	case 1:
		return Vec4{v[0], v[0], v[0], v[0]}
	case 2:
		return Vec4{v[0], v[1], v[0], v[1]}
	case 3:
		return Vec4{v[0], v[1], v[2], v[2]}
	default:
		return Vec4{v[0], v[1], v[2], v[3]}

	}
}

// Len returns the length of the vec4.
func (v Vec4) Len() Float {
	return Float(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2] + v[3]*v[3])))
}

// Add returns the addition of each element with v2 in a new vec4.
func (v Vec4) Add(v2 Vec4) Vec4 {
	return Vec4{v[0] + v2[0], v[1] + v2[1], v[2] + v2[2], v[3] + v2[3]}
}

// Sub returns the subtraction of each element with v2 in a new vec4.
func (v Vec4) Sub(v2 Vec4) Vec4 {
	return Vec4{v[0] - v2[0], v[1] - v2[1], v[2] - v2[2], v[3] - v2[3]}
}

// Normalize returns the normalized vec4.
func (v Vec4) Normalize() Vec4 {
	l := 1.0 / v.Len()
	return Vec4{v[0] * l, v[1] * l, v[2] * l, v[3] * l}
}

// Dot returns the dot product of v and v2.
func (v Vec4) Dot(v2 Vec4) Float {
	return v[0]*v2[0] + v[1]*v2[1] + v[2]*v2[2] + v[3]*v2[3]
}

// Mul returns a new vec4 with the elements multiplied by c.
func (v Vec4) Mul(c Float) Vec4 {
	return Vec4{v[0] * c, v[1] * c, v[2] * c, v[3] * c}
}

// MulVec4 multiplies each element of v with each element of v2.
func (v Vec4) MulVec4(v2 Vec4) Vec4 {
	return Vec4{v[0] * v2[0], v[1] * v2[1], v[2] * v2[2], v[3] * v2[3]}
}

// Vec2 returns a vec2 based on first and second element.
func (v Vec4) Vec2() Vec2 {
	return Vec2{v[0], v[1]}
}

// Vec3 returns a vec3 based on the first elements of vec4.
func (v Vec4) Vec3() Vec3 {
	return Vec3{v[0], v[1], v[2]}
}

// Equal returns true if v and v2 are approximately equal.
func (v Vec4) Equal(v2 Vec4) bool {
	return FloatEqual(v[0], v2[0]) &&
		FloatEqual(v[1], v2[1]) &&
		FloatEqual(v[2], v2[2]) &&
		FloatEqual(v[3], v2[3])
}

// Lerp Linear interpolation between 2 vecs
func (v Vec4) Lerp(b Vec4, t Float) Vec4 {
	return Vec4{
		v[0] + t*(b[0]-v[0]),
		v[1] + t*(b[1]-v[1]),
		v[2] + t*(b[2]-v[2]),
		v[3] + t*(b[3]-v[3]),
	}
}
