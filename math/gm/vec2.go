package gm

import "math"

// Vec2 vector2
type Vec2 [2]Float

func V2(v ...Float) Vec2 {
	switch len(v) {
	case 0:
		return Vec2{}
	case 1:
		return Vec2{v[0], v[0]}
	default:
		return Vec2{v[0], v[1]}
	}
}

// Len returns the len of v.
func (v Vec2) Len() Float {
	return Float(math.Hypot(float64(v[0]), float64(v[1])))
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
func (v Vec2) Mul(c Float) Vec2 {
	return Vec2{v[0] * c, v[1] * c}
}

// MulVec2 returns a new Vec2 with the multiplication by the corresponding
// elements of v2.
func (v Vec2) MulVec2(v2 Vec2) Vec2 {
	return Vec2{
		v[0] * v2[0],
		v[1] * v2[1],
	}
}

// Vec3 returns a vec3 with the extra value z.
func (v Vec2) Vec3(z Float) Vec3 {
	return Vec3{v[0], v[1], z}
}

// Vec4 returns a vec4 with the extra values z and w.
func (v Vec2) Vec4(z, w Float) Vec4 {
	return Vec4{v[0], v[1], z, w}
}

// Clamp clamps the vec2 elements to specific min and max floats.
func (v Vec2) Clamp(min, max Vec2) Vec2 {
	return Vec2{
		Clamp(v[0], min[0], max[0]),
		Clamp(v[1], min[1], max[1]),
	}
}

// Abs returns the vec2 with abs values for each element.
func (v Vec2) Abs() Vec2 {
	return Vec2{
		Abs(v[0]),
		Abs(v[1]),
	}
}

// Lerp Linear interpolation between 2 vecs2.
func (v Vec2) Lerp(b Vec2, t Float) Vec2 {
	return Vec2{
		v[0] + t*(b[0]-v[0]),
		v[1] + t*(b[1]-v[1]),
	}
}

func (v Vec2) X() Float { return v[0] }
func (v Vec2) Y() Float { return v[1] }
