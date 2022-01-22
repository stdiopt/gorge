package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/text"
)

// TextAlign returns the text alignment
// 0 params it will align to start
// 1 param will align both to same value
// 2 params will align horizontally to first param and vertically to second
// param
func TextAlign(a ...text.Align) [2]text.Align {
	switch {
	case len(a) == 0:
		return [2]text.Align{text.AlignStart, text.AlignStart}
	case len(a) == 1:
		return [2]text.Align{a[0], a[0]}
	default:
		return [2]text.Align{a[0], a[1]}

	}
}

// Label functional.
func Label(t string) Func {
	return func(b *Builder) {
		var autoSize bool

		mesh := text.NewMesh(gorgeui.DefaultFont)
		ent := newEntity(mesh)
		// Use this instead or some common way to GET UI
		ent.Material.SetTexture("albedoMap", gorgeui.DefaultFont.Texture)
		ent.SetScale(1, -1, 1) // UI is inverted
		ent.Color = gm.Color(1)

		Alignment := [2]text.Align{text.AlignCenter, text.AlignCenter}

		// Defaults
		mesh.Overflow = text.OverflowWordWrap
		mesh.Size = 2
		mesh.Alignment = Alignment[0]

		root := b.Root()
		root.AddElement(ent)
		b.BeginContainer()
		b.ClientArea()
		b.EndContainer()

		event.Handle(root, func(gorgeui.EventUpdate) {
			r := root.Rect()
			// AutoSize is experimental and probably buggy.
			// it will only resize each side if it is not anchored
			if autoSize {
				pr := r
				if p, ok := root.Parent().(gorgeui.Entity); ok {
					pr = p.RectTransform().Rect()
				}
				if root.Anchor[0] == root.Anchor[2] {
					r[0], r[2] = pr[0], pr[2]
					root.Dim[0] = mesh.Max[0]
				}
				if root.Anchor[1] == root.Anchor[3] {
					root.Dim[1] = mesh.Max[1] - mesh.Min[1]
				}
			}

			bounds := gm.Vec2{r[2] - r[0], r[3] - r[1]}
			if mesh.Boundary != bounds {
				mesh.SetBoundary(bounds[0], bounds[1])
			}

			// This is executed regardless the text change
			textHeight := float32(mesh.Lines) * mesh.Size
			switch Alignment[1] {
			case text.AlignStart:
				ent.Position[1] = r[1] + mesh.Size*0.25
			case text.AlignCenter:
				ent.Position[1] = r[3] - (bounds[1]*.5 + textHeight*.5) + mesh.Size*0.25 // top, center
			case text.AlignEnd:
				ent.Position[1] = r[3] - textHeight + mesh.Size*0.25
			}
		})

		b.Observe("autoSize", func(v bool) { autoSize = v })
		b.Observe("text", func(s string) { mesh.SetText(s) })
		b.Observe("textColor", func(c gm.Vec4) { ent.SetColorv(c) })
		b.Observe("fontScale", func(v float32) { mesh.SetSize(v) })
		b.Observe("textAlign", func(a [2]text.Align) {
			Alignment = a
			mesh.SetAlignment(Alignment[0])
		})
		b.Observe("overflow", func(o text.Overflow) { mesh.SetOverflow(o) })
		b.Observe("textOverflow", func(o text.Overflow) { mesh.SetOverflow(o) })
		b.Observe("material", func(m gorge.Materialer) { ent.SetMaterial(m) })
		b.Observe("order", func(o int) { ent.SetOrder(o) })

		b.Observe("stencil", func(s *gorge.Stencil) {
			// log.Println("Receiving stencil on label", s)
			ent.Stencil = s
		})
		b.Observe("_maskDepth", func(d int) {
			s := calcMaskOn(d)
			s.WriteMask = 0x00
			s.Fail = gorge.StencilOpKeep
			s.ZFail = gorge.StencilOpKeep
			s.ZPass = gorge.StencilOpKeep
			ent.Renderable().Stencil = s
			for _, c := range root.Children() {
				c.Set("_maskDepth", d)
			}
		})

		root.Set("text", t)
	}
}

// Label creates a label on builder.
func (b *Builder) Label(t string) *Entity {
	return b.Add(Label(t))
}
