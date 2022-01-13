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
	"runtime"
	"strings"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/layerfs"
	"github.com/stdiopt/gorge/static"
)

const gorgeStatic = "_gorge/"

type refCounter struct {
	count int
	ref   any // gorge.ResourceRef
}

// Resource the resource manager.
type Resource struct {
	gorge *gorge.Context
	fs    layerfs.FS
	refs  map[string]*refCounter
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
	gorge.Trigger(r.gorge, EventOpen{Name: name})
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
func (r *Resource) Load(v any, name string, opts ...any) error {
	gorge.Trigger(r.gorge, EventLoadStart{
		Name:     name,
		Resource: v,
	})

	err := r.load(v, name, opts...)

	gorge.Trigger(r.gorge, EventLoadComplete{
		Name:     name,
		Resource: v,
		Err:      err,
	})
	return err
}

// MustLoad loads the resource if an error occurs it will be sent to gorge as an event.
func (r *Resource) MustLoad(v any, path string, opts ...any) {
	if err := r.Load(v, path, opts...); err != nil {
		r.gorge.Error(err)
	}
}

func (r *Resource) Error(err error) {
	r.gorge.Error(err)
}

func (r *Resource) load(v any, name string, opts ...any) error {
	ext := filepath.Ext(name)
	loader := getLoader(v, ext)
	if loader == nil {
		return fmt.Errorf("no driver for type: %T with ext: %v", v, ext)
	}
	return loader(&Context{r}, v, name, opts...)
}

// Tracks resource reference and overrides ref Resource referrer if already
// exists
func (r *Resource) track(name string, ref, v any) (*refCounter, bool) {
	if r.refs == nil {
		r.refs = map[string]*refCounter{}
	}
	tracker, ok := r.refs[name]
	if !ok {
		tracker = &refCounter{
			count: 0,
			ref:   v,
		}
		r.refs[name] = tracker
	}
	tracker.count++
	// ref.res = tracker.res
	runtime.SetFinalizer(ref, func(any) {
		tracker.count--
		if tracker.count == 0 {
			log.Println("Releasing finalizer")
			delete(r.refs, name)
		}
	})
	return tracker, ok
}
