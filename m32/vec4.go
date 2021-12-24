package m32

import (
	"math"
	"math/rand"
)

type (
	// Vec4 float32 array of size 4 with vector operation methods.
	Vec4 [4]float32
	// Color alias of vec4
)

// Color returns a vec4 with the colors value
// if the number of params is:
// 0 - returns a zero Vec4
// 1 - returns a Vec4 with all 3 elements set to param 0 and alpha 1
// 2 - returns a Vec4 with all 3 elements set to param 0 and alpha to param 1
// 3 - returns a Vec4 with all 3 elements set to params and alpha 1
// 4 - returns a Vec4 with all 4 elements set to params.
func Color(v ...float32) Vec4 {
	switch len(v) {
	case 0:
		return Vec4{}
	case 1:
		return Vec4{v[0], v[0], v[0], 1}
	case 2:
		return Vec4{v[0], v[0], v[0], v[1]}
	case 3:
		return Vec4{v[0], v[1], v[2], 1}
	default:
		return Vec4{v[0], v[1], v[2], v[3]}
	}
}

// Len returns the length of the vec4.
func (v Vec4) Len() float32 {
	return float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2] + v[3]*v[3])))
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
func (v Vec4) Dot(v2 Vec4) float32 {
	return v[0]*v2[0] + v[1]*v2[1] + v[2]*v2[2] + v[3]*v2[3]
}

// Mul returns a new vec4 with the elements multiplied by c.
func (v Vec4) Mul(c float32) Vec4 {
	return Vec4{v[0] * c, v[1] * c, v[2] * c, v[3] * c}
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
func (v Vec4) Lerp(b Vec4, t float32) Vec4 {
	return Vec4{
		v[0] + t*(b[0]-v[0]),
		v[1] + t*(b[1]-v[1]),
		v[2] + t*(b[2]-v[2]),
		v[3] + t*(b[3]-v[3]),
	}
}

// V4Rand returns a vec4 with random values between [0,1]
func V4Rand() Vec4 {
	return Vec4{
		rand.Float32(),
		rand.Float32(),
		rand.Float32(),
		rand.Float32(),
	}
}
