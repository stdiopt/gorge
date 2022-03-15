package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type WList struct {
	Widget[*WList]

	layout ListLayout

	autoHeight Layouter
}

func List(c ...gorge.Entity) *WList {
	return Build(&WList{layout: *LayoutList(.1)}).Add(c...)
}

func (w *WList) Build(b *B) {
	event.Handle(w, func(gorgeui.EventUpdate) {
		w.layout.Layout(w)
		if w.autoHeight != nil {
			w.autoHeight.Layout(w)
		}
	})
}

func (w *WList) SetAutoHeight(b bool) *WList {
	if b {
		w.autoHeight = AutoHeight(0)
	} else {
		w.autoHeight = nil
	}
	return w
}

func (w *WList) SetSpacing(v float32) *WList {
	w.layout.SetSpacing(v)
	return w
}

func (w *WList) SetDirection(dir Direction) *WList {
	w.layout.SetDirection(dir)
	return w
}

func (b *B) BeginList(spacing float32) *WList {
	w := List().
		SetSpacing(spacing).
		SetDirection(Vertical)
	b.Begin(w)
	return w
}

func (b *B) EndList() {
	b.End()
}
