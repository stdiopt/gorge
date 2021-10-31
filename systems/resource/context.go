package resource

import "github.com/stdiopt/gorge"

type manager = Resource

// Context to be used in gorge systems
type Context struct {
	*manager
}

// FromContext returns a Context from a gorge Context
func FromContext(g *gorge.Context) *Context {
	var ret *Context
	if err := g.BindProps(func(c *Context) { ret = c }); err != nil {
		g.Error(err)
	}
	return ret
}

// Texture helper that returns a texture with a resourcer ref
func (r *Context) Texture(name string, opts ...interface{}) *gorge.Texture {
	tex := gorge.NewTexture(nil)
	r.LoadRef(tex, name, opts...)
	return tex
}

// Mesh helper that returns a mesh with a resourcer ref
func (r *Context) Mesh(name string, opts ...interface{}) *gorge.Mesh {
	mesh := gorge.NewMesh(nil)
	r.LoadRef(mesh, name, opts...)
	return mesh
}

// Material loads everytime right away
func (r *Context) Material(name string, opts ...interface{}) *gorge.Material {
	var data gorge.ShaderData
	if err := r.Load(&data, name, opts...); err != nil {
		r.Error(err)
		return gorge.NewMaterial()
	}

	return gorge.NewShaderMaterial(&data)
}
