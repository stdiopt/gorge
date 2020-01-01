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

// +build !js

package renderer

import (
	"reflect"
	"unsafe"
)

// to avoid editor pain

const f32size = 4

// TODO: Space for improvements, need a Redim() method to resize the float
// buffer
// TODO: wasm version is almost similar to this maybe create
// a common base overriding Redim() and Get() which are the funcs currently
// having any js code

//F32TransferBuf holds a buffer and a js reference
type F32TransferBuf struct {
	buffer []float32
}

//NewF32TransferBuf creates a go2js transfer buf
func NewF32TransferBuf(sz int) *F32TransferBuf {
	buffer := make([]float32, sz)
	return &F32TransferBuf{
		buffer,
	}
}

// WriteAt write floats at offset in buffer
func (b *F32TransferBuf) WriteAt(floats []float32, off int) {
	copy(b.buffer[off:], floats)
}

// Get copy bytes to js and return the reference
func (b *F32TransferBuf) Get() interface{} {
	return b.buffer
}

// F32Bytes Cast slice of floats to slice of bytes
func F32Bytes(values ...float32) []byte {
	// Get the slice header
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&values))
	header.Len *= f32size
	header.Cap *= f32size

	// Convert slice header to []byte
	data := *(*[]byte)(unsafe.Pointer(&header))
	return data
}
