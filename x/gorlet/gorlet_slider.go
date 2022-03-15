package gorlet

import (
	"fmt"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/math/ray"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type WSlider struct {
	Widget[*WSlider]

	background *WPane
	track      *WContainer
	handler    *WButton
	lbl        *WLabel

	value       float32
	handlerSize float32
	valFmt      string
	min, max    float32
	changefn    func(float32)
}

func Slider(min, max float32) *WSlider {
	return Build(&WSlider{
		min:         min,
		max:         max,
		valFmt:      "%.2f",
		handlerSize: 3,
	})
}

func (w *WSlider) Build(b *B) {
	w.SetAnchor(0, 0, 1, 0)
	w.SetSize(0, 3)
	w.SetDragEvents(true)

	w.background = b.BeginPane()
	{
		w.track = b.BeginContainer().
			SetName("sliderTrack").
			SetRect(w.handlerSize/2, 0, w.handlerSize/2, 0).
			SetPivot(0, .5)
		{
			w.handler = b.Button().
				SetName("sliderHandler").
				SetPivot(.5).
				SetSize(w.handlerSize, 0).
				SetAnchor(0, 0, 0, 1).
				SetMargin(.2)
			w.lbl = b.Label(fmt.Sprintf(w.valFmt, w.min)).
				SetName("sliderLabel").
				FillParent().
				SetDisableRaycast(true)
		}
		b.EndContainer()
	}
	b.EndPane()

	var dragging bool
	event.Handle(w, func(e gorgeui.EventPointerUp) {
		if dragging {
			return
		}
		w.doDrag(e.PointerData)
	})
	event.Handle(w, func(e gorgeui.EventDrag) {
		dragging = true
		w.doDrag(e.PointerData)
	})
	event.Handle(w, func(gorgeui.EventDragEnd) {
		dragging = false
	})
}

func (w *WSlider) SetMin(v float32) *WSlider {
	w.min = v
	w.SetValue(w.value)
	return w
}

func (w *WSlider) SetMax(v float32) *WSlider {
	w.max = v
	w.SetValue(w.value)
	return w
}

func (w *WSlider) doDrag(pd *gorgeui.PointerData) {
	rect := w.track.Rect()

	// We ray trace on track
	m := w.Mat4()
	v0 := m.MulV4(gm.Vec4{rect[0], rect[1], 0, 1}).Vec3() // 0
	v1 := m.MulV4(gm.Vec4{rect[2], rect[1], 0, 1}).Vec3() // right
	v2 := m.MulV4(gm.Vec4{rect[0], rect[3], 0, 1}).Vec3() // up)

	ui := gorgeui.RootUI(w)
	r := ray.FromScreen(ui.ScreenSize(), ui.Camera, pd.Position)
	res := ray.IntersectRect(r, v0, v1, v2)

	v := res.UV[0]

	v -= w.handlerSize / rect[2] / 2
	v = gm.Clamp(v, 0, 1)
	w.SetValue(w.real(v))
}

func (w *WSlider) SetValue(v float32) *WSlider {
	v = w.norm(v)
	if w.value == v {
		return w
	}
	w.value = v
	w.updateValue()
	rval := w.real(w.value)
	if w.changefn != nil {
		w.changefn(rval)
	}
	event.Trigger(w, EventValueChanged{rval})
	return w
}

func (w *WSlider) SetFormat(f string) *WSlider {
	w.valFmt = f
	w.updateValue()
	return w
}

func (w *WSlider) SetFontScale(v float32) *WSlider {
	w.lbl.SetFontScale(v)
	return w
}

func (w *WSlider) OnChange(f func(float32)) *WSlider {
	w.changefn = f
	return w
}

func (w *WSlider) updateValue() {
	rval := w.real(w.value)
	w.handler.SetAnchor(w.value, 0, w.value, 1)
	w.lbl.SetText(fmt.Sprintf(w.valFmt, rval))
}

func (w *WSlider) norm(v float32) float32 {
	return gm.Clamp((v-w.min)/(w.max-w.min), 0, 1)
}

// 0-1 to min-max
func (w *WSlider) real(v float32) float32 {
	return w.min + v*(w.max-w.min)
}

func (b *B) Slider(min, max float32) *WSlider {
	w := Slider(min, max)
	b.Add(w)
	return w
}
