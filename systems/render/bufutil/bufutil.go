// Package bufutil provides means to cache and write into underlying gpu buffers
package bufutil

import (
	"fmt"
	"unsafe"

	"github.com/stdiopt/gorge/systems/render/gl"
)

type number interface {
	~float32 | ~float64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

const maxaddr = 0x7FFFFFFF

func AsBytes(data any) []byte {
	switch v := data.(type) {
	case []float32:
		if len(v) == 0 {
			return nil
		}
		bsz := len(v) * 4
		return (*(*[maxaddr]byte)(unsafe.Pointer(&v[0])))[:bsz:bsz]
	case []uint32:
		if len(v) == 0 {
			return nil
		}
		bsz := len(v) * 4
		return (*(*[maxaddr]byte)(unsafe.Pointer(&v[0])))[:bsz:bsz]
	case []uint16:
		if len(v) == 0 {
			return nil
		}
		bsz := len(v) * 2
		return (*(*[maxaddr]byte)(unsafe.Pointer(&v[0])))[:bsz:bsz]
	case []byte:
		return v
	}
	panic(fmt.Errorf("type :%T not implemented", data))
}

// Future
/*
type number interface {
	constraints.Float | constraints.Integer
}
func AsBytes[T number](data []T) []byte {
	var z T
	bsz := int(unsafe.Sizeof(z)) * len(data)
	return (*(*[maxaddr]byte)(unsafe.Pointer(&v[0])))[:bsz:bsz]
}
*/

type buffer interface {
	ID() gl.Buffer
	Init(sz int)
	Destroy()
	Bind()
	WriteAt(data []byte, offs int)
	Flush()
	Size() int
}
