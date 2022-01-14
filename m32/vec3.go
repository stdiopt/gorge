package m32

import (
	"math"
)

// Vec3 an float32 array of size 3 with methods for vector operations.
type Vec3 [3]float32

// Add returns a new Vec3 based on the sum of the param.
func (v Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{v[0] + v2[0], v[1] + v2[1], v[2] + v2[2]}
}

// Sub returns a new Vec3 with subtracting each element of v2.
func (v Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{v[0] - v2[0], v[1] - v2[1], v[2] - v2[2]}
}

// Mul returns a new vec3 based on the multiplication of vec3*scalar.
func (v Vec3) Mul(c float32) Vec3 {
	return Vec3{v[0] * c, v[1] * c, v[2] * c}
}

// MulVec3 multiplies each element by the corresponding element of v2.
func (v Vec3) MulVec3(v2 Vec3) Vec3 {
	return Vec3{v[0] * v2[0], v[1] * v2[1], v[2] * v2[2]}
}

// Len returns the vec3 len.
func (v Vec3) Len() float32 {
	return float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))
}

// Normalize returns a new normalized vec3.
func (v Vec3) Normalize() Vec3 {
	l := 1.0 / v.Len()
	return Vec3{v[0] * l, v[1] * l, v[2] * l}
}

// Cross returns a new vec3 with the cross product of v with v2.
func (v Vec3) Cross(v2 Vec3) Vec3 {
	return Vec3{v[1]*v2[2] - v[2]*v2[1], v[2]*v2[0] - v[0]*v2[2], v[0]*v2[1] - v[1]*v2[0]}
}

// Dot returns the dot product of v with v2
func (v Vec3) Dot(v2 Vec3) float32 {
	return v[0]*v2[0] + v[1]*v2[1] + v[2]*v2[2]
}

// Reflect returns a reflected vec3 with n
func (v Vec3) Reflect(n Vec3) Vec3 {
	return v.Sub(n.Mul(2 * v.Dot(n)))
}

// Vec2 returns a Vec2 with the first 2 elements of v.
func (v Vec3) Vec2() Vec2 {
	return Vec2{v[0], v[1]}
}

// Vec4 returns a vec4 with the new element w.
func (v Vec3) Vec4(w float32) Vec4 {
	return Vec4{v[0], v[1], v[2], w}
}

// Clamp returns a vec3 with the v elements clamped to min and max.
func (v Vec3) Clamp(min, max Vec3) Vec3 {
	return Vec3{
		Clamp(v[0], min[0], max[0]),
		Clamp(v[1], min[1], max[1]),
		Clamp(v[2], min[2], max[2]),
	}
}

// Lerp Linear interpolation between 2 vecs
func (v Vec3) Lerp(b Vec3, t float32) Vec3 {
	return Vec3{
		v[0] + t*(b[0]-v[0]),
		v[1] + t*(b[1]-v[1]),
		v[2] + t*(b[2]-v[2]),
	}
}

// Up returns a up vector
func Up() Vec3 { return Vec3{0, 1, 0} }

// Down returns a down vector
func Down() Vec3 { return Vec3{0, -1, 0} }

// Forward returns a vector facing forward
func Forward() Vec3 { return Vec3{0, 0, -1} }

// Backward returns a vector facing backward
func Backward() Vec3 { return Vec3{0, 0, 1} }

// Left returns a vector pointing left
func Left() Vec3 { return Vec3{-1, 0, 0} }

// Right returns a vector pointing left
func Right() Vec3 { return Vec3{1, 0, 0} }
