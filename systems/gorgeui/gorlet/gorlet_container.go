package gorlet

// Container empty container consturctor.
func Container() Func {
	return func(b *B) {}
}

// BeginContainer begins an empty container.
func (b *B) BeginContainer(layout ...Layouter) *Entity {
	b.UseLayout(layout...)
	return b.Begin(Container())
}

// EndContainer alias for End.
func (b *B) EndContainer() {
	b.End()
}

/*
func List() Func {
	return func(b *Builder) {
		l := LayoutList(0)
		// Parent widget must redim it self too.
		b.Root().SetLayout(AutoHeight(0))
		b.BeginContainer(l, AutoHeight(0))
		b.ClientArea()
		b.EndContainer()

		b.Observe("spacing", func(v float32) {
			l.Spacing = v
		})
	}
}

func (b *Builder) BeginList() *Entity {
	return b.Begin(List())
}


// EndList alias to End().
func (b *Builder) EndList() { b.End() }
*/
