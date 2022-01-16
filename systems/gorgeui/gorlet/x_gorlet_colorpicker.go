package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/m32/ray"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

const colorTexSz = 256

var colorTex = func() *gorge.TextureData {
	tw, th := colorTexSz, colorTexSz
	halfH := th / 2
	pixData := make([]byte, tw*th*3)
	for y := 0; y < th; y++ {
		for x := 0; x < tw; x++ {
			i := (y*tw + x) * 3
			ch := float32(x) / float32(tw)
			cv, cs := float32(1), float32(1)
			if y < halfH {
				cv = float32(y) / float32(halfH)
			} else {
				cs = 1 - float32(y-halfH)/float32(halfH)
			}
			c := m32.HSV2RGB(ch, cs, cv)
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

func ColorPicker(fn func(m32.Vec4)) Func {
	if fn == nil {
		fn = func(m32.Vec4) {}
	}
	return func(b *Builder) {
		var val m32.Vec4

		root := b.Root()
		b.BeginPanel(LayoutFlexHorizontal(1, 7))

		b.UseDragEvents(true)
		b.UseRect(0, 0, .3, 0)
		outColor := b.Quad()

		b.UseProps(Props{
			"color":   m32.Color(1),
			"texture": gorge.NewTexture(colorTex),
		})
		b.UseDragEvents(true)
		picker := b.Quad()

		b.EndPanel()

		b.Observe("spacing", ObsFunc(func(v float32) {
			outColor.SetRect(0, 0, v, 0)
		}))
		b.Observe("color", ObsFunc(func(v m32.Vec4) {
			val = v
			outColor.Set("color", v)
			fn(v)
		}))
		pickColor := func(pd *gorgeui.PointerData) {
			rect := picker.Rect()
			m := picker.Mat4()
			v0 := m.MulV4(m32.Vec4{rect[0], rect[1], 0, 1}).Vec3()
			v1 := m.MulV4(m32.Vec4{rect[2], rect[1], 0, 1}).Vec3() // right
			v2 := m.MulV4(m32.Vec4{rect[0], rect[3], 0, 1}).Vec3() // up)

			ui := gorgeui.RootUI(picker)
			if ui == nil {
				return
			}
			r := ray.FromScreen(ui.ScreenSize(), ui.Camera, pd.Position)
			res := ray.IntersectRect(r, v0, v1, v2)

			n := res.UV.Clamp(m32.V2(), m32.V2(1))

			x := int(n[0] * float32(colorTex.Width-1))
			y := int((1 - n[1]) * float32(colorTex.Height-1))

			i := (y*colorTex.Width + x) * 3
			val[0] = float32(colorTex.PixelData[i]) / 255
			val[1] = float32(colorTex.PixelData[i+1]) / 255
			val[2] = float32(colorTex.PixelData[i+2]) / 255

			root.Set("color", val)
		}
		gorge.HandleFunc(picker, func(e gorgeui.EventPointerDown) { pickColor(e.PointerData) })
		gorge.HandleFunc(picker, func(e gorgeui.EventDragBegin) { pickColor(e.PointerData) })
		gorge.HandleFunc(picker, func(e gorgeui.EventDrag) { pickColor(e.PointerData) })
		gorge.HandleFunc(outColor, func(e gorgeui.EventDrag) {
			val[3] = m32.Clamp(val[3]-e.Delta[1]*.01, 0, 1)
			root.Set("color", val)
		})
		root.Set("color", m32.Color(1))
	}
}

// ColorPicker returns adds a color picker to current builder.
// props:
// - color: sets the current color value
// - spacing: spacing between color picker and output color
func (b *Builder) ColorPicker(fn func(m32.Vec4)) *Entity {
	return b.Add(ColorPicker(fn))
}