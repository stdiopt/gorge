package pipeline

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/systems/render"
	"github.com/stdiopt/gorge/systems/render/gl"
)

func (pl *PL) CaptureBRDF(target string) PipelineFunc {
	return func(r *render.Context, next PassFunc) PassFunc {
		const size = 512

		brdfSD := &gorge.ShaderData{Src: static.MustData("shaders/ibl/brdf.glsl")}
		brdfShader := r.NewShader(brdfSD)

		brdfLUTTex := &render.Texture{
			ID:   gl.CreateTexture(),
			Type: gl.TEXTURE_2D,
		}
		gl.BindTexture(gl.TEXTURE_2D, brdfLUTTex.ID)

		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RG16F, size, size, gl.RG, gl.FLOAT, nil)
		// be sure to set wrapping mode to GL_CLAMP_TO_EDGE
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		return func(p *Pass) {
			gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, pl.captureFBO)

			gl.BindRenderbuffer(gl.RENDERBUFFER, pl.captureRBO)
			gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, size, size)
			gl.FramebufferTexture2D(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, brdfLUTTex.ID, 0)

			gl.Viewport(0, 0, size, size)
			brdfShader.Bind()
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			pl.renderQuad()
			gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, gl.Null)

			p.Samplers[target] = brdfLUTTex
			next(p)
		}
	}
}
