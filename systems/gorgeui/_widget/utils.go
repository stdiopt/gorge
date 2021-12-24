package widget

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/static"
)

// Quad renderable entity used in some widgets.
type Quad struct {
	gorge.TransformComponent
	gorge.ColorableComponent
	*gorge.RenderableComponent
}

// QuadEntity returns a quad meshEntity based on primitive.MeshEntity.
func QuadEntity() *Quad {
	mat := gorge.NewShaderMaterial(static.Shaders.UI)
	mat.Queue = 100
	mat.Depth = gorge.DepthNone

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
	return &Quad{
		TransformComponent:  *gorge.NewTransformComponent(),
		ColorableComponent:  *gorge.NewColorableComponent(1, 1, 1, 1),
		RenderableComponent: gorge.NewRenderableComponent(mesh, mat),
	}
}

// QuadEntity returns a quad meshEntity based on primitive.MeshEntity.
/*func QuadEntity() *Quad {
	mat := gorge.NewShaderMaterial(static.Shaders.UI)
	mat.Queue = 100
	mat.Depth = gorge.DepthNone

	mesh := gorge.NewMesh(&gorge.MeshData{
		Format: gorge.VertexFormatPTN(),
		Vertices: []float32{
			0, 1, 0, 0, 0, 0, 0, 1,
			1, 1, 0, 1, 0, 0, 0, 1,
			1, 0, 0, 1, 1, 0, 0, 1,
			0, 0, 0, 0, 1, 0, 0, 1,
		},
		Indices: []uint32{
			0, 1, 2,
			2, 3, 0,
		},
	})
	return &Quad{
		TransformComponent:  *gorge.NewTransformComponent(),
		ColorableComponent:  *gorge.NewColorableComponent(1, 1, 1, 1),
		RenderableComponent: gorge.NewRenderableComponent(mesh, mat),
	}
}*/

func v4Color(v ...float32) m32.Vec4 {
	switch len(v) {
	case 0:
		return m32.Vec4{}
	case 1:
		return m32.Vec4{v[0], v[0], v[0], 1}
	case 2:
		return m32.Vec4{v[0], v[0], v[0], v[1]}
	case 3:
		return m32.Vec4{v[0], v[1], v[2], 1}
	default:
		return m32.Vec4{v[0], v[1], v[2], v[3]}
	}
}
