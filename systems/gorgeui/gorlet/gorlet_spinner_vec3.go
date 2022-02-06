package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
)

// SpinnerVec3 creates 3 spinners to control a gm.Vec3
func SpinnerVec3(fn func(gm.Vec3)) Func {
	return func(b *B) {
		var (
			// props
			fontScale   = b.Prop("fontScale", 2)
			background  = b.Prop("background", nil)
			border      = b.Prop("border")
			borderColor = b.Prop("borderColor")

			labelColorX = b.Prop("x.labelColor", gm.Color(.5, 0, 0))
			labelColorY = b.Prop("y.labelColor", gm.Color(0, .5, 0))
			labelColorZ = b.Prop("z.labelColor", gm.Color(0, 0, .5))

			// spinners
			x   *Entity
			y   *Entity
			z   *Entity
			val gm.Vec3
		)
		b.Push("fontScale", fontScale)

		root := b.Root()
		b.UseProps(Props{
			"color":       background,
			"border":      border,
			"borderColor": borderColor,
		})
		b.BeginPanel(LayoutFlexHorizontal(1))
		{

			obsFn := func(i int) func(v float32) {
				return func(v float32) {
					cp := val
					cp[i] = v
					root.Set("value", cp)
				}
			}
			b.Use("labelColor", labelColorX)
			x = b.Spinner("X", obsFn(0))

			b.Use("labelColor", labelColorY)
			y = b.Spinner("Y", obsFn(1))

			b.Use("labelColor", labelColorZ)
			z = b.Spinner("Z", obsFn(2))

		}
		b.EndPanel()

		Observe(b, "value", func(v gm.Vec3) {
			if val == v {
				return
			}
			val = v
			x.Set("value", v[0])
			y.Set("value", v[1])
			z.Set("value", v[2])
			if fn != nil {
				fn(val)
			}
			event.Trigger(root, EventValueChanged{Value: val})
		})
	}
}

// SpinnerVec3 adds a spinner to builder.
func (b *B) SpinnerVec3(fn func(gm.Vec3)) *Entity {
	return b.Add(SpinnerVec3(fn))
}
