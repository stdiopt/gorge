package gorgeutil

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/systems/input"
)

// CameraRig thing
type CameraRig struct {
	*gorge.TransformComponent
	Vert   *gorge.TransformComponent
	Camera cameraEntity

	lastP *gm.Vec2

	dragging      bool
	disableEvents bool
}

// GetEntities returns the underlying camera entities.
func (r *CameraRig) GetEntities() []gorge.Entity {
	return []gorge.Entity{r.Camera}
}

// HandleEvent implements the event handler interface.
func (r *CameraRig) HandleEvent(e event.Event) {
	switch e := e.(type) {
	case gorgeui.EventDragging:
		r.disableEvents = e.Entity != nil
	case input.EventPointer:
		if r.disableEvents {
			return
		}

		if r.lastP == nil {
			r.lastP = &gm.Vec2{}
			*r.lastP = e.Pointers[0].Pos
			return
		}
		delta := e.Pointers[0].Pos.Sub(*r.lastP)
		*r.lastP = e.Pointers[0].Pos

		// Commented until we find a way to avoid scrolling and zoom out simultaneously.
		/*
			if e.Type == input.MouseWheel {
				t := r.Camera.Transform()
				dist := t.WorldPosition().Len()
				multiplier := dist * 0.005
				t.Translate(0, 0, e.Pointers[0].ScrollDelta[1]*multiplier)
				if t.Position[2] < 0 {
					t.Position[2] = 0
				}
			}
		*/

		if e.Type == input.MouseDown {
			r.dragging = true
		}
		if e.Type == input.MouseUp {
			r.dragging = false
		}

		// If dragging or pointer move
		if r.dragging || e.Type == input.PointerMove {
			if len(e.Pointers) == 1 {
				scale := float32(0.005)
				v := gm.Vec2{delta[1], -delta[0]}.Mul(scale)
				r.Vert.Rotate(-v[0], 0, 0)
				r.Transform().Rotate(0, v[1], 0)
			}
		}
	}
}

type cameraEntity interface {
	Transform() *gorge.TransformComponent
	Mat4() gm.Mat4
	Camera() *gorge.CameraComponent
}

// NewTrackballCamera attaches events and all to make a trackball
func NewTrackballCamera(c cameraEntity) *CameraRig {
	if c == nil {
		c = NewCamera()
	}
	camRig := &CameraRig{
		TransformComponent: gorge.NewTransformComponent(),
		Vert:               gorge.NewTransformComponent(),
		Camera:             c,
	}
	t := c.Transform()
	t.SetParent(camRig.Vert)
	t.SetEuler(0, 0, 0)
	t.SetPosition(0, 0, 0)

	camRig.Vert.SetParent(camRig)
	camRig.Vert.Rotate(-0.7, 0, 0)

	return camRig
}

func AddTrackballCamera(a entityAdder) *CameraRig {
	c := NewTrackballCamera(nil)
	a.Add(c)
	a.AddHandler(c)
	return c
}
