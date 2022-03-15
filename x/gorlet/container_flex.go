package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type WFlex struct {
	Widget[*WFlex]

	layout FlexLayout
}

func Flex(vs ...float32) *WFlex {
	return Build(&WFlex{layout: *layoutFlex(Horizontal, vs...)})
}

func (w *WFlex) Build(b *B) {
	event.Handle(w, func(gorgeui.EventUpdate) {
		w.layout.Layout(w)
	})
}

func (w *WFlex) SetSizes(vs ...float32) *WFlex {
	w.layout.SetSizes(vs...)
	return w
}

func (w *WFlex) SetDirection(v Direction) *WFlex {
	w.layout.SetDirection(v)
	return w
}

func (w *WFlex) SetSpacing(v float32) *WFlex {
	w.layout.SetSpacing(v)
	return w
}

// Builder

func (b *B) BeginFlex(sizes ...float32) *WFlex {
	w := Flex(sizes...)
	b.Begin(w)
	return w
}

func (b *B) EndFlex() {
	b.End()
}
