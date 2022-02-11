package gorgeutil

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/systems/input"
)

const step = .01

// Basic default scene handler.
// TODO: {lpf} it should be moved to gorgeutil.
type Basic struct {
	gorge.Container

	gorge   *gorge.Context
	input   *input.Context
	gorgeui *gorgeui.Context

	PointLight *Light
	DirLight   *Light
	CamRig     *CameraRig

	camRot float32
	camVec gm.Vec3

	lightRoot *gorge.TransformComponent
}

// HandleEvent on basic scene
func (b *Basic) HandleEvent(v event.Event) {
	if b.gorgeui.Dragging() == nil {
		b.CamRig.HandleEvent(v)
	}
	e, ok := v.(gorge.EventUpdate)
	if !ok {
		return
	}
	dt := float32(e)
	power := float32(.1)

	if b.input.KeyDown(input.KeyLeftShift) {
		power = 1
	}
	if b.input.KeyDown(input.KeyLeftControl) {
		power = .01
	}
	// tot += dt
	if b.input.KeyPress(input.KeyV) {
		b.PointLight.DisableShadow = !b.PointLight.DisableShadow
		// b.PointLight.CastShadows = !b.PointLight.CastShadows
		// b.DirLight.CastShadows = !b.DirLight.CastShadows
	}
	if b.input.KeyDown(input.KeyZ) {
		b.lightRoot.Rotate(0, -dt, 0)
	}
	if b.input.KeyDown(input.KeyX) {
		b.lightRoot.Rotate(0, dt, 0)
	}

	if b.input.KeyDown(input.KeyA) {
		b.camVec = b.camVec.Add(b.CamRig.Left().Mul(power))
	}
	if b.input.KeyDown(input.KeyD) {
		b.camVec = b.camVec.Add(b.CamRig.Right().Mul(power))
	}
	if b.input.KeyDown(input.KeyE) {
		b.camVec = b.camVec.Add(b.CamRig.Up().Mul(power))
	}
	if b.input.KeyDown(input.KeyQ) {
		b.camVec = b.camVec.Add(b.CamRig.Down().Mul(power))
	}
	b.CamRig.Translatev(b.camVec.Mul(dt))
	b.CamRig.Rotate(0, b.camRot*dt, 0)
	b.camVec = b.camVec.Mul(.92)
	b.camRot *= .92

	if b.input.KeyDown(input.KeyW) {
		camTransform := b.CamRig.Camera.Transform()
		b.camVec = b.camVec.Add(camTransform.Forward().Mul(power))
	}
	if b.input.KeyDown(input.KeyS) {
		camTransform := b.CamRig.Camera.Transform()
		b.camVec = b.camVec.Add(camTransform.Backward().Mul(power))
	}
}

// NewBasic creates a default scene.
func NewBasic(g gorge.Contexter) *Basic {
	container := gorge.Container{}

	// These...
	ui := gorgeui.FromContext(g.G())
	ic := input.FromContext(g.G())

	lightRoot := gorge.NewTransformComponent()

	pointLight := NewPointLight()
	pointLight.SetParent(lightRoot)
	pointLight.SetRange(300)
	pointLight.SetPosition(1, 5, 2)

	pointLightGimbal := NewGimbal()
	pointLightGimbal.SetParent(pointLight)

	cam := NewCamera()
	cam.SetClearFlag(gorge.ClearSkybox)
	cam.SetClearColor(.4, .4, .4)

	camRig := NewTrackballCamera(cam)
	camRig.SetPosition(0, 5, 5)

	camGimbal := NewGimbal()
	camGimbal.SetParent(camRig)

	dirLight := NewDirectionalLight()
	dirLight.SetPosition(10, 10, -5)
	dirLight.LookAtPosition(gm.Vec3{0, 0, 0})

	dirLightGimbal := NewGimbal()
	dirLightGimbal.SetParent(dirLight)

	container.Add(
		dirLight,
		dirLightGimbal,
	)

	container.Add(
		pointLight,
		pointLightGimbal,
		camRig,
		camGimbal,
	)

	return &Basic{
		Container:  container,
		gorge:      g.G(),
		input:      ic,
		gorgeui:    ui,
		PointLight: pointLight,
		DirLight:   dirLight,
		CamRig:     camRig,
		lightRoot:  lightRoot,
	}
}

// AddBasic adds a basic scene to the given gorge context.
func AddBasic(g gorge.Contexter) *Basic {
	b := NewBasic(g)
	g.Add(b)
	g.G().AddHandler(b)
	return b
}
