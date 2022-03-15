package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

const colorTexSz = 256

// global to reuse same texture
var colorTex = func() *gorge.TextureData {
	tw, th := colorTexSz, colorTexSz
	halfH := th / 2
	pixData := make([]byte, tw*th*3)

	gs := 8
	// Grays
	for y := 0; y < th; y++ {
		for x := 0; x < gs; x++ {
			i := (y*tw + x) * 3
			v := float32(y) / float32(th)
			pixData[i] = byte(v * 255)
			pixData[i+1] = byte(v * 255)
			pixData[i+2] = byte(v * 255)

		}
	}
	// Color
	for y := 0; y < th; y++ {
		for x := gs; x < tw; x++ {
			i := (y*tw + x) * 3
			ch := float32(x-gs) / float32(tw-gs)
			cv, cs := float32(1), float32(1)
			if y < halfH {
				cv = float32(y) / float32(halfH)
			} else {
				cs = 1 - float32(y-halfH)/float32(halfH)
			}
			c := gm.HSV2RGB(ch, cs, cv)
			pixData[i] = byte(c[0] * 255)
			pixData[i+1] = byte(c[1] * 255)
			pixData[i+2] = byte(c[2] * 255)
		}
	}
	return &gorge.TextureData{
		Width:     tw,
		Height:    th,
		Format:    gorge.TextureFormatRGB,
		PixelData: pixData,
	}
}()

type WColorPicker struct {
	Widget[*WColorPicker]
	spacing  float32
	value    gm.Vec4
	outColor *WPane
	picker   *WPane

	alpha *WSlider

	changefn func(gm.Vec4)
}

func ColorPicker() *WColorPicker {
	return Build(&WColorPicker{
		spacing: 0,
		value:   gm.Color(1),
	})
}

func (w *WColorPicker) Build(b *B) {
	w.SetAnchor(0, 0, 1, 0)
	w.SetSize(0, 5)

	tex := gorge.NewTexture(colorTex)
	tex.SetWrapUVW(gorge.TextureWrapClamp)

	b.BeginPane()
	{
		w.outColor = b.Pane().
			SetColor(1).
			SetMargin(0, 0, .1, 0).
			SetAnchor(0, 0, .13, 1).
			SetDragEvents(true)

		b.BeginContainer().
			SetMargin(.1, 0, 0, 0).
			SetAnchor(.13, 0, 1, 1)
		{
			w.picker = b.Pane().
				SetColor(w.value[:]...).
				SetDragEvents(true).
				SetTexture(tex).
				SetRect(0).
				SetAnchor(0, 0, 1, .8).
				SetMargin(0, 0, 0, .1)

			w.alpha = b.Slider(0, 1).
				SetFormat("alpha: %.2f").
				SetFontScale(1.5).
				SetValue(1).
				SetRect(0).
				SetAnchor(0, .8, 1, 1).
				OnChange(func(v float32) {
					val := w.value
					val[3] = gm.Clamp(v, 0, 1)
					w.SetValue(val)
				})
		}
		b.EndContainer()
	}
	b.EndPane()

	event.Handle(w.picker, func(e gorgeui.EventPointerDown) { w.pickColor(e.PointerData) })
	event.Handle(w.picker, func(e gorgeui.EventDragBegin) { w.pickColor(e.PointerData) })
	event.Handle(w.picker, func(e gorgeui.EventDrag) { w.pickColor(e.PointerData) })
	event.Handle(w.outColor, func(e gorgeui.EventDrag) {
		val := w.value
		val[3] = gm.Clamp(w.value[3]-e.Delta[1]*.01, 0, 1)
		w.SetValue(val)
	})
}

func (w *WColorPicker) pickColor(pd *gorgeui.PointerData) {
	res := w.picker.IntersectFromScreen(pd.Position)

	n := res.UV.Clamp(gm.V2(), gm.V2(1))

	x := int(n[0] * float32(colorTex.Width-1))
	y := int((1 - n[1]) * float32(colorTex.Height-1))

	i := (y*colorTex.Width + x) * 3

	val := w.value
	val[0] = float32(colorTex.PixelData[i]) / 255
	val[1] = float32(colorTex.PixelData[i+1]) / 255
	val[2] = float32(colorTex.PixelData[i+2]) / 255

	w.SetValue(val)
}

func (w *WColorPicker) SetValue(c gm.Vec4) *WColorPicker {
	if w.value == c {
		return w
	}
	w.value = c
	w.outColor.SetColor(c[:]...)
	w.alpha.SetValue(c[3])
	if w.changefn != nil {
		w.changefn(c)
	}
	return w
}

func (w *WColorPicker) OnChange(fn func(gm.Vec4)) *WColorPicker {
	w.changefn = fn
	return w
}

func (b *B) ColorPicker() *WColorPicker {
	w := ColorPicker()
	b.Add(w)
	return w
}
