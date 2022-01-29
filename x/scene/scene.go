package scene

import (
	"fmt"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
)

type Scene struct {
	event.Bus
	gorge.Container

	Name string

	initfn func(*gorge.Context)

	gorge       *gorge.Context
	initialized bool
}

func New(name string) *Scene {
	return &Scene{Name: name}
}

func (s Scene) String() string {
	return fmt.Sprintf("Scene(%s) elements: %v", s.Name, s.Container)
}
func (s *Scene) GetScene() *Scene { return s }

func (s *Scene) OnInit(fn func(*gorge.Context)) {
	s.initfn = fn
}

func (s *Scene) G() *gorge.Context {
	return s.gorge
}

func (s *Scene) Add(e ...gorge.Entity) {
	s.Container.Add(e...)

	if s.gorge != nil {
		s.gorge.Add(e...)
	}
}

func (s *Scene) Remove(e ...gorge.Entity) {
	s.Container.Remove(e...)

	if s.gorge != nil {
		s.gorge.Remove(e...)
	}
}

func (s *Scene) initScene(g *gorge.Context) {
	s.gorge = g
	g.AddBus(s)

	if s.initfn != nil {
		s.initfn(g)
	}
	s.initialized = true
}

func (s *Scene) destroyScene(g *gorge.Context) {
	g.RemoveBus(s)
	s.gorge = nil
}
