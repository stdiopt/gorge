package gorlet

// New: attempt layouting with containers instead

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

func List() Func {
	return func(b *Builder) {
		var (
			spacing float32
			dir     = Vertical
		)

		root := b.Root()
		Observe(b, "spacing", func(v float32) { spacing = v })
		Observe(b, "direction", func(v Direction) { dir = v })

		event.Handle(root, func(gorgeui.EventUpdate) {
			cur := float32(0)
			children := root.Children()
			for _, e := range children {
				rt := e.RectTransform()
				r := rt.Rect()
				switch dir {
				case Vertical:
					rt.SetAnchor(0, 0, 1, 0)
					d := r[3] - r[1] + rt.Margin[1] + rt.Margin[3]
					rt.Position[1] = cur
					cur += d + spacing
				case Horizontal:
					rt.SetAnchor(0, 0, 0, 1)
					d := r[2] - r[0] + rt.Margin[0] + rt.Margin[2]
					rt.Position[0] = cur
					cur += d + spacing
				}
			}
		})
	}
}

func (b *Builder) BeginList() *Entity {
	return b.Begin(List())
}

func (b *Builder) EndList() {
	b.End()
}
