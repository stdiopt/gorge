// Package gl ...
// nolint
package gl

var wrapper *Wrapper

func Init(glw *Wrapper) {
	wrapper = glw
}

// This function redirects to wrapper which implements the platform specific GL
// There could be an overhead but hopefully they will be inlined

func ActiveTexture(texture Enum) { wrapper.ActiveTexture(texture) }

func AttachShader(p Program, s Shader) { wrapper.AttachShader(p, s) }

func BindAttribLocation(p Program, a Attrib, name string) { wrapper.BindAttribLocation(p, a, name) }

func BindBuffer(target Enum, b Buffer) { wrapper.BindBuffer(target, b) }

func BindFramebuffer(target Enum, fb Framebuffer) { wrapper.BindFramebuffer(target, fb) }

func BindRenderbuffer(target Enum, rb Renderbuffer) { wrapper.BindRenderbuffer(target, rb) }

func BindTexture(target Enum, t Texture) { wrapper.BindTexture(target, t) }

func BindVertexArray(rb VertexArray) { wrapper.BindVertexArray(rb) }

func BlendColor(red, green, blue, alpha float32) { wrapper.BlendColor(red, green, blue, alpha) }

func BlendEquation(mode Enum) { wrapper.BlendEquation(mode) }

func BlendEquationSeparate(modeRGB, modeAlpha Enum) {
	wrapper.BlendEquationSeparate(modeRGB, modeAlpha)
}

func BlendFunc(sfactor, dfactor Enum) { wrapper.BlendFunc(sfactor, dfactor) }

func BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha Enum) {
	wrapper.BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha)
}

func BufferInit(target Enum, size int, usage Enum) { wrapper.BufferInit(target, size, usage) }

func BufferData(target Enum, src interface{}, usage Enum) { wrapper.BufferData(target, src, usage) }

func BufferSubData(target Enum, offset int, src interface{}) {
	wrapper.BufferSubData(target, offset, src)
}

func CheckFramebufferStatus(target Enum) Enum { return wrapper.CheckFramebufferStatus(target) }

func Clear(mask Enum) { wrapper.Clear(mask) }

func ClearColor(red, green, blue, alpha float32) { wrapper.ClearColor(red, green, blue, alpha) }

func ClearDepthf(d float32) { wrapper.ClearDepthf(d) }

func ClearStencil(s int) { wrapper.ClearStencil(s) }

func ColorMask(red, green, blue, alpha bool) { wrapper.ColorMask(red, green, blue, alpha) }

func CompileShader(s Shader) { wrapper.CompileShader(s) }

func CompressedTexImage2D(target Enum, level int, internalformat Enum, width, height, border int, data []byte) {
	wrapper.CompressedTexImage2D(target, level, internalformat, width, height, border, data)
}

func CompressedTexSubImage2D(target Enum, level, xoffset, yoffset, width, height int, format Enum, data []byte) {
	wrapper.CompressedTexSubImage2D(target, level, xoffset, yoffset, width, height, format, data)
}

func CopyTexImage2D(target Enum, level int, internalformat Enum, x, y, width, height, border int) {
	wrapper.CopyTexImage2D(target, level, internalformat, x, y, width, height, border)
}

func CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int) {
	wrapper.CopyTexSubImage2D(target, level, xoffset, yoffset, x, y, width, height)
}

func CreateBuffer() Buffer { return wrapper.CreateBuffer() }

func CreateFramebuffer() Framebuffer { return wrapper.CreateFramebuffer() }

func CreateProgram() Program { return wrapper.CreateProgram() }

func CreateRenderbuffer() Renderbuffer { return wrapper.CreateRenderbuffer() }

func CreateShader(ty Enum) Shader { return wrapper.CreateShader(ty) }

func CreateTexture() Texture { return wrapper.CreateTexture() }

func CreateVertexArray() VertexArray { return wrapper.CreateVertexArray() }

func CullFace(mode Enum) { wrapper.CullFace(mode) }

func DeleteBuffer(v Buffer) { wrapper.DeleteBuffer(v) }

func DeleteFramebuffer(v Framebuffer) { wrapper.DeleteFramebuffer(v) }

func DeleteProgram(p Program) { wrapper.DeleteProgram(p) }

func DeleteRenderbuffer(v Renderbuffer) { wrapper.DeleteRenderbuffer(v) }

func DeleteShader(s Shader) { wrapper.DeleteShader(s) }

func DeleteTexture(v Texture) { wrapper.DeleteTexture(v) }

func DeleteVertexArray(v VertexArray) { wrapper.DeleteVertexArray(v) }

func DepthFunc(fn Enum) { wrapper.DepthFunc(fn) }

func DepthMask(flag bool) { wrapper.DepthMask(flag) }

func DepthRangef(n, f float32) { wrapper.DepthRangef(n, f) }

func DetachShader(p Program, s Shader) { wrapper.DetachShader(p, s) }

func Disable(cap Enum) { wrapper.Disable(cap) }

func DisableVertexAttribArray(a Attrib) { wrapper.DisableVertexAttribArray(a) }

func DrawArrays(mode Enum, first, count int) { wrapper.DrawArrays(mode, first, count) }

func DrawElements(mode Enum, count int, ty Enum, offset int) {
	wrapper.DrawElements(mode, count, ty, offset)
}

func Enable(cap Enum) { wrapper.Enable(cap) }

func EnableVertexAttribArray(a Attrib) { wrapper.EnableVertexAttribArray(a) }

func Finish() { wrapper.Finish() }

func Flush() { wrapper.Flush() }

func FramebufferRenderbuffer(target, attachment, rbTarget Enum, rb Renderbuffer) {
	wrapper.FramebufferRenderbuffer(target, attachment, rbTarget, rb)
}

func FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	wrapper.FramebufferTexture2D(target, attachment, texTarget, t, level)
}

func FrontFace(mode Enum) { wrapper.FrontFace(mode) }

func GenerateMipmap(target Enum) { wrapper.GenerateMipmap(target) }

func GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum) {
	return wrapper.GetActiveAttrib(p, index)
}

func GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum) {
	return wrapper.GetActiveUniform(p, index)
}

func GetAttachedShaders(p Program) []Shader { return wrapper.GetAttachedShaders(p) }

func GetAttribLocation(p Program, name string) Attrib { return wrapper.GetAttribLocation(p, name) }

func GetBooleanv(dst []bool, pname Enum) { wrapper.GetBooleanv(dst, pname) }

func GetFloatv(dst []float32, pname Enum) { wrapper.GetFloatv(dst, pname) }

func GetIntegerv(dst []int32, pname Enum) { wrapper.GetIntegerv(dst, pname) }

func GetInteger(pname Enum) int { return wrapper.GetInteger(pname) }

func GetBufferParameteri(target, value Enum) int { return wrapper.GetBufferParameteri(target, value) }

func GetError() Enum { return wrapper.GetError() }

func GetFramebufferAttachmentParameteri(target, attachment, pname Enum) int {
	return wrapper.GetFramebufferAttachmentParameteri(target, attachment, pname)
}

func GetProgrami(p Program, pname Enum) int { return wrapper.GetProgrami(p, pname) }

func GetProgramInfoLog(p Program) string { return wrapper.GetProgramInfoLog(p) }

func GetRenderbufferParameteri(target, pname Enum) int {
	return wrapper.GetRenderbufferParameteri(target, pname)
}

func GetShaderi(s Shader, pname Enum) int { return wrapper.GetShaderi(s, pname) }

func GetShaderInfoLog(s Shader) string { return wrapper.GetShaderInfoLog(s) }

func GetShaderPrecisionFormat(shadertype, precisiontype Enum) (rangeLow, rangeHigh, precision int) {
	return wrapper.GetShaderPrecisionFormat(shadertype, precisiontype)
}

func GetShaderSource(s Shader) string { return wrapper.GetShaderSource(s) }

func GetString(pname Enum) string { return wrapper.GetString(pname) }

func GetTexParameterfv(dst []float32, target, pname Enum) {
	wrapper.GetTexParameterfv(dst, target, pname)
}

func GetTexParameteriv(dst []int32, target, pname Enum) {
	wrapper.GetTexParameteriv(dst, target, pname)
}

func GetUniformfv(dst []float32, src Uniform, p Program) { wrapper.GetUniformfv(dst, src, p) }

func GetUniformiv(dst []int32, src Uniform, p Program) { wrapper.GetUniformiv(dst, src, p) }

func GetUniformLocation(p Program, name string) Uniform { return wrapper.GetUniformLocation(p, name) }

func GetVertexAttribf(src Attrib, pname Enum) float32 { return wrapper.GetVertexAttribf(src, pname) }

func GetVertexAttribfv(dst []float32, src Attrib, pname Enum) {
	wrapper.GetVertexAttribfv(dst, src, pname)
}

func GetVertexAttribi(src Attrib, pname Enum) int32 { return wrapper.GetVertexAttribi(src, pname) }

func GetVertexAttribiv(dst []int32, src Attrib, pname Enum) {
	wrapper.GetVertexAttribiv(dst, src, pname)
}

func Hint(target, mode Enum) { wrapper.Hint(target, mode) }

func IsBuffer(b Buffer) bool { return wrapper.IsBuffer(b) }

func IsEnabled(cap Enum) bool { return wrapper.IsEnabled(cap) }

func IsFramebuffer(fb Framebuffer) bool { return wrapper.IsFramebuffer(fb) }

func IsProgram(p Program) bool { return wrapper.IsProgram(p) }

func IsRenderbuffer(rb Renderbuffer) bool { return wrapper.IsRenderbuffer(rb) }

func IsShader(s Shader) bool { return wrapper.IsShader(s) }

func IsTexture(t Texture) bool { return wrapper.IsTexture(t) }

func LineWidth(width float32) { wrapper.LineWidth(width) }

func LinkProgram(p Program) { wrapper.LinkProgram(p) }

func PixelStorei(pname Enum, param int32) { wrapper.PixelStorei(pname, param) }

func PolygonOffset(factor, units float32) { wrapper.PolygonOffset(factor, units) }

func ReadPixels(dst []byte, x, y, width, height int, format, ty Enum) {
	wrapper.ReadPixels(dst, x, y, width, height, format, ty)
}

func ReleaseShaderCompiler() { wrapper.ReleaseShaderCompiler() }

func RenderbufferStorage(target, internalFormat Enum, width, height int) {
	wrapper.RenderbufferStorage(target, internalFormat, width, height)
}

func SampleCoverage(value float32, invert bool) { wrapper.SampleCoverage(value, invert) }

func Scissor(x, y, width, height int32) { wrapper.Scissor(x, y, width, height) }

func ShaderSource(s Shader, src string) { wrapper.ShaderSource(s, src) }

func StencilFunc(fn Enum, ref int, mask uint32) { wrapper.StencilFunc(fn, ref, mask) }

func StencilFuncSeparate(face, fn Enum, ref int, mask uint32) {
	wrapper.StencilFuncSeparate(face, fn, ref, mask)
}

func StencilMask(mask uint32) { wrapper.StencilMask(mask) }

func StencilMaskSeparate(face Enum, mask uint32) { wrapper.StencilMaskSeparate(face, mask) }

func StencilOp(fail, zfail, zpass Enum) { wrapper.StencilOp(fail, zfail, zpass) }

func StencilOpSeparate(face, sfail, dpfail, dppass Enum) {
	wrapper.StencilOpSeparate(face, sfail, dpfail, dppass)
}

func TexImage2D(target Enum, level int, internalFormat int, width, height int, format Enum, ty Enum, data []byte) {
	wrapper.TexImage2D(target, level, internalFormat, width, height, format, ty, data)
}

func TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data []byte) {
	wrapper.TexSubImage2D(target, level, x, y, width, height, format, ty, data)
}

func TexParameterf(target, pname Enum, param float32) { wrapper.TexParameterf(target, pname, param) }

func TexParameterfv(target, pname Enum, params []float32) {
	wrapper.TexParameterfv(target, pname, params)
}

func TexParameteri(target, pname Enum, param int) { wrapper.TexParameteri(target, pname, param) }

func TexParameteriv(target, pname Enum, params []int32) {
	wrapper.TexParameteriv(target, pname, params)
}

func Uniform1f(dst Uniform, v float32) { wrapper.Uniform1f(dst, v) }

func Uniform1fv(dst Uniform, src []float32) { wrapper.Uniform1fv(dst, src) }

func Uniform1i(dst Uniform, v int) { wrapper.Uniform1i(dst, v) }

func Uniform1iv(dst Uniform, src []int32) { wrapper.Uniform1iv(dst, src) }

func Uniform2f(dst Uniform, v0, v1 float32) { wrapper.Uniform2f(dst, v0, v1) }

func Uniform2fv(dst Uniform, src []float32) { wrapper.Uniform2fv(dst, src) }

func Uniform2i(dst Uniform, v0, v1 int) { wrapper.Uniform2i(dst, v0, v1) }

func Uniform2iv(dst Uniform, src []int32) { wrapper.Uniform2iv(dst, src) }

func Uniform3f(dst Uniform, v0, v1, v2 float32) { wrapper.Uniform3f(dst, v0, v1, v2) }

func Uniform3fv(dst Uniform, src []float32) { wrapper.Uniform3fv(dst, src) }

func Uniform3i(dst Uniform, v0, v1, v2 int32) { wrapper.Uniform3i(dst, v0, v1, v2) }

func Uniform3iv(dst Uniform, src []int32) { wrapper.Uniform3iv(dst, src) }

func Uniform4f(dst Uniform, v0, v1, v2, v3 float32) { wrapper.Uniform4f(dst, v0, v1, v2, v3) }

func Uniform4fv(dst Uniform, src []float32) { wrapper.Uniform4fv(dst, src) }

func Uniform4i(dst Uniform, v0, v1, v2, v3 int32) { wrapper.Uniform4i(dst, v0, v1, v2, v3) }

func Uniform4iv(dst Uniform, src []int32) { wrapper.Uniform4iv(dst, src) }

func UniformMatrix2fv(dst Uniform, src []float32) { wrapper.UniformMatrix2fv(dst, src) }

func UniformMatrix3fv(dst Uniform, src []float32) { wrapper.UniformMatrix3fv(dst, src) }

func UniformMatrix4fv(dst Uniform, src []float32) { wrapper.UniformMatrix4fv(dst, src) }

func UseProgram(p Program) { wrapper.UseProgram(p) }

func ValidateProgram(p Program) { wrapper.ValidateProgram(p) }

func VertexAttrib1f(dst Attrib, x float32) { wrapper.VertexAttrib1f(dst, x) }

func VertexAttrib1fv(dst Attrib, src []float32) { wrapper.VertexAttrib1fv(dst, src) }

func VertexAttrib2f(dst Attrib, x, y float32) { wrapper.VertexAttrib2f(dst, x, y) }

func VertexAttrib2fv(dst Attrib, src []float32) { wrapper.VertexAttrib2fv(dst, src) }

func VertexAttrib3f(dst Attrib, x, y, z float32) { wrapper.VertexAttrib3f(dst, x, y, z) }

func VertexAttrib3fv(dst Attrib, src []float32) { wrapper.VertexAttrib3fv(dst, src) }

func VertexAttrib4f(dst Attrib, x, y, z, w float32) { wrapper.VertexAttrib4f(dst, x, y, z, w) }

func VertexAttrib4fv(dst Attrib, src []float32) { wrapper.VertexAttrib4fv(dst, src) }

func VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride, offset int) {
	wrapper.VertexAttribPointer(dst, size, ty, normalized, stride, offset)
}

func Viewport(x, y, width, height int) { wrapper.Viewport(x, y, width, height) }

func GetUniformBlockIndex(p Program, name string) uint32 {
	return wrapper.GetUniformBlockIndex(p, name)
}

func UniformBlockBinding(p Program, index, bind uint32) {
	wrapper.UniformBlockBinding(p, index, bind)
}

func BindBufferBase(target Enum, n uint32, b Buffer) {
	wrapper.BindBufferBase(target, n, b)
}

func GetActiveUniformi(p Program, index uint32, pname Enum) int32 {
	return wrapper.GetActiveUniformi(p, index, pname)
}

func GetActiveUniformBlockiv(p Program, index uint32, pname Enum, params []int32) {
	wrapper.GetActiveUniformBlockiv(p, index, pname, params)
}

func GetActiveUniformBlockName(p Program, index uint32) string {
	return wrapper.GetActiveUniformBlockName(p, index)
}

func DrawArraysInstanced(mode Enum, first, count, primcount uint32) {
	wrapper.DrawArraysInstanced(mode, first, count, primcount)
}

func DrawElementsInstanced(mode Enum, count uint32, typ Enum, offset, primcount uint32) {
	wrapper.DrawElementsInstanced(mode, count, typ, offset, primcount)
}

func VertexAttribDivisor(index Attrib, divisor uint32) {
	wrapper.VertexAttribDivisor(index, divisor)
}

/*
func TexImage3D(target Enum, level, internalFormat, width, height, depth int, format, ty Enum, data []byte) {
	wrapper.TexImage3D(target, level, internalFormat, width, height, depth, format, ty, data)
}

func FramebufferTextureLayer(target Enum, attachment Enum, texture Texture, level int, layer int) {
	wrapper.FramebufferTextureLayer(target, attachment, texture, level, layer)
}
*/
