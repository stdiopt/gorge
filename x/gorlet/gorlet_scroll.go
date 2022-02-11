package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type scrollent struct {
	disabled bool
	entity   *WScrollbar
	size     float32
	value    float32
}

type WScroll struct {
	Widget[WScroll]

	scrollVel gm.Vec2

	scrolls   [2]scrollent
	container *WContainer
	mask      *WMask
}

func Scroll(c ...gorge.Entity) *WScroll {
	return Build(&WScroll{}).Add(c...)
}

func (w *WScroll) Build(b *B) {
	w.mask = b.BeginMask()
	w.container = b.BeginContainer()
	b.ClientArea()
	b.EndContainer()
	b.EndMask()
	// w.SetClientArea(w.container)

	w.scrolls = [2]scrollent{
		{
			entity: Scrollbar(Horizontal).
				SetAnchor(0, 1, 1, 1).
				OnChange(w.scrollFunc(0)),
			size: 1,
		},
		{
			entity: Scrollbar(Vertical).
				SetAnchor(1, 0, 1, 1).
				OnChange(w.scrollFunc(1)),
			size: 1,
		},
	}
	b.Add(w.scrolls[0].entity)
	b.Add(w.scrolls[1].entity)

	event.Handle(w, func(e gorgeui.EventPointerWheel) {
		sz := w.mask.ContentSize()
		b := w.container.CalcMax()

		if !w.scrolls[0].disabled && sz[0] > b[0] {
			w.scrollVel[0] -= e.Wheel[0] * 0.005
			if e.Wheel[0] > 0 && w.scrolls[0].value != 1 {
				e.StopPropagation()
			} else if e.Wheel[0] < 0 && w.scrolls[0].value != 0 {
				e.StopPropagation()
			}
		}
		if !w.scrolls[1].disabled && sz[1] < b[1] {
			w.scrollVel[1] += e.Wheel[1] * 0.005
			if e.Wheel[1] > 0 && !gm.FloatEqual(w.scrolls[1].value, 1) {
				e.StopPropagation()
			} else if e.Wheel[1] < 0 && !gm.FloatEqual(w.scrolls[1].value, 0) {
				e.StopPropagation()
			}
		}
	})
	event.Handle(w, func(e gorgeui.EventUpdate) {
		sz := w.mask.ContentSize()
		b := w.container.CalcMax()

		ssz := gm.Vec2{}

		w.scrollVel = w.scrollVel.Lerp(gm.Vec2{}, e.DeltaTime()*20)
		if !w.scrolls[0].disabled {
			if sz[0]+gm.Epsilon >= b[0] {
				w.scrolls[0].entity.SetValue(0)
				w.scrolls[0].entity.SetMax(sz[0])
			} else {
				ssz[0] = w.scrolls[0].size

				hs := sz[0] * (sz[0] / b[0])
				delta := w.scrollVel[0] * .8
				curScroll := w.container.Position[0] / (sz[0] - b[0])
				w.scrolls[0].entity.SetValue(curScroll + delta)
				w.scrolls[0].entity.SetMax(hs)
			}
		}

		if !w.scrolls[1].disabled {
			// Vertical
			if sz[1]+gm.Epsilon >= b[1] {
				w.scrolls[1].entity.SetValue(0)
				w.scrolls[1].entity.SetMax(sz[1])
			} else {
				ssz[1] = w.scrolls[1].size
				hs := sz[1] * (sz[1] / b[1])
				curScroll := w.container.Position[1] / (sz[1] - b[1])
				delta := w.scrollVel[1] * .8
				w.scrolls[1].entity.SetValue(curScroll + delta)
				w.scrolls[1].entity.SetMax(hs)
			}
		}
		w.mask.SetSize(ssz[1], ssz[0])
		w.scrolls[0].entity.SetRect(0, -ssz[0], ssz[1], ssz[0])
		w.scrolls[1].entity.SetRect(-ssz[1], 0, ssz[1], ssz[0])
	})
}

func (w *WScroll) SetScrollSize(v ...float32) *WScroll {
	switch len(v) {
	case 0:
		w.scrolls[0].size = 1
		w.scrolls[1].size = 1
	case 1:
		w.scrolls[0].size = v[0]
		w.scrolls[1].size = v[0]
	default:
		w.scrolls[0].size = v[0]
		w.scrolls[1].size = v[1]
	}
	return w
}

func (w *WScroll) scrollFunc(n int) func(v float32) {
	return func(v float32) {
		if w.scrolls[n].disabled {
			return
		}
		v = gm.Clamp(v, 0, 1)
		if v == w.scrolls[n].value {
			return
		}
		sz := w.mask.ContentSize()
		b := w.container.CalcMax()
		if sz[n]/b[n] > 1 {
			return
		}
		w.container.Position[n] = v * (sz[n] - b[n])
		w.scrolls[n].value = v
	}
}

type WScrollbar struct {
	Widget[WScrollbar]

	handler  *WButton
	track    *WContainer
	changefn func(float32)

	handlerSize float32
	dir         Direction
	dragging    *gm.Vec2
	value       float32
}

func Scrollbar(dir Direction) *WScrollbar {
	return Build(&WScrollbar{dir: dir})
}

func (w *WScrollbar) Build(b *B) {
	const sp = .3
	const spx = .2

	b.BeginPane()
	if w.dir == Horizontal {
		w.track = b.BeginContainer().
			SetPosition(w.handlerSize/2, 0, 0).
			SetSize(w.handlerSize/2, 0)
		w.handler = b.Button().
			SetPivot(.5, 0).
			SetMargin(spx, sp).
			SetSize(w.handlerSize, 0).
			SetAnchor(w.value, 0, w.value, 1)
		b.EndContainer()
	} else {
		w.track = b.BeginContainer().
			SetPosition(0, w.handlerSize/2, 0).
			SetSize(0, w.handlerSize/2)
		w.handler = b.Button().
			SetPivot(0, .5).
			SetMargin(sp, spx).
			SetSize(0, w.handlerSize).
			SetAnchor(0, w.value, 1, w.value)
		b.EndContainer()
	}
	b.EndPane()
	w.SetDragEvents(true)
	event.Handle(w, func(e gorgeui.EventPointerUp) {
		if w.dragging != nil {
			return
		}
		res := w.track.IntersectFromScreen(e.Position)
		var v float32
		if w.dir == Horizontal {
			v = res.UV[0]
		} else {
			v = res.UV[1]
		}
		v = gm.Clamp(v, 0, 1)
		w.SetValue(v)
	})
	event.Handle(w, func(e gorgeui.EventDragBegin) {
		uv := w.track.IntersectFromScreen(e.Position)
		w.dragging = &uv.UV
	})
	event.Handle(w, func(e gorgeui.EventDrag) {
		if w.dragging == nil {
			return
		}

		uv := w.track.IntersectFromScreen(e.Position).UV
		delta := uv.Sub(*w.dragging)
		w.dragging = &uv
		var v float32
		if w.dir == Horizontal {
			v = w.value + delta[0]
		} else {
			v = w.value + delta[1]
		}
		w.SetValue(v)
	})
	event.Handle(w, func(gorgeui.EventDragEnd) {
		w.dragging = nil
	})
}

func (w *WScrollbar) OnChange(fn func(v float32)) *WScrollbar {
	w.changefn = fn
	return w
}

func (w *WScrollbar) SetValue(v float32) *WScrollbar {
	v = gm.Clamp(v, 0, 1)
	if w.value == v {
		return w
	}
	w.value = v
	if w.dir == Horizontal {
		w.handler.SetAnchor(w.value, 0, w.value, 1)
	} else {
		w.handler.SetAnchor(0, w.value, 1, w.value)
	}
	if w.changefn != nil {
		w.changefn(v)
	}

	return w
}

func (w *WScrollbar) SetMax(v float32) *WScrollbar {
	w.handlerSize = v
	if w.dir == Horizontal {
		w.handler.SetSize(w.handlerSize, 0)
		w.track.SetRect(w.handlerSize/2, 0, w.handlerSize/2, 0)
	} else {
		w.handler.SetSize(0, w.handlerSize)
		w.track.SetRect(0, w.handlerSize/2, 0, w.handlerSize/2)
	}
	return w
}
