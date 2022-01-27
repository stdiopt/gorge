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
		var (
			background = b.Prop("background", gm.Color(0, .3))

			container  *Entity
			scrollable *Entity
			scrolls    = [2]scroll{
				{size: 1},
				{size: 1},
			}
		)

		doscroll := func(n int) func(v float32) {
			return func(v float32) {
				if scrolls[n].disabled {
					return
				}
				v = gm.Clamp(v, 0, 1)
				if v == scrolls[n].value {
					return
				}
				sz := scrollable.CalcSize()
				b := container.CalcBounds()
				if sz[0+n]/b[2+n] > 1 {
					return
				}
				container.Position[n] = v * (sz[0+n] - b[2+n])
				scrolls[n].value = v
			}
		}
		// Update scrollsize
		updateScrolls := func() {
			sz := gm.Vec2{}
			if !scrolls[0].disabled {
				sz[0] = scrolls[0].size
			}
			if !scrolls[1].disabled {
				sz[1] = scrolls[1].size
			}
			scrollable.SetRect(0, 0, sz[1], sz[0])

			scrolls[0].entity.SetRect(0, -sz[0], sz[1], sz[0])
			scrolls[1].entity.SetRect(-sz[1], 0, sz[1], sz[0])
		}
		/*
			hscroll := func(v float32) {
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
			}

			vscroll := func(v float32) {
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
			}*/

		// TODO: Observe root while setting root.
		// we might want to skip this
		b.Use("color", background)
		b.SetRoot(Panel())
		{
			b.UseRect(0, 0, scrolls[1].size, scrolls[0].size)
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
			b.Use("color", gm.Color(0, 0, 0, .2))
			scrolls[0].entity = b.HScrollBar(doscroll(0))

			// Vertical scrollBar
			b.UseAnchor(1, 0, 1, 1)
			b.Use("color", gm.Color(0, 0, 0, .2))
			scrolls[1].entity = b.VScrollBar(doscroll(1))
		}

		Observe(b, "hscrollSize", func(v float32) {
			scrolls[0].disabled = (v == 0)
			scrolls[0].size = v
			updateScrolls()
		})
		Observe(b, "vscrollSize", func(v float32) {
			scrolls[1].disabled = (v == 0)
			scrolls[1].size = v
			updateScrolls()
		})
		Observe(b, "scrollSize", func(v float32) {
			scrolls[0].size = v
			scrolls[1].size = v
			updateScrolls()
		})

		root := b.Root()
		event.Handle(root, func(e gorgeui.EventPointerWheel) {
			if !scrolls[0].disabled {
				h := scrolls[0].value
				scrolls[0].entity.Set("value", gm.Clamp(h+e.Wheel[0]*0.01, 0, 1))
			}
			if !scrolls[1].disabled {
				v := scrolls[1].value
				scrolls[1].entity.Set("value", gm.Clamp(v+e.Wheel[1]*0.01, 0, 1))
			}
		})
		event.Handle(root, func(gorgeui.EventUpdate) {
			sz := scrollable.CalcSize()
			b := container.CalcBounds()

			if !scrolls[0].disabled {
				if sz[0] >= b[2] {
					scrolls[0].entity.Set("value", 0)
					scrolls[0].entity.Set("handlerSize", sz[0])
				} else if sz[0] < b[2] {
					hs := sz[0] * (sz[0] / b[2])
					curScroll := container.Position[0] / (sz[0] - b[2])
					scrolls[0].entity.Set("value", curScroll)
					scrolls[0].entity.Set("handlerSize", hs)
				}
			}

			if !scrolls[1].disabled {
				// Vertical
				if sz[1] >= b[3] {
					scrolls[1].entity.Set("value", 0)
					scrolls[1].entity.Set("handlerSize", sz[1])
				} else if sz[1] < b[3] {
					hs := sz[1] * (sz[1] / b[3])
					curScroll := container.Position[1] / (sz[1] - b[3])
					scrolls[1].entity.Set("value", curScroll)
					scrolls[1].entity.Set("handlerSize", hs)
				}
			}
		})
		root.Set("scrollSize", 1)
	}
}

func HScrollBar(fn func(float32)) Func {
	return ScrollBar(Horizontal, fn)
}

func VScrollBar(fn func(float32)) Func {
	return ScrollBar(Vertical, fn)
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
			if dir == Horizontal {
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
				if dir == Horizontal {
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
			if dir == Horizontal {
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
			if dir == Horizontal {
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
			if dir == Horizontal {
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
			if dir == Horizontal {
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
	return b.Add(ScrollBar(Vertical, fn))
}

// ScrollBar creates a scroll bar
func (b *Builder) HScrollBar(fn func(float32)) *Entity {
	return b.Add(ScrollBar(Horizontal, fn))
}
