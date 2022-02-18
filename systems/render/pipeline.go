package render

import (
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/render/bufutil"
	"github.com/stdiopt/gorge/systems/render/gl"
)

// StepFunc func that will render.
type StepFunc func(ri *Step)

// Step holds information about a render Step.
type Step struct {
	RenderNumber int
	Viewport     gm.Vec4

	CameraUBO *bufutil.NamedOffset

	// Current rendering camera.
	Camera Camera
	Lights []HLight

	QueuesIndex []int
	Queues      map[int]*Queue

	StencilDirty bool

	// Global specified uniforms, could be fetch from camera
	// or directly on props
	Projection gm.Mat4
	View       gm.Mat4
	CamPos     gm.Vec3
	Ambient    gm.Vec3

	// Global props that will be set in every material
	// As in, defaults
	Ubos     map[string]gl.Buffer
	Props    map[string]any
	Samplers map[string]*Texture
}

// Queue renderer queue.
type Queue struct {
	Sort        int
	Renderables []*RenderableGroup
}
