package gorlet

import "github.com/stdiopt/gorge/math/gm"

type Overflow int

const (
	OverflowVisible = Overflow(iota)
	OverflowHidden
	OverflowScroll
)

func Panel() Func {
	return func(b *B) {
		var (
			cur        Overflow
			scrollSize = gm.Vec2{1, 1}
		)
		// Need to store state about panel, mask and Scroll too

		// This way we can set props on Quad
		root := b.SetRoot(Quad())
		wrapper := b.BeginContainer()
		container := b.BeginContainer()
		b.ClientArea()
		b.EndContainer()
		b.EndContainer()

		Observe(root, "hscrollSize", Ptr(&scrollSize[0]))
		Observe(root, "vscrollSize", Ptr(&scrollSize[1]))
		Observe(root, "scrollSize", func(v float32) { scrollSize = gm.V2(v) })

		Observe(root, "overflow", func(o Overflow) {
			if cur == o {
				return
			}
			cur = o
			root.SetClientArea(nil)
			root.Remove(wrapper)
			// This will swap containers
			switch o {
			// case Resize:
			case OverflowHidden:
				wrapper = Create(Mask())
			case OverflowScroll:
				wrapper = Create(Scroll())
				wrapper.Set("hscrollSize", scrollSize[0])
				wrapper.Set("vscrollSize", scrollSize[1])
			default:
				container.Set("_maskDepth", -1)
				wrapper = Create(Container())
			}
			wrapper.Add(container)
			root.Add(wrapper)
			root.SetClientArea(container)
		})
	}
}

// BeginPanel begins a panel.
func (b *B) BeginPanel(ls ...Layouter) *Entity {
	b.UseLayout(ls...)
	return b.Begin(Panel())
}

// EndPanel alias to b.End()
func (b *B) EndPanel() {
	b.End()
}
