package gorgeutil

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/static"
)

type renderabler interface {
	Renderable() *gorge.RenderableComponent
}

// GeomEntity generic renderable entity.
type GeomEntity struct {
	gorge.TransformComponent
	gorge.RenderableComponent
}

// Geom helper to draw per frame debug information.
type Geom struct {
	*gorge.TransformComponent
	*gorge.Material
	color gm.Vec4

	lines  *gorge.MeshData
	points *gorge.MeshData

	gorge.Container
}

// NewGeom returns a Geom entity.
func NewGeom() *Geom {
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

	linesEntity := &GeomEntity{
		TransformComponent:  gorge.TransformIdent(),
		RenderableComponent: *gorge.NewRenderableComponent(linesMesh, mat),
	}
	linesEntity.SetParent(rootTrans)

	points := &gorge.MeshData{Format: posColorFormat}
	pointsMesh := gorge.NewMesh(points)
	pointsMesh.DrawMode = gorge.DrawPoints

	pointsEntity := &GeomEntity{
		TransformComponent:  gorge.TransformIdent(),
		RenderableComponent: *gorge.NewRenderableComponent(pointsMesh, mat),
	}
	pointsEntity.SetParent(rootTrans)

	return &Geom{
		TransformComponent: rootTrans,
		Material:           mat,
		lines:              lines,
		points:             points,
		Container: gorge.Container{
			linesEntity,
			pointsEntity,
		},
	}
}

// SetColor set the current color state.
func (dg *Geom) SetColor(r, g, b, a float32) {
	dg.color = gm.Vec4{r, g, b, a}
}

// Clear clear the debug information.
func (dg *Geom) Clear() {
	dg.points.Vertices = dg.points.Vertices[:0]
	dg.points.Updates++
	dg.lines.Vertices = dg.lines.Vertices[:0]
	dg.lines.Updates++
}

// SetCullMask sets the cull mask for the debug renderables.
func (dg *Geom) SetCullMask(m gorge.CullMaskFlags) {
	for _, e := range dg.Container {
		if e, ok := e.(renderabler); ok {
			r := e.Renderable()
			r.CullMask = m
		}
	}
}

// AddLine adds a line using the current color.
func (dg *Geom) AddLine(p1 gm.Vec3, p2 gm.Vec3) {
	dg.lines.Vertices = append(dg.lines.Vertices,
		p1[0], p1[1], p1[2], dg.color[0], dg.color[1], dg.color[2], dg.color[3],
		p2[0], p2[1], p2[2], dg.color[0], dg.color[1], dg.color[2], dg.color[3],
	)
	dg.lines.Updates++
}

// AddRect3 adds a rect based on 3 points and the current color.
func (dg *Geom) AddRect3(p1, p2, p3 gm.Vec3) {
	e1 := p2.Sub(p1)
	e2 := p3.Sub(p1)

	p4 := p1.Add(e1).Add(e2)

	dg.AddRect(p1, p2, p4, p3)
}

// AddRect adds a rect based on 4 points and the current color.
func (dg *Geom) AddRect(p1, p2, p3, p4 gm.Vec3) {
	dg.lines.Vertices = append(dg.lines.Vertices,
		p1[0], p1[1], p1[2], dg.color[0], dg.color[1], dg.color[2], dg.color[3],
		p2[0], p2[1], p2[2], dg.color[0], dg.color[1], dg.color[2], dg.color[3],

		p2[0], p2[1], p2[2], dg.color[0], dg.color[1], dg.color[2], dg.color[3],
		p3[0], p3[1], p3[2], dg.color[0], dg.color[1], dg.color[2], dg.color[3],

		p3[0], p3[1], p3[2], dg.color[0], dg.color[1], dg.color[2], dg.color[3],
		p4[0], p4[1], p4[2], dg.color[0], dg.color[1], dg.color[2], dg.color[3],

		p4[0], p4[1], p4[2], dg.color[0], dg.color[1], dg.color[2], dg.color[3],
		p1[0], p1[1], p1[2], dg.color[0], dg.color[1], dg.color[2], dg.color[3],
	)
	dg.lines.Updates++
}

// AddPoint adds a single point with the current color..
func (dg *Geom) AddPoint(p gm.Vec3) {
	dg.points.Vertices = append(dg.points.Vertices,
		p[0], p[1], p[2], dg.color[0], dg.color[1], dg.color[2], dg.color[3],
	)
	dg.points.Updates++
}
