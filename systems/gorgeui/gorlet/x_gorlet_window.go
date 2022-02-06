package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/math/ray"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// Window creates a draggable window with a title bar.
func Window(def string) Func {
	return func(b *B) {
		const (
			spacing = .4
		)
		var (
			// fontScale      = b.Prop("fontScale", 2)
			winColor       = b.Prop("background", gm.Color(0, .3))
			bodyColor      = b.Prop("body.color", gm.Color())
			titleFontScale = b.Prop("title.fontScale", 1.5)
			titleColor     = b.Prop("title.color", gm.Color(0, 0, .3, .3))
			titleTextColor = b.Prop("title.textcolor", gm.Color(1))
			titleText      = b.Prop("title.text", def)
		)

		root := b.Root()

		b.Use("color", winColor)
		b.BeginPanel()

		b.UseProps(Props{
			"fontScale": titleFontScale,
			"color":     titleColor,
			"textColor": titleTextColor,
		})
		// TitleBar
		b.UseAnchor(0, 0, 1, 0)
		b.UseRect(0, 0, 0, 2)
		title := b.BeginPanel()
		{
			title.SetDragEvents(true)
			event.Handle(title, func(e gorgeui.EventDrag) {
				ui := gorgeui.RootUI(root)
				wp := ray.FromScreen(ui.ScreenSize(), ui.Camera, e.Delta).GetPoint(1)
				wp = wp.Sub(ray.FromScreen(ui.ScreenSize(), ui.Camera, gm.Vec2{}).GetPoint(1))
				root.Translate(wp[0], -wp[1], 0)
			})

			b.UseAnchor(0, 0, 1, 1)
			b.Use("text", titleText)
			b.Label("")
		}
		b.EndPanel()

		// Body
		b.UseRect(0, 2, 0, 0)
		b.Use("color", bodyColor)
		b.Use("overflow", b.Prop("overflow"))
		b.BeginPanel()
		b.ClientArea()
		b.EndPanel()

		b.Use("color", nil)
		b.Use("textColor", nil)

		b.UseAnchor(1)
		b.UseRect(0, 0, 1, 1)
		b.UsePivot(.8)
		b.Use("color", titleColor)

		resizer := b.Quad()
		resizer.SetDragEvents(true)
		b.EndPanel()
		event.Handle(resizer, func(e gorgeui.EventDrag) {
			ui := gorgeui.RootUI(root)
			wp := ray.FromScreen(ui.ScreenSize(), ui.Camera, e.Delta).GetPoint(1)
			wp = wp.Sub(ray.FromScreen(ui.ScreenSize(), ui.Camera, gm.Vec2{}).GetPoint(1))
			root.Size = root.Size.Add(gm.Vec2{wp[0], -wp[1]})
		})
	}
}

// WindowWrap wraps a window directly.
func WindowWrap(title string, bodyFn Func) Func {
	return func(b *B) {
		b.SetRoot(Window(title))
		body := b.Add(bodyFn)
		b.ForwardProps("", body)
	}
}

// BeginWindow begins a window.
func (b *B) BeginWindow(titleText string) *Entity {
	return b.Begin(Window(titleText))
}

// EndWindow alias to End().
func (b *B) EndWindow() {
	b.End()
}
