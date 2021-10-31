package primitive

import (
	"github.com/stdiopt/gorge"
)

// NewCube generates a cube mesh
func NewCube() *gorge.Mesh {
	return gorge.NewMesh(CubeMeshData())
}

// CubeMeshData returns a 1x1x1 cube meshData.
func CubeMeshData() *gorge.MeshData {
	return &gorge.MeshData{
		Name:   "primitive.Cube:",
		Format: gorge.VertexFormatPTN(),
		Vertices: []float32{
			// Front face
			-.5, .5, .5, 0, 0, 0, 0, 1,
			.5, .5, .5, 1, 0, 0, 0, 1,
			.5, -.5, .5, 1, 1, 0, 0, 1,
			-.5, -.5, .5, 0, 1, 0, 0, 1,

			// Back face
			.5, -.5, -.5, 0, 1, 0, 0, -1,
			.5, .5, -.5, 0, 0, 0, 0, -1,
			-.5, .5, -.5, 1, 0, 0, 0, -1,
			-.5, -.5, -.5, 1, 1, 0, 0, -1,

			// Top face
			.5, .5, -.5, 1, 0, 0, 1, 0,
			.5, .5, .5, 1, 1, 0, 1, 0,
			-.5, .5, .5, 0, 1, 0, 1, 0,
			-.5, .5, -.5, 0, 0, 0, 1, 0,

			// Bottom face
			-.5, -.5, .5, 1, 1, 0, -1, 0,
			.5, -.5, .5, 0, 1, 0, -1, 0,
			.5, -.5, -.5, 0, 0, 0, -1, 0,
			-.5, -.5, -.5, 1, 0, 0, -1, 0,

			// Right face
			.5, -.5, .5, 0, 1, 1, 0, 0,
			.5, .5, .5, 0, 0, 1, 0, 0,
			.5, .5, -.5, 1, 0, 1, 0, 0,
			.5, -.5, -.5, 1, 1, 1, 0, 0,

			// Left face
			-.5, .5, -.5, 0, 0, -1, 0, 0,
			-.5, .5, .5, 1, 0, -1, 0, 0,
			-.5, -.5, .5, 1, 1, -1, 0, 0,
			-.5, -.5, -.5, 0, 1, -1, 0, 0,
		},
		Indices: []uint32{
			0, 1, 2, 0, 2, 3, // front
			4, 5, 6, 4, 6, 7, // back
			8, 9, 10, 8, 10, 11, // top
			12, 13, 14, 12, 14, 15, // bottom
			16, 17, 18, 16, 18, 19, // right
			20, 21, 22, 20, 22, 23, // left
		},
	}
}
