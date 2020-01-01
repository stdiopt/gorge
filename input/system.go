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

package input

import (
	"github.com/stdiopt/gorge"
)

// Input system tied to manager
type Input struct {
	manager  *gorge.Gorge
	keyState map[string]bool
}

// System returns a input system
func System(g *gorge.Gorge) {
	s := &Input{
		manager:  g,
		keyState: map[string]bool{},
	}
	g.Persist(s)
}

// KeyDown triggers key down
func (s Input) KeyDown(key string) {
	s.keyState[key] = true
	s.manager.Trigger(KeyEvent{
		KeyDown,
		key,
	})
}

// KeyUp triggers key down
func (s Input) KeyUp(key string) {
	delete(s.keyState, key)
	s.manager.Trigger(KeyEvent{
		KeyUp,
		key,
	})
}

// GetKey checks if a key was pressed
func (s Input) GetKey(key string) bool {
	return s.keyState[key]
}

type queryier interface {
	Query(fn interface{})
}

// FromECS returns a input system from gorge
func FromECS(q queryier) *Input {
	// Get from messaging store
	var ret *Input
	q.Query(func(s *Input) { ret = s })

	if ret == nil {
		panic("input system doesn't exist in gorge")
	}

	return ret
}
