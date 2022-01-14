package ray

import "github.com/stdiopt/gorge/m32"

// IntersectRect at point p1, edge e1, edge e2
//
//    v2--------+
//    |         |
//    |         |
//    |         |
//    v0--------v1
func IntersectRect(r Ray, v0, v1, v2 m32.Vec3) Result {
	planeNormal := CalcNormal(v0, v1, v2)

	planeRes := IntersectPlane(r, planeNormal, v0)
	if !planeRes.Hit {
		return planeRes
	}
	ipos := planeRes.Position
	// var vlen, plen m32.Vec3

	uv := m32.Vec2{}

	// horizontal side
	{
		vlen := v1.Sub(v0)
		plen := planeRes.Position.Sub(v0)
		uv[0] = plen.Dot(vlen) / vlen.Dot(vlen)
	}
	{
		vlen := v2.Sub(v0)
		plen := ipos.Sub(v0)
		uv[1] = plen.Dot(vlen) / vlen.Dot(vlen)
	}
	if uv[0] < 0 || uv[0] > 1 || uv[1] < 0 || uv[1] > 1 { // out of edge
		return Result{Hit: false, Position: ipos, UV: uv}
	}

	// vertical side
	// vlen = v2.Sub(v0)
	// plen = ipos.Sub(v0)
	// v := plen.Dot(vlen) / vlen.Dot(vlen)
	return Result{Hit: true, Position: ipos, UV: uv}
}
