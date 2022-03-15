package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type WMask struct {
	Widget[*WMask]

	borderRadius float32
	depthMask    int

	maskOn    *gEntity
	container *WContainer
	maskOff   *gEntity
}

func Mask(c ...gorge.Entity) *WMask {
	return Build(&WMask{
		borderRadius: .4,
	}).Add(c...)
}

func (w *WMask) Build(b *B) {
	m := newRoundedQuadMesh(gm.Vec2{}, w.borderRadius)
	// Maybe use the rounded one instead?
	w.maskOn = newEntity("maskOn", m)
	w.maskOn.SetColor(1, 0, 0, 1)
	w.maskOn.SetColorMask(false, false, false, false)
	w.maskOff = newEntity("maskOff", m)
	w.maskOff.SetColorMask(false, false, false, false)

	w.masked = true

	w.container = Container()
	w.Add(
		w.maskOn,
		w.container,
		w.maskOff,
	)
	w.SetClientArea(w.container)

	prevRect := gm.Vec4{}
	event.Handle(w, func(gorgeui.EventUpdate) {
		w.Widget.setMaskDepth(w.depthMask + 1)

		sz := w.ContentSize()

		newRect := gm.Vec4{
			-w.Border[0],
			-w.Border[1],
			sz[0] + w.Border[2] + w.Border[0],
			sz[1] + w.Border[3] + w.Border[1],
		}
		if newRect.Equal(prevRect) {
			return
		}
		prevRect = newRect
		borderMin := newRect.XY().Vec3(0)

		m.SetSize(newRect.ZW())
		m.SetRadius(w.borderRadius)
		m.update()

		w.maskOn.Position = borderMin
		w.maskOff.Position = borderMin
		w.maskOn.Stencil = calcMaskOn(w.depthMask)
		w.maskOff.Stencil = calcMaskOff(w.depthMask)
	})
}

func (w *WMask) setMaskDepth(n int) {
	if n < 0 {
		n = 0
	}
	w.depthMask = n
}

func (w *WMask) SetBorderRadius(r float32) *WMask {
	w.borderRadius = r
	return w
}

func (b *B) BeginMask() *WMask {
	w := Mask()
	b.Begin(w)
	return w
}

func (b *B) EndMask() {
	b.End()
}
