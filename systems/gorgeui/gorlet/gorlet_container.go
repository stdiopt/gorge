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

func List() Func {
	return func(b *Builder) {
		l := LayoutList(0)
		// Parent widget must redim it self too.
		b.Root().SetLayout(AutoHeight(0))
		b.BeginContainer(l, AutoHeight(0))
		b.ClientArea()
		b.EndContainer()

		b.Observe("spacing", ObsFunc(func(v float32) {
			l.Spacing = v
		}))
	}
}

func (b *Builder) BeginList() *Entity {
	return b.Begin(List())
}

// BeginList creates a container that will organize its children in a list.
/*func (b *Builder) BeginList(spacing float32) *Entity {
	l := LayoutList(spacing)

	r := Create(Container)
	r.SetLayout(l)
	r.Observe("spacing", ObsFunc(func(v float32) {
		log.Println("Set spacing")
		l.Spacing = v
	}))
	b.AddEntity(r)
	b.SetPlacement(func(e *Entity) {
		e.SetPivot(0)
		e.SetAnchor(0)
		e.SetHeight(4)
	})
	log.Println("Observe spacing")
	return r
}*/

// EndList alias to End().
func (b *Builder) EndList() { b.End() }
