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

	/*
		eps := float32(0.01)

		if gm.FloatEqualThreshold(uv[0], 0, eps) ||
			gm.FloatEqualThreshold(uv[0], 1, eps) ||
			gm.FloatEqualThreshold(uv[1], 0, eps) ||
			gm.FloatEqualThreshold(uv[1], 1, eps) { // out of edge
			return Result{Hit: true, Position: ipos, UV: uv}
		}
	*/
	hit := uv[0] >= 0 && uv[0] <= 1 && uv[1] >= 0 && uv[1] <= 1
	//if uv[0] < 0 || uv[0] > 1 || uv[1] < 0 || uv[1] > 1 { // out of edge
	//	return Result{Hit: false, Position: ipos, UV: uv}
	//}

	// vertical side
	// vlen = v2.Sub(v0)
	// plen = ipos.Sub(v0)
	// v := plen.Dot(vlen) / vlen.Dot(vlen)
	return Result{Hit: hit, Position: ipos, UV: uv}
}
