package gorlet

// Container empty container consturctor.
func Container(b *Builder) {
}

// BeginContainer begins an empty container.
func (b *Builder) BeginContainer() *Element {
	return b.Begin(Container)
}

// EndContainer alias for End.
func (b *Builder) EndContainer() {
	b.End()
}
