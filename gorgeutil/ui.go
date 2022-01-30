package gorgeutil

import "github.com/stdiopt/gorge/systems/gorgeui"

func NewUI() *gorgeui.UI {
	return gorgeui.New()
}

func AddUI(a entityAdder) *gorgeui.UI {
	ui := NewUI()
	a.Add(ui)
	return ui
}

// UI returns a gorgeui.New(gorge.Context) with the injected context.
func (c Context) UI(cam cameraEntity) *gorgeui.UI {
	ui := gorgeui.New()
	ui.SetCamera(cam)
	c.Add(ui)
	return ui
}
