// nolint
package gl

import "log"

type wrapperi interface {
	Context3
	Impl() string
}

// var glw *Wrapper
var glw wrapperi

func Init(w *Wrapper) {
	s := &cached{wrapper: w}
	s.init()
	glw = s
	log.Println("GL The wrapper:", w)
	log.Println("GL version:", w.GetString(VERSION))
	log.Println("GL Renderer:", w.GetString(RENDERER))
}

func Global() wrapperi { return glw }

// This function redirects to wrapper which implements the platform specific GL
// There could be an overhead but hopefully they will be inlined

func ActiveTexture(texture Enum) { glw.ActiveTexture(texture) }

func AttachShader(p Program, s Shader) { glw.AttachShader(p, s) }

func BindAttribLocation(p Program, a Attrib, name string) { glw.BindAttribLocation(p, a, name) }

func BindBuffer(target Enum, b Buffer) { glw.BindBuffer(target, b) }

func BindFramebuffer(target Enum, fb Framebuffer) { glw.BindFramebuffer(target, fb) }

func BindRenderbuffer(target Enum, rb Renderbuffer) { glw.BindRenderbuffer(target, rb) }

func BindTexture(target Enum, t Texture) { glw.BindTexture(target, t) }

func BindVertexArray(rb VertexArray) { glw.BindVertexArray(rb) }

func BlendColor(red, green, blue, alpha float32) { glw.BlendColor(red, green, blue, alpha) }

func BlendEquation(mode Enum) { glw.BlendEquation(mode) }

func BlendEquationSeparate(modeRGB, modeAlpha Enum) {
	glw.BlendEquationSeparate(modeRGB, modeAlpha)
}

func BlendFunc(sfactor, dfactor Enum) { glw.BlendFunc(sfactor, dfactor) }

func BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha Enum) {
	glw.BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha)
}

func BufferInit(target Enum, size int, usage Enum) { glw.BufferInit(target, size, usage) }

func BufferData(target Enum, src any, usage Enum) { glw.BufferData(target, src, usage) }

func BufferSubData(target Enum, offset int, src any) {
	glw.BufferSubData(target, offset, src)
}

func CheckFramebufferStatus(target Enum) Enum { return glw.CheckFramebufferStatus(target) }

func Clear(mask Enum) { glw.Clear(mask) }

func ClearColor(red, green, blue, alpha float32) { glw.ClearColor(red, green, blue, alpha) }

func ClearDepthf(d float32) { glw.ClearDepthf(d) }

func ClearStencil(s int) { glw.ClearStencil(s) }

func ColorMask(red, green, blue, alpha bool) { glw.ColorMask(red, green, blue, alpha) }

func CompileShader(s Shader) { glw.CompileShader(s) }

func CompressedTexImage2D(target Enum, level int, internalformat Enum, width, height, border int, data []byte) {
	glw.CompressedTexImage2D(target, level, internalformat, width, height, border, data)
}

func CompressedTexSubImage2D(target Enum, level, xoffset, yoffset, width, height int, format Enum, data []byte) {
	glw.CompressedTexSubImage2D(target, level, xoffset, yoffset, width, height, format, data)
}

func CopyTexImage2D(target Enum, level int, internalformat Enum, x, y, width, height, border int) {
	glw.CopyTexImage2D(target, level, internalformat, x, y, width, height, border)
}

func CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int) {
	glw.CopyTexSubImage2D(target, level, xoffset, yoffset, x, y, width, height)
}

func CreateBuffer() Buffer { return glw.CreateBuffer() }

func CreateFramebuffer() Framebuffer { return glw.CreateFramebuffer() }

func CreateProgram() Program { return glw.CreateProgram() }

func CreateRenderbuffer() Renderbuffer { return glw.CreateRenderbuffer() }

func CreateShader(ty Enum) Shader { return glw.CreateShader(ty) }

func CreateTexture() Texture { return glw.CreateTexture() }

func CreateVertexArray() VertexArray { return glw.CreateVertexArray() }

func CullFace(mode Enum) { glw.CullFace(mode) }

func DeleteBuffer(v Buffer) { glw.DeleteBuffer(v) }

func DeleteFramebuffer(v Framebuffer) { glw.DeleteFramebuffer(v) }

func DeleteProgram(p Program) { glw.DeleteProgram(p) }

func DeleteRenderbuffer(v Renderbuffer) { glw.DeleteRenderbuffer(v) }

func DeleteShader(s Shader) { glw.DeleteShader(s) }

func DeleteTexture(v Texture) { glw.DeleteTexture(v) }

func DeleteVertexArray(v VertexArray) { glw.DeleteVertexArray(v) }

func DepthFunc(fn Enum) { glw.DepthFunc(fn) }

func DepthMask(flag bool) { glw.DepthMask(flag) }

func DepthRangef(n, f float32) { glw.DepthRangef(n, f) }

func DetachShader(p Program, s Shader) { glw.DetachShader(p, s) }

func Disable(cap Enum) { glw.Disable(cap) }

func DisableVertexAttribArray(a Attrib) { glw.DisableVertexAttribArray(a) }

func DrawArrays(mode Enum, first, count int) { glw.DrawArrays(mode, first, count) }

func DrawElements(mode Enum, count int, ty Enum, offset int) {
	glw.DrawElements(mode, count, ty, offset)
}

func Enable(cap Enum) { glw.Enable(cap) }

func EnableVertexAttribArray(a Attrib) { glw.EnableVertexAttribArray(a) }

func Finish() { glw.Finish() }

func Flush() { glw.Flush() }

func FramebufferRenderbuffer(target, attachment, rbTarget Enum, rb Renderbuffer) {
	glw.FramebufferRenderbuffer(target, attachment, rbTarget, rb)
}

func FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	glw.FramebufferTexture2D(target, attachment, texTarget, t, level)
}

func FrontFace(mode Enum) { glw.FrontFace(mode) }

func GenerateMipmap(target Enum) { glw.GenerateMipmap(target) }

func GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum) {
	return glw.GetActiveAttrib(p, index)
}

func GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum) {
	return glw.GetActiveUniform(p, index)
}

func GetAttachedShaders(p Program) []Shader { return glw.GetAttachedShaders(p) }

func GetAttribLocation(p Program, name string) Attrib { return glw.GetAttribLocation(p, name) }

func GetBooleanv(dst []bool, pname Enum) { glw.GetBooleanv(dst, pname) }

func GetFloatv(dst []float32, pname Enum) { glw.GetFloatv(dst, pname) }

func GetIntegerv(dst []int32, pname Enum) { glw.GetIntegerv(dst, pname) }

func GetInteger(pname Enum) int { return glw.GetInteger(pname) }

func GetBufferParameteri(target, value Enum) int { return glw.GetBufferParameteri(target, value) }

func GetError() Enum { return glw.GetError() }

func GetFramebufferAttachmentParameteri(target, attachment, pname Enum) int {
	return glw.GetFramebufferAttachmentParameteri(target, attachment, pname)
}

func GetProgrami(p Program, pname Enum) int { return glw.GetProgrami(p, pname) }

func GetProgramInfoLog(p Program) string { return glw.GetProgramInfoLog(p) }

func GetRenderbufferParameteri(target, pname Enum) int {
	return glw.GetRenderbufferParameteri(target, pname)
}

func GetShaderi(s Shader, pname Enum) int { return glw.GetShaderi(s, pname) }

func GetShaderInfoLog(s Shader) string { return glw.GetShaderInfoLog(s) }

func GetShaderPrecisionFormat(shadertype, precisiontype Enum) (rangeLow, rangeHigh, precision int) {
	return glw.GetShaderPrecisionFormat(shadertype, precisiontype)
}

func GetShaderSource(s Shader) string { return glw.GetShaderSource(s) }

func GetString(pname Enum) string { return glw.GetString(pname) }

func GetTexParameterfv(dst []float32, target, pname Enum) {
	glw.GetTexParameterfv(dst, target, pname)
}

func GetTexParameteriv(dst []int32, target, pname Enum) {
	glw.GetTexParameteriv(dst, target, pname)
}

func GetUniformfv(dst []float32, src Uniform, p Program) { glw.GetUniformfv(dst, src, p) }

func GetUniformiv(dst []int32, src Uniform, p Program) { glw.GetUniformiv(dst, src, p) }

func GetUniformLocation(p Program, name string) Uniform { return glw.GetUniformLocation(p, name) }

func GetVertexAttribf(src Attrib, pname Enum) float32 { return glw.GetVertexAttribf(src, pname) }

func GetVertexAttribfv(dst []float32, src Attrib, pname Enum) {
	glw.GetVertexAttribfv(dst, src, pname)
}

func GetVertexAttribi(src Attrib, pname Enum) int32 { return glw.GetVertexAttribi(src, pname) }

func GetVertexAttribiv(dst []int32, src Attrib, pname Enum) {
	glw.GetVertexAttribiv(dst, src, pname)
}

func Hint(target, mode Enum) { glw.Hint(target, mode) }

func IsBuffer(b Buffer) bool { return glw.IsBuffer(b) }

func IsEnabled(cap Enum) bool { return glw.IsEnabled(cap) }

func IsFramebuffer(fb Framebuffer) bool { return glw.IsFramebuffer(fb) }

func IsProgram(p Program) bool { return glw.IsProgram(p) }

func IsRenderbuffer(rb Renderbuffer) bool { return glw.IsRenderbuffer(rb) }

func IsShader(s Shader) bool { return glw.IsShader(s) }

func IsTexture(t Texture) bool { return glw.IsTexture(t) }

func LineWidth(width float32) { glw.LineWidth(width) }

func LinkProgram(p Program) { glw.LinkProgram(p) }

func PixelStorei(pname Enum, param int32) { glw.PixelStorei(pname, param) }

func PolygonOffset(factor, units float32) { glw.PolygonOffset(factor, units) }

func ReadPixels(dst []byte, x, y, width, height int, format, ty Enum) {
	glw.ReadPixels(dst, x, y, width, height, format, ty)
}

func ReleaseShaderCompiler() { glw.ReleaseShaderCompiler() }

func RenderbufferStorage(target, internalFormat Enum, width, height int) {
	glw.RenderbufferStorage(target, internalFormat, width, height)
}

func SampleCoverage(value float32, invert bool) { glw.SampleCoverage(value, invert) }

func Scissor(x, y, width, height int32) { glw.Scissor(x, y, width, height) }

func ShaderSource(s Shader, src string) { glw.ShaderSource(s, src) }

func StencilFunc(fn Enum, ref int, mask uint32) { glw.StencilFunc(fn, ref, mask) }

func StencilFuncSeparate(face, fn Enum, ref int, mask uint32) {
	glw.StencilFuncSeparate(face, fn, ref, mask)
}

func StencilMask(mask uint32) { glw.StencilMask(mask) }

func StencilMaskSeparate(face Enum, mask uint32) { glw.StencilMaskSeparate(face, mask) }

func StencilOp(fail, zfail, zpass Enum) { glw.StencilOp(fail, zfail, zpass) }

func StencilOpSeparate(face, sfail, dpfail, dppass Enum) {
	glw.StencilOpSeparate(face, sfail, dpfail, dppass)
}

func TexImage2D(target Enum, level int, internalFormat int, width, height int, format Enum, ty Enum, data []byte) {
	glw.TexImage2D(target, level, internalFormat, width, height, format, ty, data)
}

func TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data []byte) {
	glw.TexSubImage2D(target, level, x, y, width, height, format, ty, data)
}

func TexParameterf(target, pname Enum, param float32) { glw.TexParameterf(target, pname, param) }

func TexParameterfv(target, pname Enum, params []float32) {
	glw.TexParameterfv(target, pname, params)
}

func TexParameteri(target, pname Enum, param int) { glw.TexParameteri(target, pname, param) }

func TexParameteriv(target, pname Enum, params []int32) {
	glw.TexParameteriv(target, pname, params)
}

func Uniform1f(dst Uniform, v float32) { glw.Uniform1f(dst, v) }

func Uniform1fv(dst Uniform, src []float32) { glw.Uniform1fv(dst, src) }

func Uniform1i(dst Uniform, v int) { glw.Uniform1i(dst, v) }

func Uniform1iv(dst Uniform, src []int32) { glw.Uniform1iv(dst, src) }

func Uniform2f(dst Uniform, v0, v1 float32) { glw.Uniform2f(dst, v0, v1) }

func Uniform2fv(dst Uniform, src []float32) { glw.Uniform2fv(dst, src) }

func Uniform2i(dst Uniform, v0, v1 int) { glw.Uniform2i(dst, v0, v1) }

func Uniform2iv(dst Uniform, src []int32) { glw.Uniform2iv(dst, src) }

func Uniform3f(dst Uniform, v0, v1, v2 float32) { glw.Uniform3f(dst, v0, v1, v2) }

func Uniform3fv(dst Uniform, src []float32) { glw.Uniform3fv(dst, src) }

func Uniform3i(dst Uniform, v0, v1, v2 int32) { glw.Uniform3i(dst, v0, v1, v2) }

func Uniform3iv(dst Uniform, src []int32) { glw.Uniform3iv(dst, src) }

func Uniform4f(dst Uniform, v0, v1, v2, v3 float32) { glw.Uniform4f(dst, v0, v1, v2, v3) }

func Uniform4fv(dst Uniform, src []float32) { glw.Uniform4fv(dst, src) }

func Uniform4i(dst Uniform, v0, v1, v2, v3 int32) { glw.Uniform4i(dst, v0, v1, v2, v3) }

func Uniform4iv(dst Uniform, src []int32) { glw.Uniform4iv(dst, src) }

func UniformMatrix2fv(dst Uniform, src []float32) { glw.UniformMatrix2fv(dst, src) }

func UniformMatrix3fv(dst Uniform, src []float32) { glw.UniformMatrix3fv(dst, src) }

func UniformMatrix4fv(dst Uniform, src []float32) { glw.UniformMatrix4fv(dst, src) }

func UseProgram(p Program) { glw.UseProgram(p) }

func ValidateProgram(p Program) { glw.ValidateProgram(p) }

func VertexAttrib1f(dst Attrib, x float32) { glw.VertexAttrib1f(dst, x) }

func VertexAttrib1fv(dst Attrib, src []float32) { glw.VertexAttrib1fv(dst, src) }

func VertexAttrib2f(dst Attrib, x, y float32) { glw.VertexAttrib2f(dst, x, y) }

func VertexAttrib2fv(dst Attrib, src []float32) { glw.VertexAttrib2fv(dst, src) }

func VertexAttrib3f(dst Attrib, x, y, z float32) { glw.VertexAttrib3f(dst, x, y, z) }

func VertexAttrib3fv(dst Attrib, src []float32) { glw.VertexAttrib3fv(dst, src) }

func VertexAttrib4f(dst Attrib, x, y, z, w float32) { glw.VertexAttrib4f(dst, x, y, z, w) }

func VertexAttrib4fv(dst Attrib, src []float32) { glw.VertexAttrib4fv(dst, src) }

func VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride, offset int) {
	glw.VertexAttribPointer(dst, size, ty, normalized, stride, offset)
}

func Viewport(x, y, width, height int32) { glw.Viewport(x, y, width, height) }

func GetUniformBlockIndex(p Program, name string) uint32 {
	return glw.GetUniformBlockIndex(p, name)
}

func UniformBlockBinding(p Program, index, bind uint32) {
	glw.UniformBlockBinding(p, index, bind)
}

func BindBufferBase(target Enum, n uint32, b Buffer) {
	glw.BindBufferBase(target, n, b)
}

func GetActiveUniformi(p Program, index uint32, pname Enum) int32 {
	return glw.GetActiveUniformi(p, index, pname)
}

func GetActiveUniformBlockiv(p Program, index uint32, pname Enum, params []int32) {
	glw.GetActiveUniformBlockiv(p, index, pname, params)
}

func GetActiveUniformBlockName(p Program, index uint32) string {
	return glw.GetActiveUniformBlockName(p, index)
}

func DrawArraysInstanced(mode Enum, first, count, primcount uint32) {
	glw.DrawArraysInstanced(mode, first, count, primcount)
}

func DrawElementsInstanced(mode Enum, count uint32, typ Enum, offset, primcount uint32) {
	glw.DrawElementsInstanced(mode, count, typ, offset, primcount)
}

func VertexAttribDivisor(index Attrib, divisor uint32) {
	glw.VertexAttribDivisor(index, divisor)
}

/*
func TexImage3D(target Enum, level, internalFormat, width, height, depth int, format, ty Enum, data []byte) {
	wrapper.TexImage3D(target, level, internalFormat, width, height, depth, format, ty, data)
}

func FramebufferTextureLayer(target Enum, attachment Enum, texture Texture, level int, layer int) {
	wrapper.FramebufferTextureLayer(target, attachment, texture, level, layer)
}
*/
