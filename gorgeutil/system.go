package gorgeutil

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/primitive"
)

// Context extends gorge context with default entity instantiators.
type Context struct {
	*gorge.Context
}

// WithContext returns a helper.
func WithContext(g *gorge.Context) *Context {
	return &Context{g}
}

// System initializes gorgeutil context
func System(g *gorge.Context) func() *Context {
	return func() *Context {
		return WithContext(g)
	}
}

// Camera creates and adds a camera to a gorge context.
func (c Context) Camera() *Camera {
	cam := NewCamera()
	c.Add(cam)
	return cam
}

// OrthoCamera creates and adds a camera with default ortho projection options
// to a gorge context.
func (c Context) OrthoCamera(size, near, far float32) *Camera {
	cam := NewOrthoCamera(size, near, far)
	c.Add(cam)
	return cam
}

// PerspectiveCamera creates and adds a camera with default perspective options to a gorge context.
func (c Context) PerspectiveCamera(fov, near, far float32) *Camera {
	cam := NewPerspectiveCamera(fov, near, far)
	c.Add(cam)
	return cam
}

// UICamera creates an ortho camera and adds to context.
func (c Context) UICamera() *Camera {
	cam := NewCamera()
	cam.SetOrtho(100, -100, 100)
	cam.SetCullMask(gorge.MaskUI)
	cam.SetOrder(100)
	cam.SetClearFlag(gorge.ClearDepthOnly)
	c.Add(cam)
	return cam
}

// Light creates and adds a light to a gorge context.
func (c Context) Light() *Light {
	light := NewLight()
	c.Add(light)
	return light
}

// PointLight creates and adds a light to gorge context.
func (c Context) PointLight() *Light {
	light := NewPointLight()
	c.Add(light)
	return light
}

// DirectionalLight creates and adds a directional light to a gorge context.
func (c Context) DirectionalLight() *Light {
	light := NewDirectionalLight()
	c.Add(light)
	return light
}

// SpotLight creates and adds a new spot light to gorge context.
func (c Context) SpotLight() *Light {
	light := NewSpotLight()
	c.Add(light)
	return light
}

// Renderable creates and adds a renderable to gorge context.
func (c Context) Renderable(mesh gorge.Mesher, mat gorge.Materialer) *Renderable {
	r := NewRenderable(mesh, mat)
	c.Add(r)
	return r
}

// This might be moved to primitives again

// Sphere Creates and adds a sphere renderable.
func (c Context) Sphere(sector, stack int) *Renderable {
	r := NewSphere(sector, stack)
	c.Add(r)
	return r
}

// Cube creates and adds a cube.
func (c Context) Cube() *Renderable {
	r := NewCube()
	c.Add(r)
	return r
}

// Plane creates and adds a plane.
func (c Context) Plane(dir primitive.PlaneDir) *Renderable {
	r := NewPlane(dir)
	c.Add(r)
	return r
}