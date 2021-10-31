package gorge

import (
	"fmt"
)

// Different than previous loaders
// we don't have the load method
// instead we fire an load event in whatever loads stuff

type resourcer struct {
	Resourcer
}

func (r *resourcer) Resource() ResourceRef {
	if r.Resourcer == nil {
		return nil
	}
	return r.Resourcer.Resource()
}

func (r *resourcer) SetResourcer(rr Resourcer) {
	r.Resourcer = rr
}

// ////////////////////////////////////////////////////////////////////////////

// ResourceRef Should pinpoint ResourceRef here
type ResourceRef interface {
	isResource()
}

// Resourcer anything that can fetch a resource
type Resourcer interface {
	Resource() ResourceRef
}

// ResourcerSetter mostly used in Texture or Mesh
type ResourcerSetter interface {
	Resourcer
	SetResourcer(r Resourcer)
}

// GPUResourcer interface for some resources that are used in gpu.
type GPUResourcer interface {
	Resourcer
	isGPU()
}

// ////////////////////////////////////////////////////////////////////////////

// GPUResource resource reference for textures, mesh, material
type gpuResource struct {
	gpu interface{}
}

// NewGPUResource returns a new gpu resource reference
func NewGPUResource() ResourceRef {
	return &gpuResource{}
}

func (r *gpuResource) isResource()          {}
func (r *gpuResource) SetGPU(v interface{}) { r.gpu = v }
func (r *gpuResource) GetGPU() interface{}  { return r.gpu }

// ResourceCopyRef mostly used by resources to load Data and copy gpu ref
func ResourceCopyRef(a Resourcer, b ResourceRef) {
	ar := a.Resource()
	switch ar := ar.(type) {
	case *gpuResource:
		ar.SetGPU(b.(interface{ GetGPU() interface{} }).GetGPU())
	default:
		panic(fmt.Sprintf("cannot copy referrer of %T", ar))
	}
}

// Resource reference
type resourceRef struct {
	res ResourceRef
}

func (r *resourceRef) Resource() ResourceRef { return r.res }

// ResourceReleaseData triggers an update resource event to sync with systems
// for the specific and releases the underlying data, it will reuse the
// resource ref which the data is already binded in gpu
func (g *Gorge) ResourceReleaseData(r ResourcerSetter) {
	curRes := r.Resource()

	ref := &resourceRef{}
	if _, ok := r.(GPUResourcer); ok {
		ref.res = NewGPUResource()
	}
	r.SetResourcer(ref)

	// Force update
	g.Trigger(EventResourceUpdate{Resource: curRes})
	ResourceCopyRef(ref, curRes)
}

// ResourceRef Upload data to gpu and return a ref
func (g *Gorge) ResourceRef(r ResourceRef) Resourcer {
	g.Trigger(EventResourceUpdate{Resource: r})
	if _, ok := r.(interface{ GetGPU() interface{} }); ok {
		ref := &resourceRef{NewGPUResource()}
		ResourceCopyRef(ref, r)
		return ref
	}
	return &resourceRef{r}
}
