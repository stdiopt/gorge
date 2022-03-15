package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type WButton struct {
	Widget[*WButton]
	color      gm.Vec4
	highlight  gm.Vec4
	pressed    gm.Vec4
	fadeFactor float32

	pane    *WPane
	clickfn func()
}

func Button(children ...gorge.Entity) *WButton {
	return Build(&WButton{
		color:      gm.Color(.7, .9),
		highlight:  gm.Color(.8, .9, .8),
		pressed:    gm.Color(.4),
		fadeFactor: 10,
	}).Add(children...)
}

func (w *WButton) Build(b *B) {
	w.SetAnchor(0, 0, 1, 0)
	w.SetSize(0, 3)

	w.pane = b.Pane().SetColor(w.color[:]...)

	type buttonState int
	const (
		statePressed = 1 << (iota + 1)
		stateHover
	)
	var state buttonState

	c := w.color
	event.Handle(w, func(e gorgeui.EventUpdate) {
		s := w.fadeFactor
		target := w.color
		switch {
		case state&statePressed != 0:
			s *= 2
			target = w.pressed
		case state&stateHover != 0:
			target = w.highlight
		}
		// Due to floating point this might run everytime but
		// it is somewhat ok since comparing with epsilon might be slower
		if target != c {
			c = c.Lerp(target, e.DeltaTime()*s)
			w.pane.SetColor(c[:]...)
		}
	})

	event.Handle(w, func(gorgeui.EventPointerDown) {
		state |= statePressed
	})
	event.Handle(w, func(gorgeui.EventPointerUp) {
		state &= ^statePressed
		if w.clickfn != nil {
			w.clickfn()
		}
		event.Trigger(w, EventClick{})
	})
	event.Handle(w, func(gorgeui.EventPointerEnter) {
		state |= stateHover
		event.Trigger(gorgeui.RootUI(w).Gorge(), gorge.EventCursor(gorge.CursorHand))
	})
	event.Handle(w, func(gorgeui.EventPointerLeave) {
		state &= ^stateHover
		event.Trigger(gorgeui.RootUI(w).Gorge(), gorge.EventCursor(gorge.CursorArrow))
	})
}

func (w *WButton) Color(vs ...float32) *WButton {
	w.color = gm.Color(vs...)
	return w
}

func (w *WButton) Highlight(vs ...float32) *WButton {
	w.highlight = gm.Color(vs...)
	return w
}

func (w *WButton) Pressed(vs ...float32) *WButton {
	w.pressed = gm.Color(vs...)
	return w
}

func (w *WButton) FadeFactor(f float32) *WButton {
	w.fadeFactor = f
	return w
}

func (w *WButton) OnClick(fn func()) *WButton {
	w.clickfn = fn
	return w
}

func (b *B) Button() *WButton {
	w := Button()
	b.Add(w)
	return w
}

func (b *B) BeginButton() *WButton {
	w := Button()
	b.Begin(w)
	return w
}

func (b *B) EndButton() {
	b.End()
}
