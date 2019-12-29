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

package gorge

import (
	"reflect"
	"unsafe"
)

// VertexFormat vertex Formats
type VertexFormat int

// Vertex formats
const (
	VertexFormatP   = VertexFormat(iota)
	VertexFormatPN  // Vertex Normal
	VertexFormatPT  // Vertex Texture
	VertexFormatPTN // Vertex Texture Normal
	VertexPNC       // Vertex Normal Color
)

// MeshLoader is a mesh loader interface
type MeshLoader interface {
	Data() *MeshData
}

// MeshData raw mesh data
type MeshData struct {
	Format   VertexFormat
	Vertices []float32
	Indices  []uint32
}

// Data returns self to satisfy meshLoader
func (m *MeshData) Data() *MeshData { return m }

// Mesh representation
type Mesh struct {
	// HashID once we change data we could change this hash?
	// Repdate Flag for dynamic stuff
	MeshLoader

	// If we update mesh we should increment the update counter
	Updates int
}

///////////////////////////////////////////////////////////////////////////////

// Helper mesh struct

// VertexPTN position tex normal vertex
type VertexPTN struct {
	Pos    vec3
	Tex    vec2
	Normal vec3
}

// MeshPTN a slice of those vertices
type MeshPTN struct {
	Vertices []VertexPTN
	Indices  []uint32
}

// Add a vertex
func (m *MeshPTN) Add(p vec3, t vec2, n vec3) {
	m.Vertices = append(m.Vertices, VertexPTN{p, t, n})
}

// Data returns the mesh data
func (m *MeshPTN) Data() *MeshData {
	vsize := 3 + 2 + 3
	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&m.Vertices))
	hdr.Len *= vsize
	hdr.Cap *= vsize
	vertices := *(*[]float32)(unsafe.Pointer(&hdr))

	return &MeshData{
		Format:   VertexFormatPTN,
		Vertices: vertices,
		Indices:  m.Indices,
	}
}
