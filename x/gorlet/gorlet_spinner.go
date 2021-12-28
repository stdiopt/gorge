package gorlet

import (
	"fmt"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// Spinner creates a new spinner.
func Spinner(lbl string, fn func(float32)) Func {
	return func(b *Builder) {
		var (
			fontScale      = b.Prop("fontScale", 2)
			labelColor     = b.Prop("labelColor", m32.Color(1))
			labelTextColor = b.Prop("labelTextColor", m32.Color(1))
			textColor      = b.Prop("textColor", m32.Color(1))
		)
		var val float32 = -1

		b.Global("fontScale", fontScale)

		root := b.Root()
		b.UseLayout(LayoutFlexHorizontal(1, 2))
		b.BeginPanel()
		{

			b.Use("color", labelColor)
			b.BeginPanel()
			{
				b.Use("color", labelTextColor)
				b.Label(lbl)
			}
			b.End()
		}
		b.Use("color", textColor)
		l := b.Label("")
		b.EndPanel()

		b.Observe("value", ObsFunc(func(v float32) {
			if val == v {
				return
			}
			val = v
			l.Set("text", fmt.Sprintf("%.2f", val))
			if fn != nil {
				fn(val)
			}
			root.Trigger(EventValueChanged{val})
		}))

		root.SetDragEvents(true)
		root.HandleFunc(func(e event.Event) {
			switch e := e.(type) {
			case gorgeui.EventDrag:
				root.Set("value", val+e.Delta[0]*0.01)
			default:
			}
		})
		root.Set("value", 0)
	}
}

// Spinner add a spinner to the child.
func (b *Builder) Spinner(t string, fn func(float32)) *Entity {
	return b.Add(Spinner(t, fn))
}
