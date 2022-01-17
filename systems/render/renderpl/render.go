package renderpl

import (
	"github.com/stdiopt/gorge/systems/render"
	"github.com/stdiopt/gorge/systems/render/gl"
)

// Render Renders the geometry (forward)
func Render(r *render.Context, next render.StepFunc) render.StepFunc {
	return func(ri *render.Step) {
		dorender(r, ri)
		next(ri)
	}
}

// Might be better on Render struct
func dorender(r *render.Context, s *render.Step) {
	for _, qi := range s.QueuesIndex {
		// Clear stencil per queue, is it costy?!
		if s.StencilDirty { // Stencil is per queue
			gl.StencilMask(0xFF)
			gl.Clear(gl.STENCIL_BUFFER_BIT)
			s.StencilDirty = false
		}
		renderables := s.Queues[qi].Renderables
		for _, re := range renderables {
			mlen := re.Instances.Len()
			if mlen == 0 {
				continue
			}

			r.SetupShader(s, re)
			vao := re.VAO(nil)
			gl.BindVertexArray(vao)
			drawMode := render.DrawMode(re.Renderable().GetDrawMode())
			r.Draw(drawMode, re.VBO(), re.Count)
			gl.BindVertexArray(gl.Null)
		}
	}
}
