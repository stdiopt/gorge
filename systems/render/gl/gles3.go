package gl

// Context3 opengles3/webgl2 funcs
// Not all
type Context3 interface {
	Context

	GetActiveUniformi(p Program, index uint32, pname Enum) int32
	// glGetUniformBlockIndex retrieves the index of a uniform block within program.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockIndex.xhtml
	GetUniformBlockIndex(p Program, name string) uint32

	// UniformBlockBinding assign a binding point to an active uniform block
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockBinding.xhtml
	UniformBlockBinding(p Program, index, bind uint32)

	// BindBufferBase bind a buffer object to an indexed buffer target
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBufferBase.xhtml
	BindBufferBase(target Enum, n uint32, b Buffer)

	// glGetActiveUniformBlockiv — query information about an active uniform block
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveUniformBlockiv.xhtml
	GetActiveUniformBlockiv(p Program, index uint32, pname Enum, params []int32)

	// glGetActiveUniformBlockName — retrieve the name of an active uniform block
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveUniformBlockName.xhtml
	GetActiveUniformBlockName(p Program, index uint32) string

	// DrawArraysInstanced draw multiple instances of a range of elements
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
	DrawArraysInstanced(mode Enum, first, count, primcount uint32)

	// DrawElementsInstanced — draw multiple instances of a set of elements
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawElementsInstanced.xhtml
	DrawElementsInstanced(mode Enum, count uint32, typ Enum, offset, primcount uint32)

	// VertexAttribDivisor  modify the rate at which generic vertex attributes
	// advance during instanced rendering
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
	VertexAttribDivisor(index Attrib, divisor uint32)
}
