package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type scroll struct {
	disabled bool
	entity   *Entity
	size     float32
	value    float32
}

func Scroll() Func {
	return func(b *Builder) {
		background := b.Prop("background", gm.Color(0, .3))
		var (
			container  *Entity
			scrollable *Entity
			scroll     = [2]scroll{
				{size: 1},
				{size: 1},
			}
		)

		b.Use("color", background)
		b.BeginPanel()
		{
			b.UseRect(0, 0, scroll[1].size, scroll[0].size)
			scrollable = b.BeginMask()
			{
				b.UseAnchor(0, 0, 1, 1)
				container = b.BeginContainer()
				b.ClientArea()
				b.EndContainer()
			}
			b.EndMask()
			// Anything added to root goes into client area

			// Horizontal
			b.UseAnchor(0, 1, 1, 1)
			// b.UseRect(0, -scroll[0].size, scroll[1].size, scroll[0].size)
			b.UseRect(0)
			b.Use("color", gm.Color(0, 0, 0, .2))
			scroll[0].entity = b.HScrollBar(func(v float32) {
				if scroll[0].disabled {
					return
				}
				v = gm.Clamp(v, 0, 1)
				if v == scroll[0].value {
					return
				}
				sz := scrollable.CalcSize()
				b := container.CalcBounds()
				if sz[0]/b[2] > 1 {
					return
				}
				container.Position[0] = v * (sz[0] - b[2])
				scroll[0].value = v
			})

			// Vertical scrollBar
			b.UseAnchor(1, 0, 1, 1)
			// b.UseRect(-scroll[1].size, 0, scroll[0].size, scroll[1].size)
			b.UseRect(0)
			b.Use("color", gm.Color(0, 0, 0, .2))
			scroll[1].entity = b.VScrollBar(func(v float32) {
				if scroll[1].disabled {
					return
				}
				v = gm.Clamp(v, 0, 1)
				if v == scroll[1].value {
					return
				}
				sz := scrollable.CalcSize()
				b := container.CalcBounds()
				if sz[1]/b[3] > 1 {
					return
				}
				container.Position[1] = v * (sz[1] - b[3])
				scroll[1].value = v
			})
		}
		b.EndPanel()

		updateScrolls := func() {
			sz := gm.Vec2{}
			if !scroll[0].disabled {
				sz[0] = scroll[0].size
			}
			if !scroll[1].disabled {
				sz[1] = scroll[1].size
			}
			scrollable.SetRect(0, 0, sz[1], sz[0])

			scroll[0].entity.SetRect(0, -sz[0], sz[1], sz[0])
			scroll[1].entity.SetRect(-sz[1], 0, sz[1], sz[0])
		}
		Observe(b, "hscroll", func(b bool) {
			scroll[0].disabled = !b
		})
		Observe(b, "vscroll", func(b bool) {
			scroll[1].disabled = !b
		})
		Observe(b, "hscrollSize", func(v float32) {
			scroll[0].size = v
			updateScrolls()
		})
		Observe(b, "vscrollSize", func(v float32) {
			scroll[1].size = v
			updateScrolls()
		})
		Observe(b, "scrollSize", func(v float32) {
			scroll[0].size, scroll[1].size = v, v
			updateScrolls()
		})

		root := b.Root()
		event.Handle(root, func(e gorgeui.EventPointerWheel) {
			if !scroll[0].disabled {
				h := scroll[0].value
				scroll[0].entity.Set("value", gm.Clamp(h+e.Wheel[0]*0.01, 0, 1))
			}
			if !scroll[1].disabled {
				v := scroll[1].value
				scroll[1].entity.Set("value", gm.Clamp(v+e.Wheel[1]*0.01, 0, 1))
			}
		})
		event.Handle(root, func(gorgeui.EventUpdate) {
			sz := scrollable.CalcSize()
			b := container.CalcBounds()

			if !scroll[0].disabled {
				if sz[0] >= b[2] {
					scroll[0].entity.Set("value", 0)
					scroll[0].entity.Set("handlerSize", sz[0])
				} else if sz[0] < b[2] {
					hs := sz[0] * (sz[0] / b[2])
					curScroll := container.Position[0] / (sz[0] - b[2])
					scroll[0].entity.Set("value", curScroll)
					scroll[0].entity.Set("handlerSize", hs)
				}
			}

			if !scroll[1].disabled {
				// Vertical
				if sz[1] > b[3] {
					scroll[1].entity.Set("value", 0)
					scroll[1].entity.Set("handlerSize", sz[1])
				} else if sz[1] < b[3] {
					hs := sz[1] * (sz[1] / b[3])
					curScroll := container.Position[1] / (sz[1] - b[3])
					scroll[1].entity.Set("value", curScroll)
					scroll[1].entity.Set("handlerSize", hs)
				}
			}
		})
		root.Set("scrollSize", 1)
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
			backgroundColor = b.Prop("backgroundColor", gm.Color(.4, .2))
			handlerColor    = b.Prop("handlerColor")
		)
		var (
			dragging    *gm.Vec2
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

		Observe(b, "handlerColor", func(c gm.Vec4) {
			handler.Set("color", c)
		})
		Observe(b, "value", func(v float32) {
			v = gm.Clamp(v, 0, 1)
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
			v = gm.Clamp(v, 0, 1)
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
