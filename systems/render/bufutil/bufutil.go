// Package bufutil provides means to cache and write into underlying gpu buffers
package bufutil

import (
	"fmt"
	"unsafe"

	"github.com/stdiopt/gorge/systems/render/gl"
)

func asBytes(data interface{}) []byte {
	const max = ^uint32(0)
	switch v := data.(type) {
	case []float32:
		bsz := len(v) * 4
		return (*(*[max]byte)(unsafe.Pointer(&v[0])))[:bsz:bsz]
	case []uint32:
		bsz := len(v) * 4
		return (*(*[max]byte)(unsafe.Pointer(&v[0])))[:bsz:bsz]
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
