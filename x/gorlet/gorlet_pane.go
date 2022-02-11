package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type WPane struct {
	Widget[WPane]

	borderRadius float32
	roundedQuad  *roundedQuadMesh
	layout       Layouter

	ent *gEntity
}

func Pane(c ...gorge.Entity) *WPane {
	return Build(&WPane{borderRadius: .4}).Add(c...)
}

func (w *WPane) Build(b *B) {
	w.roundedQuad = newRoundedQuadMesh(gm.Vec2{}, w.borderRadius)
	w.ent = newEntity("panel", w.roundedQuad)
	w.ent.SetColor(0, 0, 0, .5)
	w.Add(w.ent)

	prevRect := gm.Vec4{}
	event.Handle(w, func(gorgeui.EventUpdate) {
		if w.layout != nil {
			w.layout.Layout(w)
		}
		w.roundedQuad.SetRadius(w.borderRadius)
		// Rebuild mesh here?! if size changes
		sz := w.ContentSize()
		newRect := gm.Vec4{
			-w.Border[0],
			-w.Border[1],
			sz[0] + w.Border[0] + w.Border[2],
			sz[1] + w.Border[1] + w.Border[3],
		}
		if newRect.Equal(prevRect) {
			return
		}
		prevRect = newRect

		w.roundedQuad.SetSize(newRect.ZW())
		w.roundedQuad.update()
		w.ent.Position = newRect.XY().Vec3(0)
		if w.Border != (gm.Vec4{}) {
			w.ent.Renderable().Material.Define("HAS_BORDER")
			w.ent.Renderable().Material.Set("border", w.Border)
		}
		w.ent.Renderable().Material.Set("rect", gm.Vec4{0, 0, sz[0], sz[1]})
	})
}

func (w *WPane) SetBorderRadius(r float32) *WPane {
	w.borderRadius = r
	return w
}

func (w *WPane) SetColor(vs ...float32) *WPane {
	c := gm.Color(vs...)
	w.ent.Colorable().SetColor(c[0], c[1], c[2], c[3])
	return w
}

func (w *WPane) SetBorderColor(vs ...float32) *WPane {
	c := gm.Color(vs...)
	w.ent.Material.Set("borderColor", c)
	return w
}

func (w *WPane) SetMaterial(mat gorge.Materialer) *WPane {
	w.ent.SetMaterial(mat)
	return w
}

func (w *WPane) SetTexture(tex gorge.Texturer) *WPane {
	w.ent.Material.SetTexture("albedoMap", tex)
	return w
}

func (w *WPane) SetLayout(l ...Layouter) *WPane {
	w.layout = LayoutMulti(l...)
	return w
}

func (w *WPane) setMaskDepth(n int) {
	var s *gorge.Stencil
	if n > -1 {
		s = calcMaskOn(n)
		s.WriteMask = 0
		s.Fail = gorge.StencilOpKeep
		s.ZFail = gorge.StencilOpKeep
		s.ZPass = gorge.StencilOpKeep
	}
	w.ent.Stencil = s
	w.Widget.setMaskDepth(n)
}

func (b *B) Pane() *WPane {
	p := Pane()
	b.Add(p)
	return p
}

func (b *B) BeginPane() *WPane {
	p := Pane()
	b.Begin(p)
	return p
}

func (b *B) EndPane() {
	b.End()
}
