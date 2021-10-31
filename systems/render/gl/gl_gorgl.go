//go:build android || mobile || gorgl

package gl

/*
#cgo ios                LDFLAGS: -framework OpenGLES
#cgo darwin,amd64,!ios  LDFLAGS: -framework OpenGL
#cgo darwin,arm         LDFLAGS: -framework OpenGLES
#cgo darwin,arm64       LDFLAGS: -framework OpenGLES
#cgo linux              LDFLAGS: -lGLESv2
#cgo openbsd            LDFLAGS: -L/usr/X11R6/lib/ -lGLESv2
#cgo android            LDFLAGS: -lGLESv3

#cgo android            CFLAGS: -Dos_android
#cgo ios                CFLAGS: -Dos_ios
#cgo darwin,amd64,!ios  CFLAGS: -Dos_osx
#cgo darwin,arm         CFLAGS: -Dos_ios
#cgo darwin,arm64       CFLAGS: -Dos_ios
#cgo darwin             CFLAGS: -DGL_SILENCE_DEPRECATION
#cgo linux              CFLAGS: -Dos_linux
#cgo openbsd            CFLAGS: -Dos_openbsd

#cgo openbsd            CFLAGS: -I/usr/X11R6/include/

#include <stdint.h>
#include <stdlib.h>

#ifdef os_android
#include <GLES3/gl3.h> // {lpf} previous: GLES2/gl2.h
#elif os_linux
#include <GLES3/gl3.h> // install on Ubuntu with: sudo apt-get install libegl1-mesa-dev libgles2-mesa-dev libx11-dev
#elif os_openbsd
#include <GLES3/gl3.h>
#endif

#ifdef os_ios
#include <OpenGLES/ES2/glext.h>
#endif

#ifdef os_osx
#include <OpenGL/gl3.h>
#define GL_ES_VERSION_3_0 1
#endif

#if defined(GL_ES_VERSION_3_0) && GL_ES_VERSION_3_0
#define GLES_VERSION "GL_ES_3_0"
#else
#define GLES_VERSION "GL_ES_2_0"
#endif
*/
import "C"

import (
	"fmt"
	"unsafe"
)

const (
	// NullTexture is a nil texture
	NullTexture     = Texture(0)
	NullFramebuffer = Framebuffer(0)
	NullBuffer      = Buffer(0)
	NullVertexArray = VertexArray(0)
	// Null value varies with the wrapper
	Null = 0
)

// Types that others implements
type (
	Uint         = C.GLuint
	Int          = C.GLint
	Buffer       = Uint
	Shader       = Uint
	Program      = Uint
	Attrib       = Uint
	Framebuffer  = Uint
	Renderbuffer = Uint
	Texture      = Uint
	VertexArray  = Uint
	Uniform      = Int
	Enum         = Uint
)

// IsValid returns if any of the values above is valid
func IsValid(v Uint) bool { return v != 0 }

// Wrapper for gl funcs
type Wrapper struct{}

func (Wrapper) String() string {
	return "go_gorgl.go wrapper"
}

var _ Context3 = &Wrapper{}

// ActiveTexture sets the active texture unit.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glActiveTexture.xhtml
func (glw Wrapper) ActiveTexture(texture Enum) {
	C.glActiveTexture(texture)
}

// AttachShader attaches a shader to a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glAttachShader.xhtml
func (glw Wrapper) AttachShader(p Program, s Shader) {
	C.glAttachShader(p, s)
}

// BindAttribLocation binds a vertex attribute index with a named
// variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindAttribLocation.xhtml
func (glw Wrapper) BindAttribLocation(p Program, a Attrib, name string) {
	panic("not implemented") // TODO: Implement
}

// BindBuffer binds a buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBuffer.xhtml
func (glw Wrapper) BindBuffer(target Enum, b Buffer) {
	C.glBindBuffer(target, b)
}

// BindFramebuffer binds a framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindFramebuffer.xhtml
func (glw Wrapper) BindFramebuffer(target Enum, fb Framebuffer) {
	C.glBindFramebuffer(target, fb)
}

// BindRenderbuffer binds a render buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindRenderbuffer.xhtml
func (glw Wrapper) BindRenderbuffer(target Enum, rb Renderbuffer) {
	C.glBindRenderbuffer(target, rb)
}

// BindTexture binds a texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindTexture.xhtml
func (glw Wrapper) BindTexture(target Enum, t Texture) {
	C.glBindTexture(target, t)
}

// BindVertexArray binds a vertex array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindVertexArray.xhtml
func (glw Wrapper) BindVertexArray(rb VertexArray) {
	C.glBindVertexArray(rb)
}

// BlendColor sets the blend color.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendColor.xhtml
func (glw Wrapper) BlendColor(red float32, green float32, blue float32, alpha float32) {
	panic("not implemented") // TODO: Implement
}

// BlendEquation sets both RGB and alpha blend equations.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendEquation.xhtml
func (glw Wrapper) BlendEquation(mode Enum) {
	panic("not implemented") // TODO: Implement
}

// BlendEquationSeparate sets RGB and alpha blend equations separately.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendEquationSeparate.xhtml
func (glw Wrapper) BlendEquationSeparate(modeRGB Enum, modeAlpha Enum) {
	panic("not implemented") // TODO: Implement
}

// BlendFunc sets the pixel blending factors.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendFunc.xhtml
func (glw Wrapper) BlendFunc(sfactor Enum, dfactor Enum) {
	C.glBlendFunc(sfactor, dfactor)
}

// BlendFuncSeparate sets the pixel RGB and alpha blending factors separately.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendFuncSeparate.xhtml
func (glw Wrapper) BlendFuncSeparate(sfactorRGB Enum, dfactorRGB Enum, sfactorAlpha Enum, dfactorAlpha Enum) {
	panic("not implemented") // TODO: Implement
}

// BufferInit creates a new uninitialized data store for the bound buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferData.xhtml
func (glw Wrapper) BufferInit(target Enum, size int, usage Enum) {
	C.glBufferData(target, C.GLsizeiptr(size), nil, usage)
}

// BufferData creates a new data store for the bound buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferData.xhtml
func (glw Wrapper) BufferData(target Enum, data interface{}, usage Enum) {
	d, sz := conv(data)
	C.glBufferData(target, C.GLsizeiptr(sz), d, usage)
}

// BufferSubData sets some of data in the bound buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferSubData.xhtml
func (glw Wrapper) BufferSubData(target Enum, offset int, data interface{}) {
	d, sz := conv(data)
	C.glBufferSubData(target, C.GLsizeiptr(offset), C.GLsizeiptr(sz), d)
}

// CheckFramebufferStatus reports the completeness status of the
// active framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCheckFramebufferStatus.xhtml
func (glw Wrapper) CheckFramebufferStatus(target Enum) Enum {
	return C.glCheckFramebufferStatus(target)
}

// Clear clears the window.
//
// The behavior of Clear is influenced by the pixel ownership test,
// the scissor test, dithering, and the buffer writemasks.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClear.xhtml
func (glw Wrapper) Clear(mask Enum) {
	C.glClear(mask)
}

// ClearColor specifies the RGBA values used to clear color buffers.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearColor.xhtml
func (glw Wrapper) ClearColor(r, g, b, a float32) {
	C.glClearColor(C.GLfloat(r), C.GLfloat(g), C.GLfloat(b), C.GLfloat(a))
}

// ClearDepthf sets the depth value used to clear the depth buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearDepthf.xhtml
func (glw Wrapper) ClearDepthf(d float32) {
	C.glClearDepthf(C.GLfloat(d))
}

// ClearStencil sets the index used to clear the stencil buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearStencil.xhtml
func (glw Wrapper) ClearStencil(s int) {
	panic("not implemented") // TODO: Implement
}

// ColorMask specifies whether color components in the framebuffer
// can be written.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glColorMask.xhtml
func (glw Wrapper) ColorMask(red bool, green bool, blue bool, alpha bool) {
	panic("not implemented") // TODO: Implement
}

// CompileShader compiles the source code of s.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompileShader.xhtml
func (glw Wrapper) CompileShader(s Shader) {
	C.glCompileShader(s)
}

// CompressedTexImage2D writes a compressed 2D texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompressedTexImage2D.xhtml
func (glw Wrapper) CompressedTexImage2D(target Enum, level int, internalformat Enum, width int, height int, border int, data []byte) {
	panic("not implemented") // TODO: Implement
}

// CompressedTexSubImage2D writes a subregion of a compressed 2D texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompressedTexSubImage2D.xhtml
func (glw Wrapper) CompressedTexSubImage2D(target Enum, level int, xoffset int, yoffset int, width int, height int, format Enum, data []byte) {
	panic("not implemented") // TODO: Implement
}

// CopyTexImage2D writes a 2D texture from the current framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCopyTexImage2D.xhtml
func (glw Wrapper) CopyTexImage2D(target Enum, level int, internalformat Enum, x int, y int, width int, height int, border int) {
	panic("not implemented") // TODO: Implement
}

// CopyTexSubImage2D writes a 2D texture subregion from the
// current framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCopyTexSubImage2D.xhtml
func (glw Wrapper) CopyTexSubImage2D(target Enum, level int, xoffset int, yoffset int, x int, y int, width int, height int) {
	panic("not implemented") // TODO: Implement
}

// CreateBuffer creates a buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenBuffers.xhtml
func (glw Wrapper) CreateBuffer() Buffer {
	var b Buffer
	C.glGenBuffers(1, &b)
	return b
}

// CreateFramebuffer creates a framebuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenFramebuffers.xhtml
func (glw Wrapper) CreateFramebuffer() Framebuffer {
	var fb Framebuffer
	C.glGenFramebuffers(1, &fb)
	return fb
}

// CreateProgram creates a new empty program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateProgram.xhtml
func (glw Wrapper) CreateProgram() Program {
	return C.glCreateProgram()
}

// CreateRenderbuffer create a renderbuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenRenderbuffers.xhtml
func (glw Wrapper) CreateRenderbuffer() Renderbuffer {
	var rb Renderbuffer
	C.glGenRenderbuffers(1, &rb)
	return rb
}

// CreateShader creates a new empty shader object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateShader.xhtml
func (glw Wrapper) CreateShader(ty Enum) Shader {
	return C.glCreateShader(ty)
}

// CreateTexture creates a texture object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenTextures.xhtml
func (glw Wrapper) CreateTexture() Texture {
	var t Texture
	C.glGenTextures(1, &t)
	return t
}

// CreateVertexArray creates a vertex array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenVertexArrays.xhtml
func (glw Wrapper) CreateVertexArray() VertexArray {
	var v VertexArray
	C.glGenVertexArrays(1, &v)
	return v
}

// CullFace specifies which polygons are candidates for culling.
//
// Valid modes: FRONT, BACK, FRONT_AND_BACK.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCullFace.xhtml
func (glw Wrapper) CullFace(mode Enum) {
	C.glCullFace(mode)
}

// DeleteBuffer deletes the given buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteBuffers.xhtml
func (glw Wrapper) DeleteBuffer(v Buffer) {
	C.glDeleteBuffers(1, &v)
}

// DeleteFramebuffer deletes the given framebuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteFramebuffers.xhtml
func (glw Wrapper) DeleteFramebuffer(v Framebuffer) {
	panic("not implemented") // TODO: Implement
}

// DeleteProgram deletes the given program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteProgram.xhtml
func (glw Wrapper) DeleteProgram(p Program) {
	C.glDeleteProgram(p)
}

// DeleteRenderbuffer deletes the given render buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteRenderbuffers.xhtml
func (glw Wrapper) DeleteRenderbuffer(v Renderbuffer) {
	panic("not implemented") // TODO: Implement
}

// DeleteShader deletes shader s.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteShader.xhtml
func (glw Wrapper) DeleteShader(s Shader) {
	C.glDeleteShader(s)
}

// DeleteTexture deletes the given texture object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteTextures.xhtml
func (glw Wrapper) DeleteTexture(v Texture) {
	C.glDeleteTextures(1, &v)
}

// DeleteVertexArray deletes the given render buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteVertexArrays.xhtml
func (glw Wrapper) DeleteVertexArray(v VertexArray) {
	C.glDeleteVertexArrays(1, &v)
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
func (glw Wrapper) DepthFunc(fn Enum) {
	C.glDepthFunc(fn)
}

// DepthMask sets the depth buffer enabled for writing.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDepthMask.xhtml
func (glw Wrapper) DepthMask(flag bool) {
	C.glDepthMask(glBoolean(flag))
}

// DepthRangef sets the mapping from normalized device coordinates to
// window coordinates.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDepthRangef.xhtml
func (glw Wrapper) DepthRangef(n float32, f float32) {
	panic("not implemented") // TODO: Implement
}

// DetachShader detaches the shader s from the program p.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDetachShader.xhtml
func (glw Wrapper) DetachShader(p Program, s Shader) {
	panic("not implemented") // TODO: Implement
}

// Disable disables various GL capabilities.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDisable.xhtml
func (glw Wrapper) Disable(cap Enum) {
	C.glDisable(cap)
}

// DisableVertexAttribArray disables a vertex attribute array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDisableVertexAttribArray.xhtml
func (glw Wrapper) DisableVertexAttribArray(a Attrib) {
	C.glDisableVertexAttribArray(a)
}

// DrawArrays renders geometric primitives from the bound data.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArrays.xhtml
func (glw Wrapper) DrawArrays(mode Enum, first int, count int) {
	C.glDrawArrays(mode, C.GLint(first), C.GLint(count))
}

// DrawElements renders primitives from a bound buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawElements.xhtml
func (glw Wrapper) DrawElements(mode Enum, count int, ty Enum, offset int) {
	C.glDrawElements(
		mode,
		C.GLint(count),
		ty,
		unsafe.Pointer(uintptr(offset)), // nolint: govet
	)
}

// TODO(crawshaw): consider DrawElements8 / DrawElements16 / DrawElements32

// Enable enables various GL capabilities.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnable.xhtml
func (glw Wrapper) Enable(cap Enum) {
	C.glEnable(cap)
}

// EnableVertexAttribArray enables a vertex attribute array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnableVertexAttribArray.xhtml
func (glw Wrapper) EnableVertexAttribArray(a Attrib) {
	C.glEnableVertexAttribArray(a)
}

// Finish blocks until the effects of all previously called GL
// commands are complete.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFinish.xhtml
func (glw Wrapper) Finish() {
	C.glFinish()
}

// Flush empties all buffers. It does not block.
//
// An OpenGL implementation may buffer network communication,
// the command stream, or data inside the graphics accelerator.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFlush.xhtml
func (glw Wrapper) Flush() {
	C.glFlush()
}

// FramebufferRenderbuffer attaches rb to the current frame buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFramebufferRenderbuffer.xhtml
func (glw Wrapper) FramebufferRenderbuffer(target Enum, attachment Enum, rbTarget Enum, rb Renderbuffer) {
	C.glFramebufferRenderbuffer(target, attachment, rbTarget, rb)
}

// FramebufferTexture2D attaches the t to the current frame buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFramebufferTexture2D.xhtml
func (glw Wrapper) FramebufferTexture2D(target Enum, attachment Enum, texTarget Enum, t Texture, level int) {
	C.glFramebufferTexture2D(
		target,
		attachment,
		texTarget,
		t,
		C.GLint(level),
	)
}

// FrontFace defines which polygons are front-facing.
//
// Valid modes: CW, CCW.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFrontFace.xhtml
func (glw Wrapper) FrontFace(mode Enum) {
	C.glFrontFace(mode)
}

// GenerateMipmap generates mipmaps for the current texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenerateMipmap.xhtml
func (glw Wrapper) GenerateMipmap(target Enum) {
	C.glGenerateMipmap(target)
}

// GetActiveAttrib returns details about an active attribute variable.
// A value of 0 for index selects the first active attribute variable.
// Permissible values for index range from 0 to the number of active
// attribute variables minus 1.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveAttrib.xhtml
func (glw Wrapper) GetActiveAttrib(p Program, index uint32) (string, int, Enum) {
	bufSize := GetProgrami(p, ACTIVE_ATTRIBUTE_MAX_LENGTH)
	buf := C.malloc(C.size_t(bufSize))
	defer C.free(buf)

	var cSize C.GLint
	var cType C.GLenum
	C.glGetActiveAttrib(
		p,
		C.GLuint(index),
		C.GLsizei(bufSize),
		nil,
		&cSize,
		&cType,
		(*C.GLchar)(buf),
	)
	return C.GoString((*C.char)(buf)), int(cSize), Enum(cType)
}

// GetActiveUniform returns details about an active uniform variable.
// A value of 0 for index selects the first active uniform variable.
// Permissible values for index range from 0 to the number of active
// uniform variables minus 1.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveUniform.xhtml
func (glw Wrapper) GetActiveUniform(p Program, index uint32) (string, int, Enum) {
	bufSize := GetProgrami(p, ACTIVE_UNIFORM_MAX_LENGTH)
	buf := C.malloc(C.size_t(bufSize))
	defer C.free(buf)

	var cSize C.GLint
	var cType C.GLenum

	C.glGetActiveUniform(
		p,
		C.GLuint(index),
		C.GLsizei(bufSize),
		nil,
		&cSize,
		&cType,
		(*C.GLchar)(buf),
	)
	return C.GoString((*C.char)(buf)), int(cSize), Enum(cType)
}

// GetAttachedShaders returns the shader objects attached to program p.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetAttachedShaders.xhtml
func (glw Wrapper) GetAttachedShaders(p Program) []Shader {
	panic("not implemented") // TODO: Implement
}

// GetAttribLocation returns the location of an attribute variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetAttribLocation.xhtml
func (glw Wrapper) GetAttribLocation(p Program, name string) Attrib {
	s, free := glStr(name)
	defer free()

	return Attrib(C.glGetAttribLocation(p, s))
}

// GetBooleanv returns the boolean values of parameter pname.
//
// Many boolean parameters can be queried more easily using IsEnabled.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func (glw Wrapper) GetBooleanv(dst []bool, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetFloatv returns the float values of parameter pname.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func (glw Wrapper) GetFloatv(dst []float32, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetIntegerv returns the int values of parameter pname.
//
// Single values may be queried more easily using GetInteger.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func (glw Wrapper) GetIntegerv(dst []int32, pname Enum) {
	buf := make([]C.GLint, len(dst))
	C.glGetIntegerv(pname, &buf[0])
	for i, v := range buf {
		dst[i] = int32(v)
	}
}

// GetInteger returns the int value of parameter pname.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func (glw Wrapper) GetInteger(pname Enum) int {
	panic("not implemented") // TODO: Implement
}

// GetBufferParameteri returns a parameter for the active buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetBufferParameter.xhtml
func (glw Wrapper) GetBufferParameteri(target Enum, value Enum) int {
	var params C.GLint
	C.glGetBufferParameteriv(target, value, &params)
	return int(params)
}

// GetError returns the next error.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetError.xhtml
func (glw Wrapper) GetError() Enum {
	return C.glGetError()
}

// GetFramebufferAttachmentParameteri returns attachment parameters
// for the active framebuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetFramebufferAttachmentParameteriv.xhtml
func (glw Wrapper) GetFramebufferAttachmentParameteri(target Enum, attachment Enum, pname Enum) int {
	panic("not implemented") // TODO: Implement
}

// GetProgrami returns a parameter value for a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramiv.xhtml
func (glw Wrapper) GetProgrami(p Program, pname Enum) int {
	var params C.GLint
	C.glGetProgramiv(p, pname, &params)
	return int(params)
}

// GetProgramInfoLog returns the information log for a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramInfoLog.xhtml
func (glw Wrapper) GetProgramInfoLog(p Program) string {
	infoLen := glw.GetProgrami(p, INFO_LOG_LENGTH)
	buf := C.malloc(C.size_t(infoLen))
	defer C.free(buf)
	C.glGetProgramInfoLog(p, C.GLsizei(infoLen), nil, (*C.GLchar)(buf))
	return C.GoString((*C.char)(buf))
}

// GetRenderbufferParameteri returns a parameter value for a render buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetRenderbufferParameteriv.xhtml
func (glw Wrapper) GetRenderbufferParameteri(target Enum, pname Enum) int {
	panic("not implemented") // TODO: Implement
}

// GetShaderi returns a parameter value for a shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderiv.xhtml
func (glw Wrapper) GetShaderi(s Shader, pname Enum) int {
	var params C.GLint
	C.glGetShaderiv(s, pname, &params)
	return int(params)
}

// GetShaderInfoLog returns the information log for a shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderInfoLog.xhtml
func (glw Wrapper) GetShaderInfoLog(s Shader) string {
	logLen := glw.GetShaderi(s, INFO_LOG_LENGTH)
	buf := C.malloc(C.size_t(logLen))
	defer C.free(buf)

	C.glGetShaderInfoLog(s, C.GLsizei(logLen), nil, (*C.GLchar)(buf))
	return C.GoString((*C.char)(buf))
}

// GetShaderPrecisionFormat returns range and precision limits for
// shader types.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderPrecisionFormat.xhtml
func (glw Wrapper) GetShaderPrecisionFormat(shadertype Enum, precisiontype Enum) (rangeLow int, rangeHigh int, precision int) {
	panic("not implemented") // TODO: Implement
}

// GetShaderSource returns source code of shader s.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderSource.xhtml
func (glw Wrapper) GetShaderSource(s Shader) string {
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
func (glw Wrapper) GetString(pname Enum) string {
	// Bounce through unsafe.Pointer, because on some platforms
	// GetString returns an *unsigned char which doesn't convert.
	return C.GoString((*C.char)((unsafe.Pointer)(C.glGetString(pname))))
}

// GetTexParameterfv returns the float values of a texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetTexParameter.xhtml
func (glw Wrapper) GetTexParameterfv(dst []float32, target Enum, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetTexParameteriv returns the int values of a texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetTexParameter.xhtml
func (glw Wrapper) GetTexParameteriv(dst []int32, target Enum, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetUniformfv returns the float values of a uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniform.xhtml
func (glw Wrapper) GetUniformfv(dst []float32, src Uniform, p Program) {
	panic("not implemented") // TODO: Implement
}

// GetUniformiv returns the float values of a uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniform.xhtml
func (glw Wrapper) GetUniformiv(dst []int32, src Uniform, p Program) {
	panic("not implemented") // TODO: Implement
}

// GetUniformLocation returns the location of a uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniformLocation.xhtml
func (glw Wrapper) GetUniformLocation(p Program, name string) Uniform {
	s, free := glStr(name)
	defer free()
	return Uniform(C.glGetUniformLocation(p, s))
}

// GetVertexAttribf reads the float value of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func (glw Wrapper) GetVertexAttribf(src Attrib, pname Enum) float32 {
	panic("not implemented") // TODO: Implement
}

// GetVertexAttribfv reads float values of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func (glw Wrapper) GetVertexAttribfv(dst []float32, src Attrib, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetVertexAttribi reads the int value of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func (glw Wrapper) GetVertexAttribi(src Attrib, pname Enum) int32 {
	panic("not implemented") // TODO: Implement
}

// GetVertexAttribiv reads int values of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func (glw Wrapper) GetVertexAttribiv(dst []int32, src Attrib, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// TODO(crawshaw): glGetVertexAttribPointerv

// Hint sets implementation-specific modes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glHint.xhtml
func (glw Wrapper) Hint(target Enum, mode Enum) {
	C.glHint(target, mode)
}

// IsBuffer reports if b is a valid buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsBuffer.xhtml
func (glw Wrapper) IsBuffer(b Buffer) bool {
	panic("not implemented") // TODO: Implement
}

// IsEnabled reports if cap is an enabled capability.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsEnabled.xhtml
func (glw Wrapper) IsEnabled(cap Enum) bool {
	panic("not implemented") // TODO: Implement
}

// IsFramebuffer reports if fb is a valid frame buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsFramebuffer.xhtml
func (glw Wrapper) IsFramebuffer(fb Framebuffer) bool {
	panic("not implemented") // TODO: Implement
}

// IsProgram reports if p is a valid program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsProgram.xhtml
func (glw Wrapper) IsProgram(p Program) bool {
	panic("not implemented") // TODO: Implement
}

// IsRenderbuffer reports if rb is a valid render buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsRenderbuffer.xhtml
func (glw Wrapper) IsRenderbuffer(rb Renderbuffer) bool {
	panic("not implemented") // TODO: Implement
}

// IsShader reports if s is valid shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsShader.xhtml
func (glw Wrapper) IsShader(s Shader) bool {
	panic("not implemented") // TODO: Implement
}

// IsTexture reports if t is a valid texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsTexture.xhtml
func (glw Wrapper) IsTexture(t Texture) bool {
	panic("not implemented") // TODO: Implement
}

// LineWidth specifies the width of lines.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glLineWidth.xhtml
func (glw Wrapper) LineWidth(width float32) {
	C.glLineWidth(C.GLfloat(width))
}

// LinkProgram links the specified program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glLinkProgram.xhtml
func (glw Wrapper) LinkProgram(p Program) {
	C.glLinkProgram(p)
}

// PixelStorei sets pixel storage parameters.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glPixelStorei.xhtml
func (glw Wrapper) PixelStorei(pname Enum, param int32) {
	C.glPixelStorei(pname, C.GLint(param))
}

// PolygonOffset sets the scaling factors for depth offsets.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glPolygonOffset.xhtml
func (glw Wrapper) PolygonOffset(factor float32, units float32) {
	panic("not implemented") // TODO: Implement
}

// ReadPixels returns pixel data from a buffer.
//
// In GLES 3, the source buffer is controlled with ReadBuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glReadPixels.xhtml
func (glw Wrapper) ReadPixels(dst []byte, x int, y int, width int, height int, format Enum, ty Enum) {
	panic("not implemented") // TODO: Implement
}

// ReleaseShaderCompiler frees resources allocated by the shader compiler.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glReleaseShaderCompiler.xhtml
func (glw Wrapper) ReleaseShaderCompiler() {
	panic("not implemented") // TODO: Implement
}

// RenderbufferStorage establishes the data storage, format, and
// dimensions of a renderbuffer object's image.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glRenderbufferStorage.xhtml
func (glw Wrapper) RenderbufferStorage(target, internalFormat Enum, width, height int) {
	C.glRenderbufferStorage(target, internalFormat, C.GLint(width), C.GLint(height))
}

// SampleCoverage sets multisample coverage parameters.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glSampleCoverage.xhtml
func (glw Wrapper) SampleCoverage(value float32, invert bool) {
	panic("not implemented") // TODO: Implement
}

// Scissor defines the scissor box rectangle, in window coordinates.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glScissor.xhtml
func (glw Wrapper) Scissor(x int32, y int32, width int32, height int32) {
	C.glScissor(C.GLint(x), C.GLint(y), C.GLint(width), C.GLint(height))
}

// TODO(crawshaw): ShaderBinary
// ShaderSource sets the source code of s to the given source code.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glShaderSource.xhtml
func (glw Wrapper) ShaderSource(s Shader, src string) {
	csources, free := glStr(src)
	defer free()
	C.glShaderSource(s, 1, &csources, nil)
}

// StencilFunc sets the front and back stencil test reference value.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilFunc.xhtml
func (glw Wrapper) StencilFunc(fn Enum, ref int, mask uint32) {
	panic("not implemented") // TODO: Implement
}

// StencilFunc sets the front or back stencil test reference value.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilFuncSeparate.xhtml
func (glw Wrapper) StencilFuncSeparate(face Enum, fn Enum, ref int, mask uint32) {
	panic("not implemented") // TODO: Implement
}

// StencilMask controls the writing of bits in the stencil planes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilMask.xhtml
func (glw Wrapper) StencilMask(mask uint32) {
	panic("not implemented") // TODO: Implement
}

// StencilMaskSeparate controls the writing of bits in the stencil planes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilMaskSeparate.xhtml
func (glw Wrapper) StencilMaskSeparate(face Enum, mask uint32) {
	panic("not implemented") // TODO: Implement
}

// StencilOp sets front and back stencil test actions.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilOp.xhtml
func (glw Wrapper) StencilOp(fail Enum, zfail Enum, zpass Enum) {
	panic("not implemented") // TODO: Implement
}

// StencilOpSeparate sets front or back stencil tests.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilOpSeparate.xhtml
func (glw Wrapper) StencilOpSeparate(face Enum, sfail Enum, dpfail Enum, dppass Enum) {
	panic("not implemented") // TODO: Implement
}

// TexImage2D writes a 2D texture image.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexImage2D.xhtml
func (glw Wrapper) TexImage2D(target Enum, level, internalFormat, width, height int, format, ty Enum, data []byte) {
	var ptr unsafe.Pointer
	if data != nil {
		ptr = unsafe.Pointer(&data[0])
	} else {
		ptr = unsafe.Pointer(nil)
	}

	C.glTexImage2D(
		target,
		C.GLint(level),
		C.GLint(internalFormat),
		C.GLint(width), C.GLint(height),
		0, // border
		format,
		ty,
		ptr,
	)
}

// TexSubImage2D writes a subregion of a 2D texture image.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexSubImage2D.xhtml
func (glw Wrapper) TexSubImage2D(target Enum, level int, x int, y int, width int, height int, format Enum, ty Enum, data []byte) {
	panic("not implemented") // TODO: Implement
}

// TexParameterf sets a float texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (glw Wrapper) TexParameterf(target Enum, pname Enum, param float32) {
	C.glTexParameterf(target, pname, C.GLfloat(param))
}

// TexParameterfv sets a float texture parameter array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (glw Wrapper) TexParameterfv(target Enum, pname Enum, params []float32) {
	// XXX: Check safety when we pass params
	C.glTexParameterfv(target, pname, (*C.GLfloat)(&params[0]))
}

// TexParameteri sets an integer texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (glw Wrapper) TexParameteri(target Enum, pname Enum, param int) {
	C.glTexParameteri(target, pname, C.GLint(param))
}

// TexParameteriv sets an integer texture parameter array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (glw Wrapper) TexParameteriv(target Enum, pname Enum, params []int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform1f writes a float uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform1f(dst Uniform, v float32) {
	C.glUniform1f(dst, C.GLfloat(v))
}

// Uniform1fv writes a [len(src)]float uniform array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform1fv(dst Uniform, src []float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform1i writes an int uniform variable.
//
// Uniform1i and Uniform1iv are the only two functions that may be used
// to load uniform variables defined as sampler types. Loading samplers
// with any other function will result in a INVALID_OPERATION error.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform1i(dst Uniform, v int) {
	C.glUniform1i(dst, C.GLint(v))
}

// Uniform1iv writes a int uniform array of len(src) elements.
//
// Uniform1i and Uniform1iv are the only two functions that may be used
// to load uniform variables defined as sampler types. Loading samplers
// with any other function will result in a INVALID_OPERATION error.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform1iv(dst Uniform, src []int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform2f writes a vec2 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform2f(dst Uniform, v0 float32, v1 float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform2fv writes a vec2 uniform array of len(src)/2 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform2fv(dst Uniform, src []float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform2i writes an ivec2 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform2i(dst Uniform, v0 int, v1 int) {
	panic("not implemented") // TODO: Implement
}

// Uniform2iv writes an ivec2 uniform array of len(src)/2 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform2iv(dst Uniform, src []int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform3f writes a vec3 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform3f(dst Uniform, v0, v1, v2 float32) {
	C.glUniform3f(dst, C.GLfloat(v0), C.GLfloat(v1), C.GLfloat(v2))
}

// Uniform3fv writes a vec3 uniform array of len(src)/3 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform3fv(dst Uniform, src []float32) {
	C.glUniform3fv(dst, 1, (*C.GLfloat)(&src[0]))
}

// Uniform3i writes an ivec3 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform3i(dst Uniform, v0 int32, v1 int32, v2 int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform3iv writes an ivec3 uniform array of len(src)/3 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform3iv(dst Uniform, src []int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform4f writes a vec4 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform4f(dst Uniform, v0 float32, v1 float32, v2 float32, v3 float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform4fv writes a vec4 uniform array of len(src)/4 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform4fv(dst Uniform, src []float32) {
	C.glUniform4fv(dst, 1, (*C.GLfloat)(&src[0]))
}

// Uniform4i writes an ivec4 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform4i(dst Uniform, v0 int32, v1 int32, v2 int32, v3 int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform4iv writes an ivec4 uniform array of len(src)/4 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) Uniform4iv(dst Uniform, src []int32) {
	panic("not implemented") // TODO: Implement
}

// UniformMatrix2fv writes 2x2 matrices. Each matrix uses four
// float32 values, so the number of matrices written is len(src)/4.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) UniformMatrix2fv(dst Uniform, src []float32) {
	panic("not implemented") // TODO: Implement
}

// UniformMatrix3fv writes 3x3 matrices. Each matrix uses nine
// float32 values, so the number of matrices written is len(src)/9.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) UniformMatrix3fv(dst Uniform, src []float32) {
	C.glUniformMatrix3fv(
		dst,
		1,
		glBoolean(false),
		(*C.GLfloat)(&src[0]),
	)
}

// UniformMatrix4fv writes 4x4 matrices. Each matrix uses 16
// float32 values, so the number of matrices written is len(src)/16.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (glw Wrapper) UniformMatrix4fv(dst Uniform, src []float32) {
	C.glUniformMatrix4fv(
		dst,
		1,
		glBoolean(false),
		(*C.GLfloat)(&src[0]),
	)
}

// UseProgram sets the active program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUseProgram.xhtml
func (glw Wrapper) UseProgram(p Program) {
	C.glUseProgram(p)
}

// ValidateProgram checks to see whether the executables contained in
// program can execute given the current OpenGL state.
//
// Typically only used for debugging.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glValidateProgram.xhtml
func (glw Wrapper) ValidateProgram(p Program) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib1f writes a float vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (glw Wrapper) VertexAttrib1f(dst Attrib, x float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib1fv writes a float vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (glw Wrapper) VertexAttrib1fv(dst Attrib, src []float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib2f writes a vec2 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (glw Wrapper) VertexAttrib2f(dst Attrib, x float32, y float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib2fv writes a vec2 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (glw Wrapper) VertexAttrib2fv(dst Attrib, src []float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib3f writes a vec3 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (glw Wrapper) VertexAttrib3f(dst Attrib, x float32, y float32, z float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib3fv writes a vec3 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (glw Wrapper) VertexAttrib3fv(dst Attrib, src []float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib4f writes a vec4 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (glw Wrapper) VertexAttrib4f(dst Attrib, x, y, z, w float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib4fv writes a vec4 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (glw Wrapper) VertexAttrib4fv(dst Attrib, src []float32) {
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
func (glw Wrapper) VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride int, offset int) {
	C.glVertexAttribPointer(
		dst,
		C.GLint(size),
		ty,
		glBoolean(normalized),
		C.GLint(stride),
		unsafe.Pointer(uintptr(offset)), // nolint: govet
	)
}

// Viewport sets the viewport, an affine transformation that
// normalizes device coordinates to window coordinates.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glViewport.xhtml
func (glw Wrapper) Viewport(x, y, width, height int) {
	C.glViewport(C.GLint(x), C.GLint(y), C.GLint(width), C.GLint(height))
}

// GetUniformBlockIndex retrieves the index of a uniform block within program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockIndex.xhtml
func (glw Wrapper) GetUniformBlockIndex(p Program, name string) uint32 {
	s, free := glStr(name)
	defer free()

	return uint32(C.glGetUniformBlockIndex(p, s))
}

// UniformBlockBinding assign a binding point to an active uniform block
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockBinding.xhtml
func (glw Wrapper) UniformBlockBinding(p Program, index, bind uint32) {
	C.glUniformBlockBinding(p, C.GLuint(index), C.GLuint(bind))
}

// BindBufferBase bind a buffer object to an indexed buffer target
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBufferBase.xhtml
func (glw Wrapper) BindBufferBase(target Enum, index uint32, b Buffer) {
	C.glBindBufferBase(target, C.GLuint(index), b)
}

func (glw Wrapper) GetActiveUniformBlockName(p Program, index uint32) string {
	bufSize := 256
	buf := C.malloc(C.size_t(bufSize))
	defer C.free(buf)
	var cSize C.GLint
	// var cType C.GLenum

	C.glGetActiveUniformBlockName(
		p,
		C.GLuint(index),
		C.GLsizei(bufSize),
		&cSize,
		(*C.GLchar)(buf),
	)

	return C.GoString((*C.char)(buf))
}

func (glw Wrapper) GetActiveUniformBlockiv(p Program, index uint32, pname Enum, params []int32) {
	buf := make([]C.GLint, len(params))

	C.glGetActiveUniformBlockiv(p, C.GLuint(index), pname, &buf[0])
	for i, v := range buf {
		params[i] = int32(v)
	}
}

func (glw Wrapper) GetActiveUniformi(p Program, index uint32, pname Enum) int32 {
	var params C.GLint
	idx := C.GLuint(index)
	C.glGetActiveUniformsiv(p, 1, &idx, pname, &params)
	return int32(params)
}

// DrawArraysInstanced draw multiple instances of a range of elements
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
func (glw Wrapper) DrawArraysInstanced(mode Enum, first, count, primcount uint32) {
	C.glDrawArraysInstanced(mode, C.GLint(first), C.GLint(count), C.GLint(primcount))
}

// DrawElementsInstanced — draw multiple instances of a set of elements
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawElementsInstanced.xhtml
func (glw Wrapper) DrawElementsInstanced(mode Enum, count uint32, typ Enum, offset, primcount uint32) {
	// off := unsafe.Pointer(uintptr(offset))
	C.glDrawElementsInstanced(
		mode,
		C.GLint(count),
		typ,
		unsafe.Pointer(uintptr(offset)), // nolint: govet
		C.GLint(primcount),
	)
}

// VertexAttribDivisor  modify the rate at which generic vertex attributes
// advance during instanced rendering
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
func (glw Wrapper) VertexAttribDivisor(index Attrib, divisor uint32) {
	C.glVertexAttribDivisor(index, C.GLuint(divisor))
}

// conv will convert a slice to a typedarray
//
//  Use js Copy bytes here and a temporary byte array + dataview
//
//  []float32 -> Float32Array
//  []float64 -> Float32Array (for glw Wrapper purposes)
func conv(data interface{}) (unsafe.Pointer, int) {
	switch v := data.(type) {
	case []byte:
		return unsafe.Pointer(&v[0]), len(v)
	case []uint16:
		return unsafe.Pointer(&v[0]), len(v) * 2
	case []uint32:
		return unsafe.Pointer(&v[0]), len(v) * 4
	case []float32:
		return unsafe.Pointer(&v[0]), len(v) * 4
	default:
		panic(fmt.Sprintf("Buffer type not implemented: %T", data))
	}
}

func glStr(s string) (*C.GLchar, func()) {
	zs := s + "\x00"
	str := unsafe.Pointer(C.CString(zs))
	free := func() { C.free(str) }
	return (*C.GLchar)(str), free
}

func glBoolean(v bool) C.GLboolean {
	if v {
		return TRUE
	}
	return FALSE
}
