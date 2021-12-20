package render

import (
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/render/bufutil"
	"github.com/stdiopt/gorge/systems/render/gl"
)

// StepFunc func that will render.
type StepFunc func(ri *Step)

// Step holds information about a render Step.
type Step struct {
	RenderNumber int
	Viewport     m32.Vec4

	CameraUBO *bufutil.NamedOffset

	// Current rendering camera.
	Camera Camera
	Lights []Light

	QueuesIndex []int
	Queues      map[int]*Queue

	// Global specified uniforms, could be fetch from camera
	// or directly on props
	Projection m32.Mat4
	View       m32.Mat4
	CamPos     m32.Vec3
	Ambient    m32.Vec3

	// Global props that will be set in every material
	// As in, defaults
	Ubos     map[string]gl.Buffer
	Props    map[string]interface{}
	Samplers map[string]*Texture
}

// Queue renderer queue.
type Queue struct {
	Sort        int
	Renderables []*RenderableGroup
}
