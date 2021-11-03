// Package input normalizes inputs from systems
package input

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
)

// This might better be mouse?
type mouseManager struct {
	gorge       *gorge.Context
	buttonState map[MouseButton]ActionState
	// Need to handle current mouse position and whatnots
	// MousePosition
	mpos        m32.Vec2
	deltaPos    m32.Vec2
	scrollDelta float32
}

func (m *mouseManager) update() {
	for k, v := range m.buttonState {
		switch v {
		case ActionDown:
			m.buttonState[k] = ActionHold
		case ActionUp:
			delete(m.buttonState, k)
		}
	}
	m.scrollDelta = 0
	m.deltaPos = m32.Vec2{}
}

func (m *mouseManager) SetScrollDelta(delta float32) {
	m.scrollDelta = delta

	// Legacy
	evt := EventPointer{
		Type: MouseWheel,
		Pointers: map[int]PointerData{
			0: {DeltaZ: delta, Pos: m.mpos},
		},
	}
	m.gorge.Trigger(evt) // nolint: errcheck
}

func (m *mouseManager) SetCursorPosition(x, y float32) {
	p := m32.Vec2{x, y}
	m.deltaPos = m.mpos.Sub(p)
	m.mpos = p

	// Legacy
	evt := EventPointer{
		Type: MouseMove,
		Pointers: map[int]PointerData{
			0: {Pos: m.mpos},
		},
	}
	m.gorge.Trigger(evt) // nolint: errcheck*/
	// Trigger position event
}

func (m *mouseManager) SetMouseButtonState(b MouseButton, s ActionState) {
	if m.buttonState == nil {
		m.buttonState = map[MouseButton]ActionState{}
	}

	pd := PointerData{
		DeltaZ: m.scrollDelta,
		Pos:    m.mpos,
	}

	m.buttonState[b] = s
	switch s {
	case ActionUp:
		m.gorge.Trigger(EventMouseButtonUp{b, pd})
		// legacy
		m.gorge.Trigger(EventPointer{
			Type: MouseUp,
			Pointers: map[int]PointerData{
				0: {Pos: m.mpos},
			},
		})
	case ActionDown:
		m.gorge.Trigger(EventMouseButtonDown{b, pd})
		m.gorge.Trigger(EventPointer{
			Type: MouseDown,
			Pointers: map[int]PointerData{
				0: {Pos: m.mpos},
			},
		})
	}
}

// ScrollDelta returns scrollDelta.
func (m *mouseManager) ScrollDelta() float32 {
	return m.scrollDelta
}

// CursorPosition returns the current cursor position.
func (m *mouseManager) CursorPosition() m32.Vec2 {
	return m.mpos
}

// CursorDelta returns the current cursor position.
func (m *mouseManager) CursorDelta() m32.Vec2 {
	return m.deltaPos
}

// MouseButtonUp returns true if the last mouse button state was up.
func (m *mouseManager) MouseButtonUp(k MouseButton) bool {
	return m.getMouseButton(k) == ActionUp
}

// MouseButtonClick similar to MouseButtonUp
func (m *mouseManager) MouseButtonClick(k MouseButton) bool {
	return m.getMouseButton(k) == ActionUp
}

// MouseButtonDown returns true if the state of the mouse button is Down or Hold.
func (m *mouseManager) MouseButtonDown(k MouseButton) bool {
	s := m.getMouseButton(k)
	return s == ActionDown || s == ActionHold
}

func (m *mouseManager) getMouseButton(k MouseButton) ActionState {
	if m.buttonState == nil {
		return ActionState(0)
	}
	return m.buttonState[k]
}

// MouseButton represents a mouse button.
type MouseButton int

// Mouse button constants.
const (
	MouseUnknown = MouseButton(iota)
	MouseLeft
	MouseRight
	MouseMiddle // aka wheel
	MouseThumb1
	MouseThumb2
)
