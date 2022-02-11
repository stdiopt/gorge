package gorgeutil

import (
	"fmt"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
)

// Camera basic camera entity.
type Camera struct {
	Name string
	gorge.TransformComponent
	gorge.CameraComponent
}

func (c *Camera) String() string {
	return fmt.Sprintf("<camera:%s:%s>", c.ProjectionType, c.Name)
}

// SetName sets the camera name for debugging purposes.
func (c *Camera) SetName(n string) {
	c.Name = n
}

// NewCamera returns a camera entity with transform and camera components.
func NewCamera() *Camera {
	return &Camera{
		Name:               "",
		TransformComponent: gorge.TransformIdent(),
		CameraComponent: gorge.CameraComponent{
			Fov:        90,
			Near:       0.1,
			Far:        1000,
			ClearColor: gm.Vec3{0.4, 0.4, 0.4},
			Viewport:   gm.Vec4{0, 0, 1, 1},
		},
	}
}

func AddCamera(a Contexter, fov, near, far float32) *Camera {
	c := NewPerspectiveCamera(fov, near, far)
	a.Add(c)
	return c
}

// NewOrthoCamera returns a camera defaulted to ortho projection.
func NewOrthoCamera(size, near, far float32) *Camera {
	return &Camera{
		TransformComponent: gorge.TransformIdent(),
		CameraComponent: gorge.CameraComponent{
			ProjectionType: gorge.ProjectionOrtho,
			OrthoSize:      size,
			Near:           near,
			Far:            far,
			ClearColor:     gm.Vec3{0.4, 0.4, 0.4},
			Viewport:       gm.Vec4{0, 0, 1, 1},
		},
	}
}

func AddOrthoCamera(a Contexter, size, near, far float32) *Camera {
	c := NewOrthoCamera(size, near, far)
	a.Add(c)
	return c
}

// NewPerspectiveCamera returns a default perspective camera.
func NewPerspectiveCamera(fov, near, far float32) *Camera {
	return &Camera{
		TransformComponent: gorge.TransformIdent(),
		CameraComponent: gorge.CameraComponent{
			ProjectionType: gorge.ProjectionPerspective,
			Fov:            fov,
			Near:           near,
			Far:            far,
			ClearColor:     gm.Vec3{0.4, 0.4, 0.4},
			Viewport:       gm.Vec4{0, 0, 1, 1},
		},
	}
}

func AddPerspectiveCamera(a Contexter, fov, near, far float32) *Camera {
	c := NewPerspectiveCamera(fov, near, far)
	a.Add(c)
	return c
}

// NewUICamera returns an ortho camera with a specific CullMask 1<<17 for UI.
func NewUICamera() *Camera {
	return &Camera{
		TransformComponent: gorge.TransformIdent(),
		CameraComponent: gorge.CameraComponent{
			ProjectionType: gorge.ProjectionOrtho,
			OrthoSize:      100,
			Near:           -100,
			Far:            100,
			ClearFlag:      gorge.ClearDepthOnly,
			Viewport:       gm.Vec4{0, 0, 1, 1},
			Order:          100,
			CullMask:       gorge.CullMaskUI,
		},
	}
}

func AddUICamera(a Contexter) *Camera {
	c := NewUICamera()
	a.Add(c)
	return c
}
