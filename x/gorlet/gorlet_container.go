package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type WContainer struct {
	Widget[WContainer]
	layouter Layouter
}

func Container(c ...gorge.Entity) *WContainer {
	return Build(&WContainer{}).Add(c...)
}

func (w *WContainer) Build(b *B) {
	event.Handle(w, func(gorgeui.EventUpdate) {
		if w.layouter != nil {
			w.layouter.Layout(w)
		}
	})
}

func (w *WContainer) SetLayout(l ...Layouter) *WContainer {
	w.layouter = LayoutMulti(l...)
	return w
}

// Builder

func (b *B) Container() *WContainer {
	w := Container()
	b.Add(w)
	return w
}

func (b *B) BeginContainer(l ...Layouter) *WContainer {
	w := Container().SetLayout(l...)
	b.Begin(w)
	return w
}

func (b *B) EndContainer() {
	b.End()
}
