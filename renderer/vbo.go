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

package renderer

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gl"
)

type vboManager struct {
	g gl.Context3

	vbos map[interface{}]*vbo
}

func newVBOManager(g gl.Context3) *vboManager {
	vm := &vboManager{
		g:    g,
		vbos: map[interface{}]*vbo{},
	}
	return vm
}

// Get vbo for mesh
func (vm *vboManager) Get(m *gorge.Mesh) *vbo {
	return vm.get(m)
}

func (vm *vboManager) get(m *gorge.Mesh) *vbo {
	if vb, ok := vm.vbos[m]; ok {
		vb.update() // if needed
		return vb
	}
	g := vm.g

	vb := &vbo{
		g:       g,
		mesh:    m,
		VBO:     g.CreateBuffer(),
		EBO:     g.CreateBuffer(),
		updates: -1,
	}
	vb.update()

	vm.vbos[m] = vb
	return vb
}

type vbo struct {
	g           gl.Context3
	Format      gorge.VertexFormat
	VBO         gl.Buffer
	EBO         gl.Buffer
	ElementsLen int
	VertexLen   int

	// Do we need this?
	mesh    *gorge.Mesh
	updates int
}

func (v *vbo) update() bool {
	if v.mesh.Updates == v.updates {
		return false
	}
	if v.mesh.MeshLoader == nil {
		return false
	}
	g := v.g
	// Reload mesh data
	meshData := v.mesh.Data()

	v.Format = meshData.Format
	v.VertexLen = len(meshData.Vertices) / vertSize(meshData.Format)

	g.BindBuffer(gl.ARRAY_BUFFER, v.VBO)
	g.BufferDataX(gl.ARRAY_BUFFER, meshData.Vertices, gl.STATIC_DRAW)

	if meshData.Indices != nil && len(meshData.Indices) > 0 {
		v.ElementsLen = len(meshData.Indices)
		g.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, v.EBO)
		g.BufferDataX(gl.ELEMENT_ARRAY_BUFFER, meshData.Indices, gl.STATIC_DRAW)
	}
	v.updates = v.mesh.Updates
	return true
}

// should bind for VAO normally?
func (v *vbo) bindForShader(shader *Shader) {
	g := v.g

	g.BindBuffer(gl.ARRAY_BUFFER, v.VBO)

	vsz := vertSize(v.Format)
	// Update VBO Attribute in VAO
	switch v.Format {
	// We will pass 3 vertex
	case gorge.VertexFormatP:
		if a, ok := shader.Attrib("aPosition"); ok {
			g.EnableVertexAttribArray(a)
			g.VertexAttribPointer(a, 3, gl.FLOAT, false, 0, 0)
		}
	// We will pass 3 vertex 3 normal
	case gorge.VertexFormatPN: // Interpolated
		if a, ok := shader.Attrib("aPosition"); ok {
			g.EnableVertexAttribArray(a)
			g.VertexAttribPointer(a, 3, gl.FLOAT, false, vsz*4, 0)
		}
		if a, ok := shader.Attrib("aNormal"); ok {
			g.EnableVertexAttribArray(a)
			g.VertexAttribPointer(a, 3, gl.FLOAT, false, vsz*4, 3*4)
		}
	// We will pass 3 vertex 3 normal
	case gorge.VertexFormatPT: // Interpolated
		if a, ok := shader.Attrib("aPosition"); ok {
			g.EnableVertexAttribArray(a)
			g.VertexAttribPointer(a, 3, gl.FLOAT, false, vsz*4, 0)
		}
		if a, ok := shader.Attrib("aTexCoords"); ok {
			g.EnableVertexAttribArray(a)
			g.VertexAttribPointer(a, 2, gl.FLOAT, false, vsz*4, 3*4)
		}
	// We will pass 3 vertex 2 uv 3 normal
	case gorge.VertexFormatPTN:
		if a, ok := shader.Attrib("aPosition"); ok {
			g.EnableVertexAttribArray(a)
			g.VertexAttribPointer(a, 3, gl.FLOAT, false, vsz*4, 0)
		}
		if a, ok := shader.Attrib("aTexCoords"); ok {
			g.EnableVertexAttribArray(a)
			g.VertexAttribPointer(a, 2, gl.FLOAT, false, vsz*4, 3*4)
		}
		if a, ok := shader.Attrib("aNormal"); ok {
			g.EnableVertexAttribArray(a)
			g.VertexAttribPointer(a, 3, gl.FLOAT, false, vsz*4, 5*4)
		}
	}
	v.g.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, v.EBO)
}

func vertSize(f gorge.VertexFormat) int {
	switch f {
	case gorge.VertexFormatP:
		return 3
	case gorge.VertexFormatPN:
		return 3 + 3
	case gorge.VertexFormatPT:
		return 3 + 2
	case gorge.VertexFormatPTN:
		return 3 + 2 + 3
	}
	return -1
}
