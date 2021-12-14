package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/m32/ray"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// Window creates a draggable window with a title bar.
func Window(title string, body BuildFunc) BuildFunc {
	return func(b *Builder) {
		const (
			spacing = .4
		)
		var (
			// fontScale      = b.Prop("fontScale", 2)
			winColor       = b.Prop("win.background", m32.Color(0, .3))
			titleFontScale = b.Prop("win.title.fontScale", 1.5)
			titleColor     = b.Prop("win.title.color", m32.Color(0, 0, .3, .3))
			titleTextColor = b.Prop("win.title.textcolor", m32.Color(1))
		)

		root := b.Root()
		root.FillParent(0)

		b.Set("color", winColor)
		b.Layout(gorgeui.AutoHeight(1))
		b.BeginPanel()

		b.SetProps(Props{
			"fontScale": titleFontScale,
			"color":     titleColor,
			"textColor": titleTextColor,
		})
		p := b.BeginPanel()
		p.SetAnchor(0, 0, 1, 0)
		p.SetRect(0, 0, 0, 2)
		p.SetPivot(0)
		p.SetDragEvents(true)
		p.HandleFunc(func(e event.Event) {
			switch e := e.(type) {
			case gorgeui.EventDrag:
				ui := gorgeui.RootUI(root)
				wp := ray.FromScreen(ui.ScreenSize(), ui.Camera, e.Delta).GetPoint(1)
				wp = wp.Sub(ray.FromScreen(ui.ScreenSize(), ui.Camera, m32.Vec2{}).GetPoint(1))
				root.Translate(wp[0], -wp[1], 0)
			default:
			}
		})
		b.Label(title)
		b.EndContainer()

		body := b.Add(body)
		b.ForwardProps("", body)
		body.SetAnchor(0, 0, 1, 1)
		body.SetRect(spacing, 2+spacing, spacing, spacing)
		body.SetPivot(0)
	}
}
