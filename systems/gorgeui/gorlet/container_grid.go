package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

func Grid() Func {
	return func(b *Builder) {
		var cols, rows int
		var spacing float32

		Observe(b, "cols", func(v int) { cols = v })
		Observe(b, "rows", func(v int) { rows = v })
		Observe(b, "spacing", func(v float32) { spacing = v })

		root := b.Root()
		event.Handle(root, func(gorgeui.EventUpdate) {
			sw := 1 / float32(cols)
			sh := 1 / float32(rows)
			for i, e := range root.Children() {
				cw := float32(i%cols) / float32(cols)
				ch := float32(i/cols) / float32(rows)
				e.SetAnchor(cw, ch, cw+sw, ch+sh)
				s := gm.Vec4{spacing / 2, spacing / 2, spacing / 2, spacing / 2}
				// s := gm.Vec4{0, 0, spacing, spacing}
				if cw == 0 {
					s[0] = 0
				} else if cw+sw == 1 {
					s[2] = 0
				}
				if ch == 0 {
					s[1] = 0
				} else if ch+sh == 1 {
					s[3] = 0
				}

				e.SetRect(s[:]...)
			}
		})
	}
}

func (b *Builder) BeginGrid() *Entity {
	return b.Begin(Grid())
}

func (b *Builder) EndGrid() {
	b.End()
}
