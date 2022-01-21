// Package gm math and gl math for floats32
// Contains code from github.com/go-gl/mathgl/mgl32
// BSD-3 Copyright ©2013 The go-gl Authors. All rights reserved.
package gm

import (
	"math"
)

type Float = float32

// Epsilon float error stuff.
const (
	Epsilon   = 0.000001
	MinNormal = 1.1754943508222875e-38 // 1 / 2**(127 - 1)
	MinValue  = math.SmallestNonzeroFloat32
	MaxValue  = math.MaxFloat32
	Pi        = math.Pi
)

// values from math
var (
	InfPos = Float(math.Inf(1))
	InfNeg = Float(math.Inf(-1))
	NaN    = Float(math.NaN())
)

// FloatEqualThreshold is a utility function to compare floats.
// It's Taken from http://floating-point-gui.de/errors/comparison/
//
// It is slightly altered to not call Abs when not needed.
//
// This differs from FloatEqual in that it lets you pass in your comparison
// threshold, so that you can adjust the comparison value to your specific
// needs
func FloatEqualThreshold(a, b, epsilon Float) bool {
	// Handles the case of inf or shortcuts the loop when no significant
	// error has accumulated
	if a == b {
		return true
	}

	diff := Abs(a - b)
	if a*b == 0 || diff < MinNormal { // If a or b are 0 or both are extremely close to it
		return diff < epsilon*epsilon
	}

	// Else compare difference
	return diff/(Abs(a)+Abs(b)) < epsilon
}

// FloatEqual is a safe utility function to compare floats.
// It's Taken from http://floating-point-gui.de/errors/comparison/
//
// It is slightly altered to not call Abs when not needed.
func FloatEqual(a, b Float) bool {
	return FloatEqualThreshold(a, b, Epsilon)
}

// Mod returns the floating-point remainder of x/y.
// The magnitude of the result is less than y and its
// sign agrees with that of x.
func Mod(x, y Float) Float {
	return Float(math.Mod(float64(x), float64(y)))
}

// Cos casts values to float64 and uses native math.Cos
// to return the cosine of the radian argument x.
func Cos(x Float) Float {
	return Float(math.Cos(float64(x)))
}

// Sin to return the sine of the radian argument x.
func Sin(x Float) Float {
	return Float(math.Sin(float64(x)))
}

// Sincos returns Sin(x), Cos(x).
func Sincos(x Float) (Float, Float) {
	s, c := math.Sincos(float64(x))
	return Float(s), Float(c)
}

// Hypot returns Sqrt(p*p + q*q), taking care to avoid unnecessary overflow and underflow.
func Hypot(x, y Float) Float {
	return Float(math.Hypot(float64(x), float64(y)))
}

// Atan2 returns the arc tangent of y/x, using the signs of the two to determine
// the quadrant of the return value.
func Atan2(y, x Float) Float {
	return Float(math.Atan2(float64(y), float64(x)))
}

// Asin returns the arcsine, in radians, of x.
//
// Special cases are:
//	Asin(±0) = ±0
//	Asin(x) = NaN if x < -1 or x > 1
func Asin(x Float) Float {
	return Float(math.Asin(float64(x)))
}

// Acos returns the arccosine, in radians, of x.
//
// Special case is:
//	Acos(x) = NaN if x < -1 or x > 1
func Acos(x Float) Float {
	return Float(math.Acos(float64(x)))
}

// Tan returns the tangent of the radian argument x.
func Tan(x Float) Float {
	return Float(math.Tan(float64(x)))
}

// Sqrt returns the square root of x.
func Sqrt(x Float) Float {
	return Float(math.Sqrt(float64(x)))
}

// Abs returns the absolute value of x.
func Abs(x Float) Float {
	return Float(math.Abs(float64(x)))
}

// Copysign returns a value with the magnitude of x and the sign of y.
func Copysign(x, y Float) Float {
	return Float(math.Copysign(float64(x), float64(y)))
}

// Max returns the greatest value between a or b
func Max(a, b Float) Float {
	if a > b {
		return a
	}
	return b
}

// Min returns the lowest value between a or b
func Min(a, b Float) Float {
	if a < b {
		return a
	}
	return b
}

// Ceil returns the .. ceil
func Ceil(x Float) Float {
	return Float(math.Ceil(float64(x)))
}

// Clamp v
func Clamp(v, min, max Float) Float {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// Cbrt returns the cube root of x.
func Cbrt(x Float) Float {
	return Float(math.Cbrt(float64(x)))
}

// Lerp Linear interpolation between 2 scalar
func Lerp(v0, v1, t Float) Float {
	return v0 + t*(v1-v0)
}
