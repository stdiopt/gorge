package input

import (
	"log"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
)

// Input system tied to manager
type Input struct {
	keyManager
}

// HandleEvent implements event handler.
func (s *Input) HandleEvent(v event.Event) {
	if _, ok := v.(gorge.EventPostUpdate); !ok {
		return
	}
	for k, v := range s.keyState {
		switch v {
		case KeyStateDown:
			s.keyState[k] = KeyStateHold
		case KeyStateUp:
			delete(s.keyState, k)
		}
	}
}

// System initialize the key manager on gorge.
func System(g *gorge.Context) {
	log.Println("Initializing system")
	s := &Input{
		keyManager: keyManager{
			gorge: g,
		},
	}
	g.Handle(s)
	g.PutProp(func() *Context {
		return &Context{s}
	})
}
