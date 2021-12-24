package gorlet

import "github.com/stdiopt/gorge/m32"

func TextButton(t string, clickfn func()) Func {
	return func(b *Builder) {
		var (
			text      = b.Prop("text", t)
			fontScale = b.Prop("fontScale", 2)
			textColor = b.Prop("textColor", m32.Color(0))
		)

		b.SetRoot(Button(clickfn))
		b.SetProps(Props{
			"text":      text,
			"fontScale": fontScale,
			"textColor": textColor,
		})
		b.Label(t)
	}
}

// TextButton add a text button child.
func (b *Builder) TextButton(t string, clickfn func()) *Entity {
	return b.Add(TextButton(t, clickfn))
}
