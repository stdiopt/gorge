// Package ray implements some ray casting math
package ray

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
)

// Result for a ray test
type Result struct {
	Hit      bool
	Position gm.Vec3
	UV       gm.Vec2
	// Normal
}

// Ray to be used on ray intersection testing.
type Ray struct {
	Position  gm.Vec3
	Direction gm.Vec3
}

// GetPoint on distance d on ray
func (r Ray) GetPoint(d float32) gm.Vec3 {
	return r.Position.Add(r.Direction.Mul(d))
}

// if we want to use Ray within gorge logic, might be better to move this to gorge
type cameraEntity interface {
	gorge.Matrixer
	Camera() *gorge.CameraComponent
}

// https://antongerdelan.net/opengl/raycasting.html
func FromScreen(screenSize gm.Vec2, camEnt cameraEntity, pos gm.Vec2) Ray {
	cam := camEnt.Camera()

	vp := cam.CalcViewport(screenSize)
	width, height := vp[2], vp[3]
	pos = pos.Sub(vp.Vec2())

	proj := cam.Projection(screenSize)
	cm := camEnt.Mat4()

	nds := gm.Vec3{
		(2*pos[0])/width - 1,
		1 - (2*pos[1])/height,
		1,
	}

	clip := gm.Vec4{nds[0], nds[1], -1, 1}
	eye := proj.Inv().MulV4(clip).Vec2().Vec4(-1, 0)
	dir := cm.MulV4(eye).Vec3()

	if cam.ProjectionType == gorge.ProjectionOrtho {
		return Ray{
			Position:  dir,
			Direction: cm.MulV4(gm.Forward().Vec4(0)).Vec3(),
		}
	}
	// Ray from camera Entity func somewhere
	return Ray{
		Position:  cm.Col(3).Vec3(),
		Direction: dir.Normalize(),
	}
}
