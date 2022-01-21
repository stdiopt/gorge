package gorlet

import "github.com/stdiopt/gorge/math/gm"

func TextButton(t string, clickfn func()) Func {
	return func(b *Builder) {
		var (
			text         = b.Prop("text", t)
			fontScale    = b.Prop("fontScale", 2)
			textColor    = b.Prop("textColor", gm.Color(0))
			textOverflow = b.Prop("textOverflow")
		)

		b.SetRoot(Button(clickfn))
		b.UseProps(Props{
			"text":      text,
			"fontScale": fontScale,
			"textColor": textColor,
			"overflow":  textOverflow,
		})
		b.Label(t)
	}
}

// TextButton add a text button child.
func (b *Builder) TextButton(t string, clickfn func()) *Entity {
	return b.Add(TextButton(t, clickfn))
}
