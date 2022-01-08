package gorgeutil

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/primitive"
)

// Entity entity.
type Entity struct {
	Name string
	gorge.TransformComponent
	*gorge.RenderableComponent
	*gorge.ColorableComponent
}

// RenderableWithOptions
/*func RenderableWithOptions(opts ...gorge.EntityFunc) *Renderable {
	mat := gorge.NewMaterial()
	m := &Renderable{
		TransformComponent:  gorge.TransformIdent(),
		RenderableComponent: gorge.NewRenderableComponent(nil, mat),
		ColorableComponent:  gorge.NewColorableComponent(1, 1, 1, 1),
	}
	gorge.ApplyTo(m, opts...)
	return m
}*/

// NewRenderable returns a new renderable.
func NewRenderable(mesh gorge.Mesher, mat gorge.Materialer) *Entity {
	return &Entity{
		TransformComponent:  gorge.TransformIdent(),
		RenderableComponent: gorge.NewRenderableComponent(mesh, mat),
		ColorableComponent:  gorge.NewColorableComponent(1, 1, 1, 1),
	}
}

// SetName sets renderable name.
func (r *Entity) SetName(n string) {
	r.Name = n
}

// NewPlane returns a new plane.
func NewPlane(dir primitive.PlaneDir) *Entity {
	mat := gorge.NewMaterial()
	mesh := primitive.NewPlane(dir)
	return &Entity{
		TransformComponent:  gorge.TransformIdent(),
		RenderableComponent: gorge.NewRenderableComponent(mesh, mat),
		ColorableComponent:  gorge.NewColorableComponent(1, 1, 1, 1),
	}
}

// NewSphere returns a new renderable sphere.
func NewSphere(sector, stack int) *Entity {
	mat := gorge.NewMaterial()
	mesh := primitive.NewSphere(sector, stack)
	return &Entity{
		TransformComponent:  gorge.TransformIdent(),
		RenderableComponent: gorge.NewRenderableComponent(mesh, mat),
		ColorableComponent:  gorge.NewColorableComponent(1, 1, 1, 1),
	}
}

// NewCube returns a new renderable cube.
func NewCube() *Entity {
	mat := gorge.NewMaterial()
	mesh := primitive.NewCube()
	return &Entity{
		TransformComponent:  gorge.TransformIdent(),
		RenderableComponent: gorge.NewRenderableComponent(mesh, mat),
		ColorableComponent:  gorge.NewColorableComponent(1, 1, 1, 1),
	}
}
