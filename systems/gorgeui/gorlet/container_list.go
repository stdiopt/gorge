package gorlet

// New: attempt layouting with containers instead

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

func List(s ...float32) Func {
	spacing := float32(0)
	if len(s) > 0 {
		spacing = s[0]
	}
	return func(b *Builder) {
		list := ListLayout{
			Direction: Vertical,
			Spacing:   spacing,
		}

		root := b.Root()
		Observe(b, "spacing", func(v float32) { list.Spacing = v })
		Observe(b, "direction", func(v Direction) { list.Direction = v })

		event.Handle(root, func(gorgeui.EventUpdate) {
			list.Layout(root)
		})
	}
}

func (b *Builder) BeginList(spacing ...float32) *Entity {
	return b.Begin(List(spacing...))
}

func (b *Builder) EndList() {
	b.End()
}
