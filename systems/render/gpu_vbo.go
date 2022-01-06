package render

import (
	"runtime"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/render/gl"
)

// Maybe call it meshManager?
// since we will have buffer
type vboManager struct {
	gorge         *gorge.Context
	bufferManager *bufferManager
	count         int
}

func newVBOManager(
	g *gorge.Context,
	bm *bufferManager,
) *vboManager {
	m := &vboManager{
		gorge:         g,
		bufferManager: bm,
	}
	return m
}

func (m *vboManager) New(r gorge.MeshResource) *VBO {
	v := &VBO{
		manager: m,
		updates: -1,
		vbo:     m.bufferManager.New(gl.ARRAY_BUFFER, gl.STATIC_DRAW),
		ebo:     m.bufferManager.New(gl.ELEMENT_ARRAY_BUFFER, gl.STATIC_DRAW),
	}
	m.count++

	runtime.SetFinalizer(v, func(v *VBO) {
		m.gorge.RunInMain(func() {
			m.destroy(v)
		})
	})

	if d, ok := r.(*gorge.MeshData); ok {
		v.upload(d, true)
	}

	return v
}

func (m *vboManager) GetByRef(r gorge.MeshResource) (*VBO, bool) {
	v, ok := gorge.GetGPU(r).(*VBO)
	if !ok {
		v = m.New(r)
		gorge.SetGPU(r, v)
		return v, true
	}
	updated := false
	if d, ok := r.(*gorge.MeshData); ok {
		updated = v.update(d)
	}
	return v, updated
}

func (m *vboManager) Get(mesh *gorge.Mesh) (*VBO, bool) {
	if mesh == nil {
		return nil, false
	}
	return m.GetByRef(mesh.Resource())
}

func (m *vboManager) Update(r *gorge.MeshData) {
	v, ok := gorge.GetGPU(r).(*VBO)
	if !ok {
		v = m.New(r)
		gorge.SetGPU(r, v)
	}
	// Force an update
	v.updates--
	v.update(r)
}

func (m *vboManager) destroy(v *VBO) {
	v.destroy()
}

// VBO vertexBufferObject.
type VBO struct {
	manager *vboManager
	vbo     *buffer
	ebo     *buffer

	Format      gorge.VertexFormat
	FrontFacing gorge.FrontFacing

	ElementsLen  uint32
	ElementsType gl.Enum
	VertexLen    uint32

	// updates indicates updates for dynamic MeshData
	updates      int
	shouldUpdate bool
}

func (v *VBO) destroy() {
	v.vbo.Destroy()
	v.ebo.Destroy()
	v.vbo = nil
	v.ebo = nil
	v.manager.count--
}

func (v *VBO) upload(data *gorge.MeshData, dynamic bool) {
	bufUsage := gl.Enum(gl.STATIC_DRAW)
	if dynamic {
		bufUsage = gl.DYNAMIC_DRAW
	}
	v.Format = data.Format
	v.FrontFacing = data.FrontFacing

	v.VertexLen = uint32(len(data.Vertices) / data.Format.Size())

	v.vbo.Init(len(data.Vertices) * 4)
	v.vbo.usage = bufUsage
	v.vbo.WriteAt(data.Vertices, 0)
	v.vbo.Flush()

	bsz := 1
	v.ElementsLen = 0
	switch ind := data.Indices.(type) {
	case []byte:
		bsz = 1
		v.ElementsType = gl.UNSIGNED_BYTE
		v.ElementsLen = uint32(len(ind))
	case []uint16:
		bsz = 2
		v.ElementsType = gl.UNSIGNED_SHORT
		v.ElementsLen = uint32(len(ind))
	case []uint32:
		bsz = 4
		v.ElementsType = gl.UNSIGNED_INT
		v.ElementsLen = uint32(len(ind))
	}
	if v.ElementsLen > 0 {
		v.ebo.usage = bufUsage
		v.ebo.Init(int(v.ElementsLen) * bsz) // *4 depends on type
		v.ebo.WriteAt(data.Indices, 0)
		v.ebo.Flush()
	}
	v.updates = data.Updates
}

// Load data if
// 1. State is Loading
// 2. data IsMeshData && updates arent equal

func (v *VBO) update(data *gorge.MeshData) bool {
	if v.shouldUpdate {
		v.shouldUpdate = false
		return true
	}
	if v.updates == data.Updates {
		return false
	}
	v.upload(data, true)
	return true
}

// Bind binds both vertex and element buffer
func (v *VBO) Bind() {
	v.ebo.Bind()
	v.vbo.Bind()
}

// BindAttribs Back named attribs please
func (v *VBO) BindAttribs(s *Shader) {
	if v == nil || v.VertexLen == 0 {
		return
	}
	v.vbo.Bind()
	fsz := v.Format.Size()
	offs, next := 0, 0
	for _, f := range v.Format {
		offs = next
		next += f.Size * 4
		loc, ok := shaderAttrib(s, f.Attrib)
		if !ok {
			// gl.DisableVertexAttribArray(loc)
			continue
		}
		gl.EnableVertexAttribArray(loc)
		gl.VertexAttribPointer(loc, f.Size, gl.FLOAT, false, fsz*4, offs)
	}
	v.ebo.Bind()
}
