package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

func Grid() Func {
	return func(b *B) {
		grid := GridLayout{}

		Observe(b, "cols", func(v int) { grid.Cols = v })
		Observe(b, "rows", func(v int) { grid.Rows = v })
		Observe(b, "spacing", func(v float32) { grid.Spacing = v })

		root := b.Root()
		event.Handle(root, func(gorgeui.EventUpdate) {
			grid.Layout(root)
		})
	}
}

func (b *B) BeginGrid() *Entity {
	return b.Begin(Grid())
}

func (b *B) EndGrid() {
	b.End()
}
