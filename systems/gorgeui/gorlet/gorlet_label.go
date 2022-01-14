package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/static"
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

		mat := gorge.NewShaderMaterial(static.Shaders.Unlit)
		mat.SetQueue(100)
		// Use this instead or some common way to GET UI
		// mat.SetDepth(gorge.DepthRead)
		mat.SetDepth(gorge.DepthNone)
		mat.SetTexture("albedoMap", gorgeui.DefaultFont.Texture)

		Alignment := [2]text.Align{text.AlignCenter, text.AlignCenter}

		ent := text.New(gorgeui.DefaultFont)
		ent.SetMaterial(mat)
		ent.SetScale(1, -1, 1) // UI is inverted

		// Defaults
		ent.Overflow = text.OverflowWordWrap
		ent.Size = 2
		ent.Alignment = Alignment[0]
		ent.Color = m32.Color(1)

		root := b.Root()
		root.AddElement(ent)
		b.BeginContainer()
		b.ClientArea()
		b.EndContainer()

		gorge.HandleFunc(root, func(gorgeui.EventUpdate) {
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
					root.Dim[0] = ent.Max[0]
				}
				if root.Anchor[1] == root.Anchor[3] {
					root.Dim[1] = ent.Max[1] - ent.Min[1]
				}
			}

			bounds := m32.Vec2{r[2] - r[0], r[3] - r[1]}
			if ent.Boundary != bounds {
				ent.SetBoundary(bounds[0], bounds[1])
			}

			// This is executed regardless the text change
			textHeight := float32(ent.Lines) * ent.Size
			switch Alignment[1] {
			case text.AlignStart:
				ent.Position[1] = r[1] + ent.Size*0.25
			case text.AlignCenter:
				ent.Position[1] = r[3] - (bounds[1]*.5 + textHeight*.5) + ent.Size*0.25 // top, center
			case text.AlignEnd:
				ent.Position[1] = r[3] - textHeight + ent.Size*0.25
			}
		})

		b.Observe("autoSize", ObsFunc(func(v bool) { autoSize = v }))
		b.Observe("text", ObsFunc(func(s string) { ent.SetText(s) }))
		b.Observe("textColor", ObsFunc(func(c m32.Vec4) { ent.SetColorv(c) }))
		b.Observe("fontScale", ObsFunc(func(v float32) { ent.SetSize(v) }))
		b.Observe("textAlign", ObsFunc(func(a [2]text.Align) {
			Alignment = a
			ent.SetAlignment(Alignment[0])
		}))
		b.Observe("overflow", ObsFunc(func(o text.Overflow) { ent.SetOverflow(o) }))
		b.Observe("textOverflow", ObsFunc(func(o text.Overflow) { ent.SetOverflow(o) }))
		b.Observe("material", ObsFunc(func(m gorge.Materialer) { ent.SetMaterial(m) }))
		b.Observe("order", ObsFunc(func(o int) { ent.SetOrder(o) }))

		root.Set("text", t)
	}
}

// Label creates a label on builder.
func (b *Builder) Label(t string) *Entity {
	return b.Add(Label(t))
}
