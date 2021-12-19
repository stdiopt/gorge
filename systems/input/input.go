package input

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
)

// Input system tied to manager
type Input struct {
	keyManager
	mouseManager
}

// HandleEvent implements event handler.
func (s *Input) HandleEvent(v event.Event) {
	if _, ok := v.(gorge.EventPostUpdate); !ok {
		return
	}

	s.keyManager.update()
	s.mouseManager.update()
}

// System initialize the key manager on gorge.
/*func System(g *gorge.Context) {
	log.Println("Initializing system")
	s := &Input{
		keyManager:   keyManager{gorge: g},
		mouseManager: mouseManager{gorge: g},
	}
	g.Handle(s)
	g.PutProp(func() *Context {
		return &Context{s}
	})
}*/

// ActionState type
type ActionState int

func (k ActionState) String() string {
	switch k {
	case ActionDown:
		return "ActionDown"
	case ActionHold:
		return "ActionHold"
	case ActionUp:
		return "ActionUp"
	default:
		return "<no state>"
	}
}

// Action states used in keys and buttons.
const (
	_ = ActionState(iota)
	ActionDown
	ActionHold
	ActionUp
)
