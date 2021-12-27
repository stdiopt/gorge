package gorlet

import (
	"github.com/stdiopt/gorge/text"
)

// Labeled returns a gorlet that labels an entity.
func Labeled(lbl string) Func {
	return func(b *Builder) {
		var (
			fontScale = b.Prop("fontScale", 2)
			textAlign = TextAlign(text.AlignEnd, text.AlignCenter)
			text      = b.Prop("text", lbl)
		)

		b.UseAnchor(0, 0, 1, 1)
		b.UseRect(0)
		b.BeginContainer(LayoutFlexHorizontal(1, 3))
		{

			b.Use("fontScale", fontScale)
			b.Use("textAlign", textAlign)
			b.Use("text", text)
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

// Labeled creates a labeled entity by passing the body
// it returns the entity created by fn.
func (b *Builder) Labeled(lbl string, fn Func) *Entity {
	b.BeginLabeled(lbl)
	e := b.Add(fn)
	b.EndLabeled()
	return e
}

// BeginLabeled creates a labeled Entity.
func (b *Builder) BeginLabeled(lbl string) *Entity {
	return b.Begin(Labeled(lbl))
}

// EndLabeled is an alias to End.
func (b *Builder) EndLabeled() {
	b.End()
}
