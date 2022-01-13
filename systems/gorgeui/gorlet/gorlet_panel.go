package gorlet

// Panel creates a panel.
func Panel() Func {
	return Quad()
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
