package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

func Flex() Func {
	return func(b *Builder) {
		var (
			sizes   []float32
			sum     float32
			smaller float32

			dir     Direction
			spacing float32
		)
		root := b.Root()

		Observe(b, "sizes", func(v []float32) {
			sizes = v
			sum = 0
			smaller = sizes[0]
			for _, f := range sizes {
				sum += f
				smaller = gm.Min(smaller, f)

			}
		})
		Observe(b, "direction", func(v Direction) {
			dir = v
		})
		Observe(b, "spacing", func(v float32) {
			spacing = v
		})

		event.Handle(root, func(gorgeui.EventUpdate) {
			children := root.Children()
			esum := sum // effective sum
			if d := len(children) - len(sizes); d > 0 {
				esum = sum + float32(d)*smaller
			}
			var start float32
			for i, e := range children {
				sz := smaller
				if i < len(sizes) {
					sz = sizes[i]
				}

				end := start + sz/esum
				switch dir {
				case Horizontal:
					if i != 0 {
						e.SetRect(spacing, 0, 0, 0)
					}
					e.SetAnchor(start, 0, end, 1)
				case Vertical:
					if i != 0 {
						e.SetRect(0, spacing, 0, 0)
					}
					e.SetAnchor(0, start, 1, end)
				}
				start = end
			}
		})
	}
}

func (b *Builder) BeginFlex() *Entity {
	return b.Begin(Flex())
}

func (b *Builder) EndFlex() {
	b.End()
}
