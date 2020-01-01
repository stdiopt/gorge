// Copyright 2019 Luis Figueiredo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package m32 math and gl math for floats32 it uses go-gl/mathgl/mgl32 for
// certain things
package m32

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type (
	vec2 = mgl32.Vec2
	vec3 = mgl32.Vec3
	vec4 = mgl32.Vec4
	mat3 = mgl32.Mat3
	mat4 = mgl32.Mat4
	quat = mgl32.Quat
)

/*var (
	up      = vec3{0, 1, 0}
	forward = vec3{0, 0, 1}
)*/

// Up returns a up vector
func Up() vec3 { return vec3{0, 1, 0} }

// Down returns a down vector
func Down() vec3 { return vec3{0, -1, 0} }

// Forward returns a vector facing forward
func Forward() vec3 { return vec3{0, 0, 1} }

// Backward returns a vector facing backward
func Backward() vec3 { return vec3{0, 0, -1} }

// Left returns a vector pointing left
func Left() vec3 { return vec3{-1, 0, 0} }

// Right returns a vector pointing left
func Right() vec3 { return vec3{1, 0, 0} }

// Whatever I cast I will put the func here

// Cos casts values to float64 and uses native math.Cos
// to return the cosine of the radian argument x.
func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

// Sin casts values to float64 and uses native math.Sin
// to return the sine of the radian argument x.
func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

// Sincos returns Sin(x), Cos(x).
func Sincos(x float32) (float32, float32) {
	s, c := math.Sincos(float64(x))
	return float32(s), float32(c)
}

// Hypot returns Sqrt(p*p + q*q), taking care to avoid unnecessary overflow and underflow.
func Hypot(x, y float32) float32 {
	return float32(math.Hypot(float64(x), float64(y)))
}

// Atan2 returns the arc tangent of y/x, using the signs of the two to determine
//the quadrant of the return value.
func Atan2(y, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

// Sqrt returns the square root of x.
func Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

// Abs returns the absolute value of x.
func Abs(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

// Copysign returns a value with the magnitude of x and the sign of y.
func Copysign(x, y float32) float32 {
	return float32(math.Copysign(float64(x), float64(y)))
}

// Max returns the greatest value between a or b
func Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// Min returns the lowest value between a or b
func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// Ceil returns the .. ceil
func Ceil(x float32) float32 {
	return float32(math.Ceil(float64(x)))
}

// Limit maintains v between min and max
func Limit(v, mn, mx float32) float32 {
	if v < mn {
		return mn
	} else if v > mx {
		return mx
	}
	return v
}

// Cbrt returns the cube root of x.
func Cbrt(x float32) float32 {
	return float32(math.Cbrt(float64(x)))
}
