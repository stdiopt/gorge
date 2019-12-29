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

// Cube generates a cube thing
func Cube() *MeshEntity {

	// All borked
	vertices := []float32{
		// Front face
		-1.0, 1.0, 1.0, 1, 0, 0, 0, 1,
		1.0, 1.0, 1.0, 0, 0, 0, 0, 1,
		1.0, -1.0, 1.0, 0, 1, 0, 0, 1,
		-1.0, -1.0, 1.0, 1, 1, 0, 0, 1,

		// Back face
		1.0, -1.0, -1.0, 1, 1, 0, 0, -1,
		1.0, 1.0, -1.0, 1, 0, 0, 0, -1,
		-1.0, 1.0, -1.0, 0, 0, 0, 0, -1,
		-1.0, -1.0, -1.0, 0, 1, 0, 0, -1,

		// Top face
		1.0, 1.0, -1.0, 1, 1, 0, 1, 0,
		1.0, 1.0, 1.0, 1, 0, 0, 1, 0,
		-1.0, 1.0, 1.0, 0, 0, 0, 1, 0,
		-1.0, 1.0, -1.0, 0, 1, 0, 1, 0,

		// Bottom face
		-1.0, -1.0, 1.0, 0, 1, 0, -1, 0,
		1.0, -1.0, 1.0, 1, 1, 0, -1, 0,
		1.0, -1.0, -1.0, 1, 0, 0, -1, 0,
		-1.0, -1.0, -1.0, 0, 0, 0, -1, 0,

		// Right face
		1.0, -1.0, 1.0, 1, 1, 1, 0, 0,
		1.0, 1.0, 1.0, 1, 0, 1, 0, 0,
		1.0, 1.0, -1.0, 0, 0, 1, 0, 0,
		1.0, -1.0, -1.0, 0, 1, 1, 0, 0,

		// Left face
		-1.0, 1.0, -1.0, 1, 0, -1, 0, 0,
		-1.0, 1.0, 1.0, 0, 0, -1, 0, 0,
		-1.0, -1.0, 1.0, 0, 1, -1, 0, 0,
		-1.0, -1.0, -1.0, 1, 1, -1, 0, 0,
	}
	indices := []uint32{
		0, 1, 2, 0, 2, 3, // front
		//2, 1, 0, 3, 2, 0,
		4, 5, 6, 4, 6, 7, // back
		//7, 6, 4, 6, 5, 4,
		8, 9, 10, 8, 10, 11, // top
		//11, 10, 8, 10, 9, 8,
		12, 13, 14, 12, 14, 15, // bottom
		//15, 14, 12, 14, 13, 12,

		16, 17, 18, 16, 18, 19, // right
		//19, 18, 16, 18, 17, 16,
		20, 21, 22, 20, 22, 23, // left
		//23, 22, 20, 22, 21, 20,
	}
	data := &gorge.MeshData{
		Format:   gorge.VertexFormatPTN,
		Vertices: vertices,
		Indices:  indices,
	}
	material := gorge.NewMaterial("")

	mesh := gorge.Mesh{MeshLoader: data}
	return &MeshEntity{
		*gorge.NewTransform(),
		gorge.Renderable{
			Color:    vec4{1, 1, 1, 1},
			Mesh:     &mesh,
			Material: material,
		},
		//material,
	}
}
