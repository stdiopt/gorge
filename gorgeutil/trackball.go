package gorgeutil

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/input"
)

// CameraRig thing
type CameraRig struct {
	*gorge.TransformComponent
	Vert   *gorge.TransformComponent
	Camera cameraEntity

	lastP *m32.Vec2

	dragging bool
}

// GetEntities returns the underlying camera entities.
func (r *CameraRig) GetEntities() []gorge.Entity {
	return []gorge.Entity{r.Camera}
}

// HandleEvent implements the event handler interface.
func (r *CameraRig) HandleEvent(ee event.Event) {
	e, ok := ee.(input.EventPointer)
	if !ok {
		return
	}
	if r.lastP == nil {
		r.lastP = &m32.Vec2{}
		*r.lastP = e.Pointers[0].Pos
		return
	}
	delta := e.Pointers[0].Pos.Sub(*r.lastP)
	*r.lastP = e.Pointers[0].Pos
	if e.Type == input.MouseWheel {
		t := r.Camera.Transform()
		dist := t.WorldPosition().Len()
		multiplier := dist * 0.005
		t.Translate(0, 0, e.Pointers[0].DeltaZ*multiplier)
		if t.Position[2] < 0 {
			t.Position[2] = 0
		}
	}
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
			v := m32.Vec2{delta[1], -delta[0]}.Mul(scale)
			r.Vert.Rotate(-v[0], 0, 0)
			r.Transform().Rotate(0, v[1], 0)
		}
	}
}

type cameraEntity interface {
	Transform() *gorge.TransformComponent
	Mat4() m32.Mat4
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
