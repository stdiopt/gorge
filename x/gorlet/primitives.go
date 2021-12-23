package gorlet

import (
	"math"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// Quad renderable entity used in some widgets.
type gEntity struct {
	gorge.TransformComponent
	gorge.ColorableComponent
	*gorge.RenderableComponent
}

// quadEntity returns a quad meshEntity based on primitive.MeshEntity.
func quadEntity() *gEntity {
	mat := gorge.NewShaderMaterial(static.Shaders.UI)
	mat.Queue = 100
	mat.Depth = gorge.DepthNone // Fix this

	mesh := gorge.NewMesh(&gorge.MeshData{
		Format: gorge.VertexFormatPTN(),
		Vertices: []float32{
			0, 1, 0, 0, 0, 0, 0, 1,
			1, 1, 0, 1, 0, 0, 0, 1,
			1, 0, 0, 1, 1, 0, 0, 1,
			0, 0, 0, 0, 1, 0, 0, 1,
		},
		Indices: []uint32{
			0, 2, 1,
			2, 0, 3,
		},
	})
	return &gEntity{
		TransformComponent:  *gorge.NewTransformComponent(),
		ColorableComponent:  *gorge.NewColorableComponent(1, 1, 1, 1),
		RenderableComponent: gorge.NewRenderableComponent(mesh, mat),
	}
}

// PolyMeshData returns a poly as meshData.
func polyEntity(n int) *gEntity {
	points := []float32{}
	p := m32.Vec3{0, .5, 0}
	theta := float32(math.Pi) / (float32(n) / 2)
	r := m32.M3Rotate(theta)
	for i := 0; i < n+1; i++ {
		o := p.Add(m32.Vec3{.5, .5, 0})
		points = append(points, o[:]...)
		p = r.MulV3(p)
	}
	meshData := &gorge.MeshData{
		Format:   gorge.VertexFormatP(),
		Vertices: points,
		Indices:  nil,
	}

	mat := gorge.NewShaderMaterial(static.Shaders.UI)
	mat.Queue = 100
	mat.Depth = gorge.DepthNone

	mesh := gorge.NewMesh(meshData)
	mesh.DrawMode = gorge.DrawTriangleFan

	return &gEntity{
		TransformComponent:  *gorge.NewTransformComponent(),
		ColorableComponent:  *gorge.NewColorableComponent(1, 1, 1, 1),
		RenderableComponent: gorge.NewRenderableComponent(mesh, mat),
	}
}

type graphicer interface {
	Transform() *gorge.TransformComponent
	Colorable() *gorge.ColorableComponent
}

func rectElement(ent graphicer) BuildFunc {
	return func(b *Builder) {
		p := b.Root()
		p.AddElement(ent)
		p.HandleFunc(func(e event.Event) {
			if _, ok := e.(gorgeui.EventUpdate); !ok {
				return
			}
			r := p.Rect()
			t := ent.Transform()
			t.Scale[0] = r[2] - r[0]
			t.Scale[1] = r[3] - r[1]
		})
		b.Observe("color", ObsFunc(func(c m32.Vec4) {
			ent.Colorable().SetColorv(c)
		}))
		// Defaults
		p.Set("color", m32.Color(0, 0, 0, .2))
	}
}

// Quad returns a quad entity starting at 0,0 to 1,1
func Quad() BuildFunc {
	return rectElement(quadEntity())
}

// Poly returns a polygon starting at 0,0 to 1,1
func Poly(n int) BuildFunc {
	return rectElement(polyEntity(n))
}
