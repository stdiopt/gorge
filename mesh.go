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
	"fmt"
	"reflect"
	"unsafe"
)

// DrawType type of draw for the renderer
type DrawType int

// independent from gl drawTypes
const (
	// Default triangles
	DrawTriangles = DrawType(iota)
	DrawTriangleStrip
	DrawTriangleFan
	DrawPoints
	DrawLines
	DrawLineLoop
	DrawLineStrip
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

// MeshLoader struct
/*type MeshLoader struct {
	meshLoader
	Updates int
}*/

// MeshData raw mesh data
type MeshData struct {
	Name     string
	Format   VertexFormat
	Vertices []float32
	Indices  []uint32
	Updates  int
}

func (m *MeshData) String() string {
	return fmt.Sprintf("MeshData: %s, %v verts: %v, ind: %v, upd: %v",
		m.Name,
		m.Format,
		len(m.Vertices), len(m.Indices), m.Updates,
	)
}

// Data returns self to satisfy meshLoader
func (m *MeshData) Data() *MeshData { return m }

// Mesh representation
type Mesh struct {
	asset
	// HashID once we change data we could change this hash?
	// Repdate Flag for dynamic stuff
	DrawType DrawType

	loader MeshLoader // Disallow changes
}

// Loader returns the mesh loader
func (m *Mesh) Loader() MeshLoader {
	return m.loader
}

// NewMesh creates a mesh based on loader
func NewMesh(m MeshLoader) *Mesh {
	return &Mesh{
		loader: m,
	}
}

///////////////////////////////////////////////////////////////////////////////

// Helper mesh struct

// VertexPTN position tex normal vertex
type VertexPTN struct {
	Pos    vec3
	Tex    vec2
	Normal vec3
}

// MeshDataPTN a slice of those vertices
type MeshDataPTN struct {
	Name     string
	Vertices []VertexPTN
	Indices  []uint32
}

// Add a vertex
func (m *MeshDataPTN) Add(p vec3, t vec2, n vec3) {
	m.Vertices = append(m.Vertices, VertexPTN{p, t, n})
}

// Data returns the mesh data
func (m *MeshDataPTN) Data() *MeshData {
	vsize := 3 + 2 + 3
	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&m.Vertices))
	hdr.Len *= vsize
	hdr.Cap *= vsize
	vertices := *(*[]float32)(unsafe.Pointer(&hdr))

	return &MeshData{
		Name:     m.Name,
		Format:   VertexFormatPTN,
		Vertices: vertices,
		Indices:  m.Indices,
	}
}
