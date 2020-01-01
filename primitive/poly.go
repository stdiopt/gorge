// Copyright 2019 Luis Figueiredo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package primitive

import (
	"math"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
)

// Poly creates a polygon with n sides
func Poly(n int) *MeshEntity {
	mesh := PolyMesh(n)
	mat := gorge.NewMaterial(nil)
	// Should be on mesh perhaps?
	return &MeshEntity{
		*gorge.NewTransform(),
		gorge.Renderable{
			Color:    vec4{1, 1, 1, 1},
			Mesh:     mesh,
			Material: mat,
		},
	}
}

// PolyMesh creates a poly mesh
func PolyMesh(n int) *gorge.Mesh {
	points := []float32{}
	p := vec3{0, 1, 0}
	theta := float32(math.Pi) / (float32(n) / 2)
	r := m32.HomogRotate2D(theta)
	for i := 0; i < n+1; i++ {
		points = append(points, p[:]...)
		p = r.Mul3x1(p)
	}
	data := &gorge.MeshData{
		Name:     "primitive.Poly",
		Format:   gorge.VertexFormatP,
		Vertices: points,
		Indices:  nil,
	}
	m := gorge.NewMesh(data)
	m.DrawType = gorge.DrawTriangleFan

	return m
}
