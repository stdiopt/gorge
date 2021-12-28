package gorlet

// Container empty container consturctor.
func Container(b *Builder) {}

// BeginContainer begins an empty container.
func (b *Builder) BeginContainer(layout ...Layouter) *Entity {
	b.UseLayout(layout...)
	return b.Begin(Container)
}

// EndContainer alias for End.
func (b *Builder) EndContainer() {
	b.End()
}

// BeginList creates a container that will organize its children in a list.
func (b *Builder) BeginList(spacing float32) *Entity {
	b.UsePlacement(func(e *Entity) {
		e.SetPivot(0)
		e.SetAnchor(0)
		e.SetHeight(4)
	})
	b.UseLayout(LayoutList(spacing))
	return b.Begin(Container)
}

// EndList alias to End().
func (b *Builder) EndList() { b.End() }
