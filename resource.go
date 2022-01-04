package gorge

// GPU reference for resources binded in renderer (texture,mesh)
type GPU struct {
	gpu interface{}
}

func (r *GPU) isGPU()               {}
func (r *GPU) setGPU(v interface{}) { r.gpu = v }
func (r *GPU) getGPU() interface{}  { return r.gpu }

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
