package primitive

import (
	"github.com/stdiopt/gorge"
)

// PlaneDir plane direction
type PlaneDir int

// Plane directions
// Maybe Use Forward backward, up,down,left,right
const (
	PlaneDirX = PlaneDir(iota)
	PlaneDirY
	PlaneDirZ
	PlaneDirZInv
)

// NewPlane creates a polygon facing Z
func NewPlane(d PlaneDir) *gorge.Mesh {
	return gorge.NewMesh(PlaneMeshData(d))
}

// PlaneMeshData returns a plane meshdata.
func PlaneMeshData(d PlaneDir) *gorge.MeshData {
	var vert []float32
	switch d {
	case PlaneDirZInv:
		vert = []float32{
			-1, 1, 0, 0, 0, 0, 1, 0,
			1, 1, 0, 1, 0, 0, 1, 0,
			1, -1, 0, 1, 1, 0, 1, 0,
			-1, -1, 0, 0, 1, 0, 1, 0,
		}
	case PlaneDirZ:
		vert = []float32{
			-1, 1, 0, 0, 0, 0, 0, 1,
			1, 1, 0, 1, 0, 0, 0, 1,
			1, -1, 0, 1, 1, 0, 0, 1,
			-1, -1, 0, 0, 1, 0, 0, 1,
		}
	case PlaneDirY:
		vert = []float32{
			-1, 0, -1, 0, 0, 0, 1, 0,
			1, 0, -1, 1, 0, 0, 1, 0,
			1, 0, 1, 1, 1, 0, 1, 0,
			-1, 0, 1, 0, 1, 0, 1, 0,
		}
	default:
		panic("undefined direction")
	}

	return &gorge.MeshData{
		Name:     "primitive.Plane",
		Format:   gorge.VertexFormatPTN(),
		Vertices: vert,
		Indices: []uint32{
			0, 1, 2,
			2, 3, 0,
		},
	}
}
