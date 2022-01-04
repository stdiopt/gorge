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

	ref := &gorge.TextureRef{
		Ref: gorge.NewGPUResource(),
	}
	tex.SetResourcer(ref)

	// Bind the resource, if exists we reuse the resource ref
	counter, ok := r.track(name, ref, ref.Resource())
	if ok {
		ref.Ref = counter.ref.(gorge.ResourceRef)
		return tex
	}

	// Load into a new temporary resourcer and copy the gpu reference
	go func() {
		r.gorge.TriggerInMain(EventLoadStart{
			Name:     name,
			Resource: tex,
		})
		tmp := &gorge.TextureData{}
		if err := r.load(tmp, name, opts...); err != nil {
			r.gorge.TriggerInMain(EventLoadComplete{
				Name:     name,
				Resource: tex,
				Err:      err,
			})

			r.Error(err)
			return
		}
		r.gorge.RunInMain(func() {
			r.gorge.Trigger(gorge.EventResourceUpdate{
				Resource: tmp,
			})
			r.gorge.Trigger(EventLoadComplete{
				Name:     name,
				Resource: tex,
			})
			gorge.SetGPU(counter.ref, gorge.GetGPU(tmp))
			// gorge.ResourceCopyRef(tex.Resourcer, res)
		})
	}()
	return tex
}

// Mesh helper that returns a mesh with a resourcer ref
func (r *Context) Mesh(name string, opts ...interface{}) *gorge.Mesh {
	mesh := gorge.NewMesh(nil)

	ref := &gorge.MeshRef{
		Ref: gorge.NewGPUResource(),
	}
	mesh.SetResourcer(ref)

	// Bind the resource, if exists we reuse the resource ref
	if counter, ok := r.track(name, ref, ref.Resource()); ok {
		ref.Ref = counter.ref.(gorge.ResourceRef)
		return mesh
	}

	// Load into a new temporary resourcer and copy the gpu reference
	go func() {
		r.gorge.TriggerInMain(EventLoadStart{
			Name:     name,
			Resource: mesh,
		})
		tmp := &gorge.TextureData{}
		if err := r.load(tmp, name, opts...); err != nil {
			r.gorge.TriggerInMain(EventLoadComplete{
				Name:     name,
				Resource: mesh,
				Err:      err,
			})

			r.Error(err)
			return
		}
		r.gorge.RunInMain(func() {
			res := tmp
			r.gorge.Trigger(gorge.EventResourceUpdate{
				Resource: res,
			})
			r.gorge.Trigger(EventLoadComplete{
				Name:     name,
				Resource: mesh,
			})
			gorge.SetGPU(mesh.Resourcer.Resource(), gorge.GetGPU(res))
		})
	}()
	return mesh
}

/*
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
*/

// Material loads everytime right away
func (r *Context) Material(name string, opts ...interface{}) *gorge.Material {
	var data gorge.ShaderData
	if err := r.Load(&data, name, opts...); err != nil {
		r.Error(err)
		return gorge.NewMaterial()
	}

	return gorge.NewShaderMaterial(&data)
}
