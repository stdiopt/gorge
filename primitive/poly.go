package primitive

import (
	"math"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
)

// NewPoly creates a poly mesh
func NewPoly(n int) *gorge.Mesh {
	m := gorge.NewMesh(PolyMeshData(n))
	m.DrawMode = gorge.DrawTriangleFan
	return m
}

// PolyMeshData returns a poly as meshData.
func PolyMeshData(n int) *gorge.MeshData {
	points := []float32{}
	p := gm.Vec3{0, 1, 0}
	theta := float32(math.Pi) / (float32(n) / 2)
	r := gm.M3Rotate(theta)
	for i := 0; i < n+1; i++ {
		points = append(points, p[:]...)
		p = r.MulV3(p)
	}
	return &gorge.MeshData{
		Source:   "primitive.Poly",
		Format:   gorge.VertexFormatP(),
		Vertices: points,
		Indices:  nil,
	}
}
