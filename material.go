package gorge

import (
	"fmt"
)

// Material the material
type Material struct {
	resourcer
	Name string
	// Primitive stuff
	Queue       int
	Depth       DepthMode
	DoubleSided bool
	Blend       BlendType

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
func NewShaderMaterial(r Resourcer) *Material {
	return &Material{
		resourcer: resourcer{r},
	}
}

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

// Defines override shaderProp defines with hierarchy
func (m *Material) Defines() map[string]string {
	pm, ok := m.Resourcer.(*Material)
	if !ok {
		return m.shaderProps.Defines()
	}
	ret := map[string]string{}
	for k, v := range pm.Defines() {
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
	if p, ok := m.Resourcer.(*Material); ok {
		hash ^= p.DefinesHash()
	}
	hash ^= m.shaderProps.DefinesHash()
	return hash
}

// Get returns the material property for name or if the material has a Parent
// material returns the property from parent else returns nil.
func (m *Material) Get(name string) interface{} {
	pm, ok := m.Resourcer.(*Material)
	if !ok {
		return m.shaderProps.Get(name)
	}
	if v := m.shaderProps.Get(name); v != nil {
		return v
	}
	return pm.Get(name)
}

// GetTexture returns the texture for name
func (m *Material) GetTexture(name string) *Texture {
	if t := m.shaderProps.GetTexture(name); t != nil {
		return t
	}

	if pm, ok := m.Resourcer.(*Material); ok {
		return pm.GetTexture(name)
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
