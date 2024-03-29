package debug

import (
	"math"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/static"
)

// Entity generic renderable entity.
type lineEntity struct {
	gorge.TransformComponent
	gorge.RenderableComponent
}

// Lines helper to draw per frame debug information.
type Lines struct {
	*gorge.TransformComponent
	*gorge.Material
	color gm.Vec4

	Lines  *gorge.MeshData
	Points *gorge.MeshData

	gorge.Container
}

// NewDebugLines returns a Geom entity.
func NewDebugLines() *Lines {
	mat := gorge.NewShaderMaterial(static.Shaders.DefaultNew)
	mat.Define("MATERIAL_UNLIT")

	posColorFormat := gorge.VertexFormat{
		gorge.VertexAttrib(3, "a_Position", "HAS_POSITION"),
		gorge.VertexAttrib(4, "a_Color", "HAS_VERTEX_COLOR_VEC4"),
	}

	rootTrans := gorge.NewTransformComponent()

	lines := &gorge.MeshData{Format: posColorFormat}
	linesMesh := gorge.NewMesh(lines)
	linesMesh.DrawMode = gorge.DrawLines

	linesEntity := &lineEntity{
		TransformComponent:  gorge.TransformIdent(),
		RenderableComponent: *gorge.NewRenderableComponent(linesMesh, mat),
	}
	linesEntity.SetParent(rootTrans)

	points := &gorge.MeshData{Format: posColorFormat}
	pointsMesh := gorge.NewMesh(points)
	pointsMesh.DrawMode = gorge.DrawPoints

	pointsEntity := &lineEntity{
		TransformComponent:  gorge.TransformIdent(),
		RenderableComponent: *gorge.NewRenderableComponent(pointsMesh, mat),
	}
	pointsEntity.SetParent(rootTrans)

	return &Lines{
		TransformComponent: rootTrans,
		Material:           mat,
		color:              gm.Vec4{1, 1, 1, 1},
		Lines:              lines,
		Points:             points,
		Container: gorge.Container{
			linesEntity,
			pointsEntity,
		},
	}
}

// SetColor set the current color state.
func (dg *Lines) SetColor(r, g, b, a float32) {
	dg.color = gm.Vec4{r, g, b, a}
}

// Clear clear the debug information.
func (dg *Lines) Clear() {
	dg.Points.Vertices = dg.Points.Vertices[:0]
	dg.Points.Updates++
	dg.Lines.Vertices = dg.Lines.Vertices[:0]
	dg.Lines.Updates++
}

// SetCullMask sets the cull mask for the debug renderables.
func (dg *Lines) SetCullMask(m gorge.CullMaskFlags) {
	type renderabler interface {
		Renderable() *gorge.RenderableComponent
	}
	for _, e := range dg.Container {
		if e, ok := e.(renderabler); ok {
			r := e.Renderable()
			r.CullMask = m
		}
	}
}

// AddLine adds a line using the current color.
func (dg *Lines) AddLine(p1 gm.Vec3, p2 gm.Vec3) {
	c := dg.color
	dg.Lines.Vertices = append(dg.Lines.Vertices,
		p1[0], p1[1], p1[2], c[0], c[1], c[2], c[3],
		p2[0], p2[1], p2[2], c[0], c[1], c[2], c[3],
	)
	dg.Lines.Updates++
}

// AddRect3 adds a rect based on 3 points and the current color.
func (dg *Lines) AddRect3(p1, p2, p3 gm.Vec3) {
	e1 := p2.Sub(p1)
	e2 := p3.Sub(p1)

	p4 := p1.Add(e1).Add(e2)

	dg.AddRect(p1, p2, p4, p3)
}

// AddRect adds a rect based on 4 points and the current color.
func (dg *Lines) AddRect(p1, p2, p3, p4 gm.Vec3) {
	c := dg.color
	dg.Lines.Vertices = append(dg.Lines.Vertices,
		p1[0], p1[1], p1[2], c[0], c[1], c[2], c[3],
		p2[0], p2[1], p2[2], c[0], c[1], c[2], c[3],

		p2[0], p2[1], p2[2], c[0], c[1], c[2], c[3],
		p3[0], p3[1], p3[2], c[0], c[1], c[2], c[3],

		p3[0], p3[1], p3[2], c[0], c[1], c[2], c[3],
		p4[0], p4[1], p4[2], c[0], c[1], c[2], c[3],

		p4[0], p4[1], p4[2], c[0], c[1], c[2], c[3],
		p1[0], p1[1], p1[2], c[0], c[1], c[2], c[3],
	)
	dg.Lines.Updates++
}

// AddCircle adds a circle based on a center point and a radius.
func (dg *Lines) AddCircle(p gm.Vec3, r float32) {
	c := dg.color

	sz := 90
	var vert []float32
	m := gm.M3Rotate((math.Pi * 2) / float32(sz))
	vp := gm.Vec2{r, 0}
	for i := 0; i < sz; i++ {
		vert = append(vert,
			p[0]+vp[0]*r, p[1]+vp[1]*r, p[2], c[0], c[1], c[2], c[3],
		)
		vp = m.MulV3(vp.Vec3(1)).Vec2()
		vert = append(vert,
			p[0]+vp[0]*r, p[1]+vp[1]*r, p[2], c[0], c[1], c[2], c[3],
		)
	}
	dg.Lines.Vertices = append(dg.Lines.Vertices, vert...)
	dg.Lines.Updates++
}

// AddPoint adds a single point with the current color..
func (dg *Lines) AddPoint(p gm.Vec3) {
	dg.Points.Vertices = append(dg.Points.Vertices,
		p[0], p[1], p[2], dg.color[0], dg.color[1], dg.color[2], dg.color[3],
	)
	dg.Points.Updates++
}
