package renderpl

import (
	"sort"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/systems/render"
	"github.com/stdiopt/gorge/systems/render/gl"
)

// EachCamera executes a pipeline per existing system cameras.
func EachCamera(pipes ...PipelineFunc) PipelineFunc {
	return func(r *render.Context, next render.StepFunc) render.StepFunc {
		eachCamera := Pipeline(r, pipes...)
		return func(p *render.Step) {
			sort.Sort(cameraSorter(r.Cameras.Items()))
			for _, c := range r.Cameras.Items() {
				p.Camera = c
				eachCamera(p)
			}
			next(p)
		}
	}
}

// PrepareCamera should prepare the inner renderable list for this specific camera
// - check if we will render the instance
// - append instance to render state
// - prepare upload buffer
// - upload transform and color attribs
func PrepareCamera(r *render.Context, next render.StepFunc) render.StepFunc {
	return func(p *render.Step) {
		cam := p.Camera.Camera()

		p.Viewport = cam.CalcViewport(r.Gorge().ScreenSize())
		width := p.Viewport[2]
		height := p.Viewport[3]

		// Defaults for default material
		p.Lights = r.Lights.Items()

		mat := p.Camera.Mat4()
		p.Projection = cam.ProjectionWithAspect(width / height)
		p.View = mat.Inv()
		p.CamPos = mat.Col(3).Vec3()
		p.Ambient = cam.ClearColor
		VP := p.Projection.Mul(p.View)

		p.Props["VP"] = VP
		p.Props["ambient"] = p.View
		p.Props["viewPos"] = p.CamPos

		p.CameraUBO.WriteOffset("VP", VP)
		p.CameraUBO.WriteOffset("ambient", cam.ClearColor)
		p.CameraUBO.WriteOffset("viewPos", p.CamPos)
		p.CameraUBO.Flush()

		p.Ubos["Camera"] = p.CameraUBO.ID()

		// ri.renderables = ri.renderables[:0]

		{ // queues
			// we can't reset just yet if we don't delete from the thingie
			// ri.QueuesIndex = ri.QueuesIndex[:0]
			for _, q := range p.Queues {
				q.Renderables = q.Renderables[:0]
			}
		}

		// ri.renderables = ri.renderables[:0]
		camMask := render.CullMask(cam.CullMask)
		for _, re := range r.Renderables {
			// Ignore if there is no instances
			if re.Instances.Len() == 0 {
				continue
			}

			// Mask check to see if we will render this on this camera
			reMask := render.CullMask(re.Renderable().CullMask)
			if reMask&camMask == 0 {
				continue
			}
			// Check if we already processed this in some previous camera
			// we don't need to reupload transform attribute
			if p.RenderNumber != re.RenderNumber {
				re.Update(p)
			}

			// If VBO is nil we skip
			if v := re.VBO(); v == nil || v.VertexLen == 0 {
				continue
			}

			// queue index
			qi := re.Renderable().Queue
			// Select queue to insert
			q, ok := p.Queues[qi]
			if !ok {
				q = &render.Queue{
					Renderables: []*render.RenderableGroup{},
				}
				p.Queues[qi] = q
				p.QueuesIndex = append(p.QueuesIndex, qi)
			}
			// Sort insert
			q.Renderables = append(q.Renderables, re)
		}
		// NEW: Should sort By Order renderable
		// we could eventually add Zsorter here too
		for _, q := range p.Queues {
			sort.Sort(renderableGroupSorter(q.Renderables))
		}

		sort.Ints(p.QueuesIndex)
		next(p)
	}
}

// ClearCamera returns the stage that clears the renderer based on camera.
func ClearCamera(r *render.Context, next render.StepFunc) render.StepFunc {
	// Get skybox AND procedural skybox renderer
	skyBox := CameraSkybox(r, "envMap")
	return func(ri *render.Step) {
		// Render a SkyboxMaterial quad

		cam := ri.Camera.Camera()
		gl.Viewport(
			int(ri.Viewport[0]),
			int(ri.Viewport[1]),
			int(ri.Viewport[2]),
			int(ri.Viewport[3]),
		)

		gl.Enable(gl.SCISSOR_TEST)
		gl.Scissor(
			int32(ri.Viewport[0]),
			int32(ri.Viewport[1]),
			int32(ri.Viewport[2]),
			int32(ri.Viewport[3]),
		)
		// Reset maskings
		gl.ColorMask(true, true, true, true)
		gl.DepthMask(true)

		// XXX: New stencil test clearing, defaults to 0
		// This might be optional on camera, not sure
		// gl.StencilMask(0xFF)
		// gl.ClearStencil(0)
		// gl.Clear(gl.STENCIL_BUFFER_BIT)

		switch cam.ClearFlag {
		case gorge.ClearSkybox:
			// Based on camera material
			skyBox(ri)
		case gorge.ClearColor:
			gl.ClearColor(cam.ClearColor[0], cam.ClearColor[1], cam.ClearColor[2], 1)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		case gorge.ClearDepthOnly:
			gl.Clear(gl.DEPTH_BUFFER_BIT)
		case gorge.ClearNothing:
			// Nothing duh
		}
		gl.Disable(gl.SCISSOR_TEST)
		next(ri)
	}
}

// CameraSkybox Regular skybox using ri "envMap" cube sample
func CameraSkybox(r *render.Context, srcMap string) render.StepFunc {
	shader := r.NewShader(static.Shaders.CubeEnv)

	cubeVBO := r.GetVBO(gorge.NewMesh(&gorge.MeshData{
		Format:   gorge.VertexFormatP(),
		Vertices: skyboxVert,
	}))

	skyboxVAO := gl.CreateVertexArray()

	gl.BindVertexArray(skyboxVAO)
	cubeVBO.BindAttribs(shader)
	gl.BindVertexArray(gl.Null)

	return func(ri *render.Step) {
		tex := ri.Samplers[srcMap]
		if tex == nil {
			return
		}
		VP := ri.Projection.Mul(ri.View.Mat3().Mat4())

		// Odd to be here but ... we could run SetupMaterial here.
		gl.Disable(gl.STENCIL_TEST)
		gl.DepthMask(true)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		gl.DepthMask(false)
		gl.Disable(gl.BLEND)
		{
			gl.FrontFace(gl.CCW)
			gl.BindVertexArray(skyboxVAO)
			shader.Bind()
			shader.Set("VP", VP)
			shader.Set("skybox", 0)

			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_CUBE_MAP, tex.ID)
			r.Draw(gl.TRIANGLES, cubeVBO, 1)
			gl.BindVertexArray(gl.Null)
			gl.FrontFace(gl.CW)
		}
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthMask(true)
	}
}
