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

	// For debug
	vboCount int
	vbos     map[interface{}]*vbo
}

func newVBOManager(g gl.Context3) *vboManager {
	vm := &vboManager{
		g:    g,
		vbos: map[interface{}]*vbo{},
	}
	return vm
}

func (vm *vboManager) Load(m *gorge.Mesh) *vbo {
	k := m.Loader()
	if vb, ok := vm.vbos[k]; ok {
		vb.refCount++
		return vb
	}

	vb := &vbo{
		g:       vm.g,
		loader:  m.Loader(),
		VBO:     vm.g.CreateBuffer(),
		EBO:     vm.g.CreateBuffer(),
		updates: -1,
	}
	vb.update(true)

	vm.vbos[k] = vb
	return vb
}

// Get vbo for mesh
func (vm *vboManager) Get(m *gorge.Mesh) *vbo {
	k := m.Loader()
	if vb, ok := vm.vbos[k]; ok {
		return vb
	}
	// Not found, but if it is a MeshEntity we allow it to load since it will not
	// use IO and should be already tracked in scene
	if _, ok := m.Loader().(*gorge.MeshData); ok {
		return vm.Load(m)
	}
	return nil
}

type vbo struct {
	g           gl.Context3
	Format      gorge.VertexFormat
	VBO         gl.Buffer
	EBO         gl.Buffer
	ElementsLen int
	VertexLen   int

	// Do we need this?
	//mesh     *gorge.Mesh
	// Should be a loader as mesh doesn't matter anymore here
	loader   gorge.MeshLoader
	updates  int
	refCount int
}

// Load data if
// 1. State is Loading
// 2. data IsMeshData && updates arent equal
func (v *vbo) update(initial bool) bool {
	data, isMeshData := v.loader.(*gorge.MeshData)

	if !initial {
		if !isMeshData {
			return false
		}
		if v.updates == data.Updates {
			return false
		}
	}

	bufferType := uint32(gl.STATIC_DRAW)
	if isMeshData {
		bufferType = gl.DYNAMIC_DRAW
		v.updates = data.Updates
	}

	meshData := v.loader.Data()
	g := v.g

	if meshData == nil {
		return false
	}
	v.Format = meshData.Format
	v.VertexLen = len(meshData.Vertices) / vertSize(meshData.Format)

	g.BindBuffer(gl.ARRAY_BUFFER, v.VBO)
	g.BufferDataX(gl.ARRAY_BUFFER, meshData.Vertices, bufferType)

	if meshData.Indices != nil && len(meshData.Indices) > 0 {
		v.ElementsLen = len(meshData.Indices)
		g.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, v.EBO)
		g.BufferDataX(gl.ELEMENT_ARRAY_BUFFER, meshData.Indices, bufferType)
	}
	return true
}

// should bind for VAO normally?
func (v *vbo) bindForShader(shader *shader) {
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
