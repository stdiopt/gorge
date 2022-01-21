// Package pipeline implements learnogl env mapping stuff
package pipeline

import (
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/render"
	"github.com/stdiopt/gorge/systems/render/gl"
	"github.com/stdiopt/gorge/systems/render/renderpl"
	"github.com/stdiopt/gorge/systems/resource"
)

type (
	Render       = render.Render
	Step         = render.Step
	PipelineFunc = renderpl.PipelineFunc
	StepFunc     = render.StepFunc
	VBO          = render.VBO
	Shader       = render.Shader
)

var (
	camTargets = []gm.Mat4{
		gm.LookAt(gm.Vec3{}, gm.Vec3{1, 0, 0}, gm.Vec3{0, -1, 0}),
		gm.LookAt(gm.Vec3{}, gm.Vec3{-1, 0, 0}, gm.Vec3{0, -1, 0}),

		gm.LookAt(gm.Vec3{}, gm.Vec3{0, 1, 0}, gm.Vec3{0, 0, 1}),
		gm.LookAt(gm.Vec3{}, gm.Vec3{0, -1, 0}, gm.Vec3{0, 0, -1}),

		gm.LookAt(gm.Vec3{}, gm.Vec3{0, 0, 1}, gm.Vec3{0, -1, 0}),
		gm.LookAt(gm.Vec3{}, gm.Vec3{0, 0, -1}, gm.Vec3{0, -1, 0}),
	}
	camProj = gm.Perspective(90, 1, .1, 10)
)

// PL pipeline instance that will reference VertexArray and VBOs.
type PL struct {
	render   *render.Context
	resource *resource.Context

	captureFBO gl.Framebuffer
	captureRBO gl.Renderbuffer

	cubeVAO gl.VertexArray
	cubeVBO *VBO
	quadVAO gl.VertexArray
	quadVBO *VBO
}

func New(r *render.Context, res *resource.Context) *PL {
	captureFBO := gl.CreateFramebuffer()
	captureRBO := gl.CreateRenderbuffer()

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, captureFBO)
	{
		gl.BindRenderbuffer(gl.RENDERBUFFER, captureRBO)
		gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, 512, 512)
		gl.FramebufferRenderbuffer(gl.DRAW_FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, captureRBO)
	}

	gl.BindRenderbuffer(gl.RENDERBUFFER, gl.Null)

	cubeVBO := r.GetVBO(cubeMesh)
	cubeVAO := gl.CreateVertexArray()
	gl.BindVertexArray(cubeVAO)
	cubeVBO.BindAttribs(nil)
	gl.BindVertexArray(gl.Null)

	quadVBO := r.GetVBO(quadMesh)
	quadVAO := gl.CreateVertexArray()
	gl.BindVertexArray(quadVAO)
	quadVBO.BindAttribs(nil)
	gl.BindVertexArray(gl.Null)

	return &PL{
		render:     r,
		resource:   res,
		captureFBO: captureFBO,
		captureRBO: captureRBO,
		cubeVBO:    cubeVBO,
		cubeVAO:    cubeVAO,

		quadVBO: quadVBO,
		quadVAO: quadVAO,
	}
}

func (p *PL) renderCube() {
	gl.BindVertexArray(p.cubeVAO)
	gl.FrontFace(gl.CCW)
	p.render.Draw(gl.TRIANGLES, p.cubeVBO, 1)
	gl.FrontFace(gl.CW)
	gl.BindVertexArray(gl.Null)
}

func (p *PL) renderQuad() {
	gl.BindVertexArray(p.quadVAO)
	gl.FrontFace(gl.CCW)
	p.render.Draw(gl.TRIANGLE_STRIP, p.quadVBO, 1)
	gl.FrontFace(gl.CW)
	gl.BindVertexArray(gl.Null)
}

func (p *PL) createCubeMap(sz int, mipmap bool) *render.Texture {
	const level = 0
	const internalFormat = gl.RGBA16F
	const srcFormat = gl.RGBA
	const srcType = gl.FLOAT

	cubeMap := &render.Texture{
		ID:   gl.CreateTexture(),
		Type: gl.TEXTURE_CUBE_MAP,
	}
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cubeMap.ID)
	for i := 0; i < 6; i++ {
		target := gl.Enum(gl.TEXTURE_CUBE_MAP_POSITIVE_X + i)
		gl.TexImage2D(
			target, level, internalFormat,
			sz, sz,
			srcFormat, srcType, nil,
		)
	}
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	if mipmap && isPowerOfTwo(sz) {
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.GenerateMipmap(gl.TEXTURE_CUBE_MAP)
	}

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, gl.Null)
	return cubeMap
}

func isPowerOfTwo(v int) bool {
	return (v & (v - 1)) == 0
}
