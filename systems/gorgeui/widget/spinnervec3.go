package widget

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// SpinnerVec3 is a component that displays a 3D vector as a set of spinners.
type SpinnerVec3 struct {
	Component
	Value m32.Vec3
	x     *Spinner
	y     *Spinner
	z     *Spinner
}

// NewSpinnerVec3 returns a default SpinnerVec3
func NewSpinnerVec3(value m32.Vec3) *SpinnerVec3 {
	d := float32(1) / 3

	sx := NewSpinner("X", value[0])
	sx.SetColor(m32.Vec4{.5, 0, 0, 1})
	sx.SetRect(0)
	sx.SetAnchor(0, 0, d, 1)

	sy := NewSpinner("Y", value[1])
	sy.SetColor(m32.Vec4{0, .5, 0, 1})
	sy.SetRect(0)
	sy.SetAnchor(d, 0, d*2, 1)

	sz := NewSpinner("Z", value[2])
	sz.SetColor(m32.Vec4{0, 0, .5, 1})
	sz.SetRect(0)
	sz.SetAnchor(d*2, 0, 1, 1)

	s := &SpinnerVec3{
		Component: *NewComponent(),
		Value:     value,
		x:         sx,
		y:         sy,
		z:         sz,
	}
	s.SetAnchor(0)
	s.SetPivot(0)
	gorgeui.AddChildrenTo(s, sx, sy, sz)
	return s
}

func (s *SpinnerVec3) HandleEvent(e event.Event) {
	switch e := e.(type) {
	case EventSpin:
		_ = e
	}
}
