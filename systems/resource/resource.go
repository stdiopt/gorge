// Package resource handles resources
package resource

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/layerfs"
	"github.com/stdiopt/gorge/static"
)

const gorgeStatic = "_gorge/"

// System gorge system initializer func.
func System(g *gorge.Context) {
	log.Println("Initializing system")

	lfs := layerfs.FS{}
	s, err := fs.Sub(static.Assets, "src")
	if err != nil {
		panic(fmt.Errorf("static embed not found: %w", err))
	}
	lfs.Mount(gorgeStatic, s)

	m := &Resource{
		gorge: g,
		fs:    lfs,
	}

	g.PutProp(func() *Context {
		return &Context{m}
	})
}

// Resource the resource manager.
type Resource struct {
	gorge   *gorge.Context
	fs      layerfs.FS
	tracker map[string]*resTrack
}

// AddFS adds a new file system with the prefix if a path exists it will overlay
// the existing file system.
func (r *Resource) AddFS(prefix string, fs fs.FS) {
	r.fs.Mount(prefix, fs)
}

// Gorge convinient way to return a gorge context to loaders.
func (r *Resource) Gorge() *gorge.Context {
	return r.gorge
}

// LoadBytes returns the asset as string
func (r *Resource) LoadBytes(name string) ([]byte, error) {
	rd, err := r.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rd.Close(); err != nil {
			r.Error(err)
		}
	}()

	data, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// LoadString returns asset as a string
func (r *Resource) LoadString(name string) (string, error) {
	s, err := r.LoadBytes(name)
	return string(s), err
}

// Open opens a reousrce based on the configured sourcer.
func (r *Resource) Open(name string) (io.ReadCloser, error) {
	r.gorge.Trigger(EventOpen{
		Name: name,
	})
	if strings.HasPrefix(name, gorgeStatic) {
		data, err := static.Data(name[len(gorgeStatic):])
		if err != nil {
			return nil, err
		}
		return ioutil.NopCloser(bytes.NewReader(data)), nil
	}
	return r.fs.Open(name)
}

// Load stuff
func (r *Resource) Load(v interface{}, name string, opts ...interface{}) error {
	r.gorge.Trigger(EventLoadStart{
		Name:     name,
		Resource: v,
	})
	ext := filepath.Ext(name)
	loader := getLoader(v, filepath.Ext(name))
	if loader == nil {
		return fmt.Errorf("no driver for type: %T with ext: %v", v, ext)
	}
	if err := loader(&Context{r}, v, name, opts...); err != nil {
		return err
	}
	r.gorge.Trigger(EventLoadComplete{
		Name:     name,
		Resource: v,
	})
	return nil
}

// MustLoad loads the resource if an error occurs it will be sent to gorge as an event.
func (r *Resource) MustLoad(v interface{}, path string, opts ...interface{}) {
	if err := r.Load(v, path, opts...); err != nil {
		r.gorge.Error(err)
	}
}

func (r *Resource) Error(err error) {
	r.gorge.Error(err)
}

// loadedRef will be used as a resource in Mesh or Texture
// the purpose of this is avoid load duplication
// if the resource is already loaded a new loadedRef with an existing resource ref
type loadedRef struct {
	res gorge.ResourceRef
}

func (r loadedRef) Resource() gorge.ResourceRef { return r.res }

type resTrack struct {
	count int
	res   gorge.ResourceRef
}

// Tracks resource reference and overrides ref Resource referrer if already
// exists
func (r *Resource) track(name string, ref *loadedRef) bool {
	if r.tracker == nil {
		r.tracker = map[string]*resTrack{}
	}
	tracker, ok := r.tracker[name]
	if !ok {
		tracker = &resTrack{
			count: 0,
			res:   ref.res,
		}
		r.tracker[name] = tracker
	}
	tracker.count++
	ref.res = tracker.res
	runtime.SetFinalizer(ref, func(_ interface{}) {
		tracker.count--
		if tracker.count == 0 {
			delete(r.tracker, name)
		}
	})
	return ok
}

// LoadRef sets a New loader reference and check if resource reference exists
// if the reference exists we update the loader reference with the specified
// resource else it will load in background and update the loader reference once done.
func (r *Resource) LoadRef(rs gorge.ResourcerSetter, name string, opts ...interface{}) {
	r.gorge.Trigger(EventLoadStart{
		Name:     name,
		Resource: rs,
	})

	ref := &loadedRef{}
	if _, ok := rs.(gorge.GPUResourcer); ok {
		ref.res = gorge.NewGPUResource()
	}
	rs.SetResourcer(ref)

	// Bind the resource, if exists we reuse the resource ref
	if ok := r.track(name, ref); ok {
		return
	}

	// Load into a new temporary resourcer and copy the gpu reference
	go func() {
		rr := reflect.New(reflect.TypeOf(rs).Elem()).Interface().(gorge.Resourcer)
		if err := r.Load(rr, name, opts...); err != nil {
			r.Error(err)
			return
		}
		res := rr.Resource()
		r.UpdateResource(res)

		gorge.ResourceCopyRef(rs, res)
	}()
}

// UpdateResource triggers a resource event to allow systems to aknowldge the resource.
func (r *Resource) UpdateResource(rr gorge.ResourceRef) {
	r.gorge.TriggerOnUpdate(gorge.EventResourceUpdate{
		Resource: rr,
	})
}
