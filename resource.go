package gorge

type gpuResource struct {
	gpu interface{}
}

// NewGPUResource returns a new gpu resource reference
func NewGPUResource() ResourceRef {
	return &gpuResource{}
}

func (r *gpuResource) isResource()          {}
func (r *gpuResource) setGPU(v interface{}) { r.gpu = v }
func (r *gpuResource) getGPU() interface{}  { return r.gpu }

// GetGPU returns gpu data from the resourceRef
func GetGPU(r interface{}) interface{} {
	if r, ok := r.(interface{ getGPU() interface{} }); ok {
		return r.getGPU()
	}
	return nil
}

// SetGPU sets gpu data in the resourceRef
func SetGPU(r, v interface{}) {
	if r, ok := r.(interface{ setGPU(interface{}) }); ok {
		r.setGPU(v)
	}
}

// ////////////////////////////////////////////////////////////////////////////

// ResourceRef Should pinpoint ResourceRef here
type ResourceRef interface {
	isResource()
}

// Resourcer anything that can fetch a resource
/*type Resourcer interface {
	Resource() ResourceRef
}*/

// ResourcerSetter mostly used in Texture or Mesh
/*
type ResourcerSetter interface {
	Resourcer
	SetResourcer(r Resourcer)
}
*/

// GPUResourcer interface for some resources that are used in gpu.
/*type GPUResourcer interface {
	Resourcer
	isGPU()
}*/

// ////////////////////////////////////////////////////////////////////////////

// GPUResource resource reference for textures, mesh, material

// ResourceCopyRef mostly used by resources to load Data and copy gpu ref
/*func ResourceCopyRef(a Resourcer, b ResourceRef) {
	ar := a.Resource()
	switch ar := ar.(type) {
	case *gpuResource:
		SetGPU(ar, GetGPU(b))
		// ar.setGPU(b.(interface{ getGPU() interface{} }).getGPU())
	default:
		panic(fmt.Sprintf("cannot copy referrer of %T", ar))
	}
}*/

// Resource reference
/*
type resourceRef struct {
	res ResourceRef
}

func (r *resourceRef) Resource() ResourceRef { return r.res }

// ResourceReleaseData triggers an update resource event to sync with systems
// for the specific and releases the underlying data, it will reuse the
// resource ref which the data is already binded in gpu
/*func (g *Gorge) ResourceReleaseData(r ResourcerSetter) {
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
*/
