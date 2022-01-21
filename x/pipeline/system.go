package pipeline

import (
	"github.com/stdiopt/gorge/systems/render"
	"github.com/stdiopt/gorge/systems/render/renderpl"
	"github.com/stdiopt/gorge/systems/resource"
)

func System(r *render.Context, res *resource.Context) error {
	// var pbr gorge.ShaderData
	// res.Load(&pbr, "shaders/pbr.glsl")

	p := New(r, res)

	r.SetInitStage(renderpl.Pipeline(r,
		// render.ProceduralSkyboxStage,
		p.LoadHDR("tex/hdr/newport_loft.hdr", "envMap"), // Load skybox into envMap prop
		// p.LoadSkyboxStage("envMap"),

		p.CaptureIrradiance("envMap", "irradianceMap"), // (EXTRA SLOW) Grab a cube sample from screen into irradianceMap
		p.CapturePrefilter("envMap", "prefilterMap"),   // Grab cube sample from env
		p.CaptureBRDF("brdfLUT"),                       // generate BRDF tex
	))

	r.SetRenderStage(renderpl.Pipeline(r,
		// p.RenderCube("envMap", gm.Vec4{0, 0, 1, 1}), // Draw a cube with irradiance
		// render.ProceduralSkyboxStage,
		// After procedural
		renderpl.EachCamera(
			renderpl.PrepareCamera, // Prepare/Cull
			renderpl.PrepareLights, // Render depth,shadowmaps if any
			renderpl.ClearCamera,   // Clear camera with skybox if any
			renderpl.Render,        // render Stage thing
		),
		// p.RenderCube("envMap", gm.Vec4{0, 0, .2, .2}),
		// p.RenderCube("irradianceMap", gm.Vec4{.2, 0, .2, .2}),
		// p.RenderCube("prefilterMap", gm.Vec4{.4, 0, .2, .2}),
		// p.RenderQuad("brdfLUT", gm.Vec4{.6, 0, .2, .2}),
		// p.RenderCube("depthCube[0].depthMap", gm.Vec4{.8, 0, .2, .2}),
	))
	return nil
}
