package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// TextButton creates a text button.
func TextButton(t string, clickfn func()) Func {
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

		// Setting this here is not a good idea?.
		// b.UseAnchor(0)
		// b.UseRect(0, 0, 30, 4)
		// b.UsePivot(0)
		// root := b.SetRoot(Panel())
		root := b.Root()

		b.Use("color", normal)
		b.UseAnchor(0, 0, 1, 1)
		b.UseRect(0)
		b.UsePivot(0)
		p := b.BeginPanel()
		{
			b.SetProps(Props{
				"text":      b.Prop("text", t),
				"fontScale": b.Prop("fontScale", 2),
				// "textAlign": b.Prop("textAlign", []text.Align{text.AlignCenter, text.AlignCenter}),
				"textColor": b.Prop("textColor", m32.Color(0)),
			})
			b.UseRect(0)
			b.UseAnchor(0, 0, 1, 1)
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
		// root.Set("text", t)
	}
}

// TextButton add a text button child.
func (b *Builder) TextButton(t string, clickfn func()) *Entity {
	return b.Add(TextButton(t, clickfn))
}
