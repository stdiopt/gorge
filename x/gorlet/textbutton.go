package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/text"
)

// EventClick is trigger when an event is clicked.
type EventClick struct{} // need more info

// TextButton creates a text button.
func TextButton(t string, clickfn func()) BuildFunc {
	return func(b *Builder) {
		var (
			normal     = m32.Color(.7, .9)
			highlight  = m32.Color(.8, .9, .8)
			down       = m32.Color(.4)
			fadeFactor = float32(10)
		)

		root := b.Root()
		b.SetAddMode(ElementAdd)

		b.Set("color", normal)
		p := b.BeginPanel()
		{
			b.Set("alignment", []text.AlignType{text.AlignCenter, text.AlignCenter})
			// This part is amazing, forwarding a property
			b.Set("text", Prop("text"))
			b.Set("fontScale", Prop("fontScale"))
			b.Set("textColor", Prop("textColor", m32.Color(0)))
			b.Label(t)
		}
		b.EndPanel()

		color := normal
		type buttonState int
		const (
			statePressed = 1 << (iota + 1)
			stateHover
		)
		var state buttonState
		root.HandleFunc(func(e event.Event) {
			switch e := e.(type) {
			case gorgeui.EventUpdate:
				s := fadeFactor
				target := normal
				switch {
				case state&statePressed != 0:
					s *= 10
					target = down
				case state&stateHover != 0:
					target = highlight
				}
				color = color.Lerp(target, e.DeltaTime()*s)
				if target != color {
					p.Set("color", color)
				}
			case gorgeui.EventPointerDown:
				state |= statePressed
			case gorgeui.EventPointerUp:
				state &= ^statePressed
				root.Trigger(EventClick{})
				if clickfn != nil {
					clickfn()
				}
			case gorgeui.EventPointerEnter:
				state |= stateHover
			case gorgeui.EventPointerLeave:
				state &= ^stateHover

			}
		})
		root.Set("text", t)
	}
}

// TextButton add a text button child.
func (b *Builder) TextButton(t string, clickfn func()) *Element {
	return b.Add(TextButton(t, clickfn))
}
