package debug

import (
	"log"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/gorgeutil"
	"github.com/stdiopt/gorge/m32"
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

	PointLight *gorgeutil.Light
	DirLight   *gorgeutil.Light
	CamRig     *gorgeutil.CameraRig

	camRot float32
	camVec m32.Vec3

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
		if b.PointLight.CastShadows == gorge.CastShadowEnabled {
			b.PointLight.CastShadows = gorge.CastShadowDisabled
		} else {
			b.PointLight.CastShadows = gorge.CastShadowEnabled
		}
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

	cam := b.CamRig.Camera.Camera()

	if b.input.KeyDown(input.KeyW) {
		camTransform := b.CamRig.Camera.Transform()
		b.camVec = b.camVec.Add(camTransform.Forward().Mul(power))
	}
	if b.input.KeyDown(input.KeyS) {
		camTransform := b.CamRig.Camera.Transform()
		b.camVec = b.camVec.Add(camTransform.Backward().Mul(power))
	}
	// Ambient control
	if b.input.KeyDown(input.KeyR) {
		cam.ClearColor[0] =
			m32.Min(cam.ClearColor[0]+step, 1)
	}
	if b.input.KeyDown(input.KeyF) {
		cam.ClearColor[0] =
			m32.Max(cam.ClearColor[0]-step, 0)
	}
	if b.input.KeyDown(input.KeyT) {
		cam.ClearColor[1] =
			m32.Min(cam.ClearColor[1]+step, 1)
	}
	if b.input.KeyDown(input.KeyG) {
		cam.ClearColor[1] =
			m32.Max(cam.ClearColor[1]-step, 0)
	}
	if b.input.KeyDown(input.KeyY) {
		cam.ClearColor[2] =
			m32.Min(cam.ClearColor[2]+step, 1)
	}
	if b.input.KeyDown(input.KeyH) {
		cam.ClearColor[2] =
			m32.Max(cam.ClearColor[2]-step, 0)
	}
}

// NewBasic creates a default scene.
func NewBasic(g *gorge.Context) *Basic {
	container := gorge.Container{}

	ui := gorgeui.FromContext(g)
	ic := input.FromContext(g)

	lightRoot := gorge.NewTransformComponent()

	pointLight := gorgeutil.NewPointLight()
	pointLight.SetParent(lightRoot)
	pointLight.SetRange(300)
	pointLight.SetPosition(1, 5, 2)

	pointLightGimbal := gorgeutil.NewGimbal()
	pointLightGimbal.SetParent(pointLight)

	cam := gorgeutil.NewCamera()
	cam.SetClearFlag(gorge.ClearSkybox)
	cam.SetClearColor(.4, .4, .4)

	camRig := gorgeutil.NewTrackballCamera(cam)
	camRig.SetPosition(0, 5, 5)

	camGimbal := gorgeutil.NewGimbal()
	camGimbal.SetParent(camRig)

	dirLight := gorgeutil.NewDirectionalLight()
	dirLight.SetCastShadows(gorge.CastShadowEnabled)
	dirLight.SetPosition(10, 10, -5)
	dirLight.LookAtPosition(m32.Vec3{0, 0, 0})

	dirLightGimbal := gorgeutil.NewGimbal()
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
		gorge:      g,
		input:      ic,
		gorgeui:    ui,
		PointLight: pointLight,
		DirLight:   dirLight,
		CamRig:     camRig,
		lightRoot:  lightRoot,
	}
}

// AddBasic adds a basic scene to the given gorge context.
func AddBasic(g *gorge.Context) *Basic {
	b := NewBasic(g)
	g.Add(b)
	g.Handle(b)
	return b
}

// BasicSystem initializes a default scene when used in app initializator.
func BasicSystem(g *gorge.Context) error {
	log.Println("initializing helper system")
	thing := NewBasic(g)
	// g.PutProp(thing)
	g.Add(thing)
	g.Handle(thing)
	return nil
}
