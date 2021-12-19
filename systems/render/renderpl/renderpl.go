// Package renderpl contains default rendering pipeline for gorge
package renderpl

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/render"
)

// const defCameraMask = uint32(0xFF)

// PipelineFunc middleware alike pipelining.
type PipelineFunc func(r *render.Context, next render.PassFunc) render.PassFunc

// Pipeline builds a StagerFunc from several pipelineFuncs.
func Pipeline(r *render.Context, fns ...PipelineFunc) render.PassFunc {
	if len(fns) == 0 {
		return func(ri *render.Pass) {} // End of line
	}
	next := fns[1:]
	return fns[0](r, Pipeline(r, next...))
}

// Default detault rendering pipeline
func Default(g *gorge.Context) {
	r := render.FromContext(g)
	r.SetRenderStage(Pipeline(r,
		// CubeSkyboxStage,
		ProceduralSkybox,
		EachCamera(
			PrepareCamera,
			PrepareLights,
			ClearCamera,
			Render,
		),
	))
}
