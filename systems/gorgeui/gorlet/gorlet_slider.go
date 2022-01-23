package gorlet

import (
	"fmt"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/math/ray"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/text"
)

// TODO: {lpf} find a way to turn updated on drag on and off.

// Slider guilet.
func Slider(min, max float32, fn func(float32)) Func {
	// min-max to 0-1
	norm := func(v float32) float32 {
		return (v - min) / (max - min)
	}
	// 0-1 to min-max
	real := func(v float32) float32 {
		return min + v*(max-min)
	}
	return func(b *Builder) {
		var (
			fontScale        = b.Prop("fontScale")
			backgroundColor  = b.Prop("backgroundColor", gm.Color(.4, .2))
			handlerTextColor = b.Prop("textColor")
			handlerColor     = b.Prop("handlerColor")
		)
		var (
			valFmt      = "%.2f"
			handlerSize = float32(4)
			val         float32
			track       *Entity
			handler     *Entity
		)

		root := b.Root()
		root.SetDragEvents(true)
		b.Use("color", backgroundColor)
		b.BeginPanel()
		{
			b.UseAnchor(0, 0, 1, 1)
			b.UseRect(handlerSize/2, 0, handlerSize/2, 0)
			b.UsePivot(0, .5)
			track = b.BeginContainer()
			{
				// track.SetRect(handlerSize/2, 0, handlerSize/2, 0)

				// this will be stuck forever :/
				b.Use("textColor", handlerTextColor)
				b.Use("color", handlerColor)
				b.Use("fontScale", fontScale)
				b.Use("textOverflow", text.OverflowOverlap)

				b.UsePivot(.5)
				b.UseRect(0, 0, handlerSize, 0)
				b.UseAnchor(val, 0, val, 1)
				handler = b.TextButton("0", nil)
			}
			b.End()
		}
		b.EndContainer()

		b.Observe("min", func(v float32) { min = v })
		b.Observe("max", func(v float32) { max = v })
		b.Observe("handlerColor", func(c gm.Vec4) {
			handler.Set("color", c)
		})
		b.Observe("handler", func(e *Entity) {
			// Need to remove Element first :/
			// should remove observers from handler?
			track.Remove(handler)
			handler = e
			track.Add(handler)

			handler.SetPivot(.5)
			handler.SetRect(0, 0, handlerSize, 0)
			handler.SetAnchor(val, 0, val, 1)
		})
		b.Observe("value", func(v float32) {
			v = norm(v)
			v = gm.Clamp(v, 0, 1)
			if val == v {
				return
			}
			val = v
			rval := real(val)
			handler.SetAnchor(val, 0, val, 1)
			handler.Set("text", fmt.Sprintf(valFmt, rval))
			if fn != nil {
				fn(rval)
			}
			gorge.Trigger(root, EventValueChanged{val})
		})
		b.Observe("handlerSize", func(f float32) {
			handlerSize = f
			handler.SetRect(0, 0, handlerSize/2, 0)
			track.SetRect(handlerSize/2, 0, handlerSize/2, 0)
		})
		b.Observe("textFormat", func(s string) {
			valFmt = s
			handler.Set("text", fmt.Sprintf(valFmt, real(val)))
		})

		var dragging bool
		dodrag := func(pd *gorgeui.PointerData) {
			rect := track.Rect()

			// We ray trace on track
			m := root.Mat4()
			v0 := m.MulV4(gm.Vec4{rect[0], rect[1], 0, 1}).Vec3() // 0
			v1 := m.MulV4(gm.Vec4{rect[2], rect[1], 0, 1}).Vec3() // right
			v2 := m.MulV4(gm.Vec4{rect[0], rect[3], 0, 1}).Vec3() // up)

			ui := gorgeui.RootUI(root)
			r := ray.FromScreen(ui.ScreenSize(), ui.Camera, pd.Position)
			res := ray.IntersectRect(r, v0, v1, v2)

			// wp := root.WorldPosition()
			// v := (res.Position[0] - wp[0]) / (v1[0] - v0[0])
			v := res.UV[0]

			v -= handlerSize / rect[2] / 2
			v = gm.Clamp(v, 0, 1)
			root.Set("value", real(v))
			// log.Println("Res:", res.Position[0])
		}
		event.Handle(root, func(e gorgeui.EventPointerUp) {
			if dragging {
				return
			}
			dodrag(e.PointerData)
		})
		event.Handle(root, func(e gorgeui.EventDrag) {
			dragging = true
			dodrag(e.PointerData)
		})
		event.Handle(root, func(gorgeui.EventDragEnd) {
			dragging = false
		})

		root.Set("value", float32(0))
	}
}

// Slider adds a slider to current element.
func (b *Builder) Slider(min, max float32, fn func(float32)) *Entity {
	return b.Add(Slider(min, max, fn))
}
