package gorlet

import (
	"github.com/stdiopt/gorge/text"
)

// Labeled returns a gorlet that labels an entity.
func Labeled(lbl string) BuildFunc {
	return func(b *Builder) {
		var (
			fontScale = b.Prop("fontScale", 2)
			textAlign = []text.AlignType{text.AlignEnd, text.AlignCenter}
			text      = b.Prop("text", lbl)
		)

		b.UseAnchor(0, 0, 1, 1)
		b.UseRect(0)
		b.BeginContainer(LayoutFlexHorizontal(1, 3))
		{

			b.Set("fontScale", fontScale)
			b.Set("textAlign", textAlign)
			b.Set("text", text)
			b.UseRect(.3)
			b.Label("")

			b.UseAnchor(0, 0, 1, 1)
			b.UseRect(.3)
			b.BeginContainer()
			b.ClientArea()
			b.EndContainer()
		}

		b.EndContainer()
		// b.AddEntity(e)
	}
}

// BeginLabeled creates a labeled Entity.
func (b *Builder) BeginLabeled(lbl string) *Entity {
	return b.Begin(Labeled(lbl))
}

// EndLabeled is an alias to End.
func (b *Builder) EndLabeled() {
	b.End()
}
