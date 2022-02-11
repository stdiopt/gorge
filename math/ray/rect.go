package ray

import "github.com/stdiopt/gorge/math/gm"

// IntersectRect at point p1, edge e1, edge e2
//
//    v2--------+
//    |         |
//    |         |
//    |         |
//    v0--------v1
func IntersectRect(r Ray, v0, v1, v2 gm.Vec3) Result {
	planeNormal := CalcNormal(v0, v1, v2)

	planeRes := IntersectPlane(r, planeNormal, v0)
	if !planeRes.Hit {
		return planeRes
	}
	ipos := planeRes.Position
	// var vlen, plen gm.Vec3

	uv := gm.Vec2{}

	// horizontal side
	xvec := v1.Sub(v0)
	xlen := ipos.Sub(v0)
	uv[0] = xlen.Dot(xvec) / xvec.Dot(xvec)

	// Vertical side
	yvec := v2.Sub(v0)
	ylen := ipos.Sub(v0)
	uv[1] = ylen.Dot(yvec) / yvec.Dot(yvec)

	hit := uv[0] >= 0 && uv[0] <= 1 && uv[1] >= 0 && uv[1] <= 1

	return Result{Hit: hit, Position: ipos, UV: uv}
}
