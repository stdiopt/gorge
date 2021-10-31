package ray

import (
	"github.com/stdiopt/gorge/m32"
)

// IntersectRect at point p1, edge e1, edge e2
//
//    c---------+
//    |         |
//    |         |
//    |         |
//    a---------b
func IntersectRect(r Ray, a, b, c m32.Vec3) Result {
	planeNormal := CalcNormal(a, b, c)

	planeRes := IntersectPlane(r, planeNormal, a)
	if !planeRes.Hit {
		return planeRes
	}
	ipos := planeRes.Position
	var vlen, plen m32.Vec3
	var t float32

	vlen = b.Sub(a)
	plen = planeRes.Position.Sub(a)
	t = plen.Dot(vlen) / vlen.Dot(vlen)
	if t < 0 || t > 1 { // out of edge
		return Result{Hit: false, Position: ipos}
	}

	vlen = c.Sub(a)
	plen = ipos.Sub(a)
	t = plen.Dot(vlen) / vlen.Dot(vlen)
	if t < 0 || t > 1 {
		return Result{Hit: false, Position: ipos}
	}
	return Result{Hit: true, Position: ipos}
}

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
