package primitive

import (
	"math"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
)

// http://songho.ca/opengl/gl_sphere.html

// NewSphere returns a mesh with a sphere meshData
func NewSphere(sectorCount, stackCount int) *gorge.Mesh {
	return gorge.NewMesh(SphereMeshData(sectorCount, stackCount))
}

// SphereMeshData returns the meshData only.
func SphereMeshData(sectorCount, stackCount int) *gorge.MeshData {
	const radius = 1
	verts := []float32{}
	indices := []uint32{}

	sectorStep := 2 * math.Pi / float32(sectorCount)
	stackStep := math.Pi / float32(stackCount)

	for i := 0; i <= stackCount; i++ {
		stackAngle := math.Pi/2 - float32(i)*stackStep
		xz := radius * m32.Cos(stackAngle)
		y := radius * m32.Sin(stackAngle)

		for j := 0; j <= sectorCount; j++ {
			sectorAngle := float32(j) * sectorStep

			x := xz * m32.Cos(sectorAngle)
			z := xz * m32.Sin(sectorAngle)

			verts = append(verts, x, y, z)

			u := float32(j) / float32(sectorCount)
			v := float32(i) / float32(stackCount)
			verts = append(verts, -u, v)

			verts = append(verts, x, y, z) // normal are the same as coords
		}
	}

	// indices
	//  k1--k1+1
	//  |  / |
	//  | /  |
	//  k2--k2+1
	for i := 0; i < stackCount; i++ {
		k1 := uint32(i * (sectorCount + 1))
		k2 := k1 + uint32(sectorCount+1)
		for j := 0; j < sectorCount; j++ {
			if i != 0 {
				indices = append(indices, k1, k2, k1+1)
			}
			if i != (stackCount - 1) {
				indices = append(indices, k1+1, k2, k2+1)
			}
			k1++
			k2++
		}
	}

	return &gorge.MeshData{
		Name:     "primitive.Sphere",
		Format:   gorge.VertexFormatPTN(),
		Vertices: verts,
		Indices:  indices,
	}
}
