// Copyright 2019 Luis Figueiredo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package renderer mostly agnostic gl renderer with a couple of exceptions due
// to some limitations
//
// TODO: Render pipeline
// It will listen for ECS events
package renderer

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/asset"
	"github.com/stdiopt/gorge/gl"
	"github.com/stdiopt/gorge/m32"
)

var (
	ExperimentalSkybox = false
)

type (
	vec2 = m32.Vec2
	vec3 = m32.Vec3
	vec4 = m32.Vec4
	mat3 = m32.Mat3
	mat4 = m32.Mat4
	quat = m32.Quat
)
type rLight interface {
	TransformComponent() *gorge.Transform
	LightComponent() *gorge.Light
}
type rCamera interface {
	TransformComponent() *gorge.Transform
	CameraComponent() *gorge.Camera
}

type rMesh interface {
	TransformComponent() *gorge.Transform
	RenderableComponent() *gorge.Renderable
}

type instanceKey struct {
	Material *gorge.Material
	Mesh     *gorge.Mesh
}

// Batch on add
// Mesh manager
// meshInstance controller
// a renderable instance SHOULD be a combination of {Mesh,Material}
type renderableInstance struct {
	//material *material
	Material *gorge.Material
	Mesh     *gorge.Mesh
	shader   *Shader

	VAO gl.VertexArray

	vbo *vbo
	// for instancing
	TRO gl.Buffer

	// Instances
	meshes []rMesh

	attribBuf *F32TransferBuf
}

// Renderer thing
type Renderer struct {
	gorge *gorge.Gorge
	// Camera projection
	camera rCamera
	light  rLight
	// TODO: should remove this in favour of delegated funcs
	assets    *asset.System
	instances []*renderableInstance

	instancesMap map[instanceKey]*renderableInstance
	//lastInstance *renderableInstance

	textures *textureManager
	shaders  *shaderManager
	vbos     *vboManager

	g         gl.Context3
	totalTime float32

	skyboxTex    *texture
	skyboxShader *Shader
	skyboxVAO    gl.VertexArray
	skyboxEBO    []uint32
}

// System initializes gl context and attatch the handlers on manager
func System(gm *gorge.Gorge) {
	var g *gl.Wrapper
	gm.Query(func(glctx *gl.Wrapper) {
		g = glctx
	})
	if g == nil {
		panic("renderer requires persisted glctx")
	}
	// Left hand
	g.ClearDepthf(-1)
	g.DepthFunc(gl.GEQUAL)
	g.LineWidth(3)
	g.FrontFace(gl.CCW)
	g.Hint(gl.GENERATE_MIPMAP_HINT, gl.NICEST)
	// should be on material

	// Material
	g.Enable(gl.BLEND)
	g.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	g.CullFace(gl.BACK)

	// Don't use assets
	assets := asset.FromECS(gm)
	if assets == nil {
		panic("renderer requires assets system to be initialized first")
	}

	rs := &Renderer{
		gorge:        gm,
		g:            g,
		assets:       assets,
		instances:    []*renderableInstance{},
		instancesMap: map[instanceKey]*renderableInstance{},
		textures:     newTextureManager(g),
		vbos:         newVBOManager(g),
		shaders:      &shaderManager{g: g, assets: assets},
	}

	gm.Handle(rs.watchResize)
	gm.Handle(rs.handleEntityAdd)
	gm.Handle(rs.handlePostUpdate)
	gm.Handle(rs.handleAssetAdd)

	if ExperimentalSkybox {
		rs.PrepareSkybox()
	}

}

// AddMesh adds a mesh
// Mesh Manager?
func (rs *Renderer) AddMesh(m rMesh) {
	renderable := m.RenderableComponent()
	g := rs.g
	key := instanceKey{
		renderable.Material,
		renderable.Mesh,
	}

	if instanced, ok := rs.instancesMap[key]; ok {
		instanced.meshes = append(instanced.meshes, m)
		return
	}
	instanced := &renderableInstance{}

	mat := renderable.Material
	if mat == nil {
		// set some default material
		panic("no material")
	}
	shader, err := rs.shaders.Get(mat.Name)
	if err != nil {
		panic(err)
	}

	instanced.shader = shader
	instanced.Material = mat

	// Preload textures on add
	for _, t := range mat.Textures {
		rs.textures.Get(t) // or by Name?
		/*if err != nil {
			rs.error(err)
			panic(err)
		}*/
		//mat.SetTexture(k, tex)
	}
	instanced.Mesh = renderable.Mesh
	instanced.VAO = g.CreateVertexArray()
	g.BindVertexArray(instanced.VAO)
	vb := rs.vbos.Get(renderable.Mesh)
	instanced.vbo = vb

	vb.bindForShader(shader)

	//gles3
	// Bind a transform location
	instanced.TRO = g.CreateBuffer()
	g.BindBuffer(gl.ARRAY_BUFFER, instanced.TRO)
	// Hack to upload a mat4 into attr (4*4 vec4)
	vec4size := 4 * 4 // 4 floats in bytes size
	if a, ok := shader.Attrib("aTransform"); ok {
		for i := uint32(0); i < 4; i++ {
			aa := a + gl.Attrib(i)
			g.EnableVertexAttribArray(aa)
			g.VertexAttribPointer(aa, 4, gl.FLOAT, false, 4*vec4size+vec4size, int(i)*vec4size)
			g.VertexAttribDivisor(aa, 1)
		}
	}

	if a, ok := shader.Attrib("aColor"); ok {
		g.EnableVertexAttribArray(a)
		g.VertexAttribPointer(a, 4, gl.FLOAT, false, 4*vec4size+vec4size, 4*vec4size)
		g.VertexAttribDivisor(a, 1)
	}

	instanced.meshes = append(instanced.meshes, m)

	rs.instances = append(rs.instances, instanced)
	rs.instancesMap[key] = instanced

}

func (rs *Renderer) watchResize(evt gorge.ResizeEvent) {
	sz := m32.Vec2(evt)
	rs.g.Viewport(0, 0, int(sz[0]), int(sz[1]))
}
func (rs *Renderer) handleEntityAdd(entities gorge.EntitiesAddEvent) {
	for _, e := range entities {
		switch v := e.(type) {
		case rCamera:
			rs.camera = v
		case rLight: // light list
			rs.light = v
		case rMesh:
			rs.AddMesh(v)
		}
	}
}

// HandlePostUpdate will render stuff uppon post update called
func (rs *Renderer) handlePostUpdate(evt gorge.PostUpdateEvent) {
	rs.totalTime += float32(evt)
	rs.Render()
}

// handleAssetAdd when an event from asset manager occurs we upload to gpu right away?
func (rs *Renderer) handleAssetAdd(a asset.AddEvent) {
	switch v := a.Asset.(type) {
	case *gorge.Texture:
		rs.textures.Get(v)
	case *gorge.Mesh:
		rs.vbos.Get(v)
	}
}

///////////////////////////////////////////////////////////////////////////////
// RENDER
///////////////////////////////////////////////////////////////////////////////

// RenderPass struct

// Render Scene
func (rs *Renderer) Render() {
	if rs.camera == nil {
		return
	}

	// Upload instance data
	rs.Prepare()

	// Prepare render info from renderer cameras and state
	// Prepare projection and view matrix
	cam := rs.camera.CameraComponent()
	projection := cam.Projection()
	view := rs.camera.TransformComponent().Inv()

	/*type mat4er interface{ Mat4() m32.Mat4 }
	if m, ok := rs.camera.(mat4er); ok {
		projection = m.Mat4()
	}*/
	camPos := rs.camera.TransformComponent().WorldPosition()
	ambient := cam.Ambient

	// Lights
	var lightPos vec3
	var lightColor = vec3{1, 1, 1}
	if rs.light != nil {
		lightPos = rs.light.TransformComponent().WorldPosition()
		lightColor = rs.light.LightComponent().Color
	}

	ri := renderInfo{
		// The common uniform stuff that we could pass to shader manager
		projection,
		view,
		camPos,
		ambient,
		// More possible lighting with more info
		lightPos,
		lightColor,
	}

	if ExperimentalSkybox {
		rs.Skybox(&ri)
	} else {
		rs.g.ClearColor(ri.ambient[0], ri.ambient[1], ri.ambient[2], 1)
		rs.g.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	}
	rs.Pass(&ri)
	// build transforms and whatnots

}

type renderInfo struct {
	projection mat4
	view       mat4
	camPos     vec3
	ambient    vec3
	lightPos   vec3
	lightColor vec3

	// Target framebuffer
	// Material overrides for specific cases, shadows etc
}

// Prepare will update instance attribs to gpu
func (rs *Renderer) Prepare() {
	g := rs.g
	// Per material
	for _, ins := range rs.instances {

		if len(ins.meshes) == 0 {
			continue
		}
		if ins.attribBuf == nil {
			ins.attribBuf = NewF32TransferBuf(
				len(ins.meshes)*16 + len(ins.meshes)*4,
			)
		}
		g.BindVertexArray(ins.VAO)
		if ins.vbo.update() {
			ins.vbo.bindForShader(ins.shader)
		}
		// Do we need this here at all?

		// Instancing upload all transforms into a float array
		for i, r := range ins.meshes {
			// Do the transformations
			transform := r.TransformComponent()
			r := r.RenderableComponent()
			m := transform.Mat4()

			totSize := 16 + 4
			ins.attribBuf.WriteAt(m[:], i*totSize)
			ins.attribBuf.WriteAt(r.Color[:], i*totSize+16)
		}
		g.BindBuffer(gl.ARRAY_BUFFER, ins.TRO)
		g.BufferDataX(gl.ARRAY_BUFFER, ins.attribBuf.Get(), gl.DYNAMIC_DRAW)
	}
}

// Pass update uniforms per instance and do the actually rendering
func (rs *Renderer) Pass(ri *renderInfo) {
	g := rs.g

	for _, ins := range rs.instances {
		if len(ins.meshes) == 0 {
			continue
		}
		shader := ins.shader
		rs.useMaterial(shader, ins.Material)

		shader.Set("projection", ri.projection)
		shader.Set("view", ri.view)
		shader.Set("viewPos", ri.camPos) // Maybe we don't need this if we take from view
		if rs.light != nil {
			shader.Set("lightPos[0]", ri.lightPos)
			shader.Set("lightColors[0]", ri.lightColor.Mul(1000))
		}
		// Pass sampler here to shader

		// Ambient sampler
		// Maybe we can cache this too

		// Skybox
		if ExperimentalSkybox {
			// Last active texture available
			texI := uint32(len(shader.samplers))

			g.ActiveTexture(gl.Enum(gl.TEXTURE0 + texI))
			g.BindTexture(gl.TEXTURE_CUBE_MAP, rs.skyboxTex.id)
			shader.Set("envMap", texI)
		}

		// We can avoid this
		//shader.Set("time", rs.totalTime)

		g.BindVertexArray(ins.VAO)
		drawType := gl.Enum(ins.Material.DrawType)

		if ins.vbo.ElementsLen > 0 {
			////if ins.ElementsLen > 0 {
			g.DrawElementsInstanced(drawType, ins.vbo.ElementsLen, gl.UNSIGNED_INT, 0, len(ins.meshes))
		} else {
			g.DrawArraysInstanced(drawType, 0, ins.vbo.VertexLen, len(ins.meshes))
		}
	}
}

func (rs *Renderer) error(err error) {
	rs.gorge.Trigger(gorge.ErrorEvent(err))
}

// useMaterial this function sets the shader uniforms from a gorge.Material
// it will check values if they are different it will set
func (rs *Renderer) useMaterial(shader *Shader, mat *gorge.Material) {
	g := rs.g
	// Rebind stuff here?
	// This is from material
	if mat.Depth {
		g.Enable(gl.DEPTH_TEST)
	} else {
		g.Disable(gl.DEPTH_TEST)
	}
	if !mat.DoubleSided {
		g.Enable(gl.CULL_FACE)
	} else {
		g.Disable(gl.CULL_FACE)
	}
	// Get shader
	shader.bind()

	// Update gorge Textures here too
	for i, k := range shader.samplers {
		var tex = rs.textures.gray
		if gtex, ok := mat.Textures[k]; ok {
			tex = rs.textures.Get(gtex)
		}
		g.ActiveTexture(gl.TEXTURE0 + gl.Enum(i))
		g.BindTexture(gl.TEXTURE_2D, tex.id)
		shader.Set(k, i)
	}

	// This way we can prepare stuff?
	props := mat.Props()
	for k, v := range props {
		shader.Set(k, v)
	}
}
