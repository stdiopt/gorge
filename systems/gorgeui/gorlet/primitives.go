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
	mat.DisableShadow = true
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
		Observe(b, "color", ent.Colorable().SetColorv)
		Observe(b, "material", ent.Renderable().SetMaterial)
		Observe(b, "texture", func(tex gorge.Texturer) {
			ent.Renderable().Material.SetTexture("albedoMap", tex)
		})
		Observe(b, "stencil", ent.Renderable().SetStencil)
		Observe(b, "colorMask", func(b *[4]bool) {
			ent.Renderable().ColorMask = b
		})
		// Forget order here
		Observe(b, "order", ent.Renderable().SetOrder)
		Observe(b, "_maskDepth", func(n int) {
			s := calcMaskOn(n)
			s.WriteMask = 0
			s.Fail = gorge.StencilOpKeep
			s.ZFail = gorge.StencilOpKeep
			s.ZPass = gorge.StencilOpKeep

			ent.Renderable().Stencil = s
			for _, c := range root.Children() {
				c.Set("_maskDepth", n)
			}
		})
		Observe(b, "border", root.SetBorderv)
		Observe(b, "borderColor", func(v gm.Vec4) {
			ent.Renderable().Material.Set("borderColor", v)
		})
		event.Handle(root, func(gorgeui.EventUpdate) {
			r := root.Rect()
			t := ent.Transform()
			w := r[2] - r[0]
			h := r[3] - r[1]
			t.Scale[0] = w + root.Border[2] + root.Border[0]
			t.Scale[1] = h + root.Border[3] + root.Border[1]
			t.Position[0] = -root.Border[0]
			t.Position[1] = -root.Border[1]
			if root.Border != (gm.Vec4{}) {
				ent.Renderable().Material.Define("HAS_BORDER")
				ent.Renderable().Material.Set("border", root.Border)
			}
			ent.Renderable().Material.Set("rect", r)
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
