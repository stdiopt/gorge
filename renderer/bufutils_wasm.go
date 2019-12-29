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

// +build js,wasm

// TODO: Implement growing on these arrays

package renderer

import (
	"reflect"
	"syscall/js"
	"unsafe"
)

const f32size = 4

//F32TransferBuf holds a buffer and a js reference
type F32TransferBuf struct {
	// Buffer is the local go buffer that we write
	buffer  []float32
	jsBytes js.Value
}

//NewF32TransferBuf creates a go2js transfer buf
func NewF32TransferBuf(sz int) *F32TransferBuf {
	buffer := make([]float32, sz)
	jsBytes := js.Global().Get("Uint8Array").New(sz * f32size)
	return &F32TransferBuf{
		buffer:  buffer,
		jsBytes: jsBytes,
	}
}

// WriteAt write floats at offset in buffer
func (b *F32TransferBuf) WriteAt(floats []float32, off int) {
	copy(b.buffer[off:], floats)
}

// Get copy bytes to js and return the reference
func (b *F32TransferBuf) Get() interface{} {
	js.CopyBytesToJS(b.jsBytes, F32Bytes(b.buffer...))
	return b.jsBytes
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
