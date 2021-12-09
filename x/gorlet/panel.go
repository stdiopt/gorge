package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/systems/gorgeui/widget"
)

// Panel creates a panel.
func Panel() BuildFunc {
	return func(b *Builder) {
		p := b.Root()
		ent := widget.QuadEntity()
		gorgeui.AddElementTo(p, ent)
		p.HandleFunc(func(e event.Event) {
			if _, ok := e.(gorgeui.EventUpdate); !ok {
				return
			}
			r := p.Rect()
			ent.Scale[0] = r[2] - r[0]
			ent.Scale[1] = r[3] - r[1]
		})
		b.Observe("color", func(c m32.Vec4) {
			ent.SetColorv(c)
		})

		// Defaults
		p.Set("color", m32.Color(0, 0, 0, .2))
	}
}

// BeginPanel begins a panel.
func (b *Builder) BeginPanel() *Entity {
	return b.Begin(Panel())
}

// EndPanel alias to b.End()
func (b *Builder) EndPanel() {
	b.End()
}
