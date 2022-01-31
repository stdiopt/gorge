package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

func Flex(sz ...float32) Func {
	return func(b *Builder) {
		flex := FlexLayout{}
		root := b.Root()

		Observe(b, "sizes", func(v []float32) { flex.SetSizes(v...) })
		Observe(b, "direction", func(v Direction) { flex.Direction = v })
		Observe(b, "spacing", func(v float32) { flex.Spacing = v })

		// This way we can add layouters without messing with flex.
		event.Handle(root, func(gorgeui.EventUpdate) {
			flex.Layout(root)
		})
		if len(sz) > 0 {
			flex.SetSizes(sz...)
		}
	}
}

func (b *Builder) BeginFlex(sz ...float32) *Entity {
	return b.Begin(Flex(sz...))
}

func (b *Builder) EndFlex() {
	b.End()
}
