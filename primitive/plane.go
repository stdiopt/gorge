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
	"github.com/stdiopt/gorge"
)

var planeMesh = gorge.NewMesh(&gorge.MeshData{
	Name:   "primitive.Plane",
	Format: gorge.VertexFormatPTN,
	Vertices: []float32{
		-1, 0, -1, 0, 1, 0, 1, 0,
		1, 0, -1, 1, 1, 0, 1, 0,
		1, 0, 1, 1, 0, 0, 1, 0,
		-1, 0, 1, 0, 0, 0, 1, 0,
	},
	Indices: []uint32{
		0, 1, 2,
		2, 3, 0,
	},
})

// Plane returns a mesh entity ready to add
func Plane() *MeshEntity {
	mat := gorge.NewMaterial(nil)
	return &MeshEntity{
		*gorge.NewTransform(),
		gorge.Renderable{
			Color:    vec4{1, 1, 1, 1},
			Mesh:     planeMesh,
			Material: mat,
		},
	}
}

// PlaneMesh creates a polygon facing Z
func PlaneMesh() *gorge.Mesh {
	return planeMesh
}
