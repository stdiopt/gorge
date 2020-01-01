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

// Package gorge contains mostly data only components
package gorge

import (
	"fmt"
	glog "log"

	"github.com/stdiopt/gorge/m32"
)

var (
	log = glog.New(glog.Writer(), "(gorge) ", 0)
)

type (
	vec2 = m32.Vec2
	vec3 = m32.Vec3
	vec4 = m32.Vec4
	mat4 = m32.Mat4
	quat = m32.Quat
)

// State type for gorge
type State int

func (s State) String() string {
	switch s {
	case StateZero:
		return "zero"
	case StateInitialized:
		return "initialized"
	case StateStarted:
		return "started"
	case StateClosed:
		return "closed"
	}
	return "<undefined>"
}

// gorge States
const (
	StateZero = iota
	StateInitialized
	StateStarted
	StateClosed
)

// Entity can be any thing
type Entity interface{}

// SystemFunc initializers
type SystemFunc = func(*Gorge)

// Gorge main state manager and message bus
type Gorge struct {
	Messaging
	done  chan struct{}
	tick  chan float32
	state State

	inits []SystemFunc
}

// New create a new manager
func New(systems ...SystemFunc) *Gorge {
	g := &Gorge{
		inits: systems,
		done:  make(chan struct{}),
		tick:  make(chan float32),
	}
	g.init()
	return g
}

// Init systems
func (g *Gorge) init() {
	if g.state != StateZero {
		return
	}
	g.state = StateInitialized
	for _, fn := range g.inits {
		fn(g)
	}
}

// Start the systems
func (g *Gorge) Start() {
	if g.state != StateInitialized {
		panic(fmt.Sprintf("cannot start, current state is: %v", g.state))
	}
	g.Persist(StartEvent{})
	g.Trigger(AfterStartEvent{})
}

// Run until close is called on done
func (g *Gorge) Run() {
	g.Start()
	for {
		select {
		case dt := <-g.tick:
			g.UpdateNow(dt)
		case <-g.done:
			return
		}
	}
}

// Update triggers update events
func (g *Gorge) Update(dt float32) {
	g.tick <- dt
}

// UpdateNow just updates stuff right awaym does not use the channel thing
func (g *Gorge) UpdateNow(dt float32) {

	g.Trigger(PreUpdateEvent(dt))
	g.Trigger(UpdateEvent(dt))
	g.Trigger(PostUpdateEvent(dt))

	g.Trigger(RenderEvent(dt))

	// XXX: Profiling
	/*g.Range(func(k, v interface{}) bool {
		hg := v.(*HandlerGroup)
		if hg.CallEnd.Sub(hg.CallStart) > 30*time.Millisecond {
			log.Printf("group: %v took too long %v", hg.Type, hg.CallEnd.Sub(hg.CallStart))
			hg.CallStart = time.Time{}
			hg.CallEnd = time.Time{}
		}
		return true
	})*/
}

// Close closes manager
func (g *Gorge) Close() {
	g.state = StateClosed
	close(g.done)
}

// AddEntity adds an entity
/*func (g *Gorge) AddEntity(e ...Entity) {
	ents := make([]Entity, 0, len(e)) // at least len(e)
	for _, e := range e {
		switch e := e.(type) {
		case []Entity:
			ents = append(ents, e...)
		default:
			ents = append(ents, e)
		}
	}
	g.Trigger(EntitiesAddEvent(ents))
}

// RemoveEntity an entity
func (g *Gorge) RemoveEntity(ents ...Entity) {
	g.Trigger(EntitiesRemoveEvent(ents))
}*/

//////////////////////////
// Experiment handler funcs
/////////////

// Better typed handlers, scene doesn't work if we hide the data types

// HandleStart helper listens for a gorge.StartEvent
func (g *Gorge) HandleStart(fn func()) *Handler {
	return g.Handle(func(e StartEvent) {
		fn()
	})
}

// HandleUpdate helper litens for a gorge.UpdateEvent
func (g *Gorge) HandleUpdate(fn func(dt float32)) *Handler {
	return g.Handle(func(e UpdateEvent) {
		fn(float32(e))
	})
}

// HandlePostUpdate helper listens for a gorge.PostUpdateEvent
func (g *Gorge) HandlePostUpdate(fn func(dt float32)) *Handler {
	return g.Handle(func(e PostUpdateEvent) {
		fn(float32(e))
	})
}

// Error persists an error in the event system
func (g *Gorge) Error(err error) {
	g.Persist(ErrorEvent{err})
}

// Warn persists a warning msg in the event system
func (g *Gorge) Warn(s string) {
	g.Persist(WarnEvent(s))
}
