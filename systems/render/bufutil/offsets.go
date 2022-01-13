package bufutil

import (
	"fmt"

	"github.com/stdiopt/gorge/m32"
)

// OffsetSpec specification for buffer
type OffsetSpec map[string]int

// NamedOffset buffer that has named offsets.
type NamedOffset struct {
	buffer
	offsets OffsetSpec
}

// NewNamedOffset returns a NamedOffset buffer.
func NewNamedOffset(b buffer, sz int, spec OffsetSpec) *NamedOffset {
	b.Init(sz)
	return &NamedOffset{
		buffer:  b,
		offsets: spec,
	}
}

// WriteOffset write something at a named offset
func (b *NamedOffset) WriteOffset(name string, v any) {
	// Figure out data from v
	offs, ok := b.offsets[name]
	if !ok {
		// Fix panic, just return instead
		panic(fmt.Errorf("offset name: %q not found", name))
	}

	// move this to a common write func
	// var blen int
	var data []byte
	switch v := v.(type) {
	case m32.Mat4:
		data = AsBytes(v[:]) // []float32
	case m32.Vec3:
		data = AsBytes(v[:]) // []float32
	case bool:
		if !v {
			data = []byte{0, 0, 0, 0}
		} else {
			data = []byte{0, 0, 0, 1}
		}
	case int:
		data = []byte{ // int32
			byte(v),
			byte(v >> 8),
			byte(v >> 16),
			byte(v >> 24),
		}
	case int32:
		data = []byte{ // int32
			byte(v),
			byte(v >> 8),
			byte(v >> 16),
			byte(v >> 24),
		}
	case float32: // just pass a single float32
		data = AsBytes([]float32{v})
	default:
		panic(fmt.Errorf("unhandled type: %T", v))
	}

	b.WriteAt(data, offs)
}
