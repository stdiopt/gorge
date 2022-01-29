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

// AddPlane to entityAdder.
func AddPlane(a entityAdder, dir primitive.PlaneDir) *Entity {
	p := NewPlane(dir)
	a.Add(p)
	return p
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

// AddSphere to entityAdder.
func AddSphere(a entityAdder, sector, stack int) *Entity {
	s := NewSphere(sector, stack)
	a.Add(s)
	return s
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

// AddCube to entityAdder.
func AddCube(a entityAdder) *Entity {
	c := NewCube()
	a.Add(c)
	return c
}
