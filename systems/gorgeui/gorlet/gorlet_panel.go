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

		root := b.SetRoot(Quad())

		container := b.Create(Container())
		root.Add(container)
		root.SetClientArea(container)

		Observe(root, "overflow", func(o Overflow) {
			if cur == o {
				return
			}
			cur = o
			// This will swap containers
			switch o {
			// case Resize:
			case OverflowVisible:
				children := container.Children()
				root.SetClientArea(nil)
				root.Remove(container)
				container = Create(Container())
				for _, c := range children {
					// Reset mask depth
					c.Set("_maskDepth", -1)
					container.Add(c)
				}
				root.Add(container)
				root.SetClientArea(container)
			case OverflowHidden:
				children := container.Children()
				root.SetClientArea(nil)
				root.Remove(container)
				container = Create(Mask())
				for _, c := range children {
					container.Add(c)
				}
				root.Add(container)
				root.SetClientArea(container)
			case OverflowScroll:
				children := container.Children()
				root.SetClientArea(nil)
				root.Remove(container)
				container = Create(Scroll())
				for _, c := range children {
					container.Add(c)
				}
				root.Add(container)
				root.SetClientArea(container)
			}
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
