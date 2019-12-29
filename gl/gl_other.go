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

// +build !js, !wasm
// +build !android

package gl

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

const (
	// NullTexture is a nil texture
	NullTexture = Texture(0)
)

// Types that others implements
type (
	Uint         = uint32
	Buffer       = Uint
	Shader       = Uint
	Program      = Uint
	Attrib       = Uint
	Framebuffer  = Uint
	Renderbuffer = Uint
	Texture      = Uint
	VertexArray  = Uint
	Uniform      = Uint
	Enum         = Uint
)

// Wrapper for gl funcs
type Wrapper struct {
	undef
}

// ActiveTexture sets the active texture unit.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glActiveTexture.xhtml
func (g *Wrapper) ActiveTexture(texture Enum) {
	gl.ActiveTexture(texture)
}

// AttachShader attaches a shader to a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glAttachShader.xhtml
func (g *Wrapper) AttachShader(p Program, s Shader) {
	gl.AttachShader(p, s)
}

// BindAttribLocation binds a vertex attribute index with a named
// variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindAttribLocation.xhtml
func (g *Wrapper) BindAttribLocation(p Program, a Attrib, name string) {
	panic("not implemented") // TODO: Implement
}

// BindBuffer binds a buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBuffer.xhtml
func (g *Wrapper) BindBuffer(target Enum, b Buffer) {
	gl.BindBuffer(target, b)
}

// BindFramebuffer binds a framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindFramebuffer.xhtml
func (g *Wrapper) BindFramebuffer(target Enum, fb Framebuffer) {
	panic("not implemented") // TODO: Implement
}

// BindRenderbuffer binds a render buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindRenderbuffer.xhtml
func (g *Wrapper) BindRenderbuffer(target Enum, rb Renderbuffer) {
	panic("not implemented") // TODO: Implement
}

// BindTexture binds a texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindTexture.xhtml
func (g *Wrapper) BindTexture(target Enum, t Texture) {
	gl.BindTexture(target, t)
}

// BindVertexArray binds a vertex array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindVertexArray.xhtml
func (g *Wrapper) BindVertexArray(rb VertexArray) {
	gl.BindVertexArray(rb)
}

// BlendColor sets the blend color.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendColor.xhtml
func (g *Wrapper) BlendColor(red float32, green float32, blue float32, alpha float32) {
	panic("not implemented") // TODO: Implement
}

// BlendEquation sets both RGB and alpha blend equations.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendEquation.xhtml
func (g *Wrapper) BlendEquation(mode Enum) {
	panic("not implemented") // TODO: Implement
}

// BlendEquationSeparate sets RGB and alpha blend equations separately.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendEquationSeparate.xhtml
func (g *Wrapper) BlendEquationSeparate(modeRGB Enum, modeAlpha Enum) {
	panic("not implemented") // TODO: Implement
}

// BlendFunc sets the pixel blending factors.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendFunc.xhtml
func (g *Wrapper) BlendFunc(sfactor Enum, dfactor Enum) {
	gl.BlendFunc(sfactor, dfactor)
}

// BlendFunc sets the pixel RGB and alpha blending factors separately.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendFuncSeparate.xhtml
func (g *Wrapper) BlendFuncSeparate(sfactorRGB Enum, dfactorRGB Enum, sfactorAlpha Enum, dfactorAlpha Enum) {
	panic("not implemented") // TODO: Implement
}

// BufferData creates a new data store for the bound buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferData.xhtml
func (g *Wrapper) BufferData(target Enum, src []byte, usage Enum) {
	panic("not implemented") // TODO: Implement
}

// BufferInit creates a new uninitialized data store for the bound buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferData.xhtml
func (g *Wrapper) BufferInit(target Enum, size int, usage Enum) {
	panic("not implemented") // TODO: Implement
}

// BufferSubData sets some of data in the bound buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferSubData.xhtml
func (g *Wrapper) BufferSubData(target Enum, offset int, data []byte) {
	panic("not implemented") // TODO: Implement
}

// CheckFramebufferStatus reports the completeness status of the
// active framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCheckFramebufferStatus.xhtml
func (g *Wrapper) CheckFramebufferStatus(target Enum) Enum {
	panic("not implemented") // TODO: Implement
}

// Clear clears the window.
//
// The behavior of Clear is influenced by the pixel ownership test,
// the scissor test, dithering, and the buffer writemasks.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClear.xhtml
func (g *Wrapper) Clear(mask Enum) {
	gl.Clear(mask)
}

// ClearColor specifies the RGBA values used to clear color buffers.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearColor.xhtml
func (g *Wrapper) ClearColor(red float32, green float32, blue float32, alpha float32) {
	gl.ClearColor(red, green, blue, alpha)
}

// ClearDepthf sets the depth value used to clear the depth buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearDepthf.xhtml
func (g *Wrapper) ClearDepthf(d float32) {
	gl.ClearDepthf(d)
}

// ClearStencil sets the index used to clear the stencil buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearStencil.xhtml
func (g *Wrapper) ClearStencil(s int) {
	panic("not implemented") // TODO: Implement
}

// ColorMask specifies whether color components in the framebuffer
// can be written.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glColorMask.xhtml
func (g *Wrapper) ColorMask(red bool, green bool, blue bool, alpha bool) {
	panic("not implemented") // TODO: Implement
}

// CompileShader compiles the source code of s.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompileShader.xhtml
func (g *Wrapper) CompileShader(s Shader) {
	gl.CompileShader(s)
}

// CompressedTexImage2D writes a compressed 2D texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompressedTexImage2D.xhtml
func (g *Wrapper) CompressedTexImage2D(target Enum, level int, internalformat Enum, width int, height int, border int, data []byte) {
	panic("not implemented") // TODO: Implement
}

// CompressedTexSubImage2D writes a subregion of a compressed 2D texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompressedTexSubImage2D.xhtml
func (g *Wrapper) CompressedTexSubImage2D(target Enum, level int, xoffset int, yoffset int, width int, height int, format Enum, data []byte) {
	panic("not implemented") // TODO: Implement
}

// CopyTexImage2D writes a 2D texture from the current framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCopyTexImage2D.xhtml
func (g *Wrapper) CopyTexImage2D(target Enum, level int, internalformat Enum, x int, y int, width int, height int, border int) {
	panic("not implemented") // TODO: Implement
}

// CopyTexSubImage2D writes a 2D texture subregion from the
// current framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCopyTexSubImage2D.xhtml
func (g *Wrapper) CopyTexSubImage2D(target Enum, level int, xoffset int, yoffset int, x int, y int, width int, height int) {
	panic("not implemented") // TODO: Implement
}

// CreateBuffer creates a buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenBuffers.xhtml
func (g *Wrapper) CreateBuffer() Buffer {
	var b Buffer
	gl.GenBuffers(1, &b)
	return b
}

// CreateFramebuffer creates a framebuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenFramebuffers.xhtml
func (g *Wrapper) CreateFramebuffer() Framebuffer {
	panic("not implemented") // TODO: Implement
}

// CreateProgram creates a new empty program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateProgram.xhtml
func (g *Wrapper) CreateProgram() Program {
	return gl.CreateProgram()
}

// CreateRenderbuffer create a renderbuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenRenderbuffers.xhtml
func (g *Wrapper) CreateRenderbuffer() Renderbuffer {
	panic("not implemented") // TODO: Implement
}

// CreateShader creates a new empty shader object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateShader.xhtml
func (g *Wrapper) CreateShader(ty Enum) Shader {
	return gl.CreateShader(ty)
}

// CreateTexture creates a texture object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenTextures.xhtml
func (g *Wrapper) CreateTexture() Texture {
	var t Texture
	gl.GenTextures(1, &t)
	return t
}

// CreateVertexArray creates a vertex array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenVertexArrays.xhtml
func (g *Wrapper) CreateVertexArray() VertexArray {
	var v VertexArray
	gl.GenVertexArrays(1, &v)
	return v
}

// CullFace specifies which polygons are candidates for culling.
//
// Valid modes: FRONT, BACK, FRONT_AND_BACK.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCullFace.xhtml
func (g *Wrapper) CullFace(mode Enum) {
	gl.CullFace(mode)
}

// DeleteBuffer deletes the given buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteBuffers.xhtml
func (g *Wrapper) DeleteBuffer(v Buffer) {
	panic("not implemented") // TODO: Implement
}

// DeleteFramebuffer deletes the given framebuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteFramebuffers.xhtml
func (g *Wrapper) DeleteFramebuffer(v Framebuffer) {
	panic("not implemented") // TODO: Implement
}

// DeleteProgram deletes the given program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteProgram.xhtml
func (g *Wrapper) DeleteProgram(p Program) {
	panic("not implemented") // TODO: Implement
}

// DeleteRenderbuffer deletes the given render buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteRenderbuffers.xhtml
func (g *Wrapper) DeleteRenderbuffer(v Renderbuffer) {
	panic("not implemented") // TODO: Implement
}

// DeleteShader deletes shader s.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteShader.xhtml
func (g *Wrapper) DeleteShader(s Shader) {
	panic("not implemented") // TODO: Implement
}

// DeleteTexture deletes the given texture object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteTextures.xhtml
func (g *Wrapper) DeleteTexture(v Texture) {
	panic("not implemented") // TODO: Implement
}

// DeleteVertexArray deletes the given render buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteVertexArrays.xhtml
func (g *Wrapper) DeleteVertexArray(v VertexArray) {
	panic("not implemented") // TODO: Implement
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
func (g *Wrapper) DepthFunc(fn Enum) {
	gl.DepthFunc(fn)
}

// DepthMask sets the depth buffer enabled for writing.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDepthMask.xhtml
func (g *Wrapper) DepthMask(flag bool) {
	gl.DepthMask(flag)
}

// DepthRangef sets the mapping from normalized device coordinates to
// window coordinates.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDepthRangef.xhtml
func (g *Wrapper) DepthRangef(n float32, f float32) {
	panic("not implemented") // TODO: Implement
}

// DetachShader detaches the shader s from the program p.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDetachShader.xhtml
func (g *Wrapper) DetachShader(p Program, s Shader) {
	panic("not implemented") // TODO: Implement
}

// Disable disables various GL capabilities.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDisable.xhtml
func (g *Wrapper) Disable(cap Enum) {
	gl.Disable(cap)
}

// DisableVertexAttribArray disables a vertex attribute array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDisableVertexAttribArray.xhtml
func (g *Wrapper) DisableVertexAttribArray(a Attrib) {
	panic("not implemented") // TODO: Implement
}

// DrawArrays renders geometric primitives from the bound data.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArrays.xhtml
func (g *Wrapper) DrawArrays(mode Enum, first int, count int) {
	gl.DrawArrays(mode, int32(first), int32(count))
}

// DrawElements renders primitives from a bound buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawElements.xhtml
func (g *Wrapper) DrawElements(mode Enum, count int, ty Enum, offset int) {
	gl.DrawElements(mode, int32(count), ty, unsafe.Pointer(uintptr(offset)))
}

// TODO(crawshaw): consider DrawElements8 / DrawElements16 / DrawElements32
// Enable enables various GL capabilities.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnable.xhtml
func (g *Wrapper) Enable(cap Enum) {
	gl.Enable(cap)
}

// EnableVertexAttribArray enables a vertex attribute array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnableVertexAttribArray.xhtml
func (g *Wrapper) EnableVertexAttribArray(a Attrib) {
	gl.EnableVertexAttribArray(a)
}

// Finish blocks until the effects of all previously called GL
// commands are complete.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFinish.xhtml
func (g *Wrapper) Finish() {
	panic("not implemented") // TODO: Implement
}

// Flush empties all buffers. It does not block.
//
// An OpenGL implementation may buffer network communication,
// the command stream, or data inside the graphics accelerator.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFlush.xhtml
func (g *Wrapper) Flush() {
	panic("not implemented") // TODO: Implement
}

// FramebufferRenderbuffer attaches rb to the current frame buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFramebufferRenderbuffer.xhtml
func (g *Wrapper) FramebufferRenderbuffer(target Enum, attachment Enum, rbTarget Enum, rb Renderbuffer) {
	panic("not implemented") // TODO: Implement
}

// FramebufferTexture2D attaches the t to the current frame buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFramebufferTexture2D.xhtml
func (g *Wrapper) FramebufferTexture2D(target Enum, attachment Enum, texTarget Enum, t Texture, level int) {
	panic("not implemented") // TODO: Implement
}

// FrontFace defines which polygons are front-facing.
//
// Valid modes: CW, CCW.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFrontFace.xhtml
func (g *Wrapper) FrontFace(mode Enum) {
	gl.FrontFace(mode)
}

// GenerateMipmap generates mipmaps for the current texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenerateMipmap.xhtml
func (g *Wrapper) GenerateMipmap(target Enum) {
	gl.GenerateMipmap(target)
}

// GetActiveAttrib returns details about an active attribute variable.
// A value of 0 for index selects the first active attribute variable.
// Permissible values for index range from 0 to the number of active
// attribute variables minus 1.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveAttrib.xhtml
func (g *Wrapper) GetActiveAttrib(p Program, index uint32) (string, int, Enum) {
	var sz int32
	var ty Enum

	nameSz := int32(256)
	nameBuf := [256]byte{}

	gl.GetActiveAttrib(p, index, nameSz, nil, &sz, &ty, &nameBuf[0])

	return gl.GoStr(&nameBuf[0]), int(sz), ty
}

// GetActiveUniform returns details about an active uniform variable.
// A value of 0 for index selects the first active uniform variable.
// Permissible values for index range from 0 to the number of active
// uniform variables minus 1.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveUniform.xhtml
func (g *Wrapper) GetActiveUniform(p Program, index uint32) (string, int, Enum) {
	var sz int32
	var ty Enum

	nameSz := int32(256)
	nameBuf := [256]byte{}

	gl.GetActiveUniform(p, index, nameSz, nil, &sz, &ty, &nameBuf[0])

	return gl.GoStr(&nameBuf[0]), int(sz), ty
}

// GetAttachedShaders returns the shader objects attached to program p.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetAttachedShaders.xhtml
func (g *Wrapper) GetAttachedShaders(p Program) []Shader {
	panic("not implemented") // TODO: Implement
}

// GetAttribLocation returns the location of an attribute variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetAttribLocation.xhtml
func (g *Wrapper) GetAttribLocation(p Program, name string) Attrib {
	return Attrib(gl.GetAttribLocation(p, gl.Str(name+"\x00")))
}

// GetBooleanv returns the boolean values of parameter pname.
//
// Many boolean parameters can be queried more easily using IsEnabled.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func (g *Wrapper) GetBooleanv(dst []bool, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetFloatv returns the float values of parameter pname.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func (g *Wrapper) GetFloatv(dst []float32, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetIntegerv returns the int values of parameter pname.
//
// Single values may be queried more easily using GetInteger.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func (g *Wrapper) GetIntegerv(dst []int32, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetInteger returns the int value of parameter pname.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func (g *Wrapper) GetInteger(pname Enum) int {
	panic("not implemented") // TODO: Implement
}

// GetBufferParameteri returns a parameter for the active buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetBufferParameter.xhtml
func (g *Wrapper) GetBufferParameteri(target Enum, value Enum) int {
	panic("not implemented") // TODO: Implement
}

// GetError returns the next error.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetError.xhtml
func (g *Wrapper) GetError() Enum {
	panic("not implemented") // TODO: Implement
}

// GetFramebufferAttachmentParameteri returns attachment parameters
// for the active framebuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetFramebufferAttachmentParameteriv.xhtml
func (g *Wrapper) GetFramebufferAttachmentParameteri(target Enum, attachment Enum, pname Enum) int {
	panic("not implemented") // TODO: Implement
}

// GetProgrami returns a parameter value for a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramiv.xhtml
func (g *Wrapper) GetProgrami(p Program, pname Enum) int {
	var pi int32
	gl.GetProgramiv(p, pname, &pi)
	return int(pi)
}

// GetProgramInfoLog returns the information log for a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramInfoLog.xhtml
func (g *Wrapper) GetProgramInfoLog(p Program) string {
	panic("not implemented") // TODO: Implement
}

// GetRenderbufferParameteri returns a parameter value for a render buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetRenderbufferParameteriv.xhtml
func (g *Wrapper) GetRenderbufferParameteri(target Enum, pname Enum) int {
	panic("not implemented") // TODO: Implement
}

// GetShaderi returns a parameter value for a shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderiv.xhtml
func (g *Wrapper) GetShaderi(s Shader, pname Enum) int {
	var p int32
	gl.GetShaderiv(s, pname, &p)
	return int(p)
}

// GetShaderInfoLog returns the information log for a shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderInfoLog.xhtml
func (g *Wrapper) GetShaderInfoLog(s Shader) string {
	var logLength int32
	gl.GetShaderiv(s, gl.INFO_LOG_LENGTH, &logLength)

	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(s, logLength, nil, gl.Str(log))
	return log
}

// GetShaderPrecisionFormat returns range and precision limits for
// shader types.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderPrecisionFormat.xhtml
func (g *Wrapper) GetShaderPrecisionFormat(shadertype Enum, precisiontype Enum) (rangeLow int, rangeHigh int, precision int) {
	panic("not implemented") // TODO: Implement
}

// GetShaderSource returns source code of shader s.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderSource.xhtml
func (g *Wrapper) GetShaderSource(s Shader) string {
	panic("not implemented") // TODO: Implement
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
func (g *Wrapper) GetString(pname Enum) string {
	return gl.GoStr(gl.GetString(pname))
}

// GetTexParameterfv returns the float values of a texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetTexParameter.xhtml
func (g *Wrapper) GetTexParameterfv(dst []float32, target Enum, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetTexParameteriv returns the int values of a texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetTexParameter.xhtml
func (g *Wrapper) GetTexParameteriv(dst []int32, target Enum, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetUniformfv returns the float values of a uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniform.xhtml
func (g *Wrapper) GetUniformfv(dst []float32, src Uniform, p Program) {
	panic("not implemented") // TODO: Implement
}

// GetUniformiv returns the float values of a uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniform.xhtml
func (g *Wrapper) GetUniformiv(dst []int32, src Uniform, p Program) {
	panic("not implemented") // TODO: Implement
}

// GetUniformLocation returns the location of a uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniformLocation.xhtml
func (g *Wrapper) GetUniformLocation(p Program, name string) Uniform {
	return Uniform(gl.GetUniformLocation(p, gl.Str(name+"\x00")))
}

// GetVertexAttribf reads the float value of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func (g *Wrapper) GetVertexAttribf(src Attrib, pname Enum) float32 {
	panic("not implemented") // TODO: Implement
}

// GetVertexAttribfv reads float values of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func (g *Wrapper) GetVertexAttribfv(dst []float32, src Attrib, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetVertexAttribi reads the int value of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func (g *Wrapper) GetVertexAttribi(src Attrib, pname Enum) int32 {
	panic("not implemented") // TODO: Implement
}

// GetVertexAttribiv reads int values of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func (g *Wrapper) GetVertexAttribiv(dst []int32, src Attrib, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// TODO(crawshaw): glGetVertexAttribPointerv
// Hint sets implementation-specific modes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glHint.xhtml
func (g *Wrapper) Hint(target Enum, mode Enum) {
	gl.Hint(target, mode)
}

// IsBuffer reports if b is a valid buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsBuffer.xhtml
func (g *Wrapper) IsBuffer(b Buffer) bool {
	panic("not implemented") // TODO: Implement
}

// IsEnabled reports if cap is an enabled capability.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsEnabled.xhtml
func (g *Wrapper) IsEnabled(cap Enum) bool {
	panic("not implemented") // TODO: Implement
}

// IsFramebuffer reports if fb is a valid frame buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsFramebuffer.xhtml
func (g *Wrapper) IsFramebuffer(fb Framebuffer) bool {
	panic("not implemented") // TODO: Implement
}

// IsProgram reports if p is a valid program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsProgram.xhtml
func (g *Wrapper) IsProgram(p Program) bool {
	panic("not implemented") // TODO: Implement
}

// IsRenderbuffer reports if rb is a valid render buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsRenderbuffer.xhtml
func (g *Wrapper) IsRenderbuffer(rb Renderbuffer) bool {
	panic("not implemented") // TODO: Implement
}

// IsShader reports if s is valid shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsShader.xhtml
func (g *Wrapper) IsShader(s Shader) bool {
	panic("not implemented") // TODO: Implement
}

// IsTexture reports if t is a valid texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsTexture.xhtml
func (g *Wrapper) IsTexture(t Texture) bool {
	panic("not implemented") // TODO: Implement
}

// LineWidth specifies the width of lines.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glLineWidth.xhtml
func (g *Wrapper) LineWidth(width float32) {
	gl.LineWidth(width)
}

// LinkProgram links the specified program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glLinkProgram.xhtml
func (g *Wrapper) LinkProgram(p Program) {
	gl.LinkProgram(p)
}

// PixelStorei sets pixel storage parameters.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glPixelStorei.xhtml
func (g *Wrapper) PixelStorei(pname Enum, param int32) {
	panic("not implemented") // TODO: Implement
}

// PolygonOffset sets the scaling factors for depth offsets.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glPolygonOffset.xhtml
func (g *Wrapper) PolygonOffset(factor float32, units float32) {
	panic("not implemented") // TODO: Implement
}

// ReadPixels returns pixel data from a buffer.
//
// In GLES 3, the source buffer is controlled with ReadBuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glReadPixels.xhtml
func (g *Wrapper) ReadPixels(dst []byte, x int, y int, width int, height int, format Enum, ty Enum) {
	panic("not implemented") // TODO: Implement
}

// ReleaseShaderCompiler frees resources allocated by the shader compiler.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glReleaseShaderCompiler.xhtml
func (g *Wrapper) ReleaseShaderCompiler() {
	panic("not implemented") // TODO: Implement
}

// RenderbufferStorage establishes the data storage, format, and
// dimensions of a renderbuffer object's image.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glRenderbufferStorage.xhtml
func (g *Wrapper) RenderbufferStorage(target Enum, internalFormat Enum, width int, height int) {
	panic("not implemented") // TODO: Implement
}

// SampleCoverage sets multisample coverage parameters.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glSampleCoverage.xhtml
func (g *Wrapper) SampleCoverage(value float32, invert bool) {
	panic("not implemented") // TODO: Implement
}

// Scissor defines the scissor box rectangle, in window coordinates.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glScissor.xhtml
func (g *Wrapper) Scissor(x int32, y int32, width int32, height int32) {
	panic("not implemented") // TODO: Implement
}

// TODO(crawshaw): ShaderBinary
// ShaderSource sets the source code of s to the given source code.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glShaderSource.xhtml
func (g *Wrapper) ShaderSource(s Shader, src string) {
	csources, free := gl.Strs(src + "\x00")
	gl.ShaderSource(s, 1, csources, nil)
	free()
}

// StencilFunc sets the front and back stencil test reference value.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilFunc.xhtml
func (g *Wrapper) StencilFunc(fn Enum, ref int, mask uint32) {
	panic("not implemented") // TODO: Implement
}

// StencilFunc sets the front or back stencil test reference value.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilFuncSeparate.xhtml
func (g *Wrapper) StencilFuncSeparate(face Enum, fn Enum, ref int, mask uint32) {
	panic("not implemented") // TODO: Implement
}

// StencilMask controls the writing of bits in the stencil planes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilMask.xhtml
func (g *Wrapper) StencilMask(mask uint32) {
	panic("not implemented") // TODO: Implement
}

// StencilMaskSeparate controls the writing of bits in the stencil planes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilMaskSeparate.xhtml
func (g *Wrapper) StencilMaskSeparate(face Enum, mask uint32) {
	panic("not implemented") // TODO: Implement
}

// StencilOp sets front and back stencil test actions.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilOp.xhtml
func (g *Wrapper) StencilOp(fail Enum, zfail Enum, zpass Enum) {
	panic("not implemented") // TODO: Implement
}

// StencilOpSeparate sets front or back stencil tests.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilOpSeparate.xhtml
func (g *Wrapper) StencilOpSeparate(face Enum, sfail Enum, dpfail Enum, dppass Enum) {
	panic("not implemented") // TODO: Implement
}

// TexImage2D writes a 2D texture image.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexImage2D.xhtml
func (g *Wrapper) TexImage2D(target Enum, level int, internalFormat int, width int, height int, format Enum, ty Enum, data []byte) {
	gl.TexImage2D(
		target,
		int32(level),
		int32(internalFormat),
		int32(width), int32(height),
		0, //border
		format,
		ty,
		unsafe.Pointer(&data[0]),
	)
}

// TexSubImage2D writes a subregion of a 2D texture image.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexSubImage2D.xhtml
func (g *Wrapper) TexSubImage2D(target Enum, level int, x int, y int, width int, height int, format Enum, ty Enum, data []byte) {
	panic("not implemented") // TODO: Implement
}

// TexParameterf sets a float texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (g *Wrapper) TexParameterf(target Enum, pname Enum, param float32) {
	gl.TexParameterf(target, pname, param)
}

// TexParameterfv sets a float texture parameter array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (g *Wrapper) TexParameterfv(target Enum, pname Enum, params []float32) {
	panic("not implemented") // TODO: Implement
}

// TexParameteri sets an integer texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (g *Wrapper) TexParameteri(target Enum, pname Enum, param int) {
	gl.TexParameteri(target, pname, int32(param))
}

// TexParameteriv sets an integer texture parameter array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (g *Wrapper) TexParameteriv(target Enum, pname Enum, params []int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform1f writes a float uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform1f(dst Uniform, v float32) {
	gl.Uniform1f(int32(dst), v)
}

// Uniform1fv writes a [len(src)]float uniform array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform1fv(dst Uniform, src []float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform1i writes an int uniform variable.
//
// Uniform1i and Uniform1iv are the only two functions that may be used
// to load uniform variables defined as sampler types. Loading samplers
// with any other function will result in a INVALID_OPERATION error.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform1i(dst Uniform, v int) {
	gl.Uniform1i(int32(dst), int32(v))
}

// Uniform1iv writes a int uniform array of len(src) elements.
//
// Uniform1i and Uniform1iv are the only two functions that may be used
// to load uniform variables defined as sampler types. Loading samplers
// with any other function will result in a INVALID_OPERATION error.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform1iv(dst Uniform, src []int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform2f writes a vec2 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform2f(dst Uniform, v0 float32, v1 float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform2fv writes a vec2 uniform array of len(src)/2 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform2fv(dst Uniform, src []float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform2i writes an ivec2 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform2i(dst Uniform, v0 int, v1 int) {
	panic("not implemented") // TODO: Implement
}

// Uniform2iv writes an ivec2 uniform array of len(src)/2 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform2iv(dst Uniform, src []int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform3f writes a vec3 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform3f(dst Uniform, v0 float32, v1 float32, v2 float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform3fv writes a vec3 uniform array of len(src)/3 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform3fv(dst Uniform, src []float32) {
	gl.Uniform3fv(int32(dst), 1, &src[0])
}

// Uniform3i writes an ivec3 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform3i(dst Uniform, v0 int32, v1 int32, v2 int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform3iv writes an ivec3 uniform array of len(src)/3 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform3iv(dst Uniform, src []int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform4f writes a vec4 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform4f(dst Uniform, v0 float32, v1 float32, v2 float32, v3 float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform4fv writes a vec4 uniform array of len(src)/4 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform4fv(dst Uniform, src []float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform4i writes an ivec4 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform4i(dst Uniform, v0 int32, v1 int32, v2 int32, v3 int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform4iv writes an ivec4 uniform array of len(src)/4 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) Uniform4iv(dst Uniform, src []int32) {
	panic("not implemented") // TODO: Implement
}

// UniformMatrix2fv writes 2x2 matrices. Each matrix uses four
// float32 values, so the number of matrices written is len(src)/4.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) UniformMatrix2fv(dst Uniform, src []float32) {
	panic("not implemented") // TODO: Implement
}

// UniformMatrix3fv writes 3x3 matrices. Each matrix uses nine
// float32 values, so the number of matrices written is len(src)/9.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) UniformMatrix3fv(dst Uniform, src []float32) {
	panic("not implemented") // TODO: Implement
}

// UniformMatrix4fv writes 4x4 matrices. Each matrix uses 16
// float32 values, so the number of matrices written is len(src)/16.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g *Wrapper) UniformMatrix4fv(dst Uniform, src []float32) {
	gl.UniformMatrix4fv(
		int32(dst),
		1,
		false,
		&src[0],
	)
}

// UseProgram sets the active program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUseProgram.xhtml
func (g *Wrapper) UseProgram(p Program) {
	gl.UseProgram(p)
}

// ValidateProgram checks to see whether the executables contained in
// program can execute given the current OpenGL state.
//
// Typically only used for debugging.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glValidateProgram.xhtml
func (g *Wrapper) ValidateProgram(p Program) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib1f writes a float vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (g *Wrapper) VertexAttrib1f(dst Attrib, x float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib1fv writes a float vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (g *Wrapper) VertexAttrib1fv(dst Attrib, src []float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib2f writes a vec2 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (g *Wrapper) VertexAttrib2f(dst Attrib, x float32, y float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib2fv writes a vec2 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (g *Wrapper) VertexAttrib2fv(dst Attrib, src []float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib3f writes a vec3 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (g *Wrapper) VertexAttrib3f(dst Attrib, x float32, y float32, z float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib3fv writes a vec3 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (g *Wrapper) VertexAttrib3fv(dst Attrib, src []float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib4f writes a vec4 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (g *Wrapper) VertexAttrib4f(dst Attrib, x float32, y float32, z float32, w float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib4fv writes a vec4 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (g *Wrapper) VertexAttrib4fv(dst Attrib, src []float32) {
	panic("not implemented") // TODO: Implement
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
func (g *Wrapper) VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride int, offset int) {
	gl.VertexAttribPointer(
		dst,
		int32(size),
		ty,
		normalized,
		int32(stride),
		unsafe.Pointer(uintptr(offset)),
	)
}

// Viewport sets the viewport, an affine transformation that
// normalizes device coordinates to window coordinates.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glViewport.xhtml
func (g *Wrapper) Viewport(x int, y int, width int, height int) {
	gl.Viewport(int32(x), int32(y), int32(width), int32(height))
}

// glGetUniformBlockIndex retrieves the index of a uniform block within program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockIndex.xhtml
//GetUniformBlockIndex(p Program, name string) int
// UniformBlockBinding assign a binding point to an active uniform block
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockBinding.xhtml
//UniformBlockBinding(p Program, index, bind int)
// BindBufferBase bind a buffer object to an indexed buffer target
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBufferBase.xhtml
//BindBufferBase(target Enum, n uint32, b Buffer)
// DrawArraysInstanced draw multiple instances of a range of elements
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
func (g *Wrapper) DrawArraysInstanced(mode Enum, first int, count int, primcount int) {
	gl.DrawArraysInstanced(mode, int32(first), int32(count), int32(primcount))
}

// DrawElementsInstanced — draw multiple instances of a set of elements
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawElementsInstanced.xhtml
func (g *Wrapper) DrawElementsInstanced(mode Enum, count int, typ Enum, offset int, primcount int) {
	//off := unsafe.Pointer(uintptr(offset))
	gl.DrawElementsInstanced(
		mode,
		int32(count),
		typ,
		unsafe.Pointer(uintptr(offset)),
		int32(primcount),
	)
}

// VertexAttribDivisor  modify the rate at which generic vertex attributes
// advance during instanced rendering
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
func (g *Wrapper) VertexAttribDivisor(index Attrib, divisor int) {
	gl.VertexAttribDivisor(index, uint32(divisor))
}

// BufferDataX will type switch the interface and select the proper type
// {lpf} Custom func
// TODO: Should move this Elsewhere as an helper
func (g *Wrapper) BufferDataX(target Enum, d interface{}, usage Enum) {

	switch v := d.(type) {
	case []float32:
		gl.BufferData(target, len(v)*4, gl.Ptr(v), usage)
	case []uint32:
		gl.BufferData(target, len(v)*4, gl.Ptr(v), usage)
	default:
		panic(fmt.Sprintf("Buffer type not implemented: %T", d))
	}
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
