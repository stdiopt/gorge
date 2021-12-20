package renderpl

import (
	"github.com/stdiopt/gorge/systems/render"
	"github.com/stdiopt/gorge/systems/render/gl"
)

// Render Renders the geometry (forward)
func Render(r *render.Context, next render.StepFunc) render.StepFunc {
	return func(ri *render.Step) {
		for _, qi := range ri.QueuesIndex {
			renderables := ri.Queues[qi].Renderables
			for _, re := range renderables {
				mlen := re.Instances.Len()
				if mlen == 0 {
					continue
				}

				r.SetupShader(ri, re)

				vao := re.VAO(nil)
				gl.BindVertexArray(vao)
				drawMode := render.DrawMode(re.Renderable().Mesh.GetDrawMode())
				r.Draw(drawMode, re.VBO(), uint32(mlen))
				gl.BindVertexArray(gl.Null)
			}
		}
		next(ri)
	}
}
