package pipeline

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/systems/render"
	"github.com/stdiopt/gorge/systems/render/gl"
)

func (pl *PL) LoadHDR(src string, target string) PipelineFunc {
	return func(r *render.Context, next StepFunc) StepFunc {
		equirectangularSD := &gorge.ShaderData{Src: static.MustData("shaders/ibl/equirectangular_cube.glsl")}
		equirectangularShader := r.NewShader(equirectangularSD)

		var tex gorge.TextureData
		if err := pl.resource.Load(&tex, src); err != nil {
			panic(err)
		}

		hdrTexture := &render.Texture{
			ID:   gl.CreateTexture(),
			Type: gl.TEXTURE_2D,
		}
		{
			gl.BindTexture(gl.TEXTURE_2D, hdrTexture.ID)
			gl.TexImage2D(
				gl.TEXTURE_2D, 0,
				gl.RGB16F,
				tex.Width, tex.Height,
				gl.RGB, gl.FLOAT, tex.PixelData) // note how we specify the texture's data value to be float

			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		}

		const texSize = 512

		cubeTex := pl.createCubeMap(texSize, false)

		return func(p *Step) {
			equirectangularShader.Bind()
			equirectangularShader.Set("equirectangularMap", 0)
			equirectangularShader.Set("projection", camProj)
			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, hdrTexture.ID)

			gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, pl.captureFBO)
			gl.Viewport(0, 0, texSize, texSize)
			for i := 0; i < 6; i++ {
				equirectangularShader.Set("view", camTargets[i].Mul(gm.Scale3D(1, -1, 1)))
				gl.FramebufferTexture2D(
					gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT0,
					gl.Enum(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i),
					cubeTex.ID, 0,
				)
				gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
				pl.renderCube()
			}
			gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, gl.Null)

			p.Props["has_"+target] = true
			p.Samplers[target] = cubeTex
			next(p)
		}
	}
}

func (pl *PL) LoadSkyboxStage(target string) PipelineFunc {
	return func(_ *render.Context, next StepFunc) StepFunc {
		srcs := []string{
			"skybox/right.jpg",
			"skybox/left.jpg",
			"skybox/top.jpg",
			"skybox/bottom.jpg",
			"skybox/front.jpg",
			"skybox/back.jpg",
		}

		cubeTex := &render.Texture{
			ID:   gl.CreateTexture(),
			Type: gl.TEXTURE_CUBE_MAP,
		}
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, cubeTex.ID)

		// Default texture params
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

		for i := 0; i < 6; i++ {
			var texData gorge.TextureData
			if err := pl.resource.Load(&texData, srcs[i]); err != nil {
				panic(err)
			}
			gl.TexImage2D(
				gl.Enum(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i),
				0, gl.RGBA, texData.Width, texData.Height,
				gl.RGBA, gl.UNSIGNED_BYTE, texData.PixelData,
			)
		}

		return func(p *Step) {
			p.Props["has_"+target] = true
			p.Samplers[target] = cubeTex
			next(p)
		}
	}
}
