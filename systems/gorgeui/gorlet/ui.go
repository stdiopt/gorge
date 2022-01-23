package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gorgeutil"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type gorui = gorgeui.UI

// UI helper that wraps gorgeui.UI and provides some additional functionality
type UI struct {
	*gorui
	gorge     *gorge.Context
	defCamera *gorgeutil.Camera
}

func AddUI(g *gorge.Context) *UI {
	u := gorgeutil.FromContext(g)
	uiCam := u.UICamera()
	ui := u.UI(uiCam)

	return &UI{
		gorui:     ui,
		gorge:     g,
		defCamera: uiCam,
	}
}

func (u *UI) Create(fn Func) *Entity {
	ent := Create(fn)
	u.Add(ent)
	return ent
}

type cameraEntity interface {
	gorge.Transformer
	Camera() *gorge.CameraComponent
}

func (u *UI) SetCamera(cam cameraEntity) {
	if cam, ok := cam.(*gorgeutil.Camera); ok && cam != u.defCamera {
		u.gorge.Remove(u.defCamera)
	} else if cam == nil {
		u.gorge.Add(u.defCamera)
	}
	u.gorui.SetCamera(cam)
}
