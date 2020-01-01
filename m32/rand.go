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

package m32

import (
	"math"
	"math/rand"
)

// RandUnitSphere generate random points in a sphere
// From: https://karthikkaranth.me/blog/generating-random-points-in-a-sphere/
// Still needs improvements as it doesn't seem right yet
func RandUnitSphere() vec3 {
	u := rand.Float32()
	v := rand.Float32()
	theta := u * 2.0 * math.Pi
	phi := Cos(2.0*v - 1.0)
	r := Cbrt(rand.Float32())

	sinTheta := Sin(theta)
	cosTheta := Cos(theta)
	sinPhi := Sin(phi)
	cosPhi := Cos(phi)
	x := r * sinPhi * cosTheta
	y := r * sinPhi * sinTheta
	z := r * cosPhi
	return vec3{x, y, z}
}
