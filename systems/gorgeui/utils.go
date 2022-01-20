package gorgeui

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/m32/ray"
)

// RectTransformer transform interface for UI elements
type RectTransformer interface {
	RectTransform() *RectComponent
}

type transformer interface {
	Transform() *gorge.TransformComponent
}

// Helper functions to transform a slice of floats into a specific size
// helful for arguments
func v2f(v ...float32) m32.Vec2 {
	switch len(v) {
	case 0:
		return m32.Vec2{}
	case 1:
		return m32.Vec2{v[0], v[0]}
	default:
		return m32.Vec2{v[0], v[1]}
	}
}

func v3f(v ...float32) m32.Vec3 {
	switch len(v) {
	case 0:
		return m32.Vec3{}
	case 1:
		return m32.Vec3{v[0], v[0], v[0]}
	case 2: // could be v[0],v[1],v[0],v[1]
		return m32.Vec3{v[0], v[1], 0}
	default: // Odd case should not be used
		return m32.Vec3{v[0], v[1], v[2]}
	}
}

func v4f(v ...float32) m32.Vec4 {
	switch len(v) {
	case 0:
		return m32.Vec4{}
	case 1:
		return m32.Vec4{v[0], v[0], v[0], v[0]}
	case 2: // could be v[0],v[1],v[0],v[1]
		return m32.Vec4{v[0], v[0], v[1], v[1]}
	case 3: // Odd case should not be used
		return m32.Vec4{v[0], v[1], v[2], v[1]}
	default:
		return m32.Vec4{v[0], v[1], v[2], v[3]}
	}
}

func rayRect(r ray.Ray, e Entity) ray.Result {
	rect := e.RectTransform().Rect()
	m := e.Mat4()

	v0 := m.MulV4(m32.Vec4{rect[0], rect[1], 0, 1}).Vec3()
	v1 := m.MulV4(m32.Vec4{rect[2], rect[1], 0, 1}).Vec3() // right
	v2 := m.MulV4(m32.Vec4{rect[0], rect[3], 0, 1}).Vec3() // up)

	return ray.IntersectRect(r, v0, v1, v2)
}
