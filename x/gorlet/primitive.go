package gorlet

import (
	"fmt"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/static"
)

type gEntity struct {
	name string
	gorge.TransformComponent
	gorge.ColorableComponent
	*gorge.RenderableComponent
}

func (e gEntity) String() string {
	return fmt.Sprintf("gEntity(%s)", e.name)
}

func newEntity(name string, mesh gorge.Mesher) *gEntity {
	mat := gorge.NewShaderMaterial(static.Shaders.UI)
	mat.Queue = 3000
	mat.Depth = gorge.DepthRead
	mat.DisableShadow = true
	return &gEntity{
		name:                name,
		TransformComponent:  *gorge.NewTransformComponent(),
		ColorableComponent:  *gorge.NewColorableComponent(1, 1, 1, 1),
		RenderableComponent: gorge.NewRenderableComponent(mesh, mat),
	}
}

func quadMeshData() *gorge.MeshData {
	return &gorge.MeshData{
		Format: gorge.VertexFormatPTN(),
		Vertices: []float32{
			/*P:*/ 0, 1, 0 /*T*/, 0, 0 /*N*/, 0, 0, 1,
			/*P:*/ 1, 1, 0 /*T*/, 1, 0 /*N*/, 0, 0, 1,
			/*P:*/ 1, 0, 0 /*T*/, 1, 1 /*N*/, 0, 0, 1,
			/*P:*/ 0, 0, 0 /*T*/, 0, 1 /*N*/, 0, 0, 1,
		},
		Indices: []uint32{
			0, 2, 1,
			2, 0, 3,
		},
	}
}

func quadMesh() *gorge.Mesh {
	return gorge.NewMesh(quadMeshData())
}

type mesh = gorge.Mesh

func newRoundedQuadMesh(sz gm.Vec2, r float32) *roundedQuadMesh {
	m := &roundedQuadMesh{
		size:   sz,
		radius: r,
		path:   *NewPath(),
	}
	m.update()
	m.mesh = *gorge.NewMesh(&m.path)
	m.DrawMode = gorge.DrawTriangleFan
	return m
}

type roundedQuadMesh struct {
	mesh
	radius float32
	size   gm.Vec2
	path   Path
}

func (m *roundedQuadMesh) SetRadius(r float32) {
	m.radius = r
}

func (m *roundedQuadMesh) SetSize(sz gm.Vec2) {
	m.size = sz
}

func (m *roundedQuadMesh) update() {
	sz := m.size
	m.path.Reset()
	m.path.size = sz

	cx := gm.Clamp(m.radius, 0, sz[0]/2)
	cy := gm.Clamp(m.radius, 0, sz[1]/2)
	cxi := sz[0] - cx
	cyi := sz[1] - cy
	m.path.LineTo(cxi, 0)                                  // BottomLeft to BottomRight
	m.path.CurveTo(gm.V2(sz[0], 0), gm.V2(sz[0], cy))      // BottomRight corner
	m.path.LineTo(sz[0], cyi)                              // BottomRight to TopRight
	m.path.CurveTo(gm.V2(sz[0], sz[1]), gm.V2(cxi, sz[1])) // TopRight corner
	m.path.LineTo(cx, sz[1])                               // TopRight to TopLeft
	m.path.CurveTo(gm.V2(0, sz[1]), gm.V2(0, cyi))         // TopLeft corner
	m.path.LineTo(0, cy)                                   // TopLeft to BottomLeft
	m.path.CurveTo(gm.V2(0, 0), gm.V2(cx, 0))              // BottomLeft corner
	m.path.Updates++
}

type Path struct {
	cur  gm.Vec2
	size gm.Vec2
	gorge.MeshData
}

func NewPath() *Path {
	return &Path{
		MeshData: gorge.MeshData{
			FrontFacing: gorge.FrontFacingCW,
			Format:      gorge.VertexFormatPTN(),
		},
	}
}

func (p *Path) Reset() {
	if p.Vertices != nil {
		p.Vertices = p.Vertices[:0]
	}
}

func (p *Path) MoveTo(x, y float32) {
	p.cur = gm.Vec2{x, y}
}

func (p *Path) LineTo(x, y float32) {
	p.Vertices = append(p.Vertices /*P*/, x, y, 0 /*T*/, x/p.size[0], 1-y/p.size[1] /*N*/, 0, 0, -1)
	p.cur = gm.Vec2{x, y}
}

func (p *Path) CurveTo(b, c gm.Vec2) {
	a := p.cur
	n := 16
	t, nInv := float32(0), 1/float32(n)
	for i := 0; i < n; i++ {
		t += nInv
		ab := a.Lerp(b, t)
		bc := b.Lerp(c, t)

		abc := ab.Lerp(bc, t)
		p.LineTo(abc[0], abc[1])
	}
	p.cur = c
}
