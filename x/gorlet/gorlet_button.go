package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// Button creates a simple button.
func Button(clickfn func()) Func {
	return func(b *Builder) {
		var (
			normal     = m32.Color(.7, .9)
			highlight  = m32.Color(.8, .9, .8)
			down       = m32.Color(.4)
			fadeFactor = float32(10)
		)
		b.BindProp("color", &normal)
		b.BindProp("highlight", &highlight)
		b.BindProp("down", &down)
		b.BindProp("fadeFactor", &fadeFactor)

		root := b.Root()

		b.Use("color", normal)
		p := b.BeginPanel()
		b.ClientArea()
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
					s *= 2
					target = down
				case state&stateHover != 0:
					target = highlight
				}
				// Due to floating point this might run everytime but
				// it is somewhat ok since comparing with epsilon might be slower
				if target != color {
					color = color.Lerp(target, e.DeltaTime()*s)
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
	}
}

// BeginButton begins a button.
func (b *Builder) BeginButton(clickfn func()) *Entity {
	return b.Begin(Button(clickfn))
}

// EndButton alias to End.
func (b *Builder) EndButton() { b.End() }
