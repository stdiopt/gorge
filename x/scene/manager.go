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

import (
	"log"

	"github.com/stdiopt/gorge"
)

// System ...
func System(g *gorge.Gorge) {
	m := &Manager{gorge: g}
	g.Persist(m)
}

// Manager is a system manager
type Manager struct {
	gorge *gorge.Gorge
}

// Load stuff into gorge
func (m *Manager) Load(s *Scene) {
	log.Println("Add entities to gorge")

	log.Printf("entities: %#v", s.entities)

	m.gorge.AddEntity(s.entities...)

	s.Trigger(StartEvent{})
	// Will start to receive gorge events
	m.gorge.Link(&s.Messaging)

}

// New creates a new scene
func (m *Manager) New(fns ...Func) *Scene {
	return &Scene{
		manager: m,
		inits:   fns,
	}
}

// ManagerFromGorge retrieves a SceneManager from gorge
func ManagerFromGorge(g *gorge.Gorge) *Manager {
	var manager *Manager
	g.Query(func(m *Manager) {
		manager = m
	})
	return manager
}
