package gorge

import (
	"fmt"
)

type Blend struct {
	Disabled bool
	Src      BlendEnum
	Dst      BlendEnum
	Eq       BlendEq
}

func (b Blend) String() string {
	return fmt.Sprintf("Blend(Disabled: %v, Src: %v, Dst: %v, Eq: %v)",
		b.Disabled,
		b.Src,
		b.Dst,
		b.Eq,
	)
}

var (
	// BlendOneOneMinusSrcAlpha - gl.ONE, gl.ONE_MINUS_SRC_ALPHA
	BlendOneOneMinusSrcAlpha = &Blend{
		Src: BlendOne,
		Dst: BlendOneMinusSrcAlpha,
		Eq:  BlendEqAdd,
	}
	// BlendOneOne - gl.ONE, gl.ONE
	BlendOneOne = &Blend{
		Src: BlendOne,
		Dst: BlendOne,
		Eq:  BlendEqAdd,
	}
	// BlendDisable - disable blending
	BlendDisable = &Blend{Disabled: true}
)

type BlendEnum int

const (
	BlendZero = BlendEnum(iota)
	BlendOne
	BlendSrcColor
	BlendOneMinusSrcColor
	BlendDstColor
	BlendOneMinusDstColor
	BlendSrcAlpha
	BlendOneMinusSrcAlpha
	BlendDstAlpha
	BlendOneMinusDstAlpha
	BlendConstantColor
	BlendOneMinusConstantColor
	BlendConstantAlpha
	BlendOneMinusConstantAlpha
)

func (b BlendEnum) String() string {
	switch b {
	case BlendZero:
		return "BlendZero"
	case BlendOne:
		return "BlendOne"
	case BlendSrcColor:
		return "BlendSrcColor"
	case BlendOneMinusSrcColor:
		return "BlendOneMinusSrcColor"
	case BlendDstColor:
		return "BlendDstColor"
	case BlendOneMinusDstColor:
		return "BlendOneMinusDstColor"
	case BlendSrcAlpha:
		return "BlendSrcAlpha"
	case BlendOneMinusSrcAlpha:
		return "BlendOneMinusSrcAlpha"
	case BlendDstAlpha:
		return "BlendDstAlpha"
	case BlendOneMinusDstAlpha:
		return "BlendOneMinusDstAlpha"
	case BlendConstantColor:
		return "BlendConstantColor"
	case BlendOneMinusConstantColor:
		return "BlendOneMinusConstantColor"
	case BlendConstantAlpha:
		return "BlendConstantAlpha"
	}
	return fmt.Sprintf("BlendUnknown(%d)", b)
}

type BlendEq int

const (
	BlendEqAdd = BlendEq(iota)
	BlendEqSub
	BlendEqRevSub
	BlendEqMin
	BlendEqMax
)

func (e BlendEq) String() string {
	switch e {
	case BlendEqAdd:
		return "BlendEqAdd"
	case BlendEqSub:
		return "BlendEqSubtract"
	case BlendEqRevSub:
		return "BlendEqReverseSubtract"
	case BlendEqMin:
		return "BlendEqMin"
	case BlendEqMax:
		return "BlendEqMax"
	}
	return fmt.Sprintf("BlendEqUnknown(%d)", e)
}

/*
C¯result=C¯source∗Fsource+C¯destination∗Fdestination(1)

C¯source: the source color vector. This is the color output of the fragment shader.
C¯destination: the destination color vector. This is the color vector that is currently stored in the color buffer.
Fsource: the source factor value. Sets the impact of the alpha value on the source color.
Fdestination: the destination factor value. Sets the impact of the alpha value on the destination color.


GL_ZERO	Factor is equal to 0.
GL_ONE	Factor is equal to 1.
GL_SRC_COLOR	Factor is equal to the source color vector C¯source.
GL_ONE_MINUS_SRC_COLOR	Factor is equal to 1 minus the source color vector: 1−C¯source.
GL_DST_COLOR	Factor is equal to the destination color vector C¯destination
GL_ONE_MINUS_DST_COLOR	Factor is equal to 1 minus the destination color vector: 1−C¯destination.
GL_SRC_ALPHA	Factor is equal to the alpha component of the source color vector C¯source.
GL_ONE_MINUS_SRC_ALPHA	Factor is equal to 1−alpha of the source color vector C¯source.
GL_DST_ALPHA	Factor is equal to the alpha component of the destination color vector C¯destination.
GL_ONE_MINUS_DST_ALPHA	Factor is equal to 1−alpha of the destination color vector C¯destination.
GL_CONSTANT_COLOR	Factor is equal to the constant color vector C¯constant.
GL_ONE_MINUS_CONSTANT_COLOR	Factor is equal to 1 - the constant color vector C¯constant.
GL_CONSTANT_ALPHA	Factor is equal to the alpha component of the constant color vector C¯constant.
GL_ONE_MINUS_CONSTANT_ALPHA	Factor is equal to 1−alpha of the constant color vector C¯constant.

*/
