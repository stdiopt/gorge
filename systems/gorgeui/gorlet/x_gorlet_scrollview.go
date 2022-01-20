package gorlet

import (
	"log"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

func Scroll() Func {
	return func(b *Builder) {
		background := b.Prop("background", m32.Color(0, .3))

		var (
			container  *Entity
			scrollable *Entity
			scroll     [2]*Entity
			vSize      float32 = 1
			hSize      float32 = 1
			vScroll    float32 = 0
			hScroll    float32 = 0
		)

		b.Use("color", background)
		b.BeginPanel()
		{
			b.UseRect(0, 0, vSize, hSize)
			scrollable = b.BeginMask()
			{
				b.UseAnchor(0, 0, 1, 0)
				container = b.BeginContainer()
				b.ClientArea()
				b.EndContainer()
			}
			b.EndMask()
			// Anything added to root goes into client area

			// Horizontal
			b.UseAnchor(0, 1, 1, 1)
			b.UseRect(0, -hSize, vSize, hSize)
			b.Use("color", m32.Color(0, 0, 0, .2))
			scroll[0] = b.HScrollBar(func(v float32) {
				sz := scrollable.CalcSize()
				b := container.CalcBounds()
				container.Position[0] = v * (sz[0] - b[2])
				hScroll = v
			})

			// Vertical scrollBar
			b.UseAnchor(1, 0, 1, 1)
			b.UseRect(-vSize, 0, hSize, vSize)
			b.Use("color", m32.Color(0, 0, 0, .2))
			scroll[1] = b.VScrollBar(func(v float32) {
				sz := scrollable.CalcSize()
				b := container.CalcBounds()
				container.Position[1] = v * (sz[1] - b[3])
				vScroll = v
			})
		}
		b.EndPanel()

		updateScrolls := func() {
			scrollable.SetRect(0, 0, vSize, hSize)
			scroll[0].SetRect(0, -hSize, vSize, hSize)
			scroll[1].SetRect(-vSize, 0, vSize, hSize)
		}
		Observe(b, "scrollSize", func(v float32) {
			hSize, vSize = v, v
			updateScrolls()
		})
		Observe(b, "hscrollSize", func(v float32) {
			hSize = v
			updateScrolls()
			log.Println("sizes:", hSize, vSize)
		})
		Observe(b, "vscrollSize", func(v float32) {
			vSize = v
			updateScrolls()
		})

		root := b.Root()
		event.Handle(root, func(e gorgeui.EventPointerWheel) {
			scroll[0].Set("value", hScroll+e.Wheel[0])
			scroll[1].Set("value", vScroll+e.Wheel[1]*0.01)
		})
		event.Handle(root, func(gorgeui.EventUpdate) {
			sz := scrollable.CalcSize()
			b := container.CalcBounds()

			if sz[0] >= b[2] {
				scroll[0].Set("value", 0)
				scroll[0].Set("handlerSize", sz[0])
			} else if sz[0] < b[2] {
				hs := sz[0] * (sz[0] / b[2])
				curScroll := container.Position[0] / (sz[0] - b[2])
				scroll[0].Set("value", curScroll)
				scroll[0].Set("handlerSize", hs)
			}

			// Vertical
			if sz[1] > b[3] {
				scroll[1].Set("value", 0)
				scroll[1].Set("handlerSize", sz[1])
			} else if sz[1] < b[3] {
				hs := sz[1] * (sz[1] / b[3])
				curScroll := container.Position[1] / (sz[1] - b[3])
				scroll[1].Set("value", curScroll)
				scroll[1].Set("handlerSize", hs)
			}
		})
	}
}

func HScrollBar(fn func(float32)) Func {
	return ScrollBar(DirectionHorizontal, fn)
}

func VScrollBar(fn func(float32)) Func {
	return ScrollBar(DirectionVertical, fn)
}

// ScrollBar creates a scroll bar
// TODO: implement direction
func ScrollBar(dir Direction, fn func(float32)) Func {
	return func(b *Builder) {
		var (
			backgroundColor = b.Prop("backgroundColor", m32.Color(.4, .2))
			handlerColor    = b.Prop("handlerColor")
		)
		var (
			dragging    *m32.Vec2
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
			if dir == DirectionHorizontal {
				b.UseRect(handlerSize/2, 0, handlerSize/2, 0)
			} else {
				b.UseRect(0, handlerSize/2, 0, handlerSize/2)
			}
			b.UsePivot(.5, 0)
			track = b.BeginContainer()
			{
				// this will be stuck forever :/
				b.Use("color", handlerColor)
				b.UsePivot(.5)
				if dir == DirectionHorizontal {
					b.UseRect(0, 0, handlerSize, 0)
					b.UseAnchor(val, 0, val, 1)
				} else {
					b.UseRect(0, 0, 0, handlerSize)
					b.UseAnchor(0, val, 1, val)
				}
				handler = b.Button(nil)
			}
			b.End()
		}
		b.EndPanel()

		Observe(b, "handlerColor", func(c m32.Vec4) {
			handler.Set("color", c)
		})
		Observe(b, "value", func(v float32) {
			v = m32.Clamp(v, 0, 1)
			if val == v {
				return
			}
			val = v
			if dir == DirectionHorizontal {
				handler.SetAnchor(val, 0, val, 1)
			} else {
				handler.SetAnchor(0, val, 1, val)
			}
			if fn != nil {
				fn(val)
			}
		})
		Observe(b, "handlerSize", func(f float32) {
			handlerSize = f
			if dir == DirectionHorizontal {
				handler.SetRect(0, 0, handlerSize, 0)
				track.SetRect(handlerSize/2, 0, handlerSize/2, 0)
			} else {
				handler.SetRect(0, 0, 0, handlerSize)
				track.SetRect(0, handlerSize/2, 0, handlerSize/2)
			}
		})

		event.Handle(root, func(e gorgeui.EventPointerUp) {
			if dragging != nil {
				return
			}
			res := root.IntersectFromScreen(e.Position)
			var v float32
			if dir == DirectionHorizontal {
				v := res.UV[0]
				v -= handlerSize / root.Rect()[2] / 2
			} else {
				v := res.UV[1]
				v -= handlerSize / root.Rect()[3] / 2
			}
			v = m32.Clamp(v, 0, 1)
			root.Set("value", v)
		})
		event.Handle(root, func(e gorgeui.EventDragBegin) {
			uv := track.IntersectFromScreen(e.Position)
			dragging = &uv.UV
		})
		event.Handle(root, func(e gorgeui.EventDrag) {
			if dragging == nil {
				return
			}
			uv := track.IntersectFromScreen(e.Position).UV
			delta := uv.Sub(*dragging)
			dragging = &uv
			var v float32
			if dir == DirectionHorizontal {
				v = val + delta[0]
			} else {
				v = val + delta[1]
			}
			root.Set("value", v)
		})
		event.Handle(root, func(gorgeui.EventDragEnd) {
			dragging = nil
		})

		root.Set("value", float32(0))
	}
}

// BeginScroll creates and pushes a scrollable area.
func (b *Builder) BeginScroll() *Entity {
	return b.Begin(Scroll())
}

// EndScroll alias to .End()
func (b *Builder) EndScroll() {
	b.End()
}

// ScrollBar creates a scroll bar
func (b *Builder) VScrollBar(fn func(float32)) *Entity {
	return b.Add(ScrollBar(DirectionVertical, fn))
}

// ScrollBar creates a scroll bar
func (b *Builder) HScrollBar(fn func(float32)) *Entity {
	return b.Add(ScrollBar(DirectionHorizontal, fn))
}
