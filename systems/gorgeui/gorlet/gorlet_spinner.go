package gorlet

import (
	"fmt"

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
			border         = b.Prop("border")
			borderColor    = b.Prop("borderColor")
			labelColor     = b.Prop("labelColor", gm.Color(1))
			labelTextColor = b.Prop("labelTextColor", gm.Color(1))
			textColor      = b.Prop("textColor", gm.Color(1))
			textOverflow   = b.Prop("textOverflow", text.OverflowOverlap)

			val float32 = -1

			valueLbl *Entity
		)

		b.Push("fontScale", fontScale)

		root := b.Root()
		root.SetDragEvents(true)
		b.UseProps(Props{
			"border":      border,
			"borderColor": borderColor,
		})
		b.BeginPanel(LayoutFlexHorizontal(1, 2))
		{

			b.Use("color", labelColor)
			b.BeginPanel()
			{
				b.Use("color", labelTextColor)
				b.Label(lbl)
			}
			b.End()
			b.UseProps(Props{"color": textColor, "overflow": textOverflow})
			valueLbl = b.Label("")
		}
		b.EndPanel()

		Observe(b, "value", func(v float32) {
			if val == v {
				return
			}
			val = v
			valueLbl.Set("text", fmt.Sprintf("%.2f", val))
			if fn != nil {
				fn(val)
			}
			event.Trigger(root, EventValueChanged{val})
		})

		var last gm.Vec2
		var delta float32
		event.Handle(root, func(e gorgeui.EventDragBegin) {
			res := root.IntersectFromScreen(e.Position)
			last = res.UV
		})
		event.Handle(root, func(e gorgeui.EventDrag) {
			res := root.IntersectFromScreen(e.Position)
			delta = res.UV[0] - last[0]
			last = res.UV
			factor := .01 + gm.Abs(res.UV[1]-.5)
			root.Set("value", val+delta*factor)
		})
		root.Set("value", 0)
	}
}

// Spinner add a spinner to the child.
func (b *Builder) Spinner(t string, fn func(float32)) *Entity {
	return b.Add(Spinner(t, fn))
}
