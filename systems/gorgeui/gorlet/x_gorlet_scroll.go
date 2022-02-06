package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type scrollent struct {
	disabled  bool
	notneeded bool
	entity    *Entity
	size      float32
	value     float32
}

func Scroll() Func {
	return func(b *B) {
		var (
			contentMargin = b.Prop("contentMargin", Margin(0))
			container     *Entity
			scrollable    *Entity
			scrolls       = [2]scrollent{
				{size: 1},
				{size: 1},
			}
			scrollVel gm.Vec2
		)

		scrollFunc := func(n int) func(v float32) {
			return func(v float32) {
				if scrolls[n].disabled {
					return
				}
				v = gm.Clamp(v, 0, 1)
				if v == scrolls[n].value {
					return
				}
				sz := scrollable.ContentSize()
				b := container.CalcMax()
				if sz[n]/b[n] > 1 {
					return
				}
				container.Position[n] = v * (sz[n] - b[n])
				scrolls[n].value = v
			}
		}

		root := b.Root()
		{
			b.UseRect(0, 0, scrolls[1].size, scrolls[0].size)
			b.Use("margin", contentMargin)
			scrollable = b.BeginMask()
			{
				container = b.BeginContainer()
				b.ClientArea()
				b.EndContainer()
			}
			b.EndMask()
			// Anything added to root goes into client area

			// Horizontal
			b.UseAnchor(0, 1, 1, 1)
			b.Use("color", gm.Color(0, 0, 0, .2))
			scrolls[0].entity = b.Add(
				scrollBar(Horizontal, scrollFunc(0)),
			)
			// Vertical scrollBar
			b.UseAnchor(1, 0, 1, 1)
			b.Use("color", gm.Color(0, 0, 0, .2))
			scrolls[1].entity = b.Add(
				scrollBar(Vertical, scrollFunc(1)),
			)
		}

		Observe(b, "hscrollSize", func(v float32) {
			scrolls[0].disabled = (v == 0)
			scrolls[0].size = v
		})
		Observe(b, "vscrollSize", func(v float32) {
			scrolls[1].disabled = (v == 0)
			scrolls[1].size = v
		})
		Observe(b, "scrollSize", func(v float32) {
			scrolls[0].size = v
			scrolls[1].size = v
		})

		event.Handle(root, func(e gorgeui.EventPointerWheel) {
			if !scrolls[0].disabled {
				scrollVel[0] -= e.Wheel[0] * 0.005
				if e.Wheel[0] > 0 && scrolls[0].value != 1 {
					e.StopPropagation()
				} else if e.Wheel[0] < 0 && scrolls[0].value != 0 {
					e.StopPropagation()
				}
			}
			if !scrolls[1].disabled {
				scrollVel[1] += e.Wheel[1] * 0.005
				if e.Wheel[1] > 0 && scrolls[1].value != 1 {
					e.StopPropagation()
				} else if e.Wheel[1] < 0 && scrolls[1].value != 0 {
					e.StopPropagation()
				}
			}
		})
		event.Handle(root, func(e gorgeui.EventUpdate) {
			sz := scrollable.ContentSize()
			b := container.CalcMax()

			ssz := gm.Vec2{}

			scrollVel = scrollVel.Lerp(gm.Vec2{}, e.DeltaTime()*20)
			if !scrolls[0].disabled {
				if sz[0]+gm.Epsilon >= b[0] {
					scrolls[0].entity.Set("value", 0)
					scrolls[0].entity.Set("handlerSize", sz[0])
				} else {
					ssz[0] = scrolls[0].size

					hs := sz[0] * (sz[0] / b[0])
					delta := scrollVel[0] * .8
					curScroll := container.Position[0] / (sz[0] - b[0])
					scrolls[0].entity.Set("value", curScroll+delta)
					scrolls[0].entity.Set("handlerSize", hs)
				}
			}

			if !scrolls[1].disabled {
				// Vertical
				if sz[1]+gm.Epsilon >= b[1] {
					scrolls[1].entity.Set("value", 0)
					scrolls[1].entity.Set("handlerSize", sz[1])
				} else {
					ssz[1] = scrolls[1].size
					hs := sz[1] * (sz[1] / b[1])
					curScroll := container.Position[1] / (sz[1] - b[1])
					delta := scrollVel[1] * .8
					scrolls[1].entity.Set("value", curScroll+delta)
					scrolls[1].entity.Set("handlerSize", hs)
				}
			}
			scrollable.SetRect(0, 0, ssz[1], ssz[0])
			scrolls[0].entity.SetRect(0, -ssz[0], ssz[1], ssz[0])
			scrolls[1].entity.SetRect(-ssz[1], 0, ssz[1], ssz[0])
		})
		root.Set("scrollSize", 1)
	}
}

func scrollBar(dir Direction, fn func(float32)) Func {
	const sp = .3
	const spx = .2
	return func(b *B) {
		var (
			backgroundColor = b.Prop("backgroundColor", gm.Color(.4, .3))
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
					b.UseMargin(spx, sp)
					b.UseRect(0, 0, handlerSize, 0)
					b.UseAnchor(val, 0, val, 1)
				} else {
					b.UseMargin(sp, spx)
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

// BeginScroll starts a scrollable rect.
func (b *B) BeginScroll() *Entity {
	return b.Begin(Scroll())
}

// EndScroll ends a scrollable rect.
func (b *B) EndScroll() {
	b.End()
}
