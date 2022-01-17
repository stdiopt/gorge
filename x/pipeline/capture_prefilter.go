package pipeline

import (
	"math"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/systems/render"
	"github.com/stdiopt/gorge/systems/render/gl"
)

func (pl *PL) CapturePrefilter(src, target string) PipelineFunc {
	return func(r *render.Context, next StepFunc) StepFunc {
		size := 128

		prefilterSD := &gorge.ShaderData{Src: static.MustData("shaders/ibl/prefilter.glsl")}
		prefilterShader := r.NewShader(prefilterSD)

		const level = 0
		const internalFormat = gl.RGBA16F
		const srcFormat = gl.RGBA
		const srcType = gl.FLOAT

		prefilterMap := pl.createCubeMap(size, true)

		return func(p *Step) {
			tex := p.Samplers[src]
			if tex == nil {
				next(p)
				return
			}

			prefilterShader.Bind()
			prefilterShader.Set("environmentMap", 0)
			prefilterShader.Set("projection", camProj)
			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_CUBE_MAP, tex.ID)

			gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, pl.captureFBO)

			{ // For some reason this works on webgl
				gl.BindRenderbuffer(gl.RENDERBUFFER, pl.captureRBO)
				gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, size, size)
				gl.FramebufferRenderbuffer(gl.DRAW_FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, pl.captureRBO)
			}
			maxMipLevels := 5
			for mip := 0; mip < maxMipLevels; mip++ {
				mipWidth := int32(128 * math.Pow(0.5, float64(mip)))
				mipHeight := int32(128 * math.Pow(0.5, float64(mip)))
				gl.BindRenderbuffer(gl.RENDERBUFFER, pl.captureRBO)
				gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, int(mipWidth), int(mipHeight))
				gl.Viewport(0, 0, mipWidth, mipHeight)

				roughness := float32(mip) / float32(maxMipLevels-1)
				prefilterShader.Set("roughness", roughness)
				for i := 0; i < 6; i++ {
					prefilterShader.Set("view", camTargets[i])
					gl.FramebufferTexture2D(
						gl.DRAW_FRAMEBUFFER,
						gl.COLOR_ATTACHMENT0,
						gl.Enum(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i),
						prefilterMap.ID, mip,
					)
					gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
					gl.Disable(gl.BLEND)
					gl.Disable(gl.DEPTH_TEST)
					pl.renderCube()
				}
			}
			gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, gl.Null)

			p.Samplers[target] = prefilterMap
			next(p)
		}
	}
}
