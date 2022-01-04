package gorge

import (
	"fmt"

	"github.com/stdiopt/gorge/systems/render/gl"
)

// MaterialResourcer is the interface for material resources.
type MaterialResourcer interface {
	Resource() ResourceRef
	isMaterial()
}

// Material the material
type Material struct {
	parent    *Material
	Resourcer MaterialResourcer
	Name      string
	// Primitive stuff
	Queue       int
	Depth       DepthMode
	DoubleSided bool
	Blend       BlendType

	/* New: stencil experiment */
	// Create stencil groups?!
	Stencil     bool
	StencilMask uint32
	StencilFunc StencilFunc
	StencilOp   StencilOp

	shaderProps
}

// NewMaterial returns a material using the default gorge shader
func NewMaterial() *Material {
	return &Material{}
}

// Material implements the materialer interface.
func (m *Material) Material() *Material {
	return m
}

// NewShaderMaterial returns a new material based on shader data
// if ShaderData is nil it will use the default PBR material
func NewShaderMaterial(r MaterialResourcer) *Material {
	return &Material{
		Resourcer: r,
	}
}

// Resource implements the resourcer interface.
func (m *Material) GetResource() ResourceRef {
	if m.Resourcer == nil {
		return nil
	}
	return m.Resourcer.Resource()
}

// SetResourcer implements the resource setter interface.
func (m *Material) SetResourcer(r MaterialResourcer) { m.Resourcer = r }

func (m Material) String() string {
	return fmt.Sprintf("(material: %q)", m.Name)
}

// SetQueue sets the material target queue
// transparent materials should be in a higher queue
func (m *Material) SetQueue(v int) {
	m.Queue = v
}

// SetDepth sets if material uses depth buffer.
func (m *Material) SetDepth(v DepthMode) {
	m.Depth = v
}

// SetDoubleSided sets material double sided prop.
func (m *Material) SetDoubleSided(v bool) {
	m.DoubleSided = true
}

// SetBlend sets the blend type for material.
func (m *Material) SetBlend(v BlendType) {
	m.Blend = v
}

// SetStencil sets the stencil property for material.
func (m *Material) SetStencil(v bool) {
	m.Stencil = v
}

// SetStencilMask sets the stencil mask for material.
func (m *Material) SetStencilMask(v uint32) {
	m.StencilMask = v
}

// SetStencilFunc sets the stencil func which
// stencilFunc describes whether OpenGL should pass or discard fragments based
// on the stencil buffer's content
func (m *Material) SetStencilFunc(f gl.Enum, ref int, mask uint32) {
	m.StencilFunc = StencilFunc{f, ref, mask}
}

// SetStencilOp describes the action when updating the stencil buffer.
func (m *Material) SetStencilOp(fail, zfail, zpass gl.Enum) {
	m.StencilOp = StencilOp{fail, zfail, zpass}
}

// Defines override shaderProp defines with hierarchy
func (m *Material) Defines() map[string]string {
	if m.parent == nil {
		return m.shaderProps.Defines()
	}
	ret := map[string]string{}
	for k, v := range m.parent.Defines() {
		ret[k] = v
	}
	for k, v := range m.defines {
		ret[k] = v
	}
	return ret
}

// DefinesHash returns an hash based on the material defines.
func (m *Material) DefinesHash() uint {
	if m == nil {
		return 0
	}
	hash := uint(0)
	if m.parent != nil {
		hash ^= m.parent.DefinesHash()
	}
	hash ^= m.shaderProps.DefinesHash()
	return hash
}

// Get returns the material property for name or if the material has a Parent
// material returns the property from parent else returns nil.
func (m *Material) Get(name string) interface{} {
	if m.parent == nil {
		return m.shaderProps.Get(name)
	}
	if v := m.shaderProps.Get(name); v != nil {
		return v
	}
	return m.parent.Get(name)
}

// GetTexture returns the texture for name
func (m *Material) GetTexture(name string) *Texture {
	if t := m.shaderProps.GetTexture(name); t != nil {
		return t
	}

	if m.parent != nil {
		return m.parent.GetTexture(name)
	}

	return nil
}

// BlendType for material
// TODO: Fix this blending stuff with src and dst for Func
type BlendType uint32

const (
	// BlendOneOneMinusSrcAlpha - gl.ONE, gl.ONE_MINUS_SRC_ALPHA
	BlendOneOneMinusSrcAlpha = BlendType(iota)
	// BlendOneOne - gl.ONE, gl.ONE
	BlendOneOne
)

// DepthMode handle depth R&W types on render
type DepthMode uint32

// Default deph modes
const (
	DepthReadWrite = DepthMode(iota)
	DepthRead
	DepthNone
)

// StencilFunc stencil function params.
// "only describes whether OpenGL should pass or discard fragments based on the
// stencil buffer's content, not how we can actually update the buffer."
type StencilFunc struct {
	Func gl.Enum
	Ref  int
	Mask uint32
}

// StencilOp sets the stencil operation for the material
// contains three options of which we can specify for each option what action to take:
// Fail: action to take if the stencil test fails.
// ZFail: action to take if the stencil test passes, but the depth test fails.
// ZPass: action to take if both the stencil and the depth test pass.
type StencilOp struct {
	Fail  gl.Enum
	ZFail gl.Enum
	ZPass gl.Enum
}
