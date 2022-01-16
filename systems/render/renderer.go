// Package render mostly agnostic gl renderer with a couple of exceptions due
// to some limitations
//
// It will listen for ECS events
package render

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/setlist"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/render/bufutil"
	"github.com/stdiopt/gorge/systems/render/gl"
)

// Render thing
type Render struct {
	gorge          *gorge.Context
	Cameras        setlist.SetList[Camera]
	Lights         setlist.SetList[Light]
	Renderables    []*RenderableGroup
	renderablesMap map[*gorge.RenderableComponent]*RenderableGroup

	// TODO create global gl state here to avoid calling gl stuff every material
	// cgo, js, etc.. calls might be expensive

	buffers  *bufferManager
	vbos     *vboManager
	textures *textureManager
	shaders  *shaderManager

	DrawCalls     int
	DisableRender bool
	renderInfo    Step

	// Or Renderer Render
	InitStage   StepFunc
	RenderStage StepFunc
}

func newRenderer(g *gorge.Context) *Render {
	gl.Enable(gl.PROGRAM_POINT_SIZE)
	gl.Hint(gl.GENERATE_MIPMAP_HINT, gl.NICEST)
	gl.LineWidth(1.4) // Not implemented by all drivers

	gl.ClearDepthf(1)
	gl.DepthFunc(gl.LESS)

	gl.FrontFace(gl.CW)
	gl.CullFace(gl.BACK)

	// should be on material
	gl.Disable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)

	bm := newBufferManager(g)
	vbos := newVBOManager(g, bm)
	shaders := newShaderManager(g, vbos)
	textures := newTextureManager(g)

	cameraUBO := bufutil.NewNamedOffset(
		bufutil.NewCached[byte](bm.New(gl.UNIFORM_BUFFER, gl.DYNAMIC_DRAW)),
		96,
		bufutil.OffsetSpec{
			"VP":      0,
			"ambient": 64,
			"viewPos": 80,
		},
	)
	// Setup default pipeline
	return &Render{
		gorge:          g,
		Renderables:    []*RenderableGroup{},
		renderablesMap: map[*gorge.RenderableComponent]*RenderableGroup{},
		shaders:        shaders,
		textures:       textures,
		vbos:           vbos,
		buffers:        bm,
		renderInfo: Step{
			CameraUBO: cameraUBO,
			Viewport:  m32.Vec4{},
			Queues:    map[int]*Queue{},
			Ubos:      map[string]gl.Buffer{},
			Props:     map[string]any{},
			Samplers:  map[string]*Texture{},
		},
	}
}

// Init initializes the renderer by calling the proper stages and etc...
func (r *Render) Init() {
	if r.InitStage == nil {
		return
	}
	r.InitStage(&r.renderInfo)
}

// Render starts the render pipeline
func (r *Render) Render() {
	r.DrawCalls = 0
	r.renderInfo.RenderNumber++
	r.RenderStage(&r.renderInfo)
}

// Gorge returns the gorge context.
func (r *Render) Gorge() *gorge.Context {
	return r.gorge
}

// NewShader returns a shader based on ShaderData
func (r *Render) NewShader(gs *gorge.ShaderData) *Shader {
	return r.shaders.New(gs)
}

// NewBuffer creates and returns a new buffer.
func (r *Render) NewBuffer(target, usage gl.Enum) Bufferer {
	return r.buffers.New(target, usage)
}

// GetVBO returns a VBO based on mesher
func (r *Render) GetVBO(mesh *gorge.Mesh) *VBO {
	v, _ := r.vbos.Get(mesh)
	return v
}

// AddCamera adds a camera.
func (r *Render) AddCamera(camera Camera) {
	r.Cameras.Add(camera)
}

// RemoveCamera removes a specific camera if exists.
func (r *Render) RemoveCamera(camera Camera) {
	r.Cameras.Remove(camera)
}

// AddLight adds a light to the light list.
func (r *Render) AddLight(light Light) {
	r.Lights.Add(light)
}

// RemoveLight if exists.
func (r *Render) RemoveLight(light Light) {
	r.Lights.Remove(light)
}

// AddRenderable adds a renderable instance
func (r *Render) AddRenderable(re Renderable) {
	renderable := re.Renderable()
	key := renderable

	// Find best group for renderable
	if group, ok := r.renderablesMap[key]; ok {
		// We are adding stuff so we reset the transfer buf
		// (Grow method would be better)
		group.Add(re)

		return
	}

	// New Renderable Instance
	// NewInstance basically
	group := &RenderableGroup{
		renderer:   r,
		renderable: renderable,
	}
	if ok := group.init(); !ok {
		return
	}
	group.Add(re)

	r.Renderables = append(r.Renderables, group)
	r.renderablesMap[key] = group
}

// RemoveRenderable remove to renderable renderable
func (r *Render) RemoveRenderable(re Renderable) {
	// Go trough instances and remove from there
	renderable := re.Renderable()
	key := renderable

	instance, ok := r.renderablesMap[key]
	if !ok {
		return
	}

	instance.Remove(re)

	// Remove instance from system
	if instance.Instances.Len() == 0 {
		instance.Destroy()
		delete(r.renderablesMap, key)
		for i, is := range r.Renderables {
			if is == instance {
				t := r.Renderables
				r.Renderables = append(r.Renderables[:i], r.Renderables[i+1:]...)
				t[len(t)-1] = nil // remove last one since it was copied
				break
			}
		}
	}
}

// SetInitStage sets Init stage for renderer, used to prepare stuff.
func (r *Render) SetInitStage(s StepFunc) {
	r.InitStage = s
}

// SetRenderStage sets the renderer stage.
func (r *Render) SetRenderStage(s StepFunc) {
	r.RenderStage = s
}

// PassShadow passes geometry with a specific shader
// TODO: we are checking DisableShadow here
func (r *Render) PassShadow(ri *Step, s *Shader) {
	gl.Enable(gl.BLEND)
	for _, qi := range ri.QueuesIndex {
		renderables := ri.Queues[qi].Renderables

		// Get depthCube shader
		for _, group := range renderables {
			if group.Renderable().DisableShadow {
				continue
			}
			mlen := group.Instances.Len()
			if mlen == 0 {
				continue
			}

			vao := group.VAO(s)

			mat := group.Renderable().Material
			alphaCutoff := float32(0)
			if v, ok := mat.Get("u_AlphaCutoff").(float32); ok {
				alphaCutoff = v
			}
			s.Set("u_AlphaCutoff", alphaCutoff)

			// TODO: experimental: Pass albedo here, or the alpha channel thing so
			// the shadow will behave with alphas
			if tex := mat.GetTexture("albedoMap"); tex != nil {
				gl.ActiveTexture(gl.TEXTURE0)
				r.textures.Bind(tex)
				s.Set("albedoMap", 0)
			} else if tex := mat.GetTexture("u_BaseColorSampler"); tex != nil {
				gl.ActiveTexture(gl.TEXTURE0)
				r.textures.Bind(tex)
				s.Set("albedoMap", 0)
			} else {
				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, gl.Null)
			}

			drawMode := DrawMode(group.Renderable().GetDrawMode())
			gl.BindVertexArray(vao)
			r.Draw(drawMode, group.vbo, group.Count)
			gl.BindVertexArray(gl.Null)
		}
	}
}

// SetupShader sets the current material uniforms
//
// TODO: {lpf} Maybe set on renderable only since we could use it in render
// pipeline to draw the sky box
func (r *Render) SetupShader(
	ri *Step,
	group *RenderableGroup,
) {
	mesh := group.Renderable().Mesh
	mat := group.Renderable().Material
	// Get shader by thing
	shader := group.shader

	// avoid this calls by storing the state globally?
	switch mat.Blend {
	case gorge.BlendDisable:
		gl.Disable(gl.BLEND)
	case gorge.BlendOneOneMinusSrcAlpha:
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	case gorge.BlendOneOne:
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.ONE, gl.ONE)
	}

	switch mat.Depth {
	case gorge.DepthReadWrite: // do depth test and writes do depth mask
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthMask(true)
	case gorge.DepthRead: // doesn't write to depth but still tests
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthMask(false)
	case gorge.DepthNone:
		gl.Disable(gl.DEPTH_TEST)
		gl.DepthMask(false)
	}
	if mat.ColorMask == nil {
		gl.ColorMask(true, true, true, true)
	} else {
		gl.ColorMask(mat.ColorMask[0], mat.ColorMask[1], mat.ColorMask[2], mat.ColorMask[3])
	}

	// This can be done in shader?
	if mat.DoubleSided {
		gl.Disable(gl.CULL_FACE)
	} else {
		gl.Enable(gl.CULL_FACE)
	}

	if group.vbo.FrontFacing == gorge.FrontFacingCCW {
		gl.FrontFace(gl.CCW)
	} else {
		gl.FrontFace(gl.CW)
	}
	// New stencil test
	if mat.Stencil != nil {
		gl.Enable(gl.STENCIL_TEST)
		gl.StencilMask(mat.Stencil.WriteMask)
		gl.StencilFunc(StencilFunc(mat.Stencil.Func), mat.Stencil.Ref, mat.Stencil.ReadMask)
		gl.StencilOp(
			StencilOp(mat.Stencil.Fail),
			StencilOp(mat.Stencil.ZFail),
			StencilOp(mat.Stencil.ZPass),
		)
		ri.StencilDirty = true
	} else {
		// gl.StencilMask(0xFF)
		// gl.StencilFunc(gl.ALWAYS, 0, 0xFF)
		// gl.StencilOp(gl.KEEP, gl.KEEP, gl.KEEP)
		gl.Disable(gl.STENCIL_TEST)
	}

	// Get shader
	shader.Bind()
	// This is new since I don't have any a_Attribute left
	// Pick only the first
	re := group.Front() // Get an instance
	modelMat4 := re.Mat4()

	for k, u := range shader.uniforms {
		if u.sampler {
			continue
		}
		if v := mat.Get(k); v != nil {
			shader.Set(k, v)
			continue
		}
		if v := mesh.Get(k); v != nil {
			shader.Set(k, v)
			continue
		}
		// Extra case per model
		switch k {
		case "u_ModelMatrix":
			shader.Set(k, modelMat4)
			continue
		case "u_NormalMatrix":
			shader.Set(k, modelMat4.Inv().Transpose())
			continue
		}

		// Global props
		if v, ok := ri.Props[k]; ok {
			shader.Set(k, v)
			continue
		}

		// Might be slower but we avoid reusing last material shader values in
		// this material, should not be a sampler
		shader.Set(k, nil)
	}

	// This will automatically Set shader variable like:
	// if we set a somethingMap we will send has_somethingMap bool
	// shader will ignore it if uniform doesn't exists
	for i, k := range shader.samplers {
		gl.ActiveTexture(gl.TEXTURE0 + gl.Enum(i))
		hasK := "has_" + k // disable this
		shader.Set(k, i)

		if t := mat.GetTexture(k); t != nil {
			// TODO: {lpf} do not use uniforms for this
			// we can use Defines previously
			shader.Set(hasK, true) // disable this
			r.textures.Bind(t)
			continue
		}

		// Global samplers
		if tex, ok := ri.Samplers[k]; ok {
			shader.Set(hasK, true) // disable this
			gl.BindTexture(tex.Type, tex.ID)
			continue
		}

		// Bind the default texture
		r.textures.Bind(nil)
	}

	uboi := uint32(0)
	for k := range shader.ubos {
		id, ok := ri.Ubos[k]
		if !ok {
			continue
		}
		index := gl.GetUniformBlockIndex(shader.program, k)
		gl.UniformBlockBinding(shader.program, index, uboi)
		gl.BindBufferBase(gl.UNIFORM_BUFFER, uboi, id)
		uboi++
	}
}

// Draw the main draw function
func (r *Render) Draw(mode gl.Enum, v *VBO, count uint32) {
	// Draw renderable instance perhaps instead of VBO
	// so we can rebind stuff
	r.DrawCalls++
	if v.ElementsLen > 0 {
		gl.DrawElementsInstanced(mode, v.ElementsLen, v.ElementsType, 0, count)
	} else {
		gl.DrawArraysInstanced(mode, 0, v.VertexLen, count)
	}
}
