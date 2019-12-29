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

import "github.com/go-gl/mathgl/mgl32"

type (
	// Mat3 alias
	Mat3 = mgl32.Mat3
	// Mat4 alias for a 4x4 matrix
	Mat4 = mgl32.Mat4
)

// Mat4Mul multiply sequence of matrices
func Mat4Mul(ms ...mat4) mat4 {
	r := ms[0]
	for _, m := range ms[1:] {
		r = r.Mul4(m)
	}
	return r
}
