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

package gl

type undef struct{}

// ActiveTexture sets the active texture unit.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glActiveTexture.xhtml
func (u undef) ActiveTexture(texture Enum) {
	panic("not implemented") // TODO: Implement
}

// AttachShader attaches a shader to a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glAttachShader.xhtml
func (u undef) AttachShader(p Program, s Shader) {
	panic("not implemented") // TODO: Implement
}

// BindAttribLocation binds a vertex attribute index with a named
// variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindAttribLocation.xhtml
func (u undef) BindAttribLocation(p Program, a Attrib, name string) {
	panic("not implemented") // TODO: Implement
}

// BindBuffer binds a buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBuffer.xhtml
func (u undef) BindBuffer(target Enum, b Buffer) {
	panic("not implemented") // TODO: Implement
}

// BindFramebuffer binds a framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindFramebuffer.xhtml
func (u undef) BindFramebuffer(target Enum, fb Framebuffer) {
	panic("not implemented") // TODO: Implement
}

// BindRenderbuffer binds a render buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindRenderbuffer.xhtml
func (u undef) BindRenderbuffer(target Enum, rb Renderbuffer) {
	panic("not implemented") // TODO: Implement
}

// BindTexture binds a texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindTexture.xhtml
func (u undef) BindTexture(target Enum, t Texture) {
	panic("not implemented") // TODO: Implement
}

// BindVertexArray binds a vertex array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindVertexArray.xhtml
func (u undef) BindVertexArray(rb VertexArray) {
	panic("not implemented") // TODO: Implement
}

// BlendColor sets the blend color.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendColor.xhtml
func (u undef) BlendColor(red float32, green float32, blue float32, alpha float32) {
	panic("not implemented") // TODO: Implement
}

// BlendEquation sets both RGB and alpha blend equations.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendEquation.xhtml
func (u undef) BlendEquation(mode Enum) {
	panic("not implemented") // TODO: Implement
}

// BlendEquationSeparate sets RGB and alpha blend equations separately.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendEquationSeparate.xhtml
func (u undef) BlendEquationSeparate(modeRGB Enum, modeAlpha Enum) {
	panic("not implemented") // TODO: Implement
}

// BlendFunc sets the pixel blending factors.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendFunc.xhtml
func (u undef) BlendFunc(sfactor Enum, dfactor Enum) {
	panic("not implemented") // TODO: Implement
}

// BlendFunc sets the pixel RGB and alpha blending factors separately.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendFuncSeparate.xhtml
func (u undef) BlendFuncSeparate(sfactorRGB Enum, dfactorRGB Enum, sfactorAlpha Enum, dfactorAlpha Enum) {
	panic("not implemented") // TODO: Implement
}

// BufferData creates a new data store for the bound buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferData.xhtml
func (u undef) BufferData(target Enum, src []byte, usage Enum) {
	panic("not implemented") // TODO: Implement
}

// BufferInit creates a new uninitialized data store for the bound buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferData.xhtml
func (u undef) BufferInit(target Enum, size int, usage Enum) {
	panic("not implemented") // TODO: Implement
}

// BufferSubData sets some of data in the bound buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferSubData.xhtml
func (u undef) BufferSubData(target Enum, offset int, data []byte) {
	panic("not implemented") // TODO: Implement
}

// CheckFramebufferStatus reports the completeness status of the
// active framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCheckFramebufferStatus.xhtml
func (u undef) CheckFramebufferStatus(target Enum) Enum {
	panic("not implemented") // TODO: Implement
}

// Clear clears the window.
//
// The behavior of Clear is influenced by the pixel ownership test,
// the scissor test, dithering, and the buffer writemasks.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClear.xhtml
func (u undef) Clear(mask Enum) {
	panic("not implemented") // TODO: Implement
}

// ClearColor specifies the RGBA values used to clear color buffers.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearColor.xhtml
func (u undef) ClearColor(red float32, green float32, blue float32, alpha float32) {
	panic("not implemented") // TODO: Implement
}

// ClearDepthf sets the depth value used to clear the depth buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearDepthf.xhtml
func (u undef) ClearDepthf(d float32) {
	panic("not implemented") // TODO: Implement
}

// ClearStencil sets the index used to clear the stencil buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearStencil.xhtml
func (u undef) ClearStencil(s int) {
	panic("not implemented") // TODO: Implement
}

// ColorMask specifies whether color components in the framebuffer
// can be written.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glColorMask.xhtml
func (u undef) ColorMask(red bool, green bool, blue bool, alpha bool) {
	panic("not implemented") // TODO: Implement
}

// CompileShader compiles the source code of s.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompileShader.xhtml
func (u undef) CompileShader(s Shader) {
	panic("not implemented") // TODO: Implement
}

// CompressedTexImage2D writes a compressed 2D texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompressedTexImage2D.xhtml
func (u undef) CompressedTexImage2D(target Enum, level int, internalformat Enum, width int, height int, border int, data []byte) {
	panic("not implemented") // TODO: Implement
}

// CompressedTexSubImage2D writes a subregion of a compressed 2D texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompressedTexSubImage2D.xhtml
func (u undef) CompressedTexSubImage2D(target Enum, level int, xoffset int, yoffset int, width int, height int, format Enum, data []byte) {
	panic("not implemented") // TODO: Implement
}

// CopyTexImage2D writes a 2D texture from the current framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCopyTexImage2D.xhtml
func (u undef) CopyTexImage2D(target Enum, level int, internalformat Enum, x int, y int, width int, height int, border int) {
	panic("not implemented") // TODO: Implement
}

// CopyTexSubImage2D writes a 2D texture subregion from the
// current framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCopyTexSubImage2D.xhtml
func (u undef) CopyTexSubImage2D(target Enum, level int, xoffset int, yoffset int, x int, y int, width int, height int) {
	panic("not implemented") // TODO: Implement
}

// CreateBuffer creates a buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenBuffers.xhtml
func (u undef) CreateBuffer() Buffer {
	panic("not implemented") // TODO: Implement
}

// CreateFramebuffer creates a framebuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenFramebuffers.xhtml
func (u undef) CreateFramebuffer() Framebuffer {
	panic("not implemented") // TODO: Implement
}

// CreateProgram creates a new empty program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateProgram.xhtml
func (u undef) CreateProgram() Program {
	panic("not implemented") // TODO: Implement
}

// CreateRenderbuffer create a renderbuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenRenderbuffers.xhtml
func (u undef) CreateRenderbuffer() Renderbuffer {
	panic("not implemented") // TODO: Implement
}

// CreateShader creates a new empty shader object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateShader.xhtml
func (u undef) CreateShader(ty Enum) Shader {
	panic("not implemented") // TODO: Implement
}

// CreateTexture creates a texture object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenTextures.xhtml
func (u undef) CreateTexture() Texture {
	panic("not implemented") // TODO: Implement
}

// CreateTVertexArray creates a vertex array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenVertexArrays.xhtml
func (u undef) CreateVertexArray() VertexArray {
	panic("not implemented") // TODO: Implement
}

// CullFace specifies which polygons are candidates for culling.
//
// Valid modes: FRONT, BACK, FRONT_AND_BACK.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCullFace.xhtml
func (u undef) CullFace(mode Enum) {
	panic("not implemented") // TODO: Implement
}

// DeleteBuffer deletes the given buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteBuffers.xhtml
func (u undef) DeleteBuffer(v Buffer) {
	panic("not implemented") // TODO: Implement
}

// DeleteFramebuffer deletes the given framebuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteFramebuffers.xhtml
func (u undef) DeleteFramebuffer(v Framebuffer) {
	panic("not implemented") // TODO: Implement
}

// DeleteProgram deletes the given program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteProgram.xhtml
func (u undef) DeleteProgram(p Program) {
	panic("not implemented") // TODO: Implement
}

// DeleteRenderbuffer deletes the given render buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteRenderbuffers.xhtml
func (u undef) DeleteRenderbuffer(v Renderbuffer) {
	panic("not implemented") // TODO: Implement
}

// DeleteShader deletes shader s.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteShader.xhtml
func (u undef) DeleteShader(s Shader) {
	panic("not implemented") // TODO: Implement
}

// DeleteTexture deletes the given texture object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteTextures.xhtml
func (u undef) DeleteTexture(v Texture) {
	panic("not implemented") // TODO: Implement
}

// DeleteVertexArray deletes the given render buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteVertexArrays.xhtml
func (u undef) DeleteVertexArray(v VertexArray) {
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
func (u undef) DepthFunc(fn Enum) {
	panic("not implemented") // TODO: Implement
}

// DepthMask sets the depth buffer enabled for writing.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDepthMask.xhtml
func (u undef) DepthMask(flag bool) {
	panic("not implemented") // TODO: Implement
}

// DepthRangef sets the mapping from normalized device coordinates to
// window coordinates.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDepthRangef.xhtml
func (u undef) DepthRangef(n float32, f float32) {
	panic("not implemented") // TODO: Implement
}

// DetachShader detaches the shader s from the program p.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDetachShader.xhtml
func (u undef) DetachShader(p Program, s Shader) {
	panic("not implemented") // TODO: Implement
}

// Disable disables various GL capabilities.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDisable.xhtml
func (u undef) Disable(cap Enum) {
	panic("not implemented") // TODO: Implement
}

// DisableVertexAttribArray disables a vertex attribute array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDisableVertexAttribArray.xhtml
func (u undef) DisableVertexAttribArray(a Attrib) {
	panic("not implemented") // TODO: Implement
}

// DrawArrays renders geometric primitives from the bound data.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArrays.xhtml
func (u undef) DrawArrays(mode Enum, first int, count int) {
	panic("not implemented") // TODO: Implement
}

// DrawElements renders primitives from a bound buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawElements.xhtml
func (u undef) DrawElements(mode Enum, count int, ty Enum, offset int) {
	panic("not implemented") // TODO: Implement
}

// Enable enables various GL capabilities.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnable.xhtml
func (u undef) Enable(cap Enum) {
	panic("not implemented") // TODO: Implement
}

// EnableVertexAttribArray enables a vertex attribute array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnableVertexAttribArray.xhtml
func (u undef) EnableVertexAttribArray(a Attrib) {
	panic("not implemented") // TODO: Implement
}

// Finish blocks until the effects of all previously called GL
// commands are complete.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFinish.xhtml
func (u undef) Finish() {
	panic("not implemented") // TODO: Implement
}

// Flush empties all buffers. It does not block.
//
// An OpenGL implementation may buffer network communication,
// the command stream, or data inside the graphics accelerator.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFlush.xhtml
func (u undef) Flush() {
	panic("not implemented") // TODO: Implement
}

// FramebufferRenderbuffer attaches rb to the current frame buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFramebufferRenderbuffer.xhtml
func (u undef) FramebufferRenderbuffer(target Enum, attachment Enum, rbTarget Enum, rb Renderbuffer) {
	panic("not implemented") // TODO: Implement
}

// FramebufferTexture2D attaches the t to the current frame buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFramebufferTexture2D.xhtml
func (u undef) FramebufferTexture2D(target Enum, attachment Enum, texTarget Enum, t Texture, level int) {
	panic("not implemented") // TODO: Implement
}

// FrontFace defines which polygons are front-facing.
//
// Valid modes: CW, CCW.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFrontFace.xhtml
func (u undef) FrontFace(mode Enum) {
	panic("not implemented") // TODO: Implement
}

// GenerateMipmap generates mipmaps for the current texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenerateMipmap.xhtml
func (u undef) GenerateMipmap(target Enum) {
	panic("not implemented") // TODO: Implement
}

// GetActiveAttrib returns details about an active attribute variable.
// A value of 0 for index selects the first active attribute variable.
// Permissible values for index range from 0 to the number of active
// attribute variables minus 1.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveAttrib.xhtml
func (u undef) GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum) {
	panic("not implemented") // TODO: Implement
}

// GetActiveUniform returns details about an active uniform variable.
// A value of 0 for index selects the first active uniform variable.
// Permissible values for index range from 0 to the number of active
// uniform variables minus 1.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveUniform.xhtml
func (u undef) GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum) {
	panic("not implemented") // TODO: Implement
}

// GetAttachedShaders returns the shader objects attached to program p.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetAttachedShaders.xhtml
func (u undef) GetAttachedShaders(p Program) []Shader {
	panic("not implemented") // TODO: Implement
}

// GetAttribLocation returns the location of an attribute variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetAttribLocation.xhtml
func (u undef) GetAttribLocation(p Program, name string) Attrib {
	panic("not implemented") // TODO: Implement
}

// GetBooleanv returns the boolean values of parameter pname.
//
// Many boolean parameters can be queried more easily using IsEnabled.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func (u undef) GetBooleanv(dst []bool, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetFloatv returns the float values of parameter pname.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func (u undef) GetFloatv(dst []float32, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetIntegerv returns the int values of parameter pname.
//
// Single values may be queried more easily using GetInteger.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func (u undef) GetIntegerv(dst []int32, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetInteger returns the int value of parameter pname.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func (u undef) GetInteger(pname Enum) int {
	panic("not implemented") // TODO: Implement
}

// GetBufferParameteri returns a parameter for the active buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetBufferParameter.xhtml
func (u undef) GetBufferParameteri(target Enum, value Enum) int {
	panic("not implemented") // TODO: Implement
}

// GetError returns the next error.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetError.xhtml
func (u undef) GetError() Enum {
	panic("not implemented") // TODO: Implement
}

// GetFramebufferAttachmentParameteri returns attachment parameters
// for the active framebuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetFramebufferAttachmentParameteriv.xhtml
func (u undef) GetFramebufferAttachmentParameteri(target Enum, attachment Enum, pname Enum) int {
	panic("not implemented") // TODO: Implement
}

// GetProgrami returns a parameter value for a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramiv.xhtml
func (u undef) GetProgrami(p Program, pname Enum) int {
	panic("not implemented") // TODO: Implement
}

// GetProgramInfoLog returns the information log for a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramInfoLog.xhtml
func (u undef) GetProgramInfoLog(p Program) string {
	panic("not implemented") // TODO: Implement
}

// GetRenderbufferParameteri returns a parameter value for a render buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetRenderbufferParameteriv.xhtml
func (u undef) GetRenderbufferParameteri(target Enum, pname Enum) int {
	panic("not implemented") // TODO: Implement
}

// GetShaderi returns a parameter value for a shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderiv.xhtml
func (u undef) GetShaderi(s Shader, pname Enum) int {
	panic("not implemented") // TODO: Implement
}

// GetShaderInfoLog returns the information log for a shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderInfoLog.xhtml
func (u undef) GetShaderInfoLog(s Shader) string {
	panic("not implemented") // TODO: Implement
}

// GetShaderPrecisionFormat returns range and precision limits for
// shader types.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderPrecisionFormat.xhtml
func (u undef) GetShaderPrecisionFormat(shadertype Enum, precisiontype Enum) (rangeLow int, rangeHigh int, precision int) {
	panic("not implemented") // TODO: Implement
}

// GetShaderSource returns source code of shader s.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderSource.xhtml
func (u undef) GetShaderSource(s Shader) string {
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
func (u undef) GetString(pname Enum) string {
	panic("not implemented") // TODO: Implement
}

// GetTexParameterfv returns the float values of a texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetTexParameter.xhtml
func (u undef) GetTexParameterfv(dst []float32, target Enum, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetTexParameteriv returns the int values of a texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetTexParameter.xhtml
func (u undef) GetTexParameteriv(dst []int32, target Enum, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetUniformfv returns the float values of a uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniform.xhtml
func (u undef) GetUniformfv(dst []float32, src Uniform, p Program) {
	panic("not implemented") // TODO: Implement
}

// GetUniformiv returns the float values of a uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniform.xhtml
func (u undef) GetUniformiv(dst []int32, src Uniform, p Program) {
	panic("not implemented") // TODO: Implement
}

// GetUniformLocation returns the location of a uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniformLocation.xhtml
func (u undef) GetUniformLocation(p Program, name string) Uniform {
	panic("not implemented") // TODO: Implement
}

// GetVertexAttribf reads the float value of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func (u undef) GetVertexAttribf(src Attrib, pname Enum) float32 {
	panic("not implemented") // TODO: Implement
}

// GetVertexAttribfv reads float values of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func (u undef) GetVertexAttribfv(dst []float32, src Attrib, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// GetVertexAttribi reads the int value of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func (u undef) GetVertexAttribi(src Attrib, pname Enum) int32 {
	panic("not implemented") // TODO: Implement
}

// GetVertexAttribiv reads int values of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func (u undef) GetVertexAttribiv(dst []int32, src Attrib, pname Enum) {
	panic("not implemented") // TODO: Implement
}

// Hint sets implementation-specific modes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glHint.xhtml
func (u undef) Hint(target Enum, mode Enum) {
	panic("not implemented") // TODO: Implement
}

// IsBuffer reports if b is a valid buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsBuffer.xhtml
func (u undef) IsBuffer(b Buffer) bool {
	panic("not implemented") // TODO: Implement
}

// IsEnabled reports if cap is an enabled capability.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsEnabled.xhtml
func (u undef) IsEnabled(cap Enum) bool {
	panic("not implemented") // TODO: Implement
}

// IsFramebuffer reports if fb is a valid frame buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsFramebuffer.xhtml
func (u undef) IsFramebuffer(fb Framebuffer) bool {
	panic("not implemented") // TODO: Implement
}

// IsProgram reports if p is a valid program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsProgram.xhtml
func (u undef) IsProgram(p Program) bool {
	panic("not implemented") // TODO: Implement
}

// IsRenderbuffer reports if rb is a valid render buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsRenderbuffer.xhtml
func (u undef) IsRenderbuffer(rb Renderbuffer) bool {
	panic("not implemented") // TODO: Implement
}

// IsShader reports if s is valid shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsShader.xhtml
func (u undef) IsShader(s Shader) bool {
	panic("not implemented") // TODO: Implement
}

// IsTexture reports if t is a valid texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsTexture.xhtml
func (u undef) IsTexture(t Texture) bool {
	panic("not implemented") // TODO: Implement
}

// LineWidth specifies the width of lines.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glLineWidth.xhtml
func (u undef) LineWidth(width float32) {
	panic("not implemented") // TODO: Implement
}

// LinkProgram links the specified program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glLinkProgram.xhtml
func (u undef) LinkProgram(p Program) {
	panic("not implemented") // TODO: Implement
}

// PixelStorei sets pixel storage parameters.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glPixelStorei.xhtml
func (u undef) PixelStorei(pname Enum, param int32) {
	panic("not implemented") // TODO: Implement
}

// PolygonOffset sets the scaling factors for depth offsets.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glPolygonOffset.xhtml
func (u undef) PolygonOffset(factor float32, units float32) {
	panic("not implemented") // TODO: Implement
}

// ReadPixels returns pixel data from a buffer.
//
// In GLES 3, the source buffer is controlled with ReadBuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glReadPixels.xhtml
func (u undef) ReadPixels(dst []byte, x int, y int, width int, height int, format Enum, ty Enum) {
	panic("not implemented") // TODO: Implement
}

// ReleaseShaderCompiler frees resources allocated by the shader compiler.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glReleaseShaderCompiler.xhtml
func (u undef) ReleaseShaderCompiler() {
	panic("not implemented") // TODO: Implement
}

// RenderbufferStorage establishes the data storage, format, and
// dimensions of a renderbuffer object's image.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glRenderbufferStorage.xhtml
func (u undef) RenderbufferStorage(target Enum, internalFormat Enum, width int, height int) {
	panic("not implemented") // TODO: Implement
}

// SampleCoverage sets multisample coverage parameters.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glSampleCoverage.xhtml
func (u undef) SampleCoverage(value float32, invert bool) {
	panic("not implemented") // TODO: Implement
}

// Scissor defines the scissor box rectangle, in window coordinates.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glScissor.xhtml
func (u undef) Scissor(x int32, y int32, width int32, height int32) {
	panic("not implemented") // TODO: Implement
}

// ShaderSource sets the source code of s to the given source code.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glShaderSource.xhtml
func (u undef) ShaderSource(s Shader, src string) {
	panic("not implemented") // TODO: Implement
}

// StencilFunc sets the front and back stencil test reference value.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilFunc.xhtml
func (u undef) StencilFunc(fn Enum, ref int, mask uint32) {
	panic("not implemented") // TODO: Implement
}

// StencilFunc sets the front or back stencil test reference value.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilFuncSeparate.xhtml
func (u undef) StencilFuncSeparate(face Enum, fn Enum, ref int, mask uint32) {
	panic("not implemented") // TODO: Implement
}

// StencilMask controls the writing of bits in the stencil planes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilMask.xhtml
func (u undef) StencilMask(mask uint32) {
	panic("not implemented") // TODO: Implement
}

// StencilMaskSeparate controls the writing of bits in the stencil planes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilMaskSeparate.xhtml
func (u undef) StencilMaskSeparate(face Enum, mask uint32) {
	panic("not implemented") // TODO: Implement
}

// StencilOp sets front and back stencil test actions.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilOp.xhtml
func (u undef) StencilOp(fail Enum, zfail Enum, zpass Enum) {
	panic("not implemented") // TODO: Implement
}

// StencilOpSeparate sets front or back stencil tests.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilOpSeparate.xhtml
func (u undef) StencilOpSeparate(face Enum, sfail Enum, dpfail Enum, dppass Enum) {
	panic("not implemented") // TODO: Implement
}

// TexImage2D writes a 2D texture image.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexImage2D.xhtml
func (u undef) TexImage2D(target Enum, level int, internalFormat int, width int, height int, format Enum, ty Enum, data []byte) {
	panic("not implemented") // TODO: Implement
}

// TexSubImage2D writes a subregion of a 2D texture image.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexSubImage2D.xhtml
func (u undef) TexSubImage2D(target Enum, level int, x int, y int, width int, height int, format Enum, ty Enum, data []byte) {
	panic("not implemented") // TODO: Implement
}

// TexParameterf sets a float texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (u undef) TexParameterf(target Enum, pname Enum, param float32) {
	panic("not implemented") // TODO: Implement
}

// TexParameterfv sets a float texture parameter array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (u undef) TexParameterfv(target Enum, pname Enum, params []float32) {
	panic("not implemented") // TODO: Implement
}

// TexParameteri sets an integer texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (u undef) TexParameteri(target Enum, pname Enum, param int) {
	panic("not implemented") // TODO: Implement
}

// TexParameteriv sets an integer texture parameter array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func (u undef) TexParameteriv(target Enum, pname Enum, params []int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform1f writes a float uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform1f(dst Uniform, v float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform1fv writes a [len(src)]float uniform array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform1fv(dst Uniform, src []float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform1i writes an int uniform variable.
//
// Uniform1i and Uniform1iv are the only two functions that may be used
// to load uniform variables defined as sampler types. Loading samplers
// with any other function will result in a INVALID_OPERATION error.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform1i(dst Uniform, v int) {
	panic("not implemented") // TODO: Implement
}

// Uniform1iv writes a int uniform array of len(src) elements.
//
// Uniform1i and Uniform1iv are the only two functions that may be used
// to load uniform variables defined as sampler types. Loading samplers
// with any other function will result in a INVALID_OPERATION error.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform1iv(dst Uniform, src []int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform2f writes a vec2 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform2f(dst Uniform, v0 float32, v1 float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform2fv writes a vec2 uniform array of len(src)/2 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform2fv(dst Uniform, src []float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform2i writes an ivec2 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform2i(dst Uniform, v0 int, v1 int) {
	panic("not implemented") // TODO: Implement
}

// Uniform2iv writes an ivec2 uniform array of len(src)/2 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform2iv(dst Uniform, src []int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform3f writes a vec3 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform3f(dst Uniform, v0 float32, v1 float32, v2 float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform3fv writes a vec3 uniform array of len(src)/3 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform3fv(dst Uniform, src []float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform3i writes an ivec3 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform3i(dst Uniform, v0 int32, v1 int32, v2 int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform3iv writes an ivec3 uniform array of len(src)/3 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform3iv(dst Uniform, src []int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform4f writes a vec4 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform4f(dst Uniform, v0 float32, v1 float32, v2 float32, v3 float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform4fv writes a vec4 uniform array of len(src)/4 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform4fv(dst Uniform, src []float32) {
	panic("not implemented") // TODO: Implement
}

// Uniform4i writes an ivec4 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform4i(dst Uniform, v0 int32, v1 int32, v2 int32, v3 int32) {
	panic("not implemented") // TODO: Implement
}

// Uniform4i writes an ivec4 uniform array of len(src)/4 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) Uniform4iv(dst Uniform, src []int32) {
	panic("not implemented") // TODO: Implement
}

// UniformMatrix2fv writes 2x2 matrices. Each matrix uses four
// float32 values, so the number of matrices written is len(src)/4.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) UniformMatrix2fv(dst Uniform, src []float32) {
	panic("not implemented") // TODO: Implement
}

// UniformMatrix3fv writes 3x3 matrices. Each matrix uses nine
// float32 values, so the number of matrices written is len(src)/9.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) UniformMatrix3fv(dst Uniform, src []float32) {
	panic("not implemented") // TODO: Implement
}

// UniformMatrix4fv writes 4x4 matrices. Each matrix uses 16
// float32 values, so the number of matrices written is len(src)/16.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func (u undef) UniformMatrix4fv(dst Uniform, src []float32) {
	panic("not implemented") // TODO: Implement
}

// UseProgram sets the active program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUseProgram.xhtml
func (u undef) UseProgram(p Program) {
	panic("not implemented") // TODO: Implement
}

// ValidateProgram checks to see whether the executables contained in
// program can execute given the current OpenGL state.
//
// Typically only used for debugging.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glValidateProgram.xhtml
func (u undef) ValidateProgram(p Program) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib1f writes a float vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (u undef) VertexAttrib1f(dst Attrib, x float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib1fv writes a float vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (u undef) VertexAttrib1fv(dst Attrib, src []float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib2f writes a vec2 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (u undef) VertexAttrib2f(dst Attrib, x float32, y float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib2fv writes a vec2 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (u undef) VertexAttrib2fv(dst Attrib, src []float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib3f writes a vec3 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (u undef) VertexAttrib3f(dst Attrib, x float32, y float32, z float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib3fv writes a vec3 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (u undef) VertexAttrib3fv(dst Attrib, src []float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib4f writes a vec4 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (u undef) VertexAttrib4f(dst Attrib, x float32, y float32, z float32, w float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttrib4fv writes a vec4 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func (u undef) VertexAttrib4fv(dst Attrib, src []float32) {
	panic("not implemented") // TODO: Implement
}

// VertexAttribPointer uses a bound buffer to define vertex attribute data.
//
// Direct use of VertexAttribPointer to load data into OpenGL is not
// supported via the Go bindings. Instead, use BindBuffer with an
// ARRAY_BUFFER and then fill it using BufferData.
//
// The size argument specifies the number of components per attribute,
// between 1-4. The stride argument specifies the byte offset between
// consecutive vertex attributes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttribPointer.xhtml
func (u undef) VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride int, offset int) {
	panic("not implemented") // TODO: Implement
}

// Viewport sets the viewport, an affine transformation that
// normalizes device coordinates to window coordinates.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glViewport.xhtml
func (u undef) Viewport(x int, y int, width int, height int) {
	panic("not implemented") // TODO: Implement
}

// GetUniformBlockIndex retrieves the index of a uniform block within program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockIndex.xhtml
func (u undef) GetUniformBlockIndex(p Program, name string) int {
	panic("not implemented") // TODO: Implement
}

// UniformBlockBinding assign a binding point to an active uniform block
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockBinding.xhtml
func (u undef) UniformBlockBinding(p Program, index, bind int) {
	panic("not implemented") // TODO: Implement
}

// BindBufferBase bind a buffer object to an indexed buffer target
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBufferBase.xhtml
func (u undef) BindBufferBase(target Enum, n uint32, b Buffer) {
	panic("not implemented") // TODO: Implement
}

// DrawArraysInstanced draw multiple instances of a range of elements
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
func (u undef) DrawArraysInstanced(mode Enum, first int, count int, primcount int) {
	panic("not implemented") // TODO: Implement
}

// DrawElementsInstanced — draw multiple instances of a set of elements
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawElementsInstanced.xhtml
func (u undef) DrawElementsInstanced(mode Enum, count int, typ Enum, offset int, primcount int) {
	panic("not implemented") // TODO: Implement
}

// VertexAttribDivisor  modify the rate at which generic vertex attributes
// advance during instanced rendering
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
func (u undef) VertexAttribDivisor(index Attrib, divisor int) {
	panic("not implemented") // TODO: Implement
}

// BufferDataX will type switch the interface and select the proper type
// {lpf} Custom func
// TODO: Should move this Elsewhere as an helper
func (u undef) BufferDataX(target Enum, d interface{}, usage Enum) {
	panic("not implemented") // TODO: Implement
}
