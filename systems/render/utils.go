package render

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/render/gl"
)

type (
	HLight      struct{ Light }
	HCamera     struct{ Camera }
	HRenderable struct{ Renderable }
)

// Light interface for light composition
type Light interface {
	Mat4() gm.Mat4
	Light() *gorge.LightComponent
}

// Camera interface for accepting camera structs
type Camera interface {
	Mat4() gm.Mat4
	Camera() *gorge.CameraComponent
}

// Renderable the renderer renderable component interface.
type Renderable interface {
	Mat4() gm.Mat4
	Renderable() *gorge.RenderableComponent
}

type cameraSorter []Camera

// Len is the number of elements in the collection.
func (c cameraSorter) Len() int { return len(c) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (c cameraSorter) Less(i int, j int) bool {
	return c[i].Camera().Order < c[j].Camera().Order
}

// Swap swaps the elements with indexes i and j.
func (c cameraSorter) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

// DrawMode converts gorge DrawMode into gl
// POINTS                                       = 0x0000
// LINES                                        = 0x0001
// LINE_LOOP                                    = 0x0002
// LINE_STRIP                                   = 0x0003
// TRIANGLES                                    = 0x0004
// TRIANGLE_STRIP                               = 0x0005
// TRIANGLE_FAN                                 = 0x0006
func DrawMode(d gorge.DrawMode) gl.Enum {
	switch d {
	case gorge.DrawPoints:
		return gl.POINTS
	case gorge.DrawLines:
		return gl.LINES
	case gorge.DrawLineLoop:
		return gl.LINE_LOOP
	case gorge.DrawLineStrip:
		return gl.LINE_STRIP
	case gorge.DrawTriangles:
		return gl.TRIANGLES
	case gorge.DrawTriangleStrip:
		return gl.TRIANGLE_STRIP
	case gorge.DrawTriangleFan:
		return gl.TRIANGLE_FAN
	}
	panic("unknown drawtype")
}

// TextureWrap converts gorge textureWrap to gl.
func TextureWrap(n gorge.TextureWrap) int {
	switch n {
	case gorge.TextureWrapClamp:
		return gl.CLAMP_TO_EDGE
	case gorge.TextureWrapRepeat:
		return gl.REPEAT
	case gorge.TextureWrapMirror:
		return gl.MIRRORED_REPEAT
	}
	return gl.REPEAT
}

// TextureFormat returns the internal format and format enum for
// gl.TexImage2D
func TextureFormat(n gorge.TextureFormat) (int, gl.Enum, gl.Enum) {
	switch n {
	case gorge.TextureFormatGray:
		return gl.R8, gl.RED, gl.UNSIGNED_BYTE
	case gorge.TextureFormatGray16:
		return gl.R16UI, gl.RED_INTEGER, gl.UNSIGNED_SHORT
	case gorge.TextureFormatRGB:
		return gl.RGB, gl.RGB, gl.UNSIGNED_BYTE
	case gorge.TextureFormatRGBA:
		return gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE
	}
	// default
	return gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE
}

// CullMask returns a bit CullMask if it's 0 it will return the default mask 0xFF
func CullMask(n gorge.CullMaskFlags) gorge.CullMaskFlags {
	if n == 0 {
		return gorge.CullMaskFlags(0xFF)
	}
	return n
}

// StencilFunc converts gorge stencil to opengl enum.
func StencilFunc(n gorge.StencilFunc) gl.Enum {
	switch n {
	case gorge.StencilFuncNever:
		return gl.NEVER
	case gorge.StencilFuncLess:
		return gl.LESS
	case gorge.StencilFuncLequal:
		return gl.LEQUAL
	case gorge.StencilFuncGreater:
		return gl.GREATER
	case gorge.StencilFuncGequal:
		return gl.GEQUAL
	case gorge.StencilFuncEqual:
		return gl.EQUAL
	case gorge.StencilFuncNotequal:
		return gl.NOTEQUAL
	case gorge.StencilFuncAlways:
		return gl.ALWAYS
	}
	// default
	return gl.ALWAYS
}

// StencilOp converts gorge stencil to opengl.
func StencilOp(n gorge.StencilOp) gl.Enum {
	switch n {
	case gorge.StencilOpKeep:
		return gl.KEEP
	case gorge.StencilOpZero:
		return gl.ZERO
	case gorge.StencilOpReplace:
		return gl.REPLACE
	case gorge.StencilOpIncr:
		return gl.INCR
	case gorge.StencilOpDecr:
		return gl.DECR
	case gorge.StencilOpIncrWrap:
		return gl.INCR_WRAP
	case gorge.StencilOpDecrWrap:
		return gl.DECR_WRAP
	case gorge.StencilOpInvert:
		return gl.INVERT
	}
	return gl.KEEP
}

func BlendEnum(n gorge.BlendEnum) gl.Enum {
	switch n {
	case gorge.BlendZero:
		return gl.ZERO
	case gorge.BlendOne:
		return gl.ONE
	case gorge.BlendSrcColor:
		return gl.SRC_COLOR
	case gorge.BlendOneMinusSrcColor:
		return gl.ONE_MINUS_SRC_COLOR
	case gorge.BlendDstColor:
		return gl.DST_COLOR
	case gorge.BlendOneMinusDstColor:
		return gl.ONE_MINUS_DST_COLOR
	case gorge.BlendSrcAlpha:
		return gl.SRC_ALPHA
	case gorge.BlendOneMinusSrcAlpha:
		return gl.ONE_MINUS_SRC_ALPHA
	case gorge.BlendDstAlpha:
		return gl.DST_ALPHA
	case gorge.BlendOneMinusDstAlpha:
		return gl.ONE_MINUS_DST_ALPHA
	case gorge.BlendConstantColor:
		return gl.CONSTANT_COLOR
	case gorge.BlendOneMinusConstantColor:
		return gl.ONE_MINUS_CONSTANT_COLOR
	case gorge.BlendConstantAlpha:
		return gl.CONSTANT_ALPHA
	case gorge.BlendOneMinusConstantAlpha:
		return gl.ONE_MINUS_CONSTANT_ALPHA
	default:
		return gl.ONE
	}
}

func BlendEq(n gorge.BlendEq) gl.Enum {
	switch n {
	case gorge.BlendEqAdd:
		return gl.FUNC_ADD
	case gorge.BlendEqSub:
		return gl.FUNC_SUBTRACT
	case gorge.BlendEqRevSub:
		return gl.FUNC_REVERSE_SUBTRACT
	default:
		return gl.FUNC_ADD
	}
}

// From glTf-Sample-Viewer
func stringHash(str string) uint {
	seed := uint(0)
	for _, c := range str {
		seed = 31*seed + uint(c)
	}
	return seed
}
