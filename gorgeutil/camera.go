package gorgeutil

import (
	"fmt"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
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
			ClearColor: m32.Vec3{0.4, 0.4, 0.4},
			Viewport:   m32.Vec4{0, 0, 1, 1},
		},
	}
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
			ClearColor:     m32.Vec3{0.4, 0.4, 0.4},
			Viewport:       m32.Vec4{0, 0, 1, 1},
		},
	}
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
			ClearColor:     m32.Vec3{0.4, 0.4, 0.4},
			Viewport:       m32.Vec4{0, 0, 1, 1},
		},
	}
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
			Viewport:       m32.Vec4{0, 0, 1, 1},
			Order:          100,
			CullMask:       gorge.CullMaskUI,
		},
	}
}
