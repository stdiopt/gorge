// Package bufutil provides means to cache and write into underlying gpu buffers
package bufutil

import (
	"fmt"
	"unsafe"

	"github.com/stdiopt/gorge/systems/render/gl"
)

const maxaddr = 0x7FFFFFFF

func asBytes(data interface{}) []byte {
	switch v := data.(type) {
	case []float32:
		bsz := len(v) * 4
		return (*(*[maxaddr]byte)(unsafe.Pointer(&v[0])))[:bsz:bsz]
	case []uint32:
		bsz := len(v) * 4
		return (*(*[maxaddr]byte)(unsafe.Pointer(&v[0])))[:bsz:bsz]
	case []byte:
		return v
	}
	panic(fmt.Errorf("type :%T not implemented", data))
}

type buffer interface {
	ID() gl.Buffer
	Init(sz int)
	Destroy()
	Bind()
	WriteAt(data interface{}, offs int)
	Flush()
	Size() int
}
