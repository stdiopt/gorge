package input

import (
	"log"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
)

type input = Input

// Context to be used in gorge systems
type Context struct {
	*input
}

// FromContext returns a input.Context from a gorge.Context
func FromContext(g *gorge.Context) *Context {
	if ctx, ok := gorge.GetContext[*Context](g); ok {
		return ctx
	}

	log.Println("Initializing system")
	s := &Input{
		keyManager:   keyManager{gorge: g},
		mouseManager: mouseManager{gorge: g},
	}
	event.Handle(g, func(gorge.EventPostUpdate) {
		s.keyManager.update()
		s.mouseManager.update()
	})
	return gorge.AddContext(g, &Context{s})
}

// IsDown checks if wether a key or mouse button is pressed.
func (c Context) IsDown(v any) bool {
	var a ActionState
	switch v := v.(type) {
	case Key:
		a = c.getKey(v)
	case MouseButton:
		a = c.getMouseButton(v)
	}
	return a == ActionDown || a == ActionHold
}

// IsUp checks if wether a key or mouse button is released.
func (c Context) IsUp(v any) bool {
	switch v := v.(type) {
	case Key:
		return c.getKey(v) == ActionUp
	case MouseButton:
		return c.getMouseButton(v) == ActionUp
	}
	return false
}

// IsPressed alias to IsUp.
func (c Context) IsPressed(v any) bool {
	return c.IsUp(v)
}
