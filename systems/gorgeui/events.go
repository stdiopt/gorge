package gorgeui

import (
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/m32/ray"
)

/*
type Event struct {
	Entity          Entity
	Value           interface{}
	stopPropagation bool
}
*/

// Events used by widgets.
type (
	// EventPointerEnter triggers when a pointer enters the widget rect.
	EventPointerEnter struct{ *PointerData }
	// EventPointerLeave triggers when a pointer leaves the widget rect.
	EventPointerLeave struct{ *PointerData }
	// EventPointerDown triggers when a pointer is pressed on widget rect.
	EventPointerDown struct{ *PointerData }
	// EventPointerUp triggers when a pointer is unpressed on widget rect.
	EventPointerUp struct{ *PointerData }
	// EventDragBegin triggers when a pointer starts dragging from the widget rect.
	EventDragBegin struct{ *PointerData }
	// EventDrag triggers on the widget that started the drag.
	EventDrag struct{ *PointerData }
	// EventDragEnd triggers when the pointer drop the widget.
	EventDragEnd struct{ *PointerData }
	// EventUpdate triggers on gorge.EventUpdate
	EventUpdate float32
)

// PointerData common pointer data for pointer events.
type PointerData struct {
	Target    Entity // could be a children
	Position  m32.Vec2
	Delta     m32.Vec2
	RayResult ray.Result

	stopPropagation bool
}

// StopPropagation sets the stop propagation flag.
func (p *PointerData) StopPropagation() {
	p.stopPropagation = true
}

// DeltaTime for the EventUpdate.
func (e EventUpdate) DeltaTime() float32 { return float32(e) }
