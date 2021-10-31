// Package gltf attempt to implement glft
package gltf

import (
	"fmt"

	"github.com/stdiopt/gorge"
)

// Doc Document
type Doc struct {
	Asset       Asset         `json:"asset"`
	Scene       int           `json:"scene"`
	Scenes      []*Scene      `json:"scenes"`
	Nodes       []*Node       `json:"nodes"`
	Cameras     []*Camera     `json:"cameras"`
	Meshes      []*Mesh       `json:"meshes,omitempty"`
	Materials   []*Material   `json:"materials,omitempty"`
	Textures    []*Texture    `json:"textures"`
	Samplers    []*Sampler    `json:"samplers"`
	Accessors   []*Accessor   `json:"accessors"`
	BufferViews []*BufferView `json:"bufferViews"`
	Buffers     []*Buffer     `json:"buffers"`
	Images      []*Image      `json:"images"`
	Animations  []*Animation  `json:"animations"`
	Skins       []*Skin       `json:"skins"`

	BasePath string `json:"-"`
}

// Asset gltf asset data.
type Asset struct {
	Copyright string `json:"copyright,omitempty"`
	Generator string `json:"generator,omitempty"`
	Version   string `json:"version"`
}

// Scene gltf scene data.
type Scene struct {
	Nodes []int `json:"nodes"`
}

// Node gltf node data.
type Node struct {
	Name        *string      `json:"name,omitempty"`
	Children    []int        `json:"children"`
	Matrix      *[16]float32 `json:"matrix,omitempty"`
	Rotation    *[4]float32  `json:"rotation,omitempty"`
	Translation *[3]float32  `json:"translation,omitempty"`
	Scale       *[3]float32  `json:"scale,omitempty"`

	Camera *int `json:"camera"`
	Mesh   *int `json:"mesh"`
	Skin   *int `json:"skin"`
}

// Camera gltf camera data.
type Camera struct {
	Type        string            `json:"type"`
	Perspective CameraPerspective `json:"perspective"`
}

// CameraPerspective gltf camera perspective data.
type CameraPerspective struct {
	AspectRatio float32 `json:"aspectRatio"`
	Yfov        float32 `json:"yfov"`
	Zfar        float32 `json:"zfar"`
	Znear       float32 `json:"znear"`
}

// Mesh gltf mesh
// TODO: Add missing stuff
type Mesh struct {
	Name       string           `json:"name,omitempty"`
	Primitives []*MeshPrimitive `json:"primitives"`
	Weights    []float32        `json:"weights,omitempty"`
}

// MeshPrimitive gltf primitive data.
type MeshPrimitive struct {
	Attributes map[string]int   `json:"attributes"`
	Targets    []map[string]int `json:"targets"` // this is an array of ATTRS
	Indices    *int             `json:"indices"`
	Material   *int             `json:"material"`
	Mode       int              `json:"mode"`
}

// Material stuff
// TODO add more stuff
type Material struct {
	Name                 string                `json:"name"`
	AlphaMode            *string               `json:"alphaMode,omitempty"`
	DoubleSided          *bool                 `json:"doubleSided,omitempty"`
	AlphaCutoff          *float32              `json:"alphaCutoff,omitempty"`
	EmissiveFactor       *[3]float32           `json:"emissiveFactor,omitempty"`
	PBRMetallicRoughness *MatMetallicRoughness `json:"pbrMetallicRoughness,omitempty"`
	EmissiveTexture      *TextureInfo          `json:"emissiveTexture,omitempty"`
	NormalTexture        *NormalTextureInfo    `json:"normalTexture,omitempty"`
	OcclusionTexture     *OcclusionTextureInfo `json:"occlusionTexture,omitempty"`
	Extensions           *MaterialExt          `json:"extensions"`
}

// MaterialExt This contains the implemented extensions
type MaterialExt struct {
	Clearcoat *MatExtClearcoat `json:"KHR_materials_clearcoat,omitempty"`
	Unlit     *struct{}        `json:"KHR_materials_unlit,omitempty"`
}

// TextureInfo gltf ata for textureinfo.
type TextureInfo struct {
	Index    int `json:"index"`
	TexCoord int `json:"texCoord"`
}

// MatMetallicRoughness stuff
// XXX: is this specific for this material or every other material might apply this
// TODO: Needs more fields from spec
type MatMetallicRoughness struct {
	BaseColorFactor          *[4]float32  `json:"baseColorFactor,omitempty"`
	BaseColorTexture         *TextureInfo `json:"baseColorTexture,omitempty"`
	MetallicFactor           *float32     `json:"metallicFactor,omitempty"`
	RoughnessFactor          *float32     `json:"roughnessFactor,omitempty"`
	MetallicRoughnessTexture *TextureInfo `json:"metallicRoughnessTexture,omitempty"`
}

// OcclusionTextureInfo gltf texture info.
type OcclusionTextureInfo struct {
	TextureInfo
	Strength *float32 `json:"strength,omitempty"` // Default 1
}

// NormalTextureInfo gltf textureInfo.
type NormalTextureInfo struct {
	TextureInfo
	Scale *float32 `json:"scale"`
}

// MatExtClearcoat gltf material clearcoat extension.
type MatExtClearcoat struct {
	ClearcoatFactor           *float32           `json:"clearcoatFactor,omitempty"`
	ClearcoatTexture          *TextureInfo       `json:"clearcoatTexture,omitempty"`
	ClearcoatRoughnessFactor  *float32           `json:"clearcoatRoughnessFactor,omitempty"`
	ClearcoatRoughnessTexture *TextureInfo       `json:"clearcoatRoughnessTexture,omitempty"`
	ClearcoatNormalTexture    *NormalTextureInfo `json:"clearcoatNormalTexture,omitempty"`
}

// Sampler gltf sampler data.
type Sampler struct {
	MinFilter *SamplerFilter `json:"minFilter"`
	MagFilter *SamplerFilter `json:"magFilter"`
	WrapS     *SamplerWrap   `json:"WrapS"`
	WrapT     *SamplerWrap   `json:"WrapT"`
}

// Texture gltf texture data.
type Texture struct {
	Sampler int `json:"sampler"`
	Source  int `json:"source"`
}

// Accessor gltf accessor data.
type Accessor struct {
	BufferView    int           `json:"bufferView"`
	ByteOffset    int           `json:"byteOffset"`
	ComponentType ComponentType `json:"componentType"`
	Count         int           `json:"count"`
	Max           []float32     `json:"max"`
	Min           []float32     `json:"min"`
	Type          AccessorType  `json:"type"`
}

// BufferView gltf bufferView data.
type BufferView struct {
	Buffer     int `json:"buffer"`
	ByteLength int `json:"byteLength"`
	ByteOffset int `json:"byteOffset"`
	ByteStride int `json:"byteStride"`
	Target     int `json:"target"`
}

// Buffer buffer data.
type Buffer struct {
	ByteLength int    `json:"byteLength"`
	URI        string `json:"uri,omitempty"`

	// Shouldn't be here but return on demand?
	RawData []byte
}

// Image object
// required one of URI or BufferView
type Image struct {
	Name       string `json:"name"`
	URI        string `json:"uri,omitempty"`
	MimeType   string `json:"mimeType,omitempty"`
	BufferView *int   `json:"bufferView,omitempty"`

	TexData *gorge.TextureData
}

// Animation gltf data struct.
type Animation struct {
	Name     string              `json:"name"`
	Channels []*AnimationChannel `json:"channels"`
	Samplers []*AnimationSampler `json:"samplers"`
}

// AnimationChannel gltf data struct.
type AnimationChannel struct {
	Sampler int                    `json:"sampler"`
	Target  AnimationChannelTarget `json:"target"`
}

// AnimationChannelTarget gltf data struct.
type AnimationChannelTarget struct {
	Node int    `json:"node"`
	Path string `json:"path"`
}

// AnimationSampler gltf data struct.
type AnimationSampler struct {
	Input         int    `json:"input"`
	Interpolation string `json:"interpolation"`
	Output        int    `json:"output"`
}

// Skin gltf data struct.
type Skin struct {
	Name                *string `json:"name,omitempty"`
	InverseBindMatrices *int    `json:"inverseBindMatrices,omitempty"`
	Joints              []int   `json:"joints"`
	Skeleton            *int    `json:"skeleton"`
}

// TODO: We copy Accessor to a byte slice which is converted later to a
// specific type

// AccessorBuffer returns a copy of the buffer solving all the strides and
// offsets
// TODO: Sparse accessor
func (g *Doc) AccessorBuffer(i int) ([]byte, int, ComponentType) {
	ret := []byte{}
	accessor := g.Accessors[i]
	bv := g.BufferViews[accessor.BufferView]
	buf := g.Buffers[bv.Buffer].RawData
	aBuf := buf[bv.ByteOffset:][:bv.ByteLength][accessor.ByteOffset:]

	sz := accessor.ComponentType.ByteLen()
	sz *= accessor.Type.UnitLength()
	for count := 0; count < accessor.Count; count++ {
		ret = append(ret, aBuf[:sz]...)
		if bv.ByteStride != 0 {
			aBuf = aBuf[bv.ByteStride:]
		} else {
			aBuf = aBuf[sz:]
		}
	}
	return ret, sz, accessor.ComponentType
}

// AccessorType represents the underlying type.
type AccessorType string

// UnitLength returns the len in components for the accessor type.
func (t AccessorType) UnitLength() int {
	switch t {
	case "SCALAR":
		return 1
	case "VEC2":
		return 2
	case "VEC3":
		return 3
	case "VEC4":
		return 4
	case "MAT4":
		return 16
	default:
		panic(fmt.Errorf("%q not implemented yet", t))
	}
}

// Accessor types.
const (
	AccessorScalar = AccessorType("SCALAR")
	AccessorVec2   = AccessorType("VEC2")
	AccessorVec3   = AccessorType("VEC3")
	AccessorVec4   = AccessorType("VEC4")
	AccessorMat4   = AccessorType("MAT4")
)

// ComponentType gltf unit data type
type ComponentType int

// Component types
const (
	ComponentByte   = 5120
	ComponentUByte  = 5121
	ComponentShort  = 5122
	ComponentUShort = 5123
	ComponentUInt   = 5125
	ComponentFloat  = 5126
)

func (t ComponentType) String() string {
	switch t {
	case ComponentByte:
		return "BYTE"
	case ComponentUByte: // gl.BYTE, gl.UNSIGNED_BYTE
		return "UNSIGNED_BYTE"
	case ComponentShort:
		return "SHORT"
	case ComponentUShort: // gl.SHORT, gl.UNSIGNED_SHORT
		return "UNSIGNED_SHORT"
	case ComponentUInt:
		return "UNSIGNED_INT"
	case ComponentFloat: // gl.UNSIGNED_INT, gl.FLOAT
		return "FLOAT"
	}
	return "<unknown>"
}

// ByteLen returns the byte len for component type.
func (t ComponentType) ByteLen() int {
	switch t {
	case ComponentByte, ComponentUByte: // gl.BYTE, gl.UNSIGNED_BYTE
		return 1
	case ComponentShort, ComponentUShort: // gl.SHORT, gl.UNSIGNED_SHORT
		return 2
	case ComponentUInt, ComponentFloat: // gl.UNSIGNED_INT, gl.FLOAT
		return 4
	}
	return 0
}

// default texture filtering values.
const (
	SamplerNearest              = SamplerFilter(9728)
	SamplerLinear               = SamplerFilter(9729)
	SamplerNearestMipmapNearest = SamplerFilter(9984)
	SamplerLinearMipmapNearest  = SamplerFilter(9985)
	SamplerNearestMipmapLinear  = SamplerFilter(9986)
	SamplerLinearMipmapLinear   = SamplerFilter(9987)
)

// SamplerFilter represents gltf texture filtering.
type SamplerFilter int

func (f SamplerFilter) String() string {
	switch f {
	case SamplerNearest:
		return "NEAREST"
	case SamplerLinear:
		return "LINEAR"
	case SamplerNearestMipmapNearest:
		return "NEAREST_MIPMAP_NEAREST"
	case SamplerLinearMipmapNearest:
		return "LINEAR_MIPMAP_NEAREST"
	case SamplerNearestMipmapLinear:
		return "NEAREST_MIPMAP_LINEAR"
	case SamplerLinearMipmapLinear:
		return "LINEAR_MIPMAP_LINEAR"
	default:
		return "<sampler filter:unknown>"
	}
}

// SamplerWrap texture wrapping.
type SamplerWrap int

// Default texture wrapping values.
const (
	SamplerClamp          = SamplerWrap(33071)
	SamplerMirroredRepeat = SamplerWrap(33648)
	SamplerRepeat         = SamplerWrap(10497)
)

func (w SamplerWrap) String() string {
	switch w {
	case SamplerClamp:
		return "CLAMP_TO_EDGE"
	case SamplerMirroredRepeat:
		return "MIRRORED_REPEAT"
	case SamplerRepeat:
		return "REPEAT"
	default:
		return fmt.Sprintf("<sampler wrap:unknown:%d>", w)
	}
}
