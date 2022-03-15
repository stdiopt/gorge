package gorlet

import (
	"fmt"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/text"
)

type WLabel struct {
	Widget[*WLabel]

	alignment [2]text.Align
	autoSize  bool

	mesh *text.Mesh
	ent  *gEntity
}

func Label(t string) *WLabel {
	return Build(&WLabel{
		alignment: [2]text.Align{text.AlignCenter, text.AlignCenter},
	}).SetText(t)
}

func (w *WLabel) Build(b *B) {
	// Check the other stuff
	w.SetAnchor(0, 0, 1, 0)
	w.SetSize(0, 3)

	w.mesh = text.NewMesh(gorgeui.DefaultFont)
	w.mesh.Overflow = text.OverflowWordWrap
	w.mesh.Size = 2
	w.mesh.Alignment = w.alignment[0]

	w.ent = newEntity("labelText", w.mesh)
	w.ent.Material.SetTexture("albedoMap", gorgeui.DefaultFont.Texture)
	w.ent.SetScale(1, -1, 1)
	w.ent.SetColor(1, 1, 1, 1)

	w.Add(w.ent)

	event.Handle(w, func(gorgeui.EventUpdate) {
		sz := w.ContentSize()
		bounds := sz
		if w.autoSize {
			psz := bounds
			if p, ok := w.Parent().(gorgeui.Entity); ok {
				psz = p.RectTransform().ContentSize()
			}
			if w.Anchor[0] == w.Anchor[2] {
				bounds[0] = psz[0]
				w.Size[0] = w.mesh.Max[0]
			}
			if w.Anchor[1] == w.Anchor[3] {
				bounds[1] = psz[1]
				w.Size[1] = float32(w.mesh.Lines) * w.mesh.Size // w.mesh.Max[1] - w.mesh.Min[1]
			}
		}
		if w.mesh.Boundary != bounds {
			w.mesh.SetBoundary(bounds[0], bounds[1])
		}
		textHeight := float32(w.mesh.Lines) * w.mesh.Size
		switch w.alignment[1] {
		case text.AlignStart:
			w.ent.Position[1] = w.mesh.Size * 0.25
		case text.AlignCenter:
			w.ent.Position[1] = sz[1] - (sz[1]*.5 + textHeight*.5) + w.mesh.Size*0.25 // top, center
		case text.AlignEnd:
			w.ent.Position[1] = sz[1] - textHeight + w.mesh.Size*0.25
		}
	})
}

func (w *WLabel) SetFont(font *text.Font) *WLabel {
	w.ent.Material.Set("albedoMap", font.Texture)
	w.mesh.SetFont(font)
	return w
}

func (w *WLabel) SetText(text string) *WLabel {
	w.mesh.SetText(text)
	return w
}

func (w *WLabel) SetTextf(f string, args ...any) *WLabel {
	w.mesh.SetText(fmt.Sprintf(f, args...))
	return w
}

func (w *WLabel) SetTextAlign(a ...text.Align) *WLabel {
	switch len(a) {
	case 0:
		w.alignment = [2]text.Align{text.AlignCenter, text.AlignCenter}
	case 1:
		w.alignment = [2]text.Align{a[0], a[0]}
	default:
		w.alignment = [2]text.Align{a[0], a[1]}
	}
	w.mesh.SetAlignment(w.alignment[0])
	return w
}

func (w *WLabel) SetColor(vs ...float32) *WLabel {
	w.ent.SetColorv(gm.Color(vs...))
	return w
}

func (w *WLabel) SetFontScale(v float32) *WLabel {
	w.mesh.SetSize(v)
	return w
}

func (w *WLabel) SetOverflow(o text.Overflow) *WLabel {
	w.mesh.SetOverflow(o)
	return w
}

func (w *WLabel) SetAutoSize(b bool) *WLabel {
	w.autoSize = b
	return w
}

func (w *WLabel) setMaskDepth(n int) {
	var s *gorge.Stencil
	if n > -1 {
		s = calcMaskOn(n)
		s.WriteMask = 0
		s.Fail = gorge.StencilOpKeep
		s.ZFail = gorge.StencilOpKeep
		s.ZPass = gorge.StencilOpKeep
	}
	w.ent.Stencil = s
}

func (b *B) Label(t string) *WLabel {
	l := Label(t)
	b.Add(l)
	return l
}
