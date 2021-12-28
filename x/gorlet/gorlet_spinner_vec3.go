package gorlet

import (
	"github.com/stdiopt/gorge/m32"
)

// SpinnerVec3 creates 3 spinners to control a m32.Vec3
func SpinnerVec3(fn func(m32.Vec3)) Func {
	return func(b *Builder) {
		var (
			fontScale   = b.Prop("fontScale", 2)
			background  = b.Prop("background", nil)
			labelColorX = b.Prop("x.labelColor", m32.Color(.5, 0, 0))
			labelColorY = b.Prop("y.labelColor", m32.Color(0, .5, 0))
			labelColorZ = b.Prop("z.labelColor", m32.Color(0, 0, .5))
		)
		var val m32.Vec3

		root := b.Root()
		obsFn := func(i int) func(v float32) {
			return func(v float32) {
				val[i] = v
				root.Set("value", val)
				if fn != nil {
					fn(val)
				}
			}
		}

		b.Global("fontScale", fontScale)

		b.Use("color", background)
		b.BeginPanel(LayoutFlexHorizontal(1))

		b.Use("labelColor", labelColorX)
		x := b.Spinner("X", obsFn(0))

		b.Use("labelColor", labelColorY)
		y := b.Spinner("Y", obsFn(1))

		b.Use("labelColor", labelColorZ)
		z := b.Spinner("Z", obsFn(2))

		b.EndPanel()

		b.Observe("value", ObsFunc(func(v m32.Vec3) {
			if val == v {
				return
			}
			val = v
			x.Set("value", v[0])
			y.Set("value", v[1])
			z.Set("value", v[2])
			fn(val)
			root.Trigger(EventValueChanged{Value: val})
		}))
	}
}

// SpinnerVec3 adds a spinner to builder.
func (b *Builder) SpinnerVec3(fn func(m32.Vec3)) *Entity {
	return b.Add(SpinnerVec3(fn))
}
