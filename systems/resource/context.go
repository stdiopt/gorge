package resource

import (
	"fmt"
	"io/fs"
	"log"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/layerfs"
	"github.com/stdiopt/gorge/static"
)

var ctxKey = struct{ string }{"resource"}

type manager = Resource

// Context to be used in gorge systems
type Context struct {
	*manager
}

// FromContext returns a Context from a gorge Context
func FromContext(g *gorge.Context) *Context {
	var ret *Context
	if ctx, ok := gorge.GetSystem(g, ctxKey).(*Context); ok {
		return ctx
	}
	log.Println("Initializing system")

	lfs := layerfs.FS{}
	s, err := fs.Sub(static.Assets, "src")
	if err != nil {
		panic(fmt.Errorf("static embed not found: %w", err))
	}
	lfs.Mount(gorgeStatic, s)

	m := &Resource{gorge: g, fs: lfs}

	ret = &Context{
		manager: m,
	}

	gorge.AddSystem(g, ctxKey, ret)
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
