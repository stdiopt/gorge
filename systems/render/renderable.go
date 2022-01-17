package render

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/setlist"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/render/bufutil"
	"github.com/stdiopt/gorge/systems/render/gl"
)

// Store this on renderable
type renderable struct {
	shader *Shader
	vbo    *VBO
}

// RenderableGroup represents instance set of renderables
type RenderableGroup struct {
	// Instances
	Instances    setlist.SetList[Renderable]
	RenderNumber int
	Count        uint32

	renderer   *Render
	renderable *gorge.RenderableComponent

	// VAO per shader attrib hash
	shaderVAO map[uint]gl.VertexArray

	tro *bufutil.Cached[float32]

	// vbo *VBO
	// shader *Shader

	// Cached material
	material *gorge.Material
	hash     uint // both mesh and material defines hash

	troResize bool
}

func (rg *RenderableGroup) init() bool {
	if rg.renderable == nil {
		return false
	}
	vbo, _ := rg.renderer.vbos.Get(rg.renderable.Mesh)
	if vbo == nil {
		return false
	}
	rg.shaderVAO = map[uint]gl.VertexArray{}

	rg.tro = bufutil.NewCached[float32](rg.renderer.buffers.New(gl.ARRAY_BUFFER, gl.DYNAMIC_DRAW))
	rg.tro.Init(4 + 16 + 16) // 1 position

	gorge.SetGPU(rg.renderable, &renderable{
		vbo: vbo,
	})
	return true
}

// Renderable returns the renderable component for this instance.
func (rg *RenderableGroup) Renderable() *gorge.RenderableComponent {
	return rg.renderable
}

// Add adds a new instance to this set.
func (rg *RenderableGroup) Add(r Renderable) {
	rg.Instances.Add(r)
	rg.troResize = true
}

// Remove removes an instance from this set.
func (rg *RenderableGroup) Remove(r Renderable) {
	rg.Instances.Remove(r)
	rg.troResize = true
}

// Destroy unreferences all resources on this group.
func (rg *RenderableGroup) Destroy() {
	gorge.SetGPU(rg.renderable, nil)

	rg.clearVAOS()
	rg.tro.Destroy()
	// Remove stuff from renderable

	rg.renderable = nil
	// rg.vbo = nil
	rg.material = nil
}

// Front returns the first Renderable on this group.
func (rg *RenderableGroup) Front() Renderable {
	return rg.Instances.Front()
}

// Update updates any related gpu buffer on the group based on pass
// we update the group with the pass render number which might be updated
func (rg *RenderableGroup) Update(s *Step) {
	rr, ok := gorge.GetGPU(rg.renderable).(*renderable)
	if !ok {
		panic("something wrong")
		// rr = &renderable{}
		// gorge.SetGPU(rg.renderable, rr)
	}

	vbo, _ := rg.renderer.vbos.Get(rg.renderable.Mesh)
	if vbo != rr.vbo {
		// Need to be aware if the the VBO format changed here
		rr.vbo = vbo
		rg.clearVAOS()
	}
	if vbo == nil || vbo.VertexLen == 0 {
		return
	}

	hash := rg.material.DefinesHash() ^ rg.renderable.Mesh.DefinesHash()
	if rg.material != rg.renderable.Material || hash != rg.hash {
		shdr := rg.renderer.shaders.GetX(rg.renderable)
		// Rebuild VAO since material or mesh changed and we need to update
		// VertexAttribs
		if rr.shader != nil && rr.shader.attribsHash != shdr.attribsHash {
			gl.DeleteVertexArray(rg.shaderVAO[rr.shader.attribsHash])
			delete(rg.shaderVAO, rr.shader.attribsHash)
		}

		// Recache stuff
		rg.material = rg.renderable.Material
		rg.hash = hash
		rr.shader = shdr
	}

	// unitSize in floats
	unitSize := 4 + 16 + 16
	if rg.troResize {
		rg.troResize = false
		sz := rg.Instances.Len() * unitSize
		rg.tro.Init(sz)
	}

	offs := 0
	rg.Count = 0
	for _, r := range rg.Instances.Items() {
		if v, ok := r.(interface{ RenderDisable() bool }); ok {
			// Could check transform Disable going upward parent
			// if some parent is disable, this is disabled
			if v.RenderDisable() {
				continue
			}
		}

		// Do the transformations
		m := r.Mat4()
		um := m.Inv().Transpose() // New: Normal Matrix
		color := m32.Vec4{1, 1, 1, 1}
		if v, ok := r.(interface{ GetColor() m32.Vec4 }); ok {
			color = v.GetColor()
		}
		totSize := unitSize
		rg.tro.WriteAt(color[:], offs)
		rg.tro.WriteAt(m[:], offs+4)
		rg.tro.WriteAt(um[:], offs+4+16)
		offs += totSize
		rg.Count++
	}
	rg.tro.Flush()
	rg.RenderNumber = s.RenderNumber
}

// VBO returns the renderable VBO.
func (rg *RenderableGroup) VBO() *VBO {
	rr, ok := gorge.GetGPU(rg.renderable).(*renderable)
	if !ok {
		return nil
	}
	return rr.vbo
}

// VAO returns an existing VAO for shader hash, if vao doesn't exists it
// will create one and bind attributes
// there might be different VAO's per shader as some shaders might not care
// for normals vertexbuffers etc.
func (rg *RenderableGroup) VAO(shader *Shader) gl.VertexArray {
	if shader == nil {
		if rr, ok := gorge.GetGPU(rg.renderable).(*renderable); ok {
			shader = rr.shader
		}
	}
	if vao, ok := rg.shaderVAO[shader.attribsHash]; ok {
		return vao
	}

	vao := gl.CreateVertexArray()
	rg.bindAttribs(vao, shader)
	rg.shaderVAO[shader.attribsHash] = vao

	return vao
}

// clearVAOS resets arrao objects for all shaders.
func (rg *RenderableGroup) clearVAOS() {
	for k, vao := range rg.shaderVAO {
		gl.DeleteVertexArray(vao)
		delete(rg.shaderVAO, k)
	}
}

// Bring Back named attribs please
func (rg *RenderableGroup) bindAttribs(vao gl.VertexArray, shader *Shader) {
	// log.Println("re bind attribs for VAO:", vao)
	gl.BindVertexArray(vao)
	// Setup TRO
	rg.tro.Bind()
	vec4size := 4 * 4 // 4 floats in bytes size
	totSizeBytes := (4 + 16 + 16) * 4

	if loc, ok := shaderAttrib(shader, aInstanceColor); ok {
		gl.EnableVertexAttribArray(loc)
		gl.VertexAttribPointer(loc, 4, gl.FLOAT, false, totSizeBytes, 0)
		gl.VertexAttribDivisor(loc, 1)
	}
	curOff := 4 * 4
	// Attribs only support vec4 at a time
	// we bind 4 times for the full model matrix
	if loc, ok := shaderAttrib(shader, aTransform); ok {
		for i := uint32(0); i < 4; i++ {
			a := loc + gl.Attrib(i) // nolint: gocritic, unconvert
			gl.EnableVertexAttribArray(a)
			gl.VertexAttribPointer(a, 4, gl.FLOAT, false, totSizeBytes, curOff+int(i)*vec4size)
			gl.VertexAttribDivisor(a, 1)
		}
	}
	curOff = (4 + 16) * 4
	// NormalMatrix
	if loc, ok := shaderAttrib(shader, aNormalTransform); ok {
		for i := uint32(0); i < 4; i++ {
			a := loc + gl.Attrib(i) // nolint: gocritic, unconvert
			gl.EnableVertexAttribArray(a)
			gl.VertexAttribPointer(a, 4, gl.FLOAT, false, totSizeBytes, curOff+int(i)*vec4size)
			gl.VertexAttribDivisor(a, 1)
		}
	}

	rr, ok := gorge.GetGPU(rg.renderable).(*renderable)
	if !ok {
		return
	}
	if rr.vbo == nil || rr.vbo.VertexLen == 0 {
		return
	}

	rr.vbo.BindAttribs(shader)
	gl.BindVertexArray(gl.Null)
}
