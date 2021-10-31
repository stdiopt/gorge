// Package input normalizes inputs from systems
package input

import "github.com/stdiopt/gorge/m32"

// PointerType pointer event type
type PointerType int

func (p PointerType) String() string {
	switch p {
	case MouseDown:
		return "MouseDown"
	case MouseUp:
		return "MouseUp"
	case MouseMove:
		return "MouseMove"
	case MouseWheel:
		return "MouseWheel"
	case PointerDown:
		return "PointerDown"
	case PointerMove:
		return "PointerMove"
	case PointerEnd:
		return "PointerEnd"
	case PointerCancel:
		return "PointerCancel"
	default:
		return "<invalid>"
	}
}

// Pointer comments
const (
	_ = PointerType(iota)
	MouseDown
	MouseUp
	MouseMove
	MouseWheel
	PointerDown
	PointerMove
	PointerEnd
	PointerCancel
)

// PointerData common
type PointerData struct {
	DeltaZ float32 // for Wheel
	Pos    m32.Vec2
}

// EventPointer on canvas
type EventPointer struct {
	Type     PointerType
	Button   int // number of button or -1 for touch?
	Pointers map[int]PointerData
}
