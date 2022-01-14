// Package ray implements some ray casting math
package ray

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
)

// Result for a ray test
type Result struct {
	Hit      bool
	Position m32.Vec3
	UV       m32.Vec2
	// Normal
}

// Ray to be used on ray intersection testing.
type Ray struct {
	Position  m32.Vec3
	Direction m32.Vec3
}

// GetPoint on distance d on ray
func (r Ray) GetPoint(d float32) m32.Vec3 {
	return r.Position.Add(r.Direction.Mul(d))
}

// if we want to use Ray within gorge logic, might be better to move this to gorge
type cameraEntity interface {
	gorge.Matrixer
	Camera() *gorge.CameraComponent
}

// FromScreen returns a ray from screen position through Camera camEnt.
func FromScreen(screenSize m32.Vec2, camEnt cameraEntity, pos m32.Vec2) Ray {
	cam := camEnt.Camera()
	t := camEnt.Mat4()

	vp := cam.CalcViewport(screenSize)
	width, height := vp[2], vp[3]
	pos = pos.Sub(vp.Vec2())
	ndc := m32.Vec4{
		2*pos[0]/width - 1,
		1 - 2*pos[1]/height,
		1, 1,
	}

	m := cam.Projection(screenSize)
	m = m.Mul(t.Inv()).Inv()
	dir := m.MulV4(ndc).Vec3()

	if cam.ProjectionType == gorge.ProjectionOrtho {
		return Ray{
			Position:  dir,
			Direction: t.MulV4(m32.Forward().Vec4(0)).Vec3(),
		}
	}
	// Ray from camera Entity func somewhere
	return Ray{
		Position:  t.Col(3).Vec3(),
		Direction: dir,
	}
}
