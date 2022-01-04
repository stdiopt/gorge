package gorlet

import (
	"fmt"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/m32/ray"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// EventF32Changed triggers on certain widdgets a value change.

// Slider guilet.
func Slider(min, max float32, fn func(float32)) Func {
	return func(b *Builder) {
		var (
			fontScale        = b.Prop("fontScale")
			backgroundColor  = b.Prop("backgroundColor", m32.Color(.4, .2))
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

		b.Use("color", backgroundColor)

		// b.UsePivot(0)
		// b.UseAnchor(0)
		// b.UseRect(0, 0, 30, 4)
		root := b.SetRoot(Panel())
		root.SetDragEvents(true)

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

				b.UsePivot(.5)
				b.UseRect(0, 0, handlerSize, 0)
				b.UseAnchor(val, 0, val, 1)
				handler = b.TextButton("0", nil)
			}
			b.End()
		}

		b.Observe("handlerColor", ObsFunc(func(c m32.Vec4) {
			handler.Set("color", c)
		}))
		b.Observe("handler", ObsFunc(func(e *Entity) {
			// Need to remove Element first :/
			// should remove observers from handler?
			track.Remove(handler)
			handler = e
			track.Add(handler)

			handler.SetPivot(.5)
			handler.SetRect(0, 0, handlerSize, 0)
			handler.SetAnchor(val, 0, val, 1)
		}))
		b.Observe("value", ObsFunc(func(v float32) {
			if val == v {
				return
			}
			val = v
			handler.SetAnchor(val, 0, val, 1)
			handler.Set("text", fmt.Sprintf(valFmt, min+(val*(max-min))))
			if fn != nil {
				fn(val)
			}
			root.Trigger(EventValueChanged{val})
		}))
		b.Observe("handlerSize", ObsFunc(func(f float32) {
			handlerSize = f
			handler.SetRect(0, 0, handlerSize/2, 0)
			track.SetRect(handlerSize/2, 0, handlerSize/2, 0)
		}))
		b.Observe("textFormat", ObsFunc(func(s string) {
			valFmt = s
			handler.Set("text", fmt.Sprintf(valFmt, min+(val*(max-min))))
		}))

		var dragging bool
		root.HandleFunc(func(e event.Event) {
			switch e := e.(type) {
			case gorgeui.EventPointerUp:
				if dragging {
					return
				}
				res := e.RayResult
				r := track.Rect()
				fullw := r[2] - r[0]

				wp := root.WorldPosition()
				v := (res.Position[0] - (wp[0] + r[0])) / fullw // Ray in thing position
				v -= handlerSize / fullw / 2
				v = m32.Clamp(v, 0, 1)
				root.Set("value", v)
			case gorgeui.EventDrag:
				dragging = true
				rect := track.Rect()
				fullw := rect[2] - rect[0]

				m := track.Mat4()
				v0 := m.MulV4(m32.Vec4{rect[0], rect[1], 0, 1}).Vec3()
				v1 := m.MulV4(m32.Vec4{rect[2], rect[1], 0, 1}).Vec3() // right
				v2 := m.MulV4(m32.Vec4{rect[0], rect[3], 0, 1}).Vec3() // up)

				ui := gorgeui.RootUI(root)
				r := ray.FromScreen(ui.ScreenSize(), ui.Camera, e.Position)
				res := ray.IntersectRect(r, v0, v1, v2)

				wp := root.WorldPosition()
				v := (res.Position[0] - (wp[0] + rect[0])) / fullw // Ray in thing position
				v -= handlerSize / fullw / 2
				v = m32.Clamp(v, 0, 1)

				root.Set("value", v)
			case gorgeui.EventDragEnd:
				dragging = false
			}
		})

		root.Set("value", float32(0))
	}
}

// Slider adds a slider to current element.
func (b *Builder) Slider(min, max float32, fn func(float32)) *Entity {
	return b.Add(Slider(min, max, fn))
}
