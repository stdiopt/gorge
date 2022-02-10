package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
)

type WSpinnerVec3 struct {
	Widget[WSpinnerVec3]

	value    gm.Vec3
	changefn func(v gm.Vec3)

	spinners [3]*WSpinner
}

func SpinnerVec3() *WSpinnerVec3 {
	return Build(&WSpinnerVec3{})
}

func (w *WSpinnerVec3) Build(b *B) {
	w.SetAnchor(0, 0, 1, 0)
	w.SetSize(0, 3)

	obsFn := func(i int) func(v float32) {
		return func(v float32) {
			cp := w.value
			cp[i] = v
			w.SetValue(cp)
		}
	}
	a := float32(1) / 3
	b.BeginPane()
	x := b.Spinner("X").
		SetAnchor(0, 0, a, 1).SetSize(0).
		OnChange(obsFn(0)).
		SetLabelColor(.5, 0, 0)
	y := b.Spinner("Y").
		SetAnchor(a, 0, a*2, 1).SetSize(0).
		OnChange(obsFn(1)).
		SetLabelColor(0, .5, 0)
	z := b.Spinner("Z").
		SetAnchor(a*2, 0, 1, 1).SetSize(0).
		OnChange(obsFn(2)).
		SetLabelColor(0, 0, .5)
	b.EndPane()

	w.spinners = [3]*WSpinner{x, y, z}
}

func (w *WSpinnerVec3) SetValue(v gm.Vec3) *WSpinnerVec3 {
	if w.value == v {
		return w
	}
	w.value = v

	w.spinners[0].SetValue(v[0])
	w.spinners[1].SetValue(v[1])
	w.spinners[2].SetValue(v[2])
	if w.changefn != nil {
		w.changefn(v)
	}
	event.Trigger(w, EventValueChanged{Value: w.value})
	return w
}

func (w *WSpinnerVec3) OnChange(fn func(v gm.Vec3)) *WSpinnerVec3 {
	w.changefn = fn
	return w
}

func (b *B) SpinnerVec3() *WSpinnerVec3 {
	w := SpinnerVec3()
	b.Add(w)
	return w
}
