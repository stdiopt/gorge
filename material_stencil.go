package gorge

import "fmt"

type Stencil struct {
	// Func
	Ref      int
	Func     StencilFunc
	ReadMask uint32
	// Mask
	WriteMask uint32
	// Op
	Fail  StencilOp
	ZFail StencilOp
	ZPass StencilOp
}

func (s *Stencil) String() string {
	return fmt.Sprintf(
		"Func: %v, Id: %d, WriteMask: %d, ReadMask: %d, Op: %v %v %v",
		s.Func, s.Ref,
		s.WriteMask,
		s.ReadMask,
		s.Fail,
		s.ZFail,
		s.ZPass,
	)
}

type StencilFunc int

const (
	StencilFuncNever = StencilFunc(iota)
	StencilFuncLess
	StencilFuncLequal
	StencilFuncGreater
	StencilFuncGequal
	StencilFuncEqual
	StencilFuncNotequal
	StencilFuncAlways
)

func (f StencilFunc) String() string {
	switch f {
	case StencilFuncNever:
		return "never"
	case StencilFuncLess:
		return "less"
	case StencilFuncLequal:
		return "lequal"
	case StencilFuncGreater:
		return "greater"
	case StencilFuncGequal:
		return "gequal"
	case StencilFuncEqual:
		return "equal"
	case StencilFuncNotequal:
		return "notequal"
	case StencilFuncAlways:
		return "always"
	}
	return "unknown"
}

type StencilOp int

const (
	StencilOpKeep = StencilOp(iota)
	StencilOpZero
	StencilOpReplace
	StencilOpIncr
	StencilOpDecr
	StencilOpInvert
	StencilOpIncrWrap
	StencilOpDecrWrap
)

func (o StencilOp) String() string {
	switch o {
	case StencilOpKeep:
		return "keep"
	case StencilOpZero:
		return "zero"
	case StencilOpReplace:
		return "replace"
	case StencilOpIncr:
		return "incr"
	case StencilOpDecr:
		return "decr"
	case StencilOpInvert:
		return "invert"
	case StencilOpIncrWrap:
		return "incr_wrap"
	case StencilOpDecrWrap:
		return "decr_wrap"
	}
	return "unknown"
}
