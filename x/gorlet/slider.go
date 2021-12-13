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
func Slider(min, max float32) BuildFunc {
	return func(b *Builder) {
		var val float32
		var handlerSize float32 = 4

		var track *Entity
		var handler *Entity

		root := b.Root()

		root.SetDragEvents(true)

		b.Set("color", m32.Color(.4, .2))
		b.SetAddMode(ElementAdd)
		b.BeginPanel()
		{
			track = b.BeginContainer()
			{
				track.SetAnchor(0, 0, 1, 1)
				track.SetPivot(0, .5)
				track.SetRect(0)
				track.SetRect(handlerSize/2, 0, handlerSize/2, 0)

				// this will be stuck forever :/
				b.Set("textColor", b.Prop("textColor"))
				b.Set("fontScale", b.Prop("fontScale", 2))
				// b.Set("color", b.Prop("handlerColor"))
				handler = b.TextButton("0", nil)
				handler.SetPivot(.5)
				handler.SetRect(0, 0, handlerSize, 0)
				handler.SetAnchor(val, 0, val, 1)
			}
			b.End()
		}
		b.End()

		b.Observe("handlerColor", func(c m32.Vec4) {
			handler.Set("color", c)
		})
		b.Observe("handler", func(e *Entity) {
			// Need to remove Element first :/
			track.RemoveElement(handler)
			handler = e
			track.AddElement(handler)

			handler.SetPivot(.5)
			handler.SetRect(0, 0, handlerSize, 0)
			handler.SetAnchor(val, 0, val, 1)
		})
		b.Observe("value", func(v float32) {
			if val == v {
				return
			}
			val = v
			root.Trigger(EventValueChanged{val})
			handler.SetAnchor(val, 0, val, 1)
			handler.Set("text", fmt.Sprintf("%.2f", min+(val*(max-min))))
		})
		b.Observe("handlerSize", func(f float32) {
			handlerSize = f
			handler.SetRect(0, 0, handlerSize/2, 0)
			track.SetRect(handlerSize/2, 0, handlerSize/2, 0)
		})

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
func (b *Builder) Slider(min, max float32) *Entity {
	return b.Add(Slider(min, max))
}
