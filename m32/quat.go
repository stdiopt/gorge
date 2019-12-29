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

// Quaternions
/*
Resources:
	https://answers.unity.com/questions/467614/what-is-the-source-code-of-quaternionlookrotation.html
*/
import (
	"github.com/go-gl/mathgl/mgl32"
)

type (
	// Quat mgl32 alias
	Quat = mgl32.Quat
)

// QuatBetweenVectors

// QuatLookAt implementation from GLM
// glm/gtc/quaternion.inl + quat_cast
func QuatLookAt(dir, up vec3) quat {

	var (
		v             = dir.Normalize()
		v2            = up.Cross(v).Normalize()
		v3            = v.Cross(v2)
		m00, m01, m02 = v2[0], v2[1], v2[2]
		m10, m11, m12 = v3[0], v3[1], v3[2]
		m20, m21, m22 = v[0], v[1], v[2]
	)

	fourXSquaredMinus1 := m00 - m11 - m22
	fourYSquaredMinus1 := m11 - m00 - m22
	fourZSquaredMinus1 := m22 - m00 - m11
	fourWSquaredMinus1 := m00 + m11 + m22

	biggestIndex := 0
	fourBiggestSquaredMinus1 := fourWSquaredMinus1
	if fourXSquaredMinus1 > fourBiggestSquaredMinus1 {
		fourBiggestSquaredMinus1 = fourXSquaredMinus1
		biggestIndex = 1
	}
	if fourYSquaredMinus1 > fourBiggestSquaredMinus1 {
		fourBiggestSquaredMinus1 = fourYSquaredMinus1
		biggestIndex = 2
	}
	if fourZSquaredMinus1 > fourBiggestSquaredMinus1 {
		fourBiggestSquaredMinus1 = fourZSquaredMinus1
		biggestIndex = 3
	}

	biggestVal := Sqrt(fourBiggestSquaredMinus1+1) * 0.5
	mult := 0.25 / biggestVal

	switch biggestIndex {
	case 0:
		return quat{
			W: biggestVal,
			V: vec3{(m12 - m21) * mult, (m20 - m02) * mult, (m01 - m10) * mult},
		}
	case 1:
		return quat{
			W: (m12 - m21) * mult,
			V: vec3{biggestVal, (m01 + m10) * mult, (m20 + m02) * mult},
		}
	case 2:
		return quat{
			W: (m20 - m02) * mult,
			V: vec3{(m01 + m10) * mult, biggestVal, (m12 + m21) * mult},
		}
	case 3:
		return quat{
			W: (m01 - m10) * mult,
			V: vec3{(m20 + m02) * mult, (m12 + m21) * mult, biggestVal},
		}
	}
	return mgl32.QuatIdent()
}

// QuatEuler returns a quaternion based on those euler angles
// https://en.wikipedia.org/wiki/Conversion_between_quaternions_and_Euler_angles
func QuatEuler(x, y, z float32) quat {

	cy, sy := Sincos(z * 0.5)
	cp, sp := Sincos(y * 0.5)
	cr, sr := Sincos(x * 0.5)

	return quat{
		W: cy*cp*cr + sy*sp*sr,
		V: vec3{
			cy*cp*sr - sy*sp*cr,
			sy*cp*sr + cy*sp*cr,
			sy*cp*cr - cy*sp*sr,
		},
	}
}
