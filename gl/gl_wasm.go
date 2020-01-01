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

// +build wasm, js

package gl

// Wasm implementation not everything is implemented as I've been implementing
// on demand

import (
	"fmt"
	"reflect"
	"syscall/js"
	"unsafe"
)

var (
	// NullTexture Nil texture just because
	NullTexture = js.Null()
)

type (
	Buffer       = js.Value
	Shader       = js.Value
	Program      = js.Value
	Attrib       = uint32
	Framebuffer  = js.Value
	Renderbuffer = js.Value
	Texture      = js.Value
	VertexArray  = js.Value
	Uniform      = js.Value
	Enum         = uint32
)

// Common to move things to gl
var (
	b8buf  = js.Global().Get("Uint8Array").New(8)
	b12buf = js.Global().Get("Uint8Array").New(12)
	b16buf = js.Global().Get("Uint8Array").New(16)
	b64buf = js.Global().Get("Uint8Array").New(64)

	// Float buffers of the byte ones
	f2buf  = js.Global().Get("Float32Array").New(b8buf.Get("buffer"))
	f3buf  = js.Global().Get("Float32Array").New(b12buf.Get("buffer"))
	f4buf  = js.Global().Get("Float32Array").New(b16buf.Get("buffer"))
	f16buf = js.Global().Get("Float32Array").New(b64buf.Get("buffer"))
)

// Wrapper exposes the methods
type Wrapper struct {
	js.Value
}

// GetWebGL Return a js.Value Wrapper context
func GetWebGL(v js.Value) Wrapper {
	return Wrapper{v}
}

// ActiveTexture sets the active texture unit.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glActiveTexture.xhtml
func (c Wrapper) ActiveTexture(texture Enum) {
	c.Call("activeTexture", texture)
}

// AttachShader attaches a shader to a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glAttachShader.xhtml
func (c Wrapper) AttachShader(p Program, s Shader) {
	c.Call("attachShader", p, s)
}

// BindAttribLocation binds a vertex attribute index with a named
// variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindAttribLocation.xhtml
func (c Wrapper) BindAttribLocation(p Program, a Attrib, name string) {
	c.Call("bindAttribLocation", p, a, name)
}

// BindBuffer binds a buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBuffer.xhtml
func (c Wrapper) BindBuffer(target Enum, b Buffer) {
	c.Call("bindBuffer", target, b)
}

// BindFramebuffer binds a framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindFramebuffer.xhtml
func (c Wrapper) BindFramebuffer(target Enum, fb Framebuffer) {
	c.Call("bindFramebuffer", target, fb)
}

// BindRenderbuffer binds a render buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindRenderbuffer.xhtml
func (c Wrapper) BindRenderbuffer(target Enum, rb Renderbuffer) {
	c.Call("bindRenderbuffer", target, rb)
}

// BindTexture binds a texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindTexture.xhtml
func (c Wrapper) BindTexture(target Enum, t Texture) {
	c.Call("bindTexture", target, t)
}

// BindVertexArray binds a vertex array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindVertexArray.xhtml
func (c Wrapper) BindVertexArray(rb VertexArray) {
	c.Call("bindVertexArray", rb)
}

// BlendColor sets the blend color.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendColor.xhtml
func (c Wrapper) BlendColor(red, green, blue, alpha float32) {
	c.Call("blendColor", red, green, blue, alpha)
}

// BlendEquation sets both RGB and alpha blend equations.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendEquation.xhtml
func (c Wrapper) BlendEquation(mode Enum) {
	c.Call("blendEquation", mode)
}

// BlendEquationSeparate sets RGB and alpha blend equations separately.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendEquationSeparate.xhtml
func (c Wrapper) BlendEquationSeparate(modeRGB, modeAlpha Enum) {
	c.Call("blendEquationSeparate", modeRGB, modeAlpha)
}

// BlendFunc sets the pixel blending factors.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendFunc.xhtml
func (c Wrapper) BlendFunc(sfactor, dfactor Enum) {
	c.Call("blendFunc", sfactor, dfactor)
}

func (c Wrapper) BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha Enum) {
	panic("not implemented")
}

// BufferData creates a new data store for the bound buffer object.
// XXX: Can be pooled
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferData.xhtml
func (c Wrapper) BufferData(target Enum, data []byte, usage Enum) {
	d := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(d, data)
	c.Call("bufferData", target, d, usage)
}

func (c Wrapper) BufferInit(target Enum, size int, usage Enum) {
	panic("not implemented")
}

// BufferSubData sets some of data in the bound buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferSubData.xhtml
func (c Wrapper) BufferSubData(target Enum, offset int, data []byte) {
	d := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(d, data)
	c.Call("bufferSubData", target, offset, d)
}

func (c Wrapper) CheckFramebufferStatus(target Enum) Enum {
	panic("not implemented")
}

// Clear clears the window.
//
// The behavior of Clear is influenced by the pixel ownership test,
// the scissor test, dithering, and the buffer writemasks.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClear.xhtml
func (c Wrapper) Clear(mask Enum) {
	c.Call("clear", mask)
}

// ClearColor specifies the RGBA values used to clear color buffers.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearColor.xhtml
func (c Wrapper) ClearColor(red, green, blue, alpha float32) {
	c.Call("clearColor", red, green, blue, alpha)
}

// ClearDepthf sets the depth value used to clear the depth buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearDepthf.xhtml
func (c Wrapper) ClearDepthf(d float32) {
	c.Call("clearDepth", d)
}

func (c Wrapper) ClearStencil(s int) {
	panic("not implemented")
}

func (c Wrapper) ColorMask(red, green, blue, alpha bool) {
	c.Call("colorMask", red, green, blue, alpha)
}

// CompileShader compiles the source code of s.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompileShader.xhtml
func (c Wrapper) CompileShader(s Shader) {
	c.Call("compileShader", s)
}

func (c Wrapper) CompressedTexImage2D(target Enum, level int, internalformat Enum, width, height, border int, data []byte) {
	panic("not implemented")
}

func (c Wrapper) CompressedTexSubImage2D(target Enum, level, xoffset, yoffset, width, height int, format Enum, data []byte) {
	panic("not implemented")
}

func (c Wrapper) CopyTexImage2D(target Enum, level int, internalformat Enum, x, y, width, height, border int) {
	panic("not implemented")
}

func (c Wrapper) CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int) {
	panic("not implemented")
}

// CreateBuffer creates a buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenBuffers.xhtml
func (c Wrapper) CreateBuffer() Buffer {
	return Buffer(c.Call("createBuffer"))
}

func (c Wrapper) CreateFramebuffer() Framebuffer {
	panic("not implemented")
}

// CreateProgram creates a new empty program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateProgram.xhtml
func (c Wrapper) CreateProgram() Program {
	return Program(c.Call("createProgram"))
}

func (c Wrapper) CreateRenderbuffer() Renderbuffer {
	panic("not implemented")
}

// CreateShader creates a new empty shader object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateShader.xhtml
func (c Wrapper) CreateShader(ty Enum) Shader {
	return Shader(c.Call("createShader", ty))
}

// CreateTexture creates a texture object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenTextures.xhtml
func (c Wrapper) CreateTexture() Texture {
	return Texture(c.Call("createTexture"))
}

// CreateVertexArray creates a vertex array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenVertexArrays.xhtml
func (c Wrapper) CreateVertexArray() VertexArray {
	return VertexArray(c.Call("createVertexArray"))
}

// CullFace specifies which polygons are candidates for culling.
//
// Valid modes: FRONT, BACK, FRONT_AND_BACK.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCullFace.xhtml
func (c Wrapper) CullFace(mode Enum) {
	c.Call("cullFace", mode)
}

// DeleteBuffer deletes the given buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteBuffers.xhtml
func (c Wrapper) DeleteBuffer(v Buffer) {
	c.Call("deleteBuffer", v)
}

// DeleteFramebuffer deletes the given framebuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteFramebuffers.xhtml
func (c Wrapper) DeleteFramebuffer(v Framebuffer) {
	c.Call("deleteFramebuffer", v)
}

// DeleteProgram deletes the given program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteProgram.xhtml
func (c Wrapper) DeleteProgram(p Program) {
	c.Call("deleteProgram", p)
}

func (c Wrapper) DeleteRenderbuffer(v Renderbuffer) {
	panic("not implemented")
}

func (c Wrapper) DeleteShader(s Shader) {
	panic("not implemented")
}

func (c Wrapper) DeleteTexture(v Texture) {
	panic("not implemented")
}

func (c Wrapper) DeleteVertexArray(v VertexArray) {
	panic("not implemented")
}

// DepthFunc sets the function used for depth buffer comparisons.
//
// Valid fn values:
//	NEVER
//	LESS
//	EQUAL
//	LEQUAL
//	GREATER
//	NOTEQUAL
//	GEQUAL
//	ALWAYS
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDepthFunc.xhtml
func (c Wrapper) DepthFunc(fn Enum) {
	c.Call("depthFunc", fn)
}

// DepthMask sets the depth buffer enabled for writing.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDepthMask.xhtml
func (c Wrapper) DepthMask(flag bool) {
	c.Call("depthMask", flag)
}

func (c Wrapper) DepthRangef(n, f float32) {
	panic("not implemented")
}

func (c Wrapper) DetachShader(p Program, s Shader) {
	panic("not implemented")
}

// Disable disables various GL capabilities.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDisable.xhtml
func (c Wrapper) Disable(cap Enum) {
	c.Call("disable", cap)
}

// DisableVertexAttribArray disables a vertex attribute array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDisableVertexAttribArray.xhtml
func (c Wrapper) DisableVertexAttribArray(a Attrib) {
	c.Call("disableVertexAttribArray", a)
}

// DrawArrays renders geometric primitives from the bound data.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArrays.xhtml
func (c Wrapper) DrawArrays(mode Enum, first, count int) {
	c.Call("drawArrays", mode, first, count)
}

// DrawElements renders primitives from a bound buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawElements.xhtml
func (c Wrapper) DrawElements(mode Enum, count int, ty Enum, offset int) {
	c.Call("drawElements", mode, count, ty, offset)
}

// Enable enables various GL capabilities.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnable.xhtml
func (c Wrapper) Enable(cp Enum) {
	c.Call("enable", cp)
}

// EnableVertexAttribArray enables a vertex attribute array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnableVertexAttribArray.xhtml
func (c Wrapper) EnableVertexAttribArray(a Attrib) {
	c.Call("enableVertexAttribArray", a)
}

// Finish blocks until the effects of all previously called GL
// commands are complete.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFinish.xhtml
func (c Wrapper) Finish() {
	c.Call("finish")
}

// Flush empties all buffers. It does not block.
//
// An OpenGL implementation may buffer network communication,
// the command stream, or data inside the graphics accelerator.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFlush.xhtml
func (c Wrapper) Flush() {
	c.Call("flush")
}

func (c Wrapper) FramebufferRenderbuffer(target, attachment, rbTarget Enum, rb Renderbuffer) {
	panic("not implemented")
}

func (c Wrapper) FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	panic("not implemented")
}

// FrontFace defines which polygons are front-facing.
//
// Valid modes: CW, CCW.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFrontFace.xhtml
func (c Wrapper) FrontFace(mode Enum) {
	c.Call("frontFace", mode)
}

// GenerateMipmap generates mipmaps for the current texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenerateMipmap.xhtml
func (c Wrapper) GenerateMipmap(target Enum) {
	c.Call("generateMipmap", target)
}

// GetActiveAttrib returns details about an active attribute variable.
// A value of 0 for index selects the first active attribute variable.
// Permissible values for index range from 0 to the number of active
// attribute variables minus 1.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveAttrib.xhtml
func (c Wrapper) GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum) {
	res := c.Call("getActiveAttrib", p, index)

	name = res.Get("name").String()
	size = res.Get("size").Int()
	ty = Enum(res.Get("type").Int())
	return name, size, ty
}

// GetActiveUniform returns details about an active uniform variable.
// A value of 0 for index selects the first active uniform variable.
// Permissible values for index range from 0 to the number of active
// uniform variables minus 1.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveUniform.xhtml
func (c Wrapper) GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum) {
	res := c.Call("getActiveUniform", p, index)

	name = res.Get("name").String()
	size = res.Get("size").Int()
	ty = Enum(res.Get("type").Int())

	return name, size, ty
}

func (c Wrapper) GetAttachedShaders(p Program) []Shader {
	panic("not implemented")
}

// GetAttribLocation returns the location of an attribute variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetAttribLocation.xhtml
func (c Wrapper) GetAttribLocation(p Program, name string) Attrib {
	return Attrib(c.Call("getAttribLocation", p, name).Int())
}

func (c Wrapper) GetBooleanv(dst []bool, pname Enum) {
	panic("not implemented")
}

func (c Wrapper) GetFloatv(dst []float32, pname Enum) {
	panic("not implemented")
}

func (c Wrapper) GetIntegerv(dst []int32, pname Enum) {
	panic("not implemented")
}

func (c Wrapper) GetInteger(pname Enum) int {
	panic("not implemented")
}

func (c Wrapper) GetBufferParameteri(target, value Enum) int {
	panic("not implemented")
}

func (c Wrapper) GetError() Enum {
	panic("not implemented")
}

func (c Wrapper) GetFramebufferAttachmentParameteri(target, attachment, pname Enum) int {
	panic("not implemented")
}

// GetProgrami returns a parameter value for a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramiv.xhtml
func (c Wrapper) GetProgrami(p Program, pname Enum) int {
	r := c.Call("getProgramParameter", p, pname)

	switch r.Type() {
	case js.TypeBoolean:
		if r.Bool() {
			return 1
		}
		return 0
	case js.TypeNumber:
		return r.Int()
	default:
		panic("unknown type")
	}
}

// GetProgramInfoLog returns the information log for a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramInfoLog.xhtml
func (c Wrapper) GetProgramInfoLog(p Program) string {
	return c.Call("getProgramInfoLog", p).String()
}

func (c Wrapper) GetRenderbufferParameteri(target, pname Enum) int {
	panic("not implemented")
}

// GetShaderi returns a parameter value for a shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderiv.xhtml
func (c Wrapper) GetShaderi(s Shader, pname Enum) int {
	if c.Call("getShaderParameter", s, pname).Bool() {
		return 1
	}
	return 0
}

// GetShaderInfoLog returns the information log for a shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderInfoLog.xhtml
func (c Wrapper) GetShaderInfoLog(s Shader) string {
	return c.Call("getShaderInfoLog", s).String()
}

func (c Wrapper) GetShaderPrecisionFormat(shadertype, precisiontype Enum) (rangeLow, rangeHigh, precision int) {
	panic("not implemented")

}

func (c Wrapper) GetShaderSource(s Shader) string {
	panic("not implemented")
}

// GetString reports current GL state.
//
// Valid name values:
//	EXTENSIONS
//	RENDERER
//	SHADING_LANGUAGE_VERSION
//	VENDOR
//	VERSION
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetString.xhtml
func (c Wrapper) GetString(pname Enum) string {
	return c.Call("getParameter", pname).String()
}

func (c Wrapper) GetTexParameterfv(dst []float32, target, pname Enum) {
	panic("not implemented")
}

func (c Wrapper) GetTexParameteriv(dst []int32, target, pname Enum) {
	panic("not implemented")
}

// GetUniformfv returns the float values of a uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniform.xhtml
func (c Wrapper) GetUniformfv(dst []float32, src Uniform, p Program) {
	c.Call("getUniformfv", dst, src, p)
}

func (c Wrapper) GetUniformiv(dst []int32, src Uniform, p Program) {
	panic("not implemented")
}

func (c Wrapper) GetUniformLocation(p Program, name string) Uniform {
	return Uniform(c.Call("getUniformLocation", p, name))
}

func (c Wrapper) GetVertexAttribf(src Attrib, pname Enum) float32 {
	panic("not implemented")
}

func (c Wrapper) GetVertexAttribfv(dst []float32, src Attrib, pname Enum) {
	panic("not implemented")
}

func (c Wrapper) GetVertexAttribi(src Attrib, pname Enum) int32 {
	panic("not implemented")
}

func (c Wrapper) GetVertexAttribiv(dst []int32, src Attrib, pname Enum) {
	panic("not implemented")
}

// Hint sets implementation-specific modes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glHint.xhtml
func (c Wrapper) Hint(target, mode Enum) {
	c.Call("hint", target, mode)
}

func (c Wrapper) IsBuffer(b Buffer) bool {
	panic("not implemented")
}

func (c Wrapper) IsEnabled(cap Enum) bool {
	panic("not implemented")

}

func (c Wrapper) IsFramebuffer(fb Framebuffer) bool {
	panic("not implemented")
}

func (c Wrapper) IsProgram(p Program) bool {
	panic("not implemented")
}

func (c Wrapper) IsRenderbuffer(rb Renderbuffer) bool {
	panic("not implemented")
}

func (c Wrapper) IsShader(s Shader) bool {
	panic("not implemented")
}

func (c Wrapper) IsTexture(t Texture) bool {
	panic("not implemented")
}

// LineWidth specifies the width of lines.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glLineWidth.xhtml
func (c Wrapper) LineWidth(width float32) {
	c.Call("lineWidth", width)
}

// LinkProgram links the specified program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glLinkProgram.xhtml
func (c Wrapper) LinkProgram(p Program) {
	c.Call("linkProgram", p)
}

func (c Wrapper) PixelStorei(pname Enum, param int32) {
	panic("not implemented")
}

func (c Wrapper) PolygonOffset(factor, units float32) {
	panic("not implemented")
}

func (c Wrapper) ReadPixels(dst []byte, x, y, width, height int, format, ty Enum) {
	panic("not implemented")
}

func (c Wrapper) ReleaseShaderCompiler() {
	panic("not implemented")
}

func (c Wrapper) RenderbufferStorage(target, internalFormat Enum, width, height int) {
	panic("not implemented")
}

func (c Wrapper) SampleCoverage(value float32, invert bool) {
	panic("not implemented")
}

func (c Wrapper) Scissor(x, y, width, height int32) {
	panic("not implemented")
}

// ShaderSource sets the source code of s to the given source code.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glShaderSource.xhtml
func (c Wrapper) ShaderSource(s Shader, src string) {
	c.Call("shaderSource", s, src)
}

func (c Wrapper) StencilFunc(fn Enum, ref int, mask uint32) {
	panic("not implemented")
}

func (c Wrapper) StencilFuncSeparate(face, fn Enum, ref int, mask uint32) {
	panic("not implemented")
}

func (c Wrapper) StencilMask(mask uint32) {
	panic("not implemented")
}

func (c Wrapper) StencilMaskSeparate(face Enum, mask uint32) {
	panic("not implemented")
}

func (c Wrapper) StencilOp(fail, zfail, zpass Enum) {
	panic("not implemented")
}

func (c Wrapper) StencilOpSeparate(face, sfail, dpfail, dppass Enum) {
	panic("not implemented")
}

// TexImage2D writes a 2D texture image.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexImage2D.xhtml
func (c Wrapper) TexImage2D(target Enum, level int, internalFormat int, width, height int, format Enum, ty Enum, data []byte) {
	jsBuf := js.Null()
	if data != nil {
		// Might leak
		jsBuf = js.Global().Get("Uint8Array").New(len(data))
		js.CopyBytesToJS(jsBuf, data)
	}
	c.Call("texImage2D",
		target, level, internalFormat,
		width, height, 0,
		format, ty, jsBuf,
	)

}

func (c Wrapper) TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data []byte) {
	panic("not implemented")

}

// TexParameterf sets a float texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (c Wrapper) TexParameterf(target, pname Enum, param float32) {
	c.Call("texParameterf", target, pname, param)
}
func (c Wrapper) TexParameterfv(target, pname Enum, params []float32) {
	panic("not implemented")
}

// TexParameteri sets an integer texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (c Wrapper) TexParameteri(target, pname Enum, param int) {
	c.Call("texParameteri", target, pname, param)
}
func (c Wrapper) TexParameteriv(target, pname Enum, params []int32) {
	panic("not implemented")
}

// Uniform1f writes a float uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (c Wrapper) Uniform1f(dst Uniform, v float32) {
	c.Call("uniform1f", dst, v)
}
func (c Wrapper) Uniform1fv(dst Uniform, src []float32) {
	panic("not implemented")
}

// Uniform1i writes an int uniform variable.
//
// Uniform1i and Uniform1iv are the only two functions that may be used
// to load uniform variables defined as sampler types. Loading samplers
// with any other function will result in a INVALID_OPERATION error.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (c Wrapper) Uniform1i(dst Uniform, v int) {
	c.Call("uniform1i", dst, v)
}
func (c Wrapper) Uniform1iv(dst Uniform, src []int32) {
	panic("not implemented")
}

// Uniform2f writes a vec2 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (c Wrapper) Uniform2f(dst Uniform, v0, v1 float32) {
	c.Call("uniform2f", dst, v0, v1)
}

// Uniform2fv writes a vec2 uniform array of len(src)/2 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (c Wrapper) Uniform2fv(dst Uniform, src []float32) {
	js.CopyBytesToJS(b8buf, F32Bytes(src...))
	c.Call("uniform2fv", dst, f2buf)
}

func (c Wrapper) Uniform2i(dst Uniform, v0, v1 int) {
	panic("not implemented")
}

func (c Wrapper) Uniform2iv(dst Uniform, src []int32) {
	panic("not implemented")
}

// Uniform3f writes a vec3 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (c Wrapper) Uniform3f(dst Uniform, v0, v1, v2 float32) {
	c.Call("uniform3f", dst, v0, v1, v2)
}

// Uniform3fv writes a vec3 uniform array of len(src)/3 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (c Wrapper) Uniform3fv(dst Uniform, src []float32) {
	js.CopyBytesToJS(b12buf, F32Bytes(src...))
	c.Call("uniform3fv", dst, f3buf)
}

func (c Wrapper) Uniform3i(dst Uniform, v0, v1, v2 int32) {
	panic("not implemented")
}

func (c Wrapper) Uniform3iv(dst Uniform, src []int32) {
	panic("not implemented")
}

// Uniform4f writes a vec4 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (c Wrapper) Uniform4f(dst Uniform, v0, v1, v2, v3 float32) {
	c.Call("uniform4f", dst, v0, v1, v2, v3)
}

// Uniform4fv writes a vec4 uniform array of len(src)/4 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (c Wrapper) Uniform4fv(dst Uniform, src []float32) {
	js.CopyBytesToJS(b16buf, F32Bytes(src...))
	c.Call("uniform4fv", dst, f4buf)
}
func (c Wrapper) Uniform4i(dst Uniform, v0, v1, v2, v3 int32) {
	panic("not implemented")
}
func (c Wrapper) Uniform4iv(dst Uniform, src []int32) {
	panic("not implemented")
}
func (c Wrapper) UniformMatrix2fv(dst Uniform, src []float32) {
	panic("not implemented")
}

// UniformMatrix3fv writes 3x3 matrices. Each matrix uses nine
// float32 values, so the number of matrices written is len(src)/9.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (c Wrapper) UniformMatrix3fv(dst Uniform, src []float32) {
	js.CopyBytesToJS(b64buf, F32Bytes(src...))
	c.Call("uniformMatrix3fv", dst, false, f16buf)
}

// UniformMatrix4fv writes 4x4 matrices. Each matrix uses 16
// float32 values, so the number of matrices written is len(src)/16.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (c Wrapper) UniformMatrix4fv(dst Uniform, src []float32) {
	js.CopyBytesToJS(b64buf, F32Bytes(src...))
	c.Call("uniformMatrix4fv", dst, false, f16buf)
}

// UseProgram sets the active program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUseProgram.xhtml
func (c Wrapper) UseProgram(p Program) {
	c.Call("useProgram", p)
}
func (c Wrapper) ValidateProgram(p Program) {
	panic("not implemented")
}
func (c Wrapper) VertexAttrib1f(dst Attrib, x float32) {
	panic("not implemented")
}
func (c Wrapper) VertexAttrib1fv(dst Attrib, src []float32) {
	panic("not implemented")
}
func (c Wrapper) VertexAttrib2f(dst Attrib, x, y float32) {
	panic("not implemented")
}
func (c Wrapper) VertexAttrib2fv(dst Attrib, src []float32) {
	panic("not implemented")
}
func (c Wrapper) VertexAttrib3f(dst Attrib, x, y, z float32) {
	panic("not implemented")
}
func (c Wrapper) VertexAttrib3fv(dst Attrib, src []float32) {
	panic("not implemented")
}
func (c Wrapper) VertexAttrib4f(dst Attrib, x, y, z, w float32) {
	panic("not implemented")
}
func (c Wrapper) VertexAttrib4fv(dst Attrib, src []float32) {
	panic("not implemented")
}

// VertexAttribPointer uses a bound buffer to define vertex attribute data.
//
// Direct use of VertexAttribPointer to load data into OpenGL is not
// supported via the Go bindings. Instead, use BindBuffer with an
// ARRAY_BUFFER and then fill it usingBufferData.
//
// The size argument specifies the number of components per attribute,
// between 1-4. The stride argument specifies the byte offset between
// consecutive vertex attributes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttribPointer.xhtml
func (c Wrapper) VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride, offset int) {
	c.Call("vertexAttribPointer", dst, size, ty, normalized, stride, offset)
}

// Viewport sets the viewport, an affine transformation that
// normalizes device coordinates to window coordinates.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glViewport.xhtml
func (c Wrapper) Viewport(x, y, width, height int) {
	c.Call("viewport", x, y, width, height)
}

///////////////////////////////////////////////////////////////////////////////
// Wrapper 2 + extra funcs
///////////////////////////////////////////////////////////////////////////////

// GetUniformBlockIndex retrieves the index of a uniform block within program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockIndex.xhtml
func (c Wrapper) GetUniformBlockIndex(p Program, name string) int {
	return c.Call("getUniformBlockIndex", p, name).Int()
}

// UniformBlockBinding assign a binding point to an active uniform block
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockBinding.xhtml
func (c Wrapper) UniformBlockBinding(p Program, index, bind int) {
	c.Call("uniformBlockBinding", p, index, bind)
}

// BindBufferBase bind a buffer object to an indexed buffer target
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBufferBase.xhtml
func (c Wrapper) BindBufferBase(target Enum, n uint32, b Buffer) {
	c.Call("bindBufferBase", target, n, b)
}

// DrawArraysInstanced draw multiple instances of a range of elements
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
func (c Wrapper) DrawArraysInstanced(mode Enum, first int, count, primcount int) {
	c.Call("drawArraysInstanced", mode, first, count, primcount)
}

// VertexAttribDivisor  modify the rate at which generic vertex attributes
// advance during instanced rendering
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
func (c Wrapper) VertexAttribDivisor(index Attrib, divisor int) {
	c.Call("vertexAttribDivisor", index, divisor)
}

// DrawElementsInstanced — draw multiple instances of a set of elements
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawElementsInstanced.xhtml
func (c Wrapper) DrawElementsInstanced(mode Enum, count int, typ Enum, offset, primcount int) {
	c.Call("drawElementsInstanced", mode, count, typ, offset, primcount)
}

// BufferDataX will type switch the interface and select the proper type
// {lpf} Custom func
// TODO: Should move this Elsewhere as an helper
func (c Wrapper) BufferDataX(target Enum, d interface{}, usage Enum) {
	c.Call("bufferData", target, conv(d), usage)
}

// F32Bytes unsafe cast list of floats to byte
func F32Bytes(values ...float32) []byte {
	// size in bytes
	f32size := 4
	// Get the slice header
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&values))
	header.Len *= f32size
	header.Cap *= f32size

	// Convert slice header to []byte
	data := *(*[]byte)(unsafe.Pointer(&header))
	return data
}

// conv will convert a slice to a typedarray
//  []float32 -> Float32Array
//  []float64 -> Float32Array (for Wrapper purposes)
func conv(data interface{}) js.Value {
	switch data := data.(type) {
	case js.Value:
		return data
	case []float32:
		d := js.Global().Get("Float32Array").New(len(data))
		for i, v := range data {
			d.SetIndex(i, v)
		}
		return d
	case []uint32:
		d := js.Global().Get("Uint32Array").New(len(data))
		for i, v := range data {
			d.SetIndex(i, v)
		}
		return d
	case []float64:
		d := js.Global().Get("Float32Array").New(len(data))
		for i, v := range data {
			d.SetIndex(i, float32(v))
		}
		return d
	default:
		panic(fmt.Sprintf("Unimplemented type: %T", data))
	}

	return js.Undefined()

}

func newFloat32Array(sz int) js.Value {
	return js.Global().Get("Float32Array").New(sz)
}
