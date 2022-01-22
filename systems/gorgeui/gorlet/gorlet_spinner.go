package gorlet

import (
	"fmt"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/text"
)

// Spinner creates a new spinner.
func Spinner(lbl string, fn func(float32)) Func {
	return func(b *Builder) {
		var (
			fontScale      = b.Prop("fontScale", 2)
			labelColor     = b.Prop("labelColor", gm.Color(1))
			labelTextColor = b.Prop("labelTextColor", gm.Color(1))
			textColor      = b.Prop("textColor", gm.Color(1))
			textOverflow   = b.Prop("textOverflow", text.OverflowOverlap)
		)
		var val float32 = -1

		b.Push("fontScale", fontScale)

		root := b.Root()
		root.SetDragEvents(true)
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
		b.Use("overflow", textOverflow)
		l := b.Label("")
		b.EndPanel()

		b.Observe("value", func(v float32) {
			if val == v {
				return
			}
			val = v
			l.Set("text", fmt.Sprintf("%.2f", val))
			if fn != nil {
				fn(val)
			}
			gorge.Trigger(root, EventValueChanged{val})
		})

		// root.SetDragEvents(true)
		event.Handle(root, func(e gorgeui.EventDrag) {
			root.Set("value", val+e.Delta[0]*0.01)
		})
		root.Set("value", 0)
	}
}

// Spinner add a spinner to the child.
func (b *Builder) Spinner(t string, fn func(float32)) *Entity {
	return b.Add(Spinner(t, fn))
}
