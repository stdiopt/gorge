package pipeline

import (
	"fmt"
	"log"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/systems/render"
	"github.com/stdiopt/gorge/systems/render/gl"
)

func (pl *PL) CaptureIrradiance(src, target string) PipelineFunc {
	srcTex := src
	dstTex := target
	return func(r *render.Context, next PassFunc) PassFunc {
		size := 32

		irradianceSD := &gorge.ShaderData{Src: static.MustData("shaders/ibl/irradiance_convolution.glsl")}

		irradianceShader := r.NewShader(irradianceSD)
		irradianceTex := pl.createCubeMap(size, false)

		// We have radiance shader now
		return func(p *Pass) {
			irradianceShader.Bind()
			irradianceShader.Set("environmentMap", 0)
			irradianceShader.Set("projection", camProj)
			tex := p.Samplers[srcTex]
			if tex == nil {
				next(p)
				return
			}

			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_CUBE_MAP, tex.ID)
			gl.Viewport(0, 0, size, size)
			gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, pl.captureFBO)

			{ // For some reason this makes it work on webgl
				gl.BindRenderbuffer(gl.RENDERBUFFER, pl.captureRBO)
				gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, size, size)
				gl.FramebufferRenderbuffer(gl.DRAW_FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, pl.captureRBO)
			}
			{
				for i := 0; i < 6; i++ {
					irradianceShader.Set("view", camTargets[i])
					gl.FramebufferTexture2D(
						gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT0,
						gl.Enum(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i), irradianceTex.ID, 0)
					// TODO REMOVE
					if v := gl.CheckFramebufferStatus(gl.DRAW_FRAMEBUFFER); v != gl.FRAMEBUFFER_COMPLETE {
						fmt.Printf("Frame buffer not complete")
						glerr := gl.GetError()
						log.Printf("frame buffer error: %v, %v", v, glerr)
					}

					gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

					pl.renderCube()
				}
			}
			gl.BindVertexArray(gl.Null)
			gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, gl.Null)

			p.Samplers[dstTex] = irradianceTex
			next(p)
		}
	}
}
