package gl

import (
	"unsafe"
)

// Make sure wrapper is a context3
var _ Context3 = &Wrapper{}

//nolint:deadcode,unused
func f32bytes(v ...float32) []byte {
	bsz := len(v) * 4
	return (*(*[^uint32(0)]byte)(unsafe.Pointer(&v[0])))[:bsz:bsz]
}

//nolint:deadcode,unused
func u32bytes(v ...uint32) []byte {
	bsz := len(v) * 4
	return (*(*[^uint32(0)]byte)(unsafe.Pointer(&v[0])))[:bsz:bsz]
}

//nolint:deadcode,unused
func u16bytes(v ...uint16) []byte {
	bsz := len(v) * 2
	return (*(*[^uint32(0)]byte)(unsafe.Pointer(&v[0])))[:bsz:bsz]
}
