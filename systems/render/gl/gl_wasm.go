//go:build js && wasm

package gl

// Wasm implementation not everything is implemented as I've been implementing
// on demand

import (
	"fmt"
	"log"
	"syscall/js"
)

var (
	// Null since this might differ between implementations
	NullFramebuffer = js.Null()
	NullBuffer      = js.Null()
	Null            = js.Null()
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

// IsValid returns if a gl value is valid or not
func IsValid(v js.Value) bool {
	// This might be wrong but can't do much now
	return v.Truthy()
}

// Common to move things to gl
var (

	// buf experiment, avoid realocation on transfers
	maxTransferBuf = 1 << 21 // 2Mb
	tbuf           = js.Global().Get("Uint8Array").New(maxTransferBuf)

	// Mostly used on uniforms
	sbuf = js.Global().Get("Uint8Array").New(64)
	// to satisfy js gl bindings we use Float32Array views with same underlying buffer
	f2buf  = js.Global().Get("Float32Array").New(sbuf.Get("buffer"), 0, 2)  // vec2
	f3buf  = js.Global().Get("Float32Array").New(sbuf.Get("buffer"), 0, 3)  // vec3
	f4buf  = js.Global().Get("Float32Array").New(sbuf.Get("buffer"), 0, 4)  // vec4
	f9buf  = js.Global().Get("Float32Array").New(sbuf.Get("buffer"), 0, 9)  // mat3 3x3
	f16buf = js.Global().Get("Float32Array").New(sbuf.Get("buffer"), 0, 16) // mat3 4x4
)

// Wrapper exposes the methods
type Wrapper struct {
	js.Value
}

func (Wrapper) String() string {
	return "renderer/gl/gl_wasm.go"
}
func (Wrapper) Impl() string { return "wasm" }

var _ Context3 = &Wrapper{}

// GetWebGL Return a js.Value Wrapper gl context
func GetWebGL(v js.Value) Wrapper {
	return Wrapper{v}
}

// ActiveTexture sets the active texture unit.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glActiveTexture.xhtml
func (g Wrapper) ActiveTexture(texture Enum) {
	g.Call("activeTexture", texture)
}

// AttachShader attaches a shader to a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glAttachShader.xhtml
func (g Wrapper) AttachShader(p Program, s Shader) {
	g.Call("attachShader", p, s)
}

// BindAttribLocation binds a vertex attribute index with a named
// variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindAttribLocation.xhtml
func (g Wrapper) BindAttribLocation(p Program, a Attrib, name string) {
	g.Call("bindAttribLocation", p, a, name)
}

// BindBuffer binds a buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBuffer.xhtml
func (g Wrapper) BindBuffer(target Enum, b Buffer) {
	g.Call("bindBuffer", target, b)
}

// BindFramebuffer binds a framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindFramebuffer.xhtml
func (g Wrapper) BindFramebuffer(target Enum, fb Framebuffer) {
	g.Call("bindFramebuffer", target, fb)
}

// BindRenderbuffer binds a render buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindRenderbuffer.xhtml
func (g Wrapper) BindRenderbuffer(target Enum, rb Renderbuffer) {
	g.Call("bindRenderbuffer", target, rb)
}

// BindTexture binds a texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindTexture.xhtml
func (g Wrapper) BindTexture(target Enum, t Texture) {
	g.Call("bindTexture", target, t)
}

// BindVertexArray binds a vertex array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindVertexArray.xhtml
func (g Wrapper) BindVertexArray(rb VertexArray) {
	g.Call("bindVertexArray", rb)
}

// BlendColor sets the blend color.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendColor.xhtml
func (g Wrapper) BlendColor(red, green, blue, alpha float32) {
	g.Call("blendColor", red, green, blue, alpha)
}

// BlendEquation sets both RGB and alpha blend equations.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendEquation.xhtml
func (g Wrapper) BlendEquation(mode Enum) {
	g.Call("blendEquation", mode)
}

// BlendEquationSeparate sets RGB and alpha blend equations separately.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendEquationSeparate.xhtml
func (g Wrapper) BlendEquationSeparate(modeRGB, modeAlpha Enum) {
	g.Call("blendEquationSeparate", modeRGB, modeAlpha)
}

// BlendFunc sets the pixel blending factors.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendFunc.xhtml
func (g Wrapper) BlendFunc(sfactor, dfactor Enum) {
	g.Call("blendFunc", sfactor, dfactor)
}

func (g Wrapper) BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha Enum) {
	panic("not implemented")
}

func (g Wrapper) BufferInit(target Enum, size int, usage Enum) {
	if size > maxTransferBuf && usage == DYNAMIC_DRAW {
		log.Printf("transfer buffer growa %d -> %d", maxTransferBuf, size)
		tbuf = js.Global().Get("Uint8Array").New(size)
		maxTransferBuf = size
	}
	g.Call("bufferData", target, size, usage)
}

// BufferData will type switch the interface and select the proper type
func (g Wrapper) BufferData(target Enum, data any, usage Enum) {
	jsData, sz := conv(data)
	// WebGL2:
	// void gl.bufferData(target, ArrayBufferView srcData, usage, srcOffset, length);
	g.Call("bufferData", target, jsData, usage, 0, sz)
}

// BufferSubData same as before but with extra step to check data type on any
func (g Wrapper) BufferSubData(target Enum, offset int, data any) {
	jsData, sz := conv(data)
	// WebGL2:
	// void gl.bufferSubData(target, dstByteOffset, ArrayBufferView srcData, srcOffset, length);
	g.Call("bufferSubData", target, offset, jsData, 0, sz)
}

// CheckFramebufferStatus reports the completeness status of the
// active framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCheckFramebufferStatus.xhtml
func (g Wrapper) CheckFramebufferStatus(target Enum) Enum {
	return Enum(g.Call("checkFramebufferStatus", target).Int())
}

// Clear clears the window.
//
// The behavior of Clear is influenced by the pixel ownership test,
// the scissor test, dithering, and the buffer writemasks.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClear.xhtml
func (g Wrapper) Clear(mask Enum) {
	g.Call("clear", mask)
}

// ClearColor specifies the RGBA values used to clear color buffers.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearColor.xhtml
func (g Wrapper) ClearColor(red, green, blue, alpha float32) {
	g.Call("clearColor", red, green, blue, alpha)
}

// ClearDepthf sets the depth value used to clear the depth buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearDepthf.xhtml
func (g Wrapper) ClearDepthf(d float32) {
	g.Call("clearDepth", d)
}

func (g Wrapper) ClearStencil(s int) {
	g.Call("clearStencil", s)
}

func (g Wrapper) ColorMask(red, green, blue, alpha bool) {
	g.Call("colorMask", red, green, blue, alpha)
}

// CompileShader compiles the source code of s.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompileShader.xhtml
func (g Wrapper) CompileShader(s Shader) {
	g.Call("compileShader", s)
}

func (g Wrapper) CompressedTexImage2D(target Enum, level int, internalformat Enum, width, height, border int, data []byte) {
	panic("not implemented")
}

func (g Wrapper) CompressedTexSubImage2D(target Enum, level, xoffset, yoffset, width, height int, format Enum, data []byte) {
	panic("not implemented")
}

func (g Wrapper) CopyTexImage2D(target Enum, level int, internalformat Enum, x, y, width, height, border int) {
	panic("not implemented")
}

func (g Wrapper) CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int) {
	panic("not implemented")
}

// CreateBuffer creates a buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenBuffers.xhtml
func (g Wrapper) CreateBuffer() Buffer {
	return Buffer(g.Call("createBuffer"))
}

// CreateFramebuffer creates a framebuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenFramebuffers.xhtml
func (g Wrapper) CreateFramebuffer() Framebuffer {
	return Framebuffer(g.Call("createFramebuffer"))
}

// CreateProgram creates a new empty program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateProgram.xhtml
func (g Wrapper) CreateProgram() Program {
	return Program(g.Call("createProgram"))
}

func (g Wrapper) CreateRenderbuffer() Renderbuffer {
	return Renderbuffer(g.Call("createRenderbuffer"))
}

// CreateShader creates a new empty shader object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateShader.xhtml
func (g Wrapper) CreateShader(ty Enum) Shader {
	return Shader(g.Call("createShader", ty))
}

// CreateTexture creates a texture object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenTextures.xhtml
func (g Wrapper) CreateTexture() Texture {
	return Texture(g.Call("createTexture"))
}

// CreateVertexArray creates a vertex array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenVertexArrays.xhtml
func (g Wrapper) CreateVertexArray() VertexArray {
	return VertexArray(g.Call("createVertexArray"))
}

// CullFace specifies which polygons are candidates for culling.
//
// Valid modes: FRONT, BACK, FRONT_AND_BACK.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCullFace.xhtml
func (g Wrapper) CullFace(mode Enum) {
	g.Call("cullFace", mode)
}

// DeleteBuffer deletes the given buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteBuffers.xhtml
func (g Wrapper) DeleteBuffer(v Buffer) {
	g.Call("deleteBuffer", v)
}

// DeleteFramebuffer deletes the given framebuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteFramebuffers.xhtml
func (g Wrapper) DeleteFramebuffer(v Framebuffer) {
	g.Call("deleteFramebuffer", v)
}

// DeleteProgram deletes the given program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteProgram.xhtml
func (g Wrapper) DeleteProgram(p Program) {
	g.Call("deleteProgram", p)
}

func (g Wrapper) DeleteRenderbuffer(v Renderbuffer) {
	panic("not implemented")
}

func (g Wrapper) DeleteShader(s Shader) {
	g.Call("deleteShader", s)
}

// DeleteTexture deletes the given texture object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteTextures.xhtml
func (g Wrapper) DeleteTexture(v Texture) {
	g.Call("deleteTexture", v)
}

// DeleteVertexArray deletes the given render buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteVertexArrays.xhtml
func (g Wrapper) DeleteVertexArray(v VertexArray) {
	g.Call("deleteVertexArray", v)
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
func (g Wrapper) DepthFunc(fn Enum) {
	g.Call("depthFunc", fn)
}

// DepthMask sets the depth buffer enabled for writing.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDepthMask.xhtml
func (g Wrapper) DepthMask(flag bool) {
	g.Call("depthMask", flag)
}

func (g Wrapper) DepthRangef(n, f float32) {
	panic("not implemented")
}

func (g Wrapper) DetachShader(p Program, s Shader) {
	panic("not implemented")
}

// Disable disables various GL capabilities.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDisable.xhtml
func (g Wrapper) Disable(cap Enum) {
	g.Call("disable", cap)
}

// DisableVertexAttribArray disables a vertex attribute array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDisableVertexAttribArray.xhtml
func (g Wrapper) DisableVertexAttribArray(a Attrib) {
	g.Call("disableVertexAttribArray", a)
}

// DrawArrays renders geometric primitives from the bound data.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArrays.xhtml
func (g Wrapper) DrawArrays(mode Enum, first, count int) {
	g.Call("drawArrays", mode, first, count)
}

// DrawElements renders primitives from a bound buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawElements.xhtml
func (g Wrapper) DrawElements(mode Enum, count int, ty Enum, offset int) {
	g.Call("drawElements", mode, count, ty, offset)
}

// Enable enables various GL capabilities.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnable.xhtml
func (g Wrapper) Enable(cp Enum) {
	g.Call("enable", cp)
}

// EnableVertexAttribArray enables a vertex attribute array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnableVertexAttribArray.xhtml
func (g Wrapper) EnableVertexAttribArray(a Attrib) {
	g.Call("enableVertexAttribArray", a)
}

// Finish blocks until the effects of all previously called GL
// commands are complete.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFinish.xhtml
func (g Wrapper) Finish() {
	g.Call("finish")
}

// Flush empties all buffers. It does not block.
//
// An OpenGL implementation may buffer network communication,
// the command stream, or data inside the graphics accelerator.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFlush.xhtml
func (g Wrapper) Flush() {
	g.Call("flush")
}

func (g Wrapper) FramebufferRenderbuffer(target, attachment, rbTarget Enum, rb Renderbuffer) {
	g.Call("framebufferRenderbuffer", target, attachment, rbTarget, rb)
}

// FramebufferTexture2D attaches the t to the current frame buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFramebufferTexture2D.xhtml
func (g Wrapper) FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	g.Call("framebufferTexture2D", target, attachment, texTarget, t, level)
}

// FrontFace defines which polygons are front-facing.
//
// Valid modes: CW, CCW.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFrontFace.xhtml
func (g Wrapper) FrontFace(mode Enum) {
	g.Call("frontFace", mode)
}

// GenerateMipmap generates mipmaps for the current texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenerateMipmap.xhtml
func (g Wrapper) GenerateMipmap(target Enum) {
	g.Call("generateMipmap", target)
}

// GetActiveAttrib returns details about an active attribute variable.
// A value of 0 for index selects the first active attribute variable.
// Permissible values for index range from 0 to the number of active
// attribute variables minus 1.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveAttrib.xhtml
func (g Wrapper) GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum) {
	res := g.Call("getActiveAttrib", p, index)

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
func (g Wrapper) GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum) {
	res := g.Call("getActiveUniform", p, index)

	name = res.Get("name").String()
	size = res.Get("size").Int()
	ty = Enum(res.Get("type").Int())

	return name, size, ty
}

func (g Wrapper) GetAttachedShaders(p Program) []Shader {
	panic("not implemented")
}

// GetAttribLocation returns the location of an attribute variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetAttribLocation.xhtml
func (g Wrapper) GetAttribLocation(p Program, name string) Attrib {
	return Attrib(g.Call("getAttribLocation", p, name).Int())
}

func (g Wrapper) GetBooleanv(dst []bool, pname Enum) {
	panic("not implemented")
}

func (g Wrapper) GetFloatv(dst []float32, pname Enum) {
	panic("not implemented")
}

func (g Wrapper) GetIntegerv(dst []int32, pname Enum) {
	panic("not implemented")
}

func (g Wrapper) GetInteger(pname Enum) int {
	panic("not implemented")
}

func (g Wrapper) GetBufferParameteri(target, value Enum) int {
	return g.Call("getBufferParameter", target, value).Int()
}

// GetError returns the next error.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetError.xhtml
func (g Wrapper) GetError() Enum {
	return Enum(g.Call("getError").Int())
}

func (g Wrapper) GetFramebufferAttachmentParameteri(target, attachment, pname Enum) int {
	panic("not implemented")
}

// GetProgrami returns a parameter value for a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramiv.xhtml
func (g Wrapper) GetProgrami(p Program, pname Enum) int {
	r := g.Call("getProgramParameter", p, pname)

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
func (g Wrapper) GetProgramInfoLog(p Program) string {
	return g.Call("getProgramInfoLog", p).String()
}

func (g Wrapper) GetRenderbufferParameteri(target, pname Enum) int {
	panic("not implemented")
}

// GetShaderi returns a parameter value for a shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderiv.xhtml
func (g Wrapper) GetShaderi(s Shader, pname Enum) int {
	if g.Call("getShaderParameter", s, pname).Bool() {
		return 1
	}
	return 0
}

// GetShaderInfoLog returns the information log for a shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderInfoLog.xhtml
func (g Wrapper) GetShaderInfoLog(s Shader) string {
	return g.Call("getShaderInfoLog", s).String()
}

func (g Wrapper) GetShaderPrecisionFormat(shadertype, precisiontype Enum) (rangeLow, rangeHigh, precision int) {
	panic("not implemented")
}

func (g Wrapper) GetShaderSource(s Shader) string {
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
func (g Wrapper) GetString(pname Enum) string {
	return g.Call("getParameter", pname).String()
}

func (g Wrapper) GetTexParameterfv(dst []float32, target, pname Enum) {
	panic("not implemented")
}

func (g Wrapper) GetTexParameteriv(dst []int32, target, pname Enum) {
	panic("not implemented")
}

// GetUniformfv returns the float values of a uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniform.xhtml
func (g Wrapper) GetUniformfv(dst []float32, src Uniform, p Program) {
	g.Call("getUniformfv", dst, src, p)
}

func (g Wrapper) GetUniformiv(dst []int32, src Uniform, p Program) {
	panic("not implemented")
}

func (g Wrapper) GetUniformLocation(p Program, name string) Uniform {
	return Uniform(g.Call("getUniformLocation", p, name))
}

func (g Wrapper) GetVertexAttribf(src Attrib, pname Enum) float32 {
	panic("not implemented")
}

func (g Wrapper) GetVertexAttribfv(dst []float32, src Attrib, pname Enum) {
	panic("not implemented")
}

func (g Wrapper) GetVertexAttribi(src Attrib, pname Enum) int32 {
	panic("not implemented")
}

func (g Wrapper) GetVertexAttribiv(dst []int32, src Attrib, pname Enum) {
	panic("not implemented")
}

// Hint sets implementation-specific modes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glHint.xhtml
func (g Wrapper) Hint(target, mode Enum) {
	g.Call("hint", target, mode)
}

func (g Wrapper) IsBuffer(b Buffer) bool {
	panic("not implemented")
}

func (g Wrapper) IsEnabled(cap Enum) bool {
	panic("not implemented")
}

func (g Wrapper) IsFramebuffer(fb Framebuffer) bool {
	panic("not implemented")
}

func (g Wrapper) IsProgram(p Program) bool {
	panic("not implemented")
}

func (g Wrapper) IsRenderbuffer(rb Renderbuffer) bool {
	panic("not implemented")
}

func (g Wrapper) IsShader(s Shader) bool {
	panic("not implemented")
}

func (g Wrapper) IsTexture(t Texture) bool {
	panic("not implemented")
}

// LineWidth specifies the width of lines.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glLineWidth.xhtml
func (g Wrapper) LineWidth(width float32) {
	g.Call("lineWidth", width)
}

// LinkProgram links the specified program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glLinkProgram.xhtml
func (g Wrapper) LinkProgram(p Program) {
	g.Call("linkProgram", p)
}

func (g Wrapper) PixelStorei(pname Enum, param int32) {
	g.Call("pixelStorei", pname, param)
}

func (g Wrapper) PolygonOffset(factor, units float32) {
	panic("not implemented")
}

func (g Wrapper) ReadPixels(dst []byte, x, y, width, height int, format, ty Enum) {
	panic("not implemented")
}

func (g Wrapper) ReleaseShaderCompiler() {
	panic("not implemented")
}

func (g Wrapper) RenderbufferStorage(target, internalFormat Enum, width, height int) {
	g.Call("renderbufferStorage", target, internalFormat, width, height)
}

func (g Wrapper) SampleCoverage(value float32, invert bool) {
	panic("not implemented")
}

func (g Wrapper) Scissor(x, y, width, height int32) {
	g.Call("scissor", x, y, width, height)
}

// ShaderSource sets the source code of s to the given source code.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glShaderSource.xhtml
func (g Wrapper) ShaderSource(s Shader, src string) {
	g.Call("shaderSource", s, src)
}

func (g Wrapper) StencilFunc(fn Enum, ref int, mask uint32) {
	g.Call("stencilFunc", fn, ref, mask)
}

func (g Wrapper) StencilFuncSeparate(face, fn Enum, ref int, mask uint32) {
	panic("not implemented")
}

func (g Wrapper) StencilMask(mask uint32) {
	g.Call("stencilMask", mask)
}

func (g Wrapper) StencilMaskSeparate(face Enum, mask uint32) {
	panic("not implemented")
}

func (g Wrapper) StencilOp(fail, zfail, zpass Enum) {
	g.Call("stencilOp", fail, zfail, zpass)
}

func (g Wrapper) StencilOpSeparate(face, sfail, dpfail, dppass Enum) {
	panic("not implemented")
}

// TexImage2D writes a 2D texture image.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexImage2D.xhtml
func (g Wrapper) TexImage2D(target Enum, level int, internalFormat int, width, height int, format Enum, ty Enum, data []byte) {
	jsBuf := js.Null()
	if data != nil {
		// Might leak
		jsBuf = js.Global().Get("Uint8Array").New(len(data))
		js.CopyBytesToJS(jsBuf, data)

		if ty == FLOAT {
			jsBuf = js.Global().Get("Float32Array").New(jsBuf.Get("buffer"))
		}
	}
	g.Call("texImage2D",
		target, level, internalFormat,
		width, height, 0,
		format, ty, jsBuf,
	)
}

func (g Wrapper) TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data []byte) {
	panic("not implemented")
}

// TexParameterf sets a float texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (g Wrapper) TexParameterf(target, pname Enum, param float32) {
	g.Call("texParameterf", target, pname, param)
}

func (g Wrapper) TexParameterfv(target, pname Enum, params []float32) {
	g.Call("texParameterfv", target, pname, params)
}

// TexParameteri sets an integer texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (g Wrapper) TexParameteri(target, pname Enum, param int) {
	g.Call("texParameteri", target, pname, param)
}

func (g Wrapper) TexParameteriv(target, pname Enum, params []int32) {
	panic("not implemented")
}

// Uniform1f writes a float uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g Wrapper) Uniform1f(dst Uniform, v float32) {
	g.Call("uniform1f", dst, v)
}

func (g Wrapper) Uniform1fv(dst Uniform, src []float32) {
	panic("not implemented")
}

// Uniform1i writes an int uniform variable.
//
// Uniform1i and Uniform1iv are the only two functions that may be used
// to load uniform variables defined as sampler types. Loading samplers
// with any other function will result in a INVALID_OPERATION error.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g Wrapper) Uniform1i(dst Uniform, v int) {
	g.Call("uniform1i", dst, v)
}

func (g Wrapper) Uniform1iv(dst Uniform, src []int32) {
	panic("not implemented")
}

// Uniform2f writes a vec2 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g Wrapper) Uniform2f(dst Uniform, v0, v1 float32) {
	g.Call("uniform2f", dst, v0, v1)
}

// Uniform2fv writes a vec2 uniform array of len(src)/2 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g Wrapper) Uniform2fv(dst Uniform, src []float32) {
	js.CopyBytesToJS(sbuf, f32bytes(src...))
	g.Call("uniform2fv", dst, f2buf)
}

func (g Wrapper) Uniform2i(dst Uniform, v0, v1 int) {
	panic("not implemented")
}

func (g Wrapper) Uniform2iv(dst Uniform, src []int32) {
	panic("not implemented")
}

// Uniform3f writes a vec3 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g Wrapper) Uniform3f(dst Uniform, v0, v1, v2 float32) {
	g.Call("uniform3f", dst, v0, v1, v2)
}

// Uniform3fv writes a vec3 uniform array of len(src)/3 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g Wrapper) Uniform3fv(dst Uniform, src []float32) {
	js.CopyBytesToJS(sbuf, f32bytes(src...))
	g.Call("uniform3fv", dst, f3buf)
}

func (g Wrapper) Uniform3i(dst Uniform, v0, v1, v2 int32) {
	panic("not implemented")
}

func (g Wrapper) Uniform3iv(dst Uniform, src []int32) {
	panic("not implemented")
}

// Uniform4f writes a vec4 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g Wrapper) Uniform4f(dst Uniform, v0, v1, v2, v3 float32) {
	g.Call("uniform4f", dst, v0, v1, v2, v3)
}

// Uniform4fv writes a vec4 uniform array of len(src)/4 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g Wrapper) Uniform4fv(dst Uniform, src []float32) {
	js.CopyBytesToJS(sbuf, f32bytes(src...))
	g.Call("uniform4fv", dst, f4buf)
}

func (g Wrapper) Uniform4i(dst Uniform, v0, v1, v2, v3 int32) {
	panic("not implemented")
}

func (g Wrapper) Uniform4iv(dst Uniform, src []int32) {
	panic("not implemented")
}

func (g Wrapper) UniformMatrix2fv(dst Uniform, src []float32) {
	panic("not implemented")
}

// UniformMatrix3fv writes 3x3 matrices. Each matrix uses nine
// float32 values, so the number of matrices written is len(src)/9.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g Wrapper) UniformMatrix3fv(dst Uniform, src []float32) {
	js.CopyBytesToJS(sbuf, f32bytes(src...))
	g.Call("uniformMatrix3fv", dst, false, f9buf)
}

// UniformMatrix4fv writes 4x4 matrices. Each matrix uses 16
// float32 values, so the number of matrices written is len(src)/16.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (g Wrapper) UniformMatrix4fv(dst Uniform, src []float32) {
	js.CopyBytesToJS(sbuf, f32bytes(src...))
	g.Call("uniformMatrix4fv", dst, false, f16buf)
}

// UseProgram sets the active program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUseProgram.xhtml
func (g Wrapper) UseProgram(p Program) {
	g.Call("useProgram", p)
}

func (g Wrapper) ValidateProgram(p Program) {
	panic("not implemented")
}

func (g Wrapper) VertexAttrib1f(dst Attrib, x float32) {
	panic("not implemented")
}

func (g Wrapper) VertexAttrib1fv(dst Attrib, src []float32) {
	panic("not implemented")
}

func (g Wrapper) VertexAttrib2f(dst Attrib, x, y float32) {
	panic("not implemented")
}

func (g Wrapper) VertexAttrib2fv(dst Attrib, src []float32) {
	panic("not implemented")
}

func (g Wrapper) VertexAttrib3f(dst Attrib, x, y, z float32) {
	panic("not implemented")
}

func (g Wrapper) VertexAttrib3fv(dst Attrib, src []float32) {
	panic("not implemented")
}

func (g Wrapper) VertexAttrib4f(dst Attrib, x, y, z, w float32) {
	panic("not implemented")
}

func (g Wrapper) VertexAttrib4fv(dst Attrib, src []float32) {
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
func (g Wrapper) VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride, offset int) {
	g.Call("vertexAttribPointer", dst, size, ty, normalized, stride, offset)
}

// Viewport sets the viewport, an affine transformation that
// normalizes device coordinates to window coordinates.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glViewport.xhtml
func (g Wrapper) Viewport(x, y, width, height int) {
	g.Call("viewport", x, y, width, height)
}

///////////////////////////////////////////////////////////////////////////////
// Wrapper 2 + extra funcs
///////////////////////////////////////////////////////////////////////////////

// GetUniformBlockIndex retrieves the index of a uniform block within program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockIndex.xhtml
func (g Wrapper) GetUniformBlockIndex(p Program, name string) uint32 {
	return uint32(g.Call("getUniformBlockIndex", p, name).Int())
}

// UniformBlockBinding assign a binding point to an active uniform block
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockBinding.xhtml
func (g Wrapper) UniformBlockBinding(p Program, index, bind uint32) {
	g.Call("uniformBlockBinding", p, index, bind)
}

// BindBufferBase bind a buffer object to an indexed buffer target
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBufferBase.xhtml
func (g Wrapper) BindBufferBase(target Enum, n uint32, b Buffer) {
	g.Call("bindBufferBase", target, n, b)
}

func (g Wrapper) GetActiveUniformBlockName(p Program, index uint32) string {
	return g.Call("getActiveUniformBlockName", p, index).String()
}

func (g Wrapper) GetActiveUniformBlockiv(p Program, index uint32, pname Enum, param []int32) {
	ret := g.Call("getActiveUniformBlockParameter", p, index, pname)

	if ret.Type() == js.TypeNumber {
		param[0] = int32(ret.Int())
		return
	}

	sz := ret.Length()

	for i := range param {
		if i > sz {
			return
		}
		param[i] = int32(ret.Index(i).Int())
	}
}

// https://www.khronos.org/registry/OpenGL-Refpages/es3.0/html/glGetActiveUniformsiv.xhtml
func (g Wrapper) GetActiveUniformi(p Program, index uint32, pname Enum) int32 {
	ret := g.Call("getActiveUniforms", p, []any{index}, pname)
	return int32(ret.Index(0).Int())
}

// DrawArraysInstanced draw multiple instances of a range of elements
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
func (g Wrapper) DrawArraysInstanced(mode Enum, first, count, primcount uint32) {
	g.Call("drawArraysInstanced", mode, first, count, primcount)
}

// VertexAttribDivisor  modify the rate at which generic vertex attributes
// advance during instanced rendering
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
func (g Wrapper) VertexAttribDivisor(index Attrib, divisor uint32) {
	g.Call("vertexAttribDivisor", index, divisor)
}

// DrawElementsInstanced — draw multiple instances of a set of elements
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawElementsInstanced.xhtml
func (g Wrapper) DrawElementsInstanced(mode Enum, count uint32, typ Enum, offset, primcount uint32) {
	g.Call("drawElementsInstanced", mode, count, typ, offset, primcount)
}

func (g *Wrapper) TexImage3D(
	target Enum,
	level, internalFormat, width, height, depth int,
	format, ty Enum,
	data []byte,
) {
	jsBuf := js.Null()
	if data != nil {
		// Might leak
		jsBuf = js.Global().Get("Uint8Array").New(len(data))
		js.CopyBytesToJS(jsBuf, data)

		if ty == FLOAT {
			jsBuf = js.Global().Get("Float32Array").New(jsBuf.Get("buffer"))
		}
	}
	g.Call("texImage3D",
		target, level, internalFormat,
		width, height, depth, 0,
		format, ty, jsBuf,
	)
}

func (g *Wrapper) FramebufferTextureLayer(
	target, attachment Enum,
	texture Texture,
	level, layer int,
) {
	g.Call("framebufferTextureLayer",
		target,
		attachment,
		texture,
		level,
		layer,
	)
}

// conv will convert a slice to a typedarray
//
//  Use js Copy bytes here and a temporary byte array + dataview
//
//  []float32 -> Float32Array
//  []float64 -> Float32Array (for Wrapper purposes)
func conv(data any) (js.Value, int) {
	var bdata []byte
	switch data := data.(type) {
	case js.Value:
		sz := data.Get("length").Int()
		return data, sz
	case []float32:
		bdata = f32bytes(data...)
	case []uint32:
		bdata = u32bytes(data...)
	case []uint16:
		bdata = u16bytes(data...)
	case []byte:
		bdata = data
	default:
		panic(fmt.Errorf("unimplemented type: %T", data))
	}
	sz := len(bdata)
	if sz > maxTransferBuf {
		buf := js.Global().Get("Uint8Array").New(sz)
		js.CopyBytesToJS(buf, bdata)
		return buf, sz
	}
	js.CopyBytesToJS(tbuf, bdata)
	return tbuf, sz
}

/*func newFloat32Array(sz int) js.Value {
	return js.Global().Get("Float32Array").New(sz)
}*/
