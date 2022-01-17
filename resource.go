package gorge

import "fmt"

// GPU reference for resources binded in renderer (texture,mesh)
type GPU struct {
	gpu any
}

func (r *GPU) isGPU()       {}
func (r *GPU) setGPU(v any) { r.gpu = v }
func (r *GPU) getGPU() any  { return r.gpu }

// GetGPU returns gpu data from the resourceRef
func GetGPU(r any) any {
	if r, ok := r.(interface{ getGPU() any }); ok {
		return r.getGPU()
	}
	return nil
}

// SetGPU sets gpu data in the resourceRef
func SetGPU(r, v any) {
	g, ok := r.(interface{ setGPU(any) })
	if !ok {
		panic(fmt.Sprintf("%T is not a gpu resource", r))
	}
	g.setGPU(v)
}

// TextureRef implements a gpu only texture resource.
type TextureRef struct{ *GPU }

func (r *TextureRef) Resource() TextureResource { return r }

// Resource implements the TextureResrouceresourcer interface.
func (r *TextureRef) isTexture() {}

// MeshRef implements a gpu only mesh resource
type MeshRef struct{ *GPU }

func (r *MeshRef) Resource() MeshResource { return r }

// Resource returns the resource ref.
func (r *MeshRef) isMesh() {}
