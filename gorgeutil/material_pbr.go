package gorgeutil

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/static"
)

type material = gorge.Material

// MaterialType type of the material (within shader)
type MaterialType int

func (m MaterialType) define() string {
	switch m {
	case MaterialMetallicRoughness:
		return "MATERIAL_METALLICROUGHNESS"
	case MaterialUnlit:
		return "MATERIAL_UNLIT"
	default:
		return ""
	}
}

// Material types
const (
	MaterialMetallicRoughness = iota
	MaterialUnlit
)

// PBRMaterial returns a material controller based on gltf default shader.
type PBRMaterial struct {
	matType MaterialType
	material
}

// NewUnlitMaterial returns a pbr material with unlit defined.
func NewUnlitMaterial() *PBRMaterial {
	mat := gorge.NewShaderMaterial(static.Shaders.DefaultNew)
	mat.Define("MATERIAL_UNLIT", "USE_HDR")
	mat.Set("u_BaseColorFactor", m32.Vec4{1, 1, 1, 1})
	return &PBRMaterial{
		matType:  MaterialUnlit,
		material: *mat,
	}
}

// NewPBRMaterial returns a new PBRMaterial with MetallicRoughness type.
func NewPBRMaterial() *PBRMaterial {
	mat := gorge.NewShaderMaterial(static.Shaders.DefaultNew)
	// Default
	mat.Define("MATERIAL_METALLICROUGHNESS", "USE_HDR")

	// default to 5 but might depent.
	mat.Set("u_MipCount", 5)

	mat.Set("u_BaseColorFactor", m32.Vec4{1, 1, 1, 1})
	mat.Set("u_Exposure", float32(1))
	mat.Set("u_MetallicFactor", float32(1))
	mat.Set("u_RoughnessFactor", float32(1))
	mat.Set("u_NormalScale", float32(1))
	mat.Set("u_OcclusionStrength", float32(1))
	return &PBRMaterial{
		matType:  MaterialMetallicRoughness,
		material: *mat,
	}
}

// Material implements gorge materialer
func (m *PBRMaterial) Material() *gorge.Material {
	return &m.material
}

// SetType sets the material mode.
func (m *PBRMaterial) SetType(t MaterialType) {
	if m.matType == t {
		return
	}
	m.Undefine(m.matType.define())
	m.matType = t
	m.Define(t.define())
}

// SetIBL sets image based lighting.
func (m *PBRMaterial) SetIBL(b bool) {
	if b {
		m.Define("USE_IBL")
		return
	}
	m.Undefine("USE_IBL")
}

// SetHDR enables hdr lighting on the shader.
func (m *PBRMaterial) SetHDR(b bool) {
	if b {
		m.Define("USE_HDR")
		return
	}
	m.Undefine("USE_HDR")
}

// SetBaseColor sets the base color factor.
func (m *PBRMaterial) SetBaseColor(v m32.Vec4) {
	m.Set("u_BaseColorFactor", v)
}

// SetBaseColorMap sets the base texture.
func (m *PBRMaterial) SetBaseColorMap(tex *gorge.Texture) {
	if tex == nil {
		m.Undefine("HAS_BASE_COLOR_MAP")
		return
	}

	m.Define("HAS_BASE_COLOR_MAP")
	m.SetTexture("u_BaseColorSampler", tex)
}

// SetMetallicFactor sets the metallic factor.
func (m *PBRMaterial) SetMetallicFactor(v float32) {
	m.Set("u_MetallicFactor", v)
}

// SetRoughnessFactor sets the metallic factor.
func (m *PBRMaterial) SetRoughnessFactor(v float32) {
	m.Set("u_RoughnessFactor", v)
}

// SetMetallicRoughnessMap sets the combined texture of metallic and roughness.
func (m *PBRMaterial) SetMetallicRoughnessMap(tex *gorge.Texture) {
	m.SetTexture("u_MetallicRoughnessSampler", tex)
	if tex == nil {
		m.Undefine("HAS_METALLIC_ROUGHNESS_MAP")
		return
	}
	m.Define("HAS_METALLIC_ROUGHNESS_MAP")
}

// SetMetallicMap separate metallic sampler.
func (m *PBRMaterial) SetMetallicMap(tex *gorge.Texture) {
	m.SetTexture("u_MetallicSampler", tex)
	if tex == nil {
		m.Undefine("HAS_METALLIC_MAP")
		return
	}
	m.Define("HAS_METALLIC_MAP")
}

// SetRoughnessMap separate roughness sampler.
func (m *PBRMaterial) SetRoughnessMap(tex *gorge.Texture) {
	m.SetTexture("u_RoughnessSampler", tex)
	if tex == nil {
		m.Undefine("HAS_ROUGHNESS_MAP")
		return
	}
	m.Define("HAS_ROUGHNESS_MAP")
}

// SetNormalMap sets the normal map.
func (m *PBRMaterial) SetNormalMap(tex *gorge.Texture) {
	m.SetTexture("u_NormalSampler", tex)
	if tex == nil {
		m.Undefine("HAS_NORMAL_MAP")
		return
	}
	m.Define("HAS_NORMAL_MAP")
}

// SetNormalScale sets the normal scale.
func (m *PBRMaterial) SetNormalScale(v float32) {
	m.Set("u_NormalScale", v)
}

// SetOcclusionMap sets the normal map.
func (m *PBRMaterial) SetOcclusionMap(tex *gorge.Texture) {
	m.SetTexture("u_OcclusionSampler", tex)
	if tex == nil {
		m.Undefine("HAS_OCCLUSION_MAP")
		return
	}
	m.Define("HAS_OCCLUSION_MAP")
}

// SetOcclusionStrength sets the occlusion strength.
func (m *PBRMaterial) SetOcclusionStrength(v float32) {
	m.Set("u_OcclusionStrength", v)
}

// SetBaseUVTransform  sets the UV transform matrix.
func (m *PBRMaterial) SetBaseUVTransform(v m32.Mat3) {
	if v == m32.M3Ident() {
		m.Undefine("HAS_BASECOLOR_UV_TRANSFORM")
		m.Set("u_BaseColorUVTransform", nil)
		return
	}
	m.Define("HAS_BASECOLOR_UV_TRANSFORM")
	m.Set("u_BaseColorUVTransform", v)
}
