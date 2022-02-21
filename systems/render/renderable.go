package render

import (
	"log"
	"runtime"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/setlist"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/render/bufutil"
	"github.com/stdiopt/gorge/systems/render/gl"
)

var vaoCount int

// This will be attached on renderableComponent to track VAO's and related
// built shaders
type renderable struct {
	// VAO per shader attrib hash
	shaderVAO map[uint]gl.VertexArray
	shader    *Shader
	vbo       *VBO
	// multiple instances transform
	tro       *bufutil.Cached[float32]
	troResize bool

	// cached stuff
	renderNumber int
	material     *gorge.Material

	// Cached material
	hash uint // both mesh and material defines hash

}

func (r *renderable) clearVAOS() {
	for k, vao := range r.shaderVAO {
		gl.DeleteVertexArray(vao)
		vaoCount--
		delete(r.shaderVAO, k)
	}
}

func (r *renderable) destroy() {
	r.clearVAOS()
	r.tro.Destroy()
	r.vbo = nil
}

// RenderableGroup represents instance set of renderables
type RenderableGroup struct {
	renderer *Render
	// Instances
	Instances setlist.SetList[Renderable]
	// Count is the number of ACTIVE renderable in this group
	Count uint32

	renderable *gorge.RenderableComponent
}

func (rg *RenderableGroup) init() bool {
	if rg.renderable == nil {
		return false
	}
	vbo, _ := rg.renderer.vbos.Get(rg.renderable.Mesh)
	if vbo == nil {
		return false
	}
	// rg.shaderVAO = map[uint]gl.VertexArray{}

	tro := bufutil.NewCached[float32](rg.renderer.buffers.New(gl.ARRAY_BUFFER, gl.DYNAMIC_DRAW))
	tro.Init(4 + 16 + 16) // 1 position

	rr := &renderable{
		shaderVAO: map[uint]gl.VertexArray{},
		vbo:       vbo,
		tro:       tro,
	}

	runtime.SetFinalizer(rr, func(r *renderable) {
		rg.renderer.gorge.RunInMain(func() {
			r.destroy()
		})
	})
	gorge.SetGPU(rg.renderable, rr)
	return true
}

// Renderable returns the renderable component for this instance.
func (rg *RenderableGroup) Renderable() *gorge.RenderableComponent {
	return rg.renderable
}

// Add adds a new instance to this set.
func (rg *RenderableGroup) Add(r Renderable) {
	rg.Instances.Add(r)
	if rr, ok := gorge.GetGPU(rg.renderable).(*renderable); ok {
		rr.troResize = true
	}
}

// Remove removes an instance from this set.
func (rg *RenderableGroup) Remove(r Renderable) {
	rg.Instances.Remove(r)
	if rr, ok := gorge.GetGPU(rg.renderable).(*renderable); ok {
		rr.troResize = true
	}
}

// Destroy unreferences all resources on this group.
func (rg *RenderableGroup) Destroy() {
	gorge.SetGPU(rg.renderable, nil)
	// Remove stuff from renderable

	rg.renderable = nil
	// rg.vbo = nil
	// rg.material = nil
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
		rr = &renderable{}
		gorge.SetGPU(rg.renderable, rr)
	}
	// No need to update for this render as it was already updated
	if s.RenderNumber == rr.renderNumber {
		return
	}

	vbo, _ := rg.renderer.vbos.Get(rg.renderable.Mesh)
	if vbo != rr.vbo {
		// Need to be aware if the the VBO format changed here
		rr.vbo = vbo
		rr.clearVAOS()
	}
	if vbo == nil || vbo.VertexLen == 0 {
		return
	}

	hash := rr.material.DefinesHash() ^ rg.renderable.Mesh.DefinesHash()
	if rr.material != rg.renderable.Material || hash != rr.hash {
		shdr := rg.renderer.shaders.GetX(rg.renderable)
		// Rebuild VAO since material or mesh changed and we need to update
		// VertexAttribs
		if rr.shader != nil && rr.shader.attribsHash != shdr.attribsHash {
			gl.DeleteVertexArray(rr.shaderVAO[rr.shader.attribsHash])
			delete(rr.shaderVAO, rr.shader.attribsHash)
		}

		// Recache stuff
		rr.hash = hash
		rr.material = rg.renderable.Material
		rr.shader = shdr
	}

	// unitSize in floats
	unitSize := 4 + 16 + 16
	if rr.troResize {
		rr.troResize = false
		sz := rg.Instances.Len() * unitSize
		rr.tro.Init(sz)
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
		color := gm.Vec4{1, 1, 1, 1}
		if v, ok := r.(interface{ GetColor() gm.Vec4 }); ok {
			color = v.GetColor()
		}
		totSize := unitSize
		rr.tro.WriteAt(color[:], offs)
		rr.tro.WriteAt(m[:], offs+4)
		rr.tro.WriteAt(um[:], offs+4+16)
		offs += totSize
		rg.Count++
	}
	rr.tro.Flush()
	rr.renderNumber = s.RenderNumber
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
	rr, ok := gorge.GetGPU(rg.renderable).(*renderable)
	if !ok {
		log.Println("[WARN] Creating empty renderable")
		rr = &renderable{}
		gorge.SetGPU(rg.renderable, rr)
	}
	if shader == nil {
		shader = rr.shader
	}
	if vao, ok := rr.shaderVAO[shader.attribsHash]; ok {
		return vao
	}

	vaoCount++
	vao := gl.CreateVertexArray()
	rg.bindAttribs(vao, shader)
	rr.shaderVAO[shader.attribsHash] = vao

	return vao
}

// Bring Back named attribs please
func (rg *RenderableGroup) bindAttribs(vao gl.VertexArray, shader *Shader) {
	rr, ok := gorge.GetGPU(rg.renderable).(*renderable)
	if !ok {
		return
	}
	// log.Println("re bind attribs for VAO:", vao)
	gl.BindVertexArray(vao)
	// Setup TRO
	rr.tro.Bind()
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

	if rr.vbo == nil || rr.vbo.VertexLen == 0 {
		return
	}

	rr.vbo.BindAttribs(shader)
	gl.BindVertexArray(gl.Null)
}
