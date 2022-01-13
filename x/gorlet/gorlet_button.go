package gorlet

import (
	"github.com/stdiopt/gorge"
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
		b.Observe("color", Ptr(&normal))
		b.Observe("highlight", Ptr(&highlight))
		b.Observe("down", Ptr(&down))
		b.Observe("fadeFactor", Ptr(&fadeFactor))

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
		gorge.HandleFunc(root, func(e gorgeui.EventUpdate) {
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
		})
		gorge.HandleFunc(root, func(gorgeui.EventPointerDown) {
			state |= statePressed
		})
		gorge.HandleFunc(root, func(gorgeui.EventPointerUp) {
			state &= ^statePressed
			gorge.Trigger(root, EventClick{})
			if clickfn != nil {
				clickfn()
			}
		})
		gorge.HandleFunc(root, func(gorgeui.EventPointerEnter) {
			state |= stateHover
		})
		gorge.HandleFunc(root, func(gorgeui.EventPointerLeave) {
			state &= ^stateHover
		})

		/*
			root.HandleFunc(func(e event.Event) {
				switch e := e.(type) {
				case gorgeui.EventUpdate:
				case gorgeui.EventPointerDown:
				case gorgeui.EventPointerUp:
				case gorgeui.EventPointerEnter:
				case gorgeui.EventPointerLeave:

				}
			})
		*/
	}
}

// BeginButton begins a button.
func (b *Builder) BeginButton(clickfn func()) *Entity {
	return b.Begin(Button(clickfn))
}

// EndButton alias to End.
func (b *Builder) EndButton() { b.End() }
