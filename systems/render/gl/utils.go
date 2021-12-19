package gl

import (
	"unsafe"
)

// Make sure wrapper is a context3
var _ Context3 = &Wrapper{}

const maxaddr = 0x7FFFFFFF

//nolint:deadcode,unused
func f32bytes(v ...float32) []byte {
	bsz := len(v) * 4
	return (*(*[maxaddr]byte)(unsafe.Pointer(&v[0])))[:bsz:bsz]
}

//nolint:deadcode,unused
func u32bytes(v ...uint32) []byte {
	bsz := len(v) * 4
	return (*(*[maxaddr]byte)(unsafe.Pointer(&v[0])))[:bsz:bsz]
}

//nolint:deadcode,unused
func u16bytes(v ...uint16) []byte {
	bsz := len(v) * 2
	return (*(*[maxaddr]byte)(unsafe.Pointer(&v[0])))[:bsz:bsz]
}
