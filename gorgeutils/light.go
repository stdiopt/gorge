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

import "github.com/stdiopt/gorge"

// Light is an entity with transform and light
type Light struct {
	gorge.Transform
	gorge.Light
}

// NewLight returns a light entity
func NewLight() *Light {
	return &Light{
		*gorge.NewTransform(),
		gorge.Light{Color: vec3{1, 1, 1}},
	}
}
