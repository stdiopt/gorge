package gorge

import (
	"fmt"
)

// TextureResourcer is an interface to handle underlying texture data.
type TextureResourcer interface {
	Resource() ResourceRef
	isTexture()
}

// TextureRef implements a texture resourcer.
type TextureRef struct {
	Ref ResourceRef
}

// Resource implements the resourcer interface.
func (r *TextureRef) Resource() ResourceRef { return r.Ref }
func (r *TextureRef) isTexture()            {}

// Texture reference
type Texture struct {
	Resourcer  TextureResourcer
	Name       string // just for reference and debugging
	Wrap       [3]TextureWrap
	FilterMode TextureFilter
}

// NewTexture returns a new texture based on resourcer.
func NewTexture(r TextureResourcer) *Texture {
	return &Texture{Resourcer: r}
}

// GetResource returns the ResourceRef for this texture.
func (t *Texture) GetResource() ResourceRef {
	if t.Resourcer == nil {
		return nil
	}
	return t.Resourcer.Resource()
}

// SetResourcer sets the resourcer for this texture.
func (t *Texture) SetResourcer(r TextureResourcer) { t.Resourcer = r }

// GetWrap get the components of TextureWrap UVW
func (t *Texture) GetWrap() (u, v, w TextureWrap) {
	return t.Wrap[0], t.Wrap[1], t.Wrap[2]
}

// GetFilterMode texture filter mode
func (t *Texture) GetFilterMode() TextureFilter {
	return t.FilterMode
}

// SetFilterMode sets the filter mode POINT,LINEAR
func (t *Texture) SetFilterMode(f TextureFilter) {
	t.FilterMode = f
}

// SetWrapUVW texture wrap for U, V, W
func (t *Texture) SetWrapUVW(uvw ...TextureWrap) {
	switch len(uvw) {
	case 1:
		t.Wrap[0], t.Wrap[1], t.Wrap[2] = uvw[0], uvw[0], uvw[0]
	default:
		copy(t.Wrap[:], uvw)
	}
}

// ReleaseData change underlying resourcer with a gpu only reference.
func (t *Texture) ReleaseData(g *Context) {
	if _, ok := t.Resourcer.(*TextureRef); ok {
		return
	}
	curRes := t.Resourcer.Resource()
	g.Trigger(EventResourceUpdate{Resource: curRes})

	gpuRef := NewGPUResource()
	SetGPU(gpuRef, GetGPU(curRes))

	t.Resourcer = &TextureRef{Ref: gpuRef}
}

// TextureFormat texture pixel format
type TextureFormat int

func (f TextureFormat) String() string {
	switch f {
	case TextureFormatRGBA:
		return "RGBA"
	case TextureFormatRGB:
		return "RGB"
	case TextureFormatGray:
		return "Gray"
	case TextureFormatGray16:
		return "Gray16"
	case TextureFormatRGB32F:
		return "RGB32F"
	default:
		return "Unknown"
	}
}

// Known texture formats
const (
	TextureFormatRGBA = TextureFormat(iota)
	TextureFormatRGB
	TextureFormatGray
	TextureFormatGray16
	TextureFormatRGB32F
)

// TextureWrap for texture
type TextureWrap int

func (m TextureWrap) String() string {
	switch m {
	case TextureWrapRepeat:
		return "TextureWrapRepeat"
	case TextureWrapClamp:
		return "TextureWrapClamp"
	case TextureWrapMirror:
		return "TextureWrapMirror"
	default:
		return "<invalid>"
	}
}

// TextureFilter texture filter mode
type TextureFilter int

// Wrapmode consts
const (
	TextureWrapRepeat = TextureWrap(iota)
	TextureWrapClamp
	TextureWrapMirror
)

// TextureFilter types
const (
	TextureFilterLinear = TextureFilter(iota)
	TextureFilterPoint
)

// TextureData is the data for the texture
type TextureData struct {
	gpuResource

	Source        string
	Format        TextureFormat
	Width, Height int
	PixelData     []byte
	Updates       int
}

func (d *TextureData) String() string {
	return fmt.Sprintf(
		"texture: (source: %q, format: %v  size: %dx%d)",
		d.Source,
		d.Format,
		d.Width, d.Height,
	)
}

func (d *TextureData) isTexture() {}

// Resource implements the resourcer interface.
func (d *TextureData) Resource() ResourceRef { return d }

//////////////////////////////////////////////////////////////////////////////
// Experiment, single pixel color texture
// /////
type texture = Texture

// ColorTexture helper for a single color texture.
type ColorTexture struct {
	TextureData
}

// NewColorTexture returns a single pixel colored texture
func NewColorTexture(r, g, b, a float32) *ColorTexture {
	t := &ColorTexture{}
	t.SetColor(r, g, b, a)
	return t
}

// SetColor sets color data for underlying texture.
func (t *ColorTexture) SetColor(r, g, b, a float32) {
	tex := &t.TextureData
	if len(tex.PixelData) == 0 {
		tex.Format = TextureFormatRGBA
		tex.Width = 1
		tex.Height = 1
		tex.PixelData = make([]byte, 4)
	}
	tex.PixelData[0] = byte(r * 255)
	tex.PixelData[1] = byte(g * 255)
	tex.PixelData[2] = byte(b * 255)
	tex.PixelData[3] = byte(a * 255)
	tex.Updates++
}

// ValueTexture helper for a single valued texture.
type ValueTexture struct {
	TextureData
}

// NewValueTexture returns a single pixel colored texture
func NewValueTexture(v float32) *ValueTexture {
	t := &ValueTexture{}
	t.SetValue(v)
	return t
}

// SetValue sets the texture Value
func (t *ValueTexture) SetValue(v float32) {
	tex := &t.TextureData
	if len(tex.PixelData) == 0 {
		tex.Format = TextureFormatGray
		tex.Width = 1
		tex.Height = 1
		tex.PixelData = make([]byte, 1)
	}
	tex.PixelData[0] = byte(v * 255)
	tex.Updates++
}
