package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// Quad renderable entity used in some widgets.
type gEntity struct {
	gorge.TransformComponent
	gorge.ColorableComponent
	*gorge.RenderableComponent
}

func newEntity(mesh gorge.Mesher) *gEntity {
	mat := gorge.NewShaderMaterial(static.Shaders.UI)
	mat.Queue = 3000
	mat.Depth = gorge.DepthRead
	return &gEntity{
		TransformComponent:  *gorge.NewTransformComponent(),
		ColorableComponent:  *gorge.NewColorableComponent(1, 1, 1, 1),
		RenderableComponent: gorge.NewRenderableComponent(mesh, mat),
	}
}

// quadEntity returns a quad meshEntity based on primitive.MeshEntity.
func quadEntity() *gEntity {
	return newEntity(gorge.NewMesh(quadMeshData()))
}

// PolyMeshData returns a poly as meshData.
func polyEntity(n int) *gEntity {
	return newEntity(gorge.NewMesh(polyMeshData(n)))
}

type graphicer interface {
	Transform() *gorge.TransformComponent
	Colorable() *gorge.ColorableComponent
	Renderable() *gorge.RenderableComponent
}

func rectElement(ent graphicer) Func {
	return func(b *Builder) {
		root := b.Root()
		root.AddElement(ent)
		// Defaults renderable, use it on label too
		b.Observe("color", ObsFunc(ent.Colorable().SetColorv))
		b.Observe("material", ObsFunc(ent.Renderable().SetMaterial))
		b.Observe("texture", ObsFunc(func(tex gorge.Texturer) {
			ent.Renderable().Material.SetTexture("albedoMap", tex)
		}))
		b.Observe("stencil", ObsFunc(func(s *gorge.Stencil) {
			ent.Renderable().Stencil = s
		}))
		b.Observe("colorMask", ObsFunc(func(b *[4]bool) {
			ent.Renderable().ColorMask = b
		}))
		// Forget order here
		b.Observe("order", ObsFunc(ent.Renderable().SetOrder))
		b.Observe("_maskDepth", ObsFunc(func(n int) {
			s := calcMaskOn(n)
			s.WriteMask = 0
			s.Fail = gorge.StencilOpKeep
			s.ZFail = gorge.StencilOpKeep
			s.ZPass = gorge.StencilOpKeep

			ent.Renderable().Stencil = s
			for _, c := range root.Children() {
				c.Set("_maskDepth", n)
			}
		}))
		event.Handle(root, func(gorgeui.EventUpdate) {
			r := root.Rect()
			t := ent.Transform()
			t.Scale[0] = r[2] - r[0]
			t.Scale[1] = r[3] - r[1]
		})
		// Defaults
		root.Set("color", gm.Color(0, 0, 0, .2))
	}
}

// Quad returns a quad entity starting at 0,0 to 1,1
func Quad() Func {
	return rectElement(quadEntity())
}

func (b *Builder) Quad() *Entity {
	return b.Add(Quad())
}

// Poly returns a polygon starting at 0,0 to 1,1
func Poly(n int) Func {
	return rectElement(polyEntity(n))
}
