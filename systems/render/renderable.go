package render

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/internal/setlist"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/render/bufutil"
	"github.com/stdiopt/gorge/systems/render/gl"
)

// RenderableGroup represents instance set of renderables
type RenderableGroup struct {
	// Instances
	Instances    setlist.SetList
	RenderNumber int

	renderer   *Render
	renderable *gorge.RenderableComponent

	// VAO per shader attrib hash
	shaderVAO map[uint]gl.VertexArray

	tro *bufutil.Cachedf32

	vbo    *VBO
	shader *Shader

	// Cached material
	material *gorge.Material
	hash     uint // both mesh and material defines hash

	// vboUpdate    bool
	troResize bool
}

func (rg *RenderableGroup) init() bool {
	vbo, _ := rg.renderer.vbos.Get(rg.renderable.Mesh)
	if vbo == nil {
		return false
	}
	rg.shaderVAO = map[uint]gl.VertexArray{}

	rg.tro = bufutil.NewCachedf32(rg.renderer.buffers.New(gl.ARRAY_BUFFER, gl.DYNAMIC_DRAW))
	rg.tro.Init(4 + 16 + 16) // 1 position

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
	rg.clearVAOS()
	rg.tro.Destroy()
	rg.renderable = nil
	rg.vbo = nil
	rg.material = nil
	rg.shader = nil
}

// Front returns the first Renderable on this group.
func (rg *RenderableGroup) Front() Renderable {
	return rg.Instances.Front().Value.(Renderable)
}

// Update updates any related gpu buffer on the group based on pass
// we update the group with the pass render number which might be updated
func (rg *RenderableGroup) Update(p *Pass) {
	vbo, _ := rg.renderer.vbos.Get(rg.renderable.Mesh)
	if vbo != rg.vbo {
		// Need to be aware if the the VBO format changed here
		rg.vbo = vbo
		rg.clearVAOS()
	}
	if vbo == nil || vbo.VertexLen == 0 {
		return
	}

	hash := rg.material.DefinesHash() ^ rg.renderable.Mesh.DefinesHash()
	if rg.material != rg.renderable.Material || hash != rg.hash {
		s := rg.renderer.shaders.GetX(rg.renderable)
		// Rebuild VAO since material or mesh changed and we need to update
		// VertexAttribs
		if rg.shader != nil && rg.shader.attribsHash != s.attribsHash {
			gl.DeleteVertexArray(rg.shaderVAO[rg.shader.attribsHash])
			delete(rg.shaderVAO, rg.shader.attribsHash)
		}

		rg.material = rg.renderable.Material
		rg.shader = s
		rg.hash = hash
	}

	// unitSize in floats
	unitSize := 4 + 16 + 16
	if rg.troResize {
		rg.troResize = false
		sz := rg.Instances.Len() * unitSize
		rg.tro.Init(sz)
	}

	offs := 0
	for e := rg.Instances.Front(); e != nil; e = e.Next() {
		mesh := e.Value.(Renderable)
		// Do the transformations
		m := mesh.Transform().Mat4()
		um := m.Inv().Transpose() // New: Normal Matrix
		color := m32.Vec4{1, 1, 1, 1}
		if v, ok := mesh.(interface{ GetColor() m32.Vec4 }); ok {
			color = v.GetColor()
		}
		totSize := unitSize
		rg.tro.WriteAt(color[:], offs)
		rg.tro.WriteAt(m[:], offs+4)
		rg.tro.WriteAt(um[:], offs+4+16)
		offs += totSize
	}
	rg.tro.Flush()
	rg.RenderNumber = p.RenderNumber
}

// VBO returns the renderable VBO.
func (rg *RenderableGroup) VBO() *VBO {
	return rg.vbo
}

// VAO returns an existing VAO for shader hash, if vao doesn't exists it
// will create one and bind attributes
// there might be different VAO's per shader as some shaders might not care
// for normals vertexbuffers etc.
func (rg *RenderableGroup) VAO(shader *Shader) gl.VertexArray {
	if shader == nil {
		shader = rg.shader
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

	if rg.vbo == nil || rg.vbo.VertexLen == 0 {
		return
	}

	rg.vbo.BindAttribs(shader)
	gl.BindVertexArray(gl.Null)
}
