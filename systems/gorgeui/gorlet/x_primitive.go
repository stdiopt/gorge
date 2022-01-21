package gorlet

import (
	"math"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
)

func quadMeshData() *gorge.MeshData {
	return &gorge.MeshData{
		Format: gorge.VertexFormatPTN(),
		Vertices: []float32{
			/*P:*/ 0, 1, 0 /*T*/, 0, 0 /*N*/, 0, 0, 1,
			/*P:*/ 1, 1, 0 /*T*/, 1, 0 /*N*/, 0, 0, 1,
			/*P:*/ 1, 0, 0 /*T*/, 1, 1 /*N*/, 0, 0, 1,
			/*P:*/ 0, 0, 0 /*T*/, 0, 1 /*N*/, 0, 0, 1,
		},
		Indices: []uint32{
			0, 2, 1,
			2, 0, 3,
		},
	}
}

func polyMeshData(n int) *gorge.MeshData {
	points := []float32{}
	p := gm.Vec3{0, .5, 0}
	theta := float32(math.Pi) / (float32(n) / 2)
	r := gm.M3Rotate(theta)
	for i := 0; i < n+1; i++ {
		o := p.Add(gm.Vec3{.5, .5, 0})
		points = append(points, o[:]...)
		p = r.MulV3(p)
	}
	return &gorge.MeshData{
		Format:   gorge.VertexFormatP(),
		Vertices: points,
		Indices:  nil,
	}
}

// Create a rectElement with a Default Material which material can be manipulated
// including the TextMesher
