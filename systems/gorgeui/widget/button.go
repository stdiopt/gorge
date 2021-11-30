package widget

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// ButtonState type
type ButtonState int

// Known button states
const (
	ButtonStatePressed = 1 << (iota + 1)
	ButtonStateHover
)

// Button widget state.
type Button struct {
	Component
	NormalColor    m32.Vec4
	HighlightColor m32.Vec4

	PressedColor  m32.Vec4
	SelectedColor m32.Vec4
	DisabledColor m32.Vec4

	FadeFactor float32
	Disabled   bool

	state ButtonState

	// Is it worth it?
	Entity graphic
}

// HandleEvent handles gorgeui events.
func (b *Button) HandleEvent(e gorgeui.Event) {
	switch e := e.Value.(type) {
	case gorgeui.EventUpdate:
		rect := b.Rect()
		et := b.Entity.Transform()
		et.Position[0] = rect[0]
		et.Position[1] = rect[1] // bottom
		et.Scale[0] = rect[2] - rect[0]
		et.Scale[1] = rect[3] - rect[1]
		// Colors based on state stuff
		color := b.Entity.Colorable().GetColor()
		target := b.NormalColor
		switch {
		case b.state&ButtonStatePressed != 0:
			target = b.PressedColor
		case b.state&ButtonStateHover != 0:
			target = b.HighlightColor
		}
		color = color.Lerp(target, e.DeltaTime()*b.FadeFactor)
		b.Entity.Colorable().SetColorv(color)

	case gorgeui.EventPointerDown:
		b.state |= ButtonStatePressed

	case gorgeui.EventPointerUp:
		b.state &= ^ButtonStatePressed
		if gorgeui.HasParent(e.Target, b) {
			gorgeui.TriggerOn(b, EventClick{b})
		}

	case gorgeui.EventPointerEnter:
		b.state |= ButtonStateHover

	case gorgeui.EventPointerLeave:
		b.state &= ^ButtonStateHover

	}
}

// NewButton returns a new button based on options.
func NewButton() *Button {
	b := &Button{
		Component:      *NewComponent(),
		NormalColor:    m32.Vec4{.7, .7, .7, .9},
		HighlightColor: m32.Vec4{.8, .9, .8, 1}, // hover
		PressedColor:   m32.Vec4{.4, .4, .4, 1},
		SelectedColor:  m32.Vec4{.9, .7, .7, 1}, // Unused
		DisabledColor:  m32.Vec4{.5, .5, .5, 1}, // Unused
		FadeFactor:     10,
		Entity:         QuadEntity(),
	}
	gorgeui.AddChildrenTo(b, b.Entity)

	return b
}

// SetFadeFactor returns fade factor.
func (b *Button) SetFadeFactor(a float32) {
	b.FadeFactor = a
}

// SetNormalColor sets the default state color.
func (b *Button) SetNormalColor(c ...float32) {
	b.NormalColor = v4Color(c...)
}

// SetHighlightColor sets the hover color.
func (b *Button) SetHighlightColor(c ...float32) {
	b.HighlightColor = v4Color(c...)
}

// SetPressedColor sets the pressed color.
func (b *Button) SetPressedColor(c ...float32) {
	b.PressedColor = v4Color(c...)
}

// SetMaterial sets the button entity material.
func (b *Button) SetMaterial(mat *gorge.Material) {
	b.Entity.Renderable().Material = mat
}
