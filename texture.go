package gorge

import (
	"fmt"

	"github.com/stdiopt/gorge/core/event"
)

// Texturer used to fetch a texture from a texture controller.
type Texturer interface {
	Texture() *Texture
}

type TextureResourcer interface {
	Resource() TextureResource
}

// TextureResource is an interface to handle underlying texture data.
type TextureResource interface {
	isTexture()
	isGPU()
}

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

func (t *Texture) Resource() TextureResource {
	if t.Resourcer == nil {
		return nil
	}
	return t.Resourcer.Resource()
}

// Texture implements Texturer
func (t *Texture) Texture() *Texture { return t }

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
	curRes := t.Resource()
	if _, ok := curRes.(*TextureData); !ok {
		return
	}
	event.Trigger(g, EventResourceUpdate{Resource: curRes})

	gpuRef := &TextureRef{&GPU{}}
	SetGPU(gpuRef, GetGPU(curRes))

	t.Resourcer = gpuRef
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
	GPU

	Source        string
	Format        TextureFormat
	Width, Height int
	PixelData     []byte
	Updates       int
}

// Resource implements TextureResourcer.
func (d *TextureData) Resource() TextureResource { return d }

// CreateRef creates a texture gpu reference.
func (d *TextureData) CreateRef(g *Context) *TextureRef {
	ref := &TextureRef{&GPU{}}
	event.Trigger(g, EventResourceUpdate{Resource: d})
	SetGPU(ref, GetGPU(d))
	return ref
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

//////////////////////////////////////////////////////////////////////////////
// Experiment, single pixel color texture
// /////
type texture = Texture

// ColorTexture helper for a single color texture.
type ColorTexture struct {
	texture
}

// NewColorTexture returns a single pixel colored texture
func NewColorTexture(r, g, b, a float32) *ColorTexture {
	t := &ColorTexture{}
	t.SetColor(r, g, b, a)
	return t
}

// SetColor sets color data for underlying texture.
func (t *ColorTexture) SetColor(r, g, b, a float32) {
	tex, ok := t.Resource().(*TextureData)
	if !ok {
		tex = &TextureData{}
		t.Resourcer = tex
	}
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
	texture
}

// NewValueTexture returns a single pixel colored texture
func NewValueTexture(v float32) *ValueTexture {
	t := &ValueTexture{}
	t.SetValue(v)
	return t
}

// SetValue sets the texture Value
func (t *ValueTexture) SetValue(v float32) {
	tex, ok := t.Resource().(*TextureData)
	if !ok {
		tex = &TextureData{}
		t.Resourcer = tex
	}
	if len(tex.PixelData) == 0 {
		tex.Format = TextureFormatGray
		tex.Width = 1
		tex.Height = 1
		tex.PixelData = make([]byte, 1)
	}
	tex.PixelData[0] = byte(v * 255)
	tex.Updates++
}
