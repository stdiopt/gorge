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

// Context3 opengles3/webgl2 funcs
// Not all
type Context3 interface {
	Context

	// glGetUniformBlockIndex retrieves the index of a uniform block within program.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockIndex.xhtml
	GetUniformBlockIndex(p Program, name string) int

	// UniformBlockBinding assign a binding point to an active uniform block
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockBinding.xhtml
	UniformBlockBinding(p Program, index, bind int)

	// BindBufferBase bind a buffer object to an indexed buffer target
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBufferBase.xhtml
	BindBufferBase(target Enum, n uint32, b Buffer)

	// DrawArraysInstanced draw multiple instances of a range of elements
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
	DrawArraysInstanced(mode Enum, first, count, primcount int)

	// DrawElementsInstanced — draw multiple instances of a set of elements
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawElementsInstanced.xhtml
	DrawElementsInstanced(mode Enum, count int, typ Enum, offset, primcount int)

	// VertexAttribDivisor  modify the rate at which generic vertex attributes
	// advance during instanced rendering
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
	VertexAttribDivisor(index Attrib, divisor int)

	// BufferDataX will type switch the interface and select the proper type
	// {lpf} Custom func
	// TODO: Should move this Elsewhere as an helper
	BufferDataX(target Enum, d interface{}, usage Enum)
}
