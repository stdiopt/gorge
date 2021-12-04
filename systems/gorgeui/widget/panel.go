package widget

import (
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// Panel stage for the panel widget.
type Panel struct {
	Component

	Color  m32.Vec4
	Entity *Quad
}

// HandleEvent handles gorgeui events.
func (p *Panel) HandleEvent(e gorgeui.Event) {
	if _, ok := e.Value.(gorgeui.EventUpdate); !ok {
		return
	}
	r := p.Rect()
	// log.Println("EVENT Current rect:", r)
	// p.Entity.Position[0] = r[0]
	// p.Entity.Position[1] = r[1] // bottom
	p.Entity.Scale[0] = r[2] - r[0]
	p.Entity.Scale[1] = r[3] - r[1]
	p.Entity.SetColorv(p.Color)
}

// NewPanel returns a new default Panel.
func NewPanel() *Panel {
	p := &Panel{
		Component: *NewComponent(),
		Color:     m32.Vec4{0, 0, 0, .5},
		Entity:    QuadEntity(),
	}
	gorgeui.AddGraphicTo(p, p.Entity)
	return p
}

// SetColor set panel options.
func (p *Panel) SetColor(c ...float32) {
	p.Color = v4Color(c...)
}
