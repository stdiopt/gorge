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

package gorgeutils

import (
	"math"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
)

// Camera entity
type Camera struct {
	gorge.Transform
	gorge.Camera
}

// NewCamera could have more stuff
func NewCamera() *Camera {
	return &Camera{
		*gorge.NewTransform().
			SetPosition(0, 4, -10).
			LookAt(vec3{0, 1, 0}, m32.Up()),
		gorge.Camera{
			Fov:         math.Pi / 4,
			AspectRatio: 1,
			Near:        0.1,
			Far:         1000,
			Ambient:     vec3{0.4, 0.4, 0.4},
		},
	}
}
