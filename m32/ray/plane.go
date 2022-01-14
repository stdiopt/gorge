package ray

import (
	"github.com/stdiopt/gorge/m32"
)

// IntersectPlane and returns where it was intersected
func IntersectPlane(r Ray, planeNormal, planePoint m32.Vec3) Result {
	if nl := planeNormal.Len(); nl != nl {
		return Result{}
	}
	diff := r.Position.Sub(planePoint)
	prod1 := diff.Dot(planeNormal)
	prod2 := r.Direction.Dot(planeNormal)
	prod3 := prod1 / prod2
	pos := r.Position.Sub(r.Direction.Mul(prod3))

	return Result{Hit: true, Position: pos}
}

// CalcNormal returns a normal from 3 points
func CalcNormal(a, b, c m32.Vec3) m32.Vec3 {
	v1 := b.Sub(a)
	v2 := c.Sub(a)
	return v1.Cross(v2).Normalize()
}
