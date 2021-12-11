package gorlet

import (
	"fmt"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// Spinner creates a new spinner.
func Spinner(lbl string) BuildFunc {
	return func(b *Builder) {
		var val float32 = -1

		b.Set("fontScale", b.Prop("fontScale", 2))
		root := b.Root()
		b.Layout(gorgeui.FlexHorizontal(1, 2))
		b.BeginPanel()
		{

			b.Set("color", b.Prop("labelColor", m32.Color(1)))
			b.BeginPanel()
			{
				b.Set("color", b.Prop("labelTextColor", m32.Color(1)))
				b.Label(lbl)
			}
			b.End()
		}
		b.Set("color", b.Prop("textColor", m32.Color(1)))

		l := b.Label("")
		b.EndPanel()

		b.Observe("value", func(v float32) {
			if val == v {
				return
			}
			val = v
			l.Set("text", fmt.Sprintf("%.2f", val))

			root.Trigger(EventValueChanged{val})
		})

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
func (b *Builder) Spinner(t string) *Entity {
	return b.Add(Spinner(t))
}
