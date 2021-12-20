package pipeline

import (
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/systems/render"
	"github.com/stdiopt/gorge/systems/render/gl"
)

// RenderCube with a specific CubeMap sampler
func (pl *PL) RenderCube(srcMap string, vp m32.Vec4) PipelineFunc {
	return func(r *render.Context, next PassFunc) PassFunc {
		shader := r.NewShader(static.Shaders.CubeEnv)

		return func(p *Pass) {
			tex := p.Samplers[srcMap]
			if tex == nil {
				next(p)
				return
			}

			VP := p.Projection.Mul(p.View.Mat3().Mat4())

			x := int(p.Viewport[2] * vp[0])
			y := int(p.Viewport[3] * vp[1])
			w := int(p.Viewport[2] * vp[2])
			h := int(p.Viewport[3] * vp[3])
			gl.Viewport(x, y, w, h)
			gl.Clear(gl.DEPTH_BUFFER_BIT)
			gl.DepthMask(false)
			{
				shader.Bind()
				shader.Set("VP", VP)
				shader.Set("skybox", 0)

				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_CUBE_MAP, tex.ID)
				pl.renderCube()
			}
			gl.DepthMask(true)
			next(p)
		}
	}
}

func (p *PL) RenderQuad(srcMap string, vp m32.Vec4) PipelineFunc {
	return func(r *render.Context, next PassFunc) PassFunc {
		emptyVAO := gl.CreateVertexArray()
		shader := r.NewShader(static.Shaders.Quad)

		return func(p *Pass) {
			tex := p.Samplers[srcMap]
			if tex == nil {
				next(p)
				return
			}
			x := int(p.Viewport[2] * vp[0])
			y := int(p.Viewport[3] * vp[1])
			w := int(p.Viewport[2] * vp[2])
			h := int(p.Viewport[3] * vp[3])
			gl.Viewport(x, y, w, h)
			gl.Clear(gl.DEPTH_BUFFER_BIT)
			gl.DepthMask(false)
			gl.FrontFace(gl.CCW)
			{
				shader.Bind()
				shader.Set("albedoMap", 0)
				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, tex.ID)
				gl.BindVertexArray(emptyVAO)
				gl.DrawArrays(gl.TRIANGLES, 0, 6)
				gl.BindVertexArray(gl.Null)
			}
			gl.DepthMask(true)
			gl.FrontFace(gl.CW)
			next(p)
		}
	}
}

func (pl *PL) RenderQuadDepth(srcMap string, vp m32.Vec4) PipelineFunc {
	return func(r *render.Context, next PassFunc) PassFunc {
		emptyVAO := gl.CreateVertexArray()
		shader := r.NewShader(static.Shaders.QuadDepth)

		return func(p *Pass) {
			tex := p.Samplers[srcMap]
			if tex == nil {
				next(p)
				return
			}

			x := int(p.Viewport[0] * vp[0])
			y := int(p.Viewport[1] * vp[1])
			w := int(p.Viewport[2] * vp[2])
			h := int(p.Viewport[3] * vp[3])
			gl.Viewport(x, y, w, h)
			gl.Scissor(int32(x), int32(y), int32(w), int32(h))
			gl.Enable(gl.SCISSOR_TEST)

			gl.Clear(gl.DEPTH_BUFFER_BIT | gl.COLOR_BUFFER_BIT)
			gl.DepthMask(false)
			gl.FrontFace(gl.CCW)
			{
				shader.Bind()
				shader.Set("near_plane", -50)
				shader.Set("far_plane", -50)
				shader.Set("perspective", 0)
				shader.Set("albedoMap", 0) // texture unit 0
				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, tex.ID)

				gl.BindVertexArray(emptyVAO)
				gl.DrawArrays(gl.TRIANGLES, 0, 6)
				gl.BindVertexArray(gl.Null)
			}
			gl.DepthMask(true)
			gl.FrontFace(gl.CW)
			gl.Disable(gl.SCISSOR_TEST)
			next(p)
		}
	}
}
