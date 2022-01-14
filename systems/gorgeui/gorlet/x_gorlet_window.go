package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/m32/ray"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// Window creates a draggable window with a title bar.
func Window(def string) Func {
	return func(b *Builder) {
		const (
			spacing = .4
		)
		var (
			// fontScale      = b.Prop("fontScale", 2)
			winColor       = b.Prop("background", m32.Color(0, .3))
			titleFontScale = b.Prop("title.fontScale", 1.5)
			titleColor     = b.Prop("title.color", m32.Color(0, 0, .3, .3))
			titleTextColor = b.Prop("title.textcolor", m32.Color(1))
			titleText      = b.Prop("title.text", def)
		)

		root := b.Root()

		b.Use("color", winColor)
		// b.Layout(gorgeui.AutoHeight(1))
		b.BeginPanel()
		// full.SetAnchor(0, 0, 1, 1)
		b.UseProps(Props{
			"fontScale": titleFontScale,
			"color":     titleColor,
			"textColor": titleTextColor,
		})
		b.UseAnchor(0, 0, 1, 0)
		b.UseRect(0, 0, 0, 2)
		title := b.BeginPanel()
		{
			title.SetDragEvents(true)
			gorge.HandleFunc(title, func(e gorgeui.EventDrag) {
				ui := gorgeui.RootUI(root)
				wp := ray.FromScreen(ui.ScreenSize(), ui.Camera, e.Delta).GetPoint(1)
				wp = wp.Sub(ray.FromScreen(ui.ScreenSize(), ui.Camera, m32.Vec2{}).GetPoint(1))
				root.Translate(wp[0], -wp[1], 0)
			})

			b.UseAnchor(0, 0, 1, 1)
			b.Use("text", titleText)
			b.Label("")
		}
		b.EndPanel()

		// Body
		b.Use("color", winColor)
		b.UseAnchor(0, 0, 1, 0)
		b.UseRect(0, 2, 0, 0)
		b.BeginPanel()
		b.ClientArea()
		b.EndPanel()

		b.Use("color", nil)
		b.Use("textColor", nil)

		// b.UseAnchor(0, 0, 1, 1)
		// b.UseRect(0, 0, 1, 0)
		// b.TextButton("status", nil)

		b.UseAnchor(1)
		b.UseRect(0, 0, 1, 1)
		b.UsePivot(.8)
		b.Use("color", titleColor)

		resizer := b.Add(Quad())
		resizer.SetDragEvents(true)
		gorge.HandleFunc(resizer, func(e gorgeui.EventDrag) {
			ui := gorgeui.RootUI(root)
			wp := ray.FromScreen(ui.ScreenSize(), ui.Camera, e.Delta).GetPoint(1)
			wp = wp.Sub(ray.FromScreen(ui.ScreenSize(), ui.Camera, m32.Vec2{}).GetPoint(1))
			root.Dim = root.Dim.Add(m32.Vec2{wp[0], -wp[1]})
		})
	}
}

// WindowWrap wraps a window directly.
func WindowWrap(title string, bodyFn Func) Func {
	return func(b *Builder) {
		b.SetRoot(Window(title))
		body := b.Add(bodyFn)
		b.ForwardProps("", body)
	}
}

// BeginWindow begins a window.
func (b *Builder) BeginWindow(titleText string) *Entity {
	return b.Begin(Window(titleText))
}

// EndWindow alias to End().
func (b *Builder) EndWindow() {
	b.End()
}
