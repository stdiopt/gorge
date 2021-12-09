package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/text"
)

// Label functional.
func Label(t string) BuildFunc {
	return func(b *Builder) {
		root := b.Root()

		mat := gorge.NewShaderMaterial(static.Shaders.Unlit)
		mat.SetQueue(100)
		mat.SetDepth(gorge.DepthNone)
		mat.SetTexture("albedoMap", gorgeui.DefaultFont.Texture)

		// Font      *text.Font
		var Size float32 = 2
		var Alignment [2]text.AlignType = [2]text.AlignType{text.AlignCenter, text.AlignCenter}

		ent := text.New(gorgeui.DefaultFont)
		ent.SetMaterial(mat)
		ent.SetScale(1, -1, 1)
		ent.SetOverflow(text.OverflowWordWrap)
		ent.Size = Size
		ent.SetAlignment(Alignment[0])
		ent.SetColor(1, 1, 1, 1)

		gorgeui.AddElementTo(root, ent)

		// Element maybe?
		b.Root().HandleFunc(func(e event.Event) {
			_, ok := e.(gorgeui.EventUpdate) // Change to PreUpdate?
			if !ok {
				return
			}
			r := b.Root().Rect()
			bounds := m32.Vec2{r[2] - r[0], r[3] - r[1]}
			ent.Position[0] = r[0] // left

			if ent.Boundary != bounds {
				ent.Boundary = bounds
				ent.Update()
			}

			// This is executed regardless the text change
			textHeight := float32(ent.Lines) * Size
			switch Alignment[1] {
			case text.AlignStart:
				ent.Position[1] = r[1] + Size*0.25
			case text.AlignCenter:
				ent.Position[1] = r[3] - (bounds[1]/2 + textHeight/2) + Size*0.25 // top, center
			case text.AlignEnd:
				ent.Position[1] = r[3] - textHeight + Size*0.25
			}
		})

		// Lose to strong type?
		b.Observe("text", func(s string) {
			ent.SetText(s)
			ent.Update()
		})
		b.Observe("textColor", func(c m32.Vec4) {
			ent.SetColorv(c)
			ent.Update()
		})
		b.Observe("fontScale", func(v float32) {
			Size = v
			ent.Size = v
			ent.Update()
		})
		b.Observe("align", func(a []text.AlignType) {
			Alignment = *(*[2]text.AlignType)(a)
			ent.Alignment = a[0]
			ent.Update()
		})
		b.Observe("overflow", func(o text.Overflow) {
			ent.Overflow = o
			ent.Update()
		})
		root.Set("text", t)
	}
}

// Label creates a label on builder.
func (b *Builder) Label(t string) *Entity {
	return b.Add(Label(t))
}
