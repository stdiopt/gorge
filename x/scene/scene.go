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

package scene

import "github.com/stdiopt/gorge"

// Func scene initializer func
type Func func(*Scene)

// Scene contains stuff
type Scene struct {
	manager *Manager
	gorge.Messaging
	entities []gorge.Entity

	inits []Func
}

// Init runs the initializers
func (s *Scene) Init() {
	for _, fn := range s.inits {
		fn(s)
	}
}

// AddEntity to scene
// XXX: should trigger on parent gorge for the renderer
func (s *Scene) AddEntity(e ...gorge.Entity) {
	if s.entities == nil {
		s.entities = []gorge.Entity{}
	}
	s.entities = append(s.entities, e...)

	// Store entities locally
	s.Trigger(gorge.EntitiesAddEvent(e))
}
