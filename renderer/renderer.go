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
	"fmt"
	glog "log"
	"runtime/debug"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gl"
	"github.com/stdiopt/gorge/m32"
)

var (
	log = glog.New(glog.Writer(), "(renderer) ", 0)
	// ExperimentalSkybox some samples uses this static skybox
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
	shader   *shader

	VAO gl.VertexArray

	vbo *vbo
	// for instancing
	TRO gl.Buffer

	// before []rMesh
	// Instances
	meshes gorge.SetList

	attribBuf *F32TransferBuf
}

// Renderer thing
type Renderer struct {
	gorge *gorge.Gorge
	// Camera projection
	camera rCamera
	light  rLight

	instances    []*renderableInstance
	instancesMap map[instanceKey]*renderableInstance

	//lastInstance *renderableInstance

	textures *textureManager
	shaders  *shaderManager
	vbos     *vboManager

	g         gl.Context3
	totalTime float32

	skyboxTex    *texture
	skyboxShader *shader
	skyboxVAO    gl.VertexArray
	skyboxEBO    []uint32
}

// System initializes gl context and attatch the handlers on manager
func System(gm *gorge.Gorge) {
	var g *gl.Wrapper
	gm.Query(func(glctx *gl.Wrapper) { g = glctx })

	if g == nil {
		panic("renderer requires persisted glctx")
	}
	// Left hand
	g.ClearDepthf(-1)
	g.DepthFunc(gl.GEQUAL)
	g.FrontFace(gl.CCW)
	g.Hint(gl.GENERATE_MIPMAP_HINT, gl.NICEST)
	g.LineWidth(3)

	// should be on material
	g.Enable(gl.BLEND)
	g.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	g.CullFace(gl.BACK)

	rs := &Renderer{
		gorge:        gm,
		g:            g,
		instances:    []*renderableInstance{},
		instancesMap: map[instanceKey]*renderableInstance{},
		textures:     newTextureManager(g),
		vbos:         newVBOManager(g),
		shaders:      newShaderManager(g),
	}

	gm.Handle(rs.handleResize)
	gm.Handle(rs.handleEntityAdd)
	gm.Handle(rs.handleEntityRemove)
	gm.Handle(rs.handlePostUpdate)

	// XXX: Two tales
	// current asset.Bundle doesn't require us to pass assets in textureManager
	// gorge.AssetBundle does as we load assets with asset.Manager here

	gm.Handle(func(bundle gorge.LoadBundleEvent) {
		for a := range bundle.Assets {
			switch v := a.(type) {
			case *gorge.Texture:
				rs.textures.Get(v)
			case *gorge.Mesh:
				rs.vbos.Load(v)
			default:
				gm.Warn(fmt.Sprintf("unknown asset: %T", a))
			}
		}
	})

	if ExperimentalSkybox {
		rs.PrepareSkybox()
	}

}

// AddMesh adds a mesh
// Mesh Manager?
func (rs *Renderer) AddMesh(m rMesh) {
	g := rs.g
	renderable := m.RenderableComponent()
	key := instanceKey{
		renderable.Material,
		renderable.Mesh,
	}

	if instanced, ok := rs.instancesMap[key]; ok {
		// We are adding stuff so we reset the transfer buf
		// (Grow method would be better)
		if instanced.attribBuf != nil {
			instanced.attribBuf = nil
		}
		instanced.meshes.Add(m)
		return
	}

	mat := renderable.Material
	if mat == nil {
		// TODO: set some default material
		panic("no material")
	}
	shader := rs.shaders.Get(mat)
	vb := rs.vbos.Get(renderable.Mesh)
	if vb == nil {
		debug.PrintStack()
		rs.gorge.Error(fmt.Errorf("asset not found in renderer: %v", renderable.Mesh))
		return
	}
	// Preload textures on add
	for _, t := range mat.Textures {
		rs.textures.Get(t)
	}

	instanced := &renderableInstance{}
	instanced.shader = shader
	instanced.Material = mat

	instanced.Mesh = renderable.Mesh
	instanced.VAO = g.CreateVertexArray()
	g.BindVertexArray(instanced.VAO)
	instanced.vbo = vb

	// Bind vertexAttribs on VAO
	vb.bindForShader(shader)

	//gles3
	// Bind a transform location
	instanced.TRO = g.CreateBuffer()
	g.BindBuffer(gl.ARRAY_BUFFER, instanced.TRO)
	// Hack to upload a mat4 into attr (4*4 vec4)
	vec4size := 4 * 4 // 4 floats in bytes size
	if a, ok := shader.Attrib("aTransform"); ok {
		// Attribs only support vec4 at a time
		// we bind 4 times for the full model matrix
		for i := uint32(0); i < 4; i++ {
			aa := a + gl.Attrib(i)
			g.EnableVertexAttribArray(aa)
			g.VertexAttribPointer(aa, 4, gl.FLOAT, false, 4*vec4size+vec4size, int(i)*vec4size)
			g.VertexAttribDivisor(aa, 1)
		}
	}

	// Bind the main color next
	if a, ok := shader.Attrib("aColor"); ok {
		g.EnableVertexAttribArray(a)
		g.VertexAttribPointer(a, 4, gl.FLOAT, false, 4*vec4size+vec4size, 4*vec4size)
		g.VertexAttribDivisor(a, 1)
	}

	instanced.meshes.Add(m)

	rs.instances = append(rs.instances, instanced)
	rs.instancesMap[key] = instanced
}

// RemoveMesh remove to renderable renderable
func (rs *Renderer) RemoveMesh(m rMesh) {
	// Go trough instances and remove from there
	renderable := m.RenderableComponent()
	key := instanceKey{
		renderable.Material,
		renderable.Mesh,
	}

	instance, ok := rs.instancesMap[key]
	if !ok {
		return
	}

	instance.meshes.Remove(m)

	// Remove instance from system
	if instance.meshes.Len() == 0 {
		delete(rs.instancesMap, key)
		for i, is := range rs.instances {
			if is != instance {
				continue
			}
			// Maintain order, slower
			rs.instances = append(rs.instances[:i], rs.instances[i+1:]...)
		}
	}

}
func (rs *Renderer) handleResize(evt gorge.ResizeEvent) {
	sz := m32.Vec2(evt)
	rs.g.Viewport(0, 0, int(sz[0]), int(sz[1]))
}
func (rs *Renderer) handleEntityAdd(ents gorge.EntitiesAddEvent) {

	for _, e := range ents {
		//for _, e := range entities {
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

// handleEntityRemove
func (rs *Renderer) handleEntityRemove(ents gorge.EntitiesRemoveEvent) {
	for _, e := range ents {
		switch v := e.(type) {
		case rCamera:
			rs.camera = v
		case rLight: // light list
			rs.light = v
		case rMesh:
			rs.RemoveMesh(v)
		}
	}
}

// HandlePostUpdate will render stuff uppon post update called
func (rs *Renderer) handlePostUpdate(e gorge.RenderEvent) {
	rs.totalTime += float32(e)
	rs.Render()
}

///////////////////////////////////////////////////////////////////////////////
// RENDER
///////////////////////////////////////////////////////////////////////////////

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
	// Waits for scene to finish renderer
	// Good to check gpu performance I guess
	rs.g.Finish()

}

// Prepare will update instance attribs to gpu
func (rs *Renderer) Prepare() {
	g := rs.g
	// Per material
	for _, ins := range rs.instances {
		mlen := ins.meshes.Len()
		if mlen == 0 {
			continue
		}
		if ins.attribBuf == nil {
			// 16- Transform attrib + 4 - Color attrib
			sz := mlen * (16 + 4)
			ins.attribBuf = NewF32TransferBuf(sz)
		}
		g.BindVertexArray(ins.VAO)
		if ins.vbo.update(false) {
			ins.vbo.bindForShader(ins.shader)
		}
		// Do we need this here at all?

		// Instancing upload all transforms into a float array

		offs := 0
		ins.meshes.Range(func(v interface{}) bool {
			mesh := v.(rMesh)
			// Do the transformations
			transform := mesh.TransformComponent()
			r := mesh.RenderableComponent()
			m := transform.Mat4()

			totSize := 16 + 4
			ins.attribBuf.WriteAt(m[:], offs)
			ins.attribBuf.WriteAt(r.Color[:], offs+16)
			offs += totSize
			return true
		})
		g.BindBuffer(gl.ARRAY_BUFFER, ins.TRO)
		g.BufferDataX(gl.ARRAY_BUFFER, ins.attribBuf.Get(), gl.DYNAMIC_DRAW)
	}
}

// Pass update uniforms per instance and do the actually rendering
func (rs *Renderer) Pass(ri *renderInfo) {
	g := rs.g

	for _, ins := range rs.instances {
		mlen := ins.meshes.Len()
		if mlen == 0 {
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
		drawType := drawType2gl(ins.Mesh.DrawType)

		if ins.vbo.ElementsLen > 0 {
			////if ins.ElementsLen > 0 {
			g.DrawElementsInstanced(drawType, ins.vbo.ElementsLen, gl.UNSIGNED_INT, 0, mlen)
		} else {
			g.DrawArraysInstanced(drawType, 0, ins.vbo.VertexLen, mlen)
		}
	}
}

func (rs *Renderer) error(err error) {
	rs.gorge.Trigger(gorge.ErrorEvent{Err: err})
}

// useMaterial this function sets the shader uniforms from a gorge.Material
// it will check values if they are different it will set
func (rs *Renderer) useMaterial(shader *shader, mat *gorge.Material) {
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

	// TODO: Bug we should run through shader props and fetch value from
	// material
	props := mat.Props()
	for k, v := range props {
		shader.Set(k, v)
	}
}

//POINTS                                       = 0x0000
//LINES                                        = 0x0001
//LINE_LOOP                                    = 0x0002
//LINE_STRIP                                   = 0x0003
//TRIANGLES                                    = 0x0004
//TRIANGLE_STRIP                               = 0x0005
//TRIANGLE_FAN                                 = 0x0006
func drawType2gl(d gorge.DrawType) gl.Enum {
	switch d {
	case gorge.DrawPoints:
		return gl.POINTS
	case gorge.DrawLines:
		return gl.LINES
	case gorge.DrawLineLoop:
		return gl.LINE_LOOP
	case gorge.DrawLineStrip:
		return gl.LINE_STRIP
	case gorge.DrawTriangles:
		return gl.TRIANGLES
	case gorge.DrawTriangleStrip:
		return gl.TRIANGLE_STRIP
	case gorge.DrawTriangleFan:
		return gl.TRIANGLE_FAN
	}
	panic("unknown drawtype")
}
