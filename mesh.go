package gorge

import (
	"fmt"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
)

// MeshResourcer is an interface for mesh resource providers.
type MeshResourcer interface {
	Resource() MeshResource
}

// MeshResource is an interface to handle underlying mesh data.
type MeshResource interface {
	isMesh()
	isGPU()
}

// Mesh representation
type Mesh struct {
	Resourcer MeshResourcer
	DrawMode  DrawMode

	// This is for shaders like material
	shaderProps
}

// Resource implements MeshResourcer.
func (m *Mesh) Resource() MeshResource {
	if m.Resourcer == nil {
		return nil
	}
	return m.Resourcer.Resource()
}

// NewMesh creates a new mesh with meshData
func NewMesh(res MeshResourcer) *Mesh {
	return &Mesh{Resourcer: res}
}

// Mesh implements mesher interface.
func (m *Mesh) Mesh() *Mesh { return m }

// Clone will clone the mesh and it's props.
func (m *Mesh) Clone() *Mesh {
	return &Mesh{
		Resourcer:   m,
		DrawMode:    m.DrawMode,
		shaderProps: m.copy(),
	}
}

// ReleaseData change the data ref to a gpu only resource.
func (m *Mesh) ReleaseData(g *Context) {
	if _, ok := m.Resource().(*MeshData); !ok {
		return
	}
	curRes := m.Resource()
	if _, ok := curRes.(*MeshData); !ok {
		return
	}
	event.Trigger(g, EventResourceUpdate{Resource: curRes})

	/*{ // free data arrays test
		r := m.Resource.(*MeshData)
		r.Vertices = nil
		r.Indices = nil
	}*/

	gpuRef := &MeshRef{&GPU{}}
	SetGPU(gpuRef, GetGPU(curRes))
	m.Resourcer = gpuRef
}

func (m Mesh) String() string {
	return fmt.Sprintf("(mesh: drawType: %v, loader: %v)",
		m.DrawMode,
		m.Resourcer,
	)
}

// GetDrawMode returns the mesh drawmode.
func (m *Mesh) GetDrawMode() DrawMode {
	return m.DrawMode
}

// Mesh draw type
const (
	DrawTriangles = DrawMode(iota)
	DrawTriangleStrip
	DrawTriangleFan
	DrawPoints
	DrawLines
	DrawLineLoop
	DrawLineStrip
)

// FrontFacing indicates the frontfacing property for the meshData vertices
const (
	FrontFacingCW  = FrontFacing(0)
	FrontFacingCCW = FrontFacing(1)
)

// VertexFormatAttrib vertex format entry for interleaving vertex data in meshData.
type VertexFormatAttrib struct {
	Size   int
	Attrib string
	Define string
}

// VertexFormat type for describing vertex formats.
type VertexFormat []VertexFormatAttrib

// Size returns the data size for this vertex
func (f VertexFormat) Size() int {
	r := 0
	for _, v := range f {
		r += v.Size
	}
	return r
}

// VertexAttrib return a vertex attribute definition
func VertexAttrib(sz int, attrib string, define string) VertexFormatAttrib {
	return VertexFormatAttrib{
		Size:   sz,
		Attrib: attrib,
		Define: define,
	}
}

// VertexFormatP default vertex with positioning only
func VertexFormatP() VertexFormat {
	return VertexFormat{
		{3, "a_Position", "HAS_POSITION"},
	}
}

// VertexFormatPN format for Position and Normal
func VertexFormatPN() VertexFormat {
	return VertexFormat{
		{3, "a_Position", "HAS_POSITION"},
		{3, "a_Normal", "HAS_NORMALS"},
	}
}

// VertexFormatPT format for Position and TexCoord
func VertexFormatPT() VertexFormat {
	return VertexFormat{
		{3, "a_Position", "HAS_POSITION"},
		{2, "a_UV1", "HAS_UV_SET1"},
	}
}

// VertexFormatPTN format for Position Texture and Normal
func VertexFormatPTN() VertexFormat {
	return VertexFormat{
		{3, "a_Position", "HAS_POSITION"},
		{2, "a_UV1", "HAS_UV_SET1"},
		{3, "a_Normal", "HAS_NORMALS"},
	}
}

// VertexFormatPNT format for Position Normal and Texture
func VertexFormatPNT() VertexFormat {
	return VertexFormat{
		{3, "a_Position", "HAS_POSITION"},
		{3, "a_Normal", "HAS_NORMALS"},
		{2, "a_UV1", "HAS_UV_SET1"},
	}
}

// MeshData raw mesh data
type MeshData struct {
	GPU

	Source string

	FrontFacing FrontFacing
	// Describe format and indexes
	Format VertexFormat

	// TODO: This might need to be pure data instead of float32
	// Indices could be a byte we just need to tell gl to read as a byte
	// so we would have a field Indices "type"
	Vertices []float32
	// Indices can be one of []byte, []uint16, []uint32
	Indices any
	Updates int
}

// Resource implements the resourcer interface so MeshData can be used directly
// in the Mesh.
func (d *MeshData) Resource() MeshResource { return d }

func (d *MeshData) isMesh() {}

// CalcBounds calculate the bounding box for this mesh (slow)
func (d *MeshData) CalcBounds() (gm.Vec3, gm.Vec3) {
	sz := d.Format.Size()
	offs := 0
	// Find renderer hardcoded aPosition attrib which is 2
	for _, f := range d.Format {
		if f.Attrib == "a_Position" {
			break
		}
		offs += f.Size
	}

	var min gm.Vec3
	var max gm.Vec3
	v := d.Vertices[offs:]
	copy(min[:], v)
	copy(max[:], v)
	for v := v[sz:]; sz < len(v); v = v[sz:] {
		min[0] = gm.Min(v[0], min[0])
		max[0] = gm.Max(v[0], max[0])
		min[1] = gm.Min(v[1], min[1])
		max[1] = gm.Max(v[1], max[1])
		min[2] = gm.Min(v[2], min[2])
		max[2] = gm.Max(v[2], max[2])

		if sz > len(v) {
			break
		}
	}
	return min, max
}

// ScaleUV manipulate meshData directly
func (d *MeshData) ScaleUV(s ...float32) {
	var scale gm.Vec2
	switch len(s) {
	case 0:
		return
	case 1:
		scale = gm.Vec2{s[0], s[0]}
	default:
		scale = gm.Vec2{s[0], s[1]}
	}

	sz := d.Format.Size()
	offs := 0
	// Find renderer hardcoded TexCoord attrib which is 2
	for _, f := range d.Format {
		if f.Attrib == "a_UV1" {
			break
		}
		offs += f.Size
	}

	for v := d.Vertices[offs:]; ; v = v[sz:] {
		v[0] *= 1 / scale[0]
		v[1] *= 1 / scale[1]
		if sz > len(v) {
			break
		}
	}
}

func (d *MeshData) String() string {
	var ind string
	switch v := d.Indices.(type) {
	case []byte:
		ind = fmt.Sprint("byte:", len(v))
	case []uint16:
		ind = fmt.Sprint("u16:", len(v))
	case []uint32:
		ind = fmt.Sprint("u32:", len(v))
	default:
		ind = "<unknown>"
	}
	return fmt.Sprintf("MeshData: %s, %v verts: %v, ind: %v, upd: %v",
		d.Source,
		d.Format,
		len(d.Vertices), ind, d.Updates,
	)
}

// ////////////////////////////////////////////////////////////////////////////

// Helper mesh struct

// DrawMode type of draw for the renderer
type DrawMode int

func (m DrawMode) String() string {
	switch m {
	case DrawTriangles:
		return "DrawTriangles"
	case DrawTriangleStrip:
		return "DrawTriangleStrip"
	case DrawTriangleFan:
		return "DrawTriangleFan"
	case DrawPoints:
		return "DrawPoints"
	case DrawLines:
		return "DrawLines"
	case DrawLineLoop:
		return "DrawLineLoop"
	case DrawLineStrip:
		return "DrawLineStrip"
	default:
		return fmt.Sprintf("DrawModeUnknown(%d)", m)
	}
}

// FrontFacing type to setup rendering cull.
type FrontFacing int

func (f FrontFacing) String() string {
	switch f {
	case FrontFacingCCW:
		return "FrontFacingCCW"
	case FrontFacingCW:
		return "FrontFacingCW"
	default:
		return fmt.Sprintf("FrontFacingUnknown(%d)", f)
	}
}
