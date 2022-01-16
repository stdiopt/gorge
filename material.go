package gorge

import (
	"fmt"
)

// Material the material
type Material struct {
	Resourcer ShaderResourcer
	Name      string
	// Primitive stuff
	Queue       int
	Depth       DepthMode
	DoubleSided bool
	Blend       BlendType

	DisableShadow bool

	/* New: stencil experiment */
	// Create stencil groups?!
	Stencil *Stencil

	// Extra stuff for enable and disable fragment rendering mask.
	ColorMask *[4]bool

	shaderProps
}

// NewMaterial returns a material using the default gorge shader
func NewMaterial() *Material {
	return &Material{}
}

// Material implements the materialer interface.
func (m *Material) Material() *Material { return m }

// Resource returns the current shader resource.
func (m *Material) Resource() ShaderResource {
	if m.Resourcer == nil {
		return nil
	}
	return m.Resourcer.Resource()
}

// NewShaderMaterial returns a new material based on shader data
// if ShaderData is nil it will use the default PBR material
func NewShaderMaterial(r ShaderResourcer) *Material {
	return &Material{Resourcer: r}
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

// SetBlend sets the blend type for material.
func (m *Material) SetBlend(v BlendType) {
	m.Blend = v
}

func (m *Material) SetDisableShadow(v bool) {
	m.DisableShadow = v
}

// SetStencil sets the stencil property for material.
func (m *Material) SetStencil(s *Stencil) {
	m.Stencil = s
}

func (m *Material) SetColorMask(r, g, b, a bool) {
	m.ColorMask = &[4]bool{r, g, b, a}
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
func (m *Material) Get(name string) any {
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

// DepthMode handle depth R&W types on render
type DepthMode uint32

// Default deph modes
const (
	DepthReadWrite = DepthMode(iota)
	DepthRead
	DepthNone
)
