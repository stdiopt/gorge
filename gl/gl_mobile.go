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

// +build android

package gl

import (
	"fmt"
	"reflect"
	"unsafe"

	"golang.org/x/mobile/gl"
)

// Types that others implements
type (
	Uint         = uint32
	Buffer       = gl.Buffer
	Shader       = gl.Shader
	Program      = gl.Program
	Attrib       = gl.Attrib
	Framebuffer  = gl.Framebuffer
	Renderbuffer = gl.Renderbuffer
	Texture      = gl.Texture
	VertexArray  = gl.VertexArray
	Uniform      = gl.Uniform
	Enum         = gl.Enum
)

// Wrapper for mobile gl.Context
type Wrapper struct {
	gl.Context3
}

// BufferDataX upload any slice type
func (g Wrapper) BufferDataX(target Enum, d interface{}, usage Enum) {

	switch v := d.(type) {
	case []float32:
		g.BufferData(target, F32Bytes(v...), usage)
	case []uint32:
		g.BufferData(target, UI32Bytes(v...), usage)
	default:
		panic(fmt.Sprintf("Buffer type not implemented: %T", d))
	}
}

// F32Bytes unsafe cast list of floats to byte
func F32Bytes(values ...float32) []byte {
	// size in bytes
	f32size := 4
	// Get the slice header
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&values))
	header.Len *= f32size
	header.Cap *= f32size

	// Convert slice header to []byte
	data := *(*[]byte)(unsafe.Pointer(&header))
	return data
}

// UI32Bytes unsafe cast list of floats to byte
func UI32Bytes(values ...uint32) []byte {
	// size in bytes
	i32size := 4
	// Get the slice header
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&values))
	header.Len *= i32size
	header.Cap *= i32size

	// Convert slice header to []byte
	data := *(*[]byte)(unsafe.Pointer(&header))
	return data
}
