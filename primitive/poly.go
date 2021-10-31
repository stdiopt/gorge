package primitive

import (
	"math"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
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
	p := m32.Vec3{0, 1, 0}
	theta := float32(math.Pi) / (float32(n) / 2)
	r := m32.M3Rotate(theta)
	for i := 0; i < n+1; i++ {
		points = append(points, p[:]...)
		p = r.MulV3(p)
	}
	return &gorge.MeshData{
		Name:     "primitive.Poly",
		Format:   gorge.VertexFormatP(),
		Vertices: points,
		Indices:  nil,
	}
}
