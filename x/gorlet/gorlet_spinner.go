package gorlet

import (
	"fmt"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/text"
)

type WSpinner struct {
	Widget[*WSpinner]

	value float32
	fmt   string

	lblbg    *WPane
	lbl      *WLabel
	valueLbl *WLabel
	changefn func(float32)
}

func Spinner(lbl string) *WSpinner {
	return Build(&WSpinner{
		fmt: "%.2f",
	}).SetLabel(lbl)
}

func (w *WSpinner) Build(b *B) {
	w.SetAnchor(0, 0, 1, 0)
	w.SetSize(0, 3)
	w.SetDragEvents(true)

	b.BeginPane()
	{
		w.lblbg = b.BeginPane().
			SetAnchor(0, 0, .3, 1)
		w.lbl = b.Label("").
			SetColor(1).
			FillParent()
		b.EndPane()

		w.valueLbl = b.Label(fmt.Sprintf(w.fmt, w.value)).
			SetAnchor(.3, 0, 1, 1).
			SetRect(0).
			SetColor(1).
			SetOverflow(text.OverflowOverlap)
	}
	b.EndPane()

	var last gm.Vec2
	var delta float32
	event.Handle(w, func(e gorgeui.EventDragBegin) {
		res := w.IntersectFromScreen(e.Position)
		last = res.UV
	})
	event.Handle(w, func(e gorgeui.EventDrag) {
		res := w.IntersectFromScreen(e.Position)
		delta = res.UV[0] - last[0]
		last = res.UV
		factor := .01 + gm.Abs(res.UV[1]-.5)
		w.SetValue(w.value + delta*factor)
	})
	w.SetValue(0)
}

func (w *WSpinner) SetValue(v float32) *WSpinner {
	if w.value == v {
		return w
	}
	w.value = v
	w.valueLbl.SetTextf(w.fmt, w.value)
	if w.changefn != nil {
		w.changefn(w.value)
	}
	event.Trigger(w, EventValueChanged{w.value})
	return w
}

func (w *WSpinner) SetLabel(t string) *WSpinner {
	w.lbl.SetText(t)
	return w
}

func (w *WSpinner) SetLabelColor(c ...float32) *WSpinner {
	w.lblbg.SetColor(c...)
	return w
}

func (w *WSpinner) SetLabelTextColor(c ...float32) *WSpinner {
	w.lbl.SetColor(c...)
	return w
}

func (w *WSpinner) SetTextColor(c ...float32) *WSpinner {
	w.valueLbl.SetColor(c...)
	return w
}

func (w *WSpinner) OnChange(fn func(float32)) *WSpinner {
	w.changefn = fn
	return w
}

func (b *B) Spinner(lbl string) *WSpinner {
	w := Spinner(lbl)
	b.Add(w)
	return w
}
