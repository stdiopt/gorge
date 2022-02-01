package gorlet

type Overflow int

const (
	OverflowVisible = Overflow(iota)
	OverflowHidden
	OverflowScroll
)

func Panel() Func {
	return func(b *Builder) {
		var cur Overflow
		// Need to store state about panel, mask and Scroll too

		// This way we can set props on Quad
		root := b.SetRoot(Quad())
		wrapper := b.BeginContainer()
		container := b.BeginContainer()
		b.ClientArea()
		b.EndContainer()
		b.EndContainer()

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
				wrapper = Create(scroll())
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
func (b *Builder) BeginPanel(ls ...Layouter) *Entity {
	b.UseLayout(ls...)
	return b.Begin(Panel())
}

// EndPanel alias to b.End()
func (b *Builder) EndPanel() {
	b.End()
}
