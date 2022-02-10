package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type WGrid struct {
	Widget[WGrid]

	layout GridLayout
}

func Grid(c ...gorge.Entity) *WGrid {
	return Build(&WGrid{
		layout: GridLayout{
			Cols:    1,
			Rows:    1,
			Spacing: .1,
		},
	}).Add(c...)
}

func (w *WGrid) Build(b *B) {
	event.Handle(w, func(gorgeui.EventUpdate) {
		w.layout.Layout(w)
	})
}

func (w *WGrid) SetDimensions(cols, rows int) *WGrid {
	w.layout.SetDimensions(cols, rows)
	return w
}

func (w *WGrid) SetRows(rows int) *WGrid {
	w.layout.Rows = rows
	return w
}

func (w *WGrid) SetCols(cols int) *WGrid {
	w.layout.Cols = cols
	return w
}

func (w *WGrid) SetSpacing(c float32) *WGrid {
	w.layout.Spacing = c
	return w
}

func (b *B) BeginGrid(cols, rows int) *WGrid {
	w := Grid().SetDimensions(cols, rows)
	b.Begin(w)
	return w
}

func (b *B) EndGrid() {
	b.End()
}
