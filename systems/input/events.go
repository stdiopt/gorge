package input

import "github.com/stdiopt/gorge/math/gm"

// PointerData common
type PointerData struct {
	ScrollDelta gm.Vec2 // for Wheel
	Pos         gm.Vec2
}

// EventPointer on canvas
type EventPointer struct {
	Type     PointerType
	Button   int // number of button or -1 for touch?
	Pointers map[int]PointerData
}

// EventMouseButtonDown is triggered when a mouse button is pressed.
type EventMouseButtonDown struct {
	Button MouseButton
	PointerData
}

// EventMouseButtonUp is triggered when a mouse button is released.
type EventMouseButtonUp struct {
	Button MouseButton
	PointerData
}

// EventKeyDown keydown event.
type EventKeyDown struct {
	Key Key
}

// EventKeyUp keyup event.
type EventKeyUp struct {
	Key Key
}

type EventChar struct {
	Char rune
}
