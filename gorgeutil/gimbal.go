package gorgeutil

import (
	"math"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/primitive"
	"github.com/stdiopt/gorge/static"
)

// MeshEntity thing
type MeshEntity struct {
	*gorge.TransformComponent
	*gorge.RenderableComponent
	*gorge.ColorableComponent
}

// Gimbal Compost object
type Gimbal struct {
	Entities []gorge.Entity
	*gorge.TransformComponent
}

// GetEntities implement the gorge entity container.
func (g Gimbal) GetEntities() []gorge.Entity {
	return g.Entities
}

// NewGimbal creates entities on manager
func NewGimbal() *Gimbal {
	// Parent thing
	root := gorge.NewTransformComponent()

	line := gorge.NewMesh(&gorge.MeshData{
		Name:   "gorgeutil.Gimbal",
		Format: gorge.VertexFormatP(),
		Vertices: []float32{
			0, 0, 0,
			0, 0, 1,
		},
		Indices: []uint32{},
	})
	line.DrawMode = gorge.DrawLines

	rot90 := float32(math.Pi / 2)

	objs := []struct {
		axis m32.Vec3
		rot  m32.Vec3
	}{
		{axis: m32.Vec3{0, 0, 1}, rot: m32.Vec3{}},
		{axis: m32.Vec3{0, 1, 0}, rot: m32.Vec3{-rot90, 0, 0}},
		{axis: m32.Vec3{1, 0, 0}, rot: m32.Vec3{0, rot90, 0}},
	}

	gm := &Gimbal{
		Entities:           []gorge.Entity{},
		TransformComponent: root,
	}

	mat := gorge.NewShaderMaterial(static.Shaders.Unlit)
	lineRenderable := gorge.NewRenderableComponent(line, mat)
	for _, o := range objs {
		color := o.axis.Vec4(1)

		l := &MeshEntity{
			gorge.NewTransformComponent(),
			lineRenderable,
			gorge.NewColorableComponent(color[0], color[1], color[2], 1),
		}
		l.SetParent(root)
		l.CastShadow = gorge.CastShadowDisabled
		l.Rotatev(o.rot)

		gm.Entities = append(gm.Entities, l)

	}
	cubeMesh := primitive.NewCube()
	cubeRenderable := gorge.NewRenderableComponent(cubeMesh, mat)
	for _, o := range objs {
		color := o.axis.Vec4(1)

		b := &MeshEntity{
			gorge.NewTransformComponent(),
			cubeRenderable,
			gorge.NewColorableComponent(color[0], color[1], color[2], 1),
		}
		b.SetParent(root)
		b.SetCastShadow(gorge.CastShadowDisabled)
		b.SetPositionv(o.axis)
		b.SetScale(0.08)

		gm.Entities = append(gm.Entities, b)
	}
	return gm
}
