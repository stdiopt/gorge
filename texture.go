// Copyright 2019 Luis Figueiredo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gorge

import (
	"github.com/stdiopt/gorge/gl"
)

// TODO: Remove gl entries

// TextureFormat texture pixel format
type TextureFormat int

// WrapMode for texture
type WrapMode int

// FilterMode texture filter mode
type FilterMode int

// Wrapmode consts
const (
	TextureWrapClamp  = WrapMode(gl.CLAMP_TO_EDGE)
	TextureWrapRepeat = WrapMode(gl.REPEAT)
	TextureWrapMirror = WrapMode(gl.MIRRORED_REPEAT)

	TextureFilterPoint  = FilterMode(gl.NEAREST)
	TextureFilterLinear = FilterMode(gl.LINEAR)

	TextureFormatRGBA = TextureFormat(gl.RGBA)
)

// TextureLoader interface for a texture loader
type TextureLoader interface {
	Data() *TextureData
}

// TextureData representes the data for the texture
type TextureData struct {
	Source        string
	Format        TextureFormat
	Width, Height int
	PixelData     []byte
}

// Data convenient func
func (d TextureData) Data() *TextureData { return &d }

// Texture reference
type Texture struct {
	asset
	Name       string // just for reference and debugging
	WrapU      WrapMode
	WrapV      WrapMode
	WrapW      WrapMode
	FilterMode FilterMode

	Updates     int
	DataUpdates int

	// Loader Should be private so we can not change?
	// Swappable texture loader
	loader TextureLoader
}

// NewTexture returns a new texture with loader
func NewTexture(loader TextureLoader) *Texture {
	return &Texture{loader: loader}
}

// Loader gets the texture Loader
func (t *Texture) Loader() TextureLoader {
	return t.loader
}

// SetFilterMode sets the filter mode POINT,LINEAR
func (t *Texture) SetFilterMode(f FilterMode) *Texture {
	t.FilterMode = f
	t.Updates++
	return t
}

// SetWrapU texture wrap for U
func (t *Texture) SetWrapU(w WrapMode) *Texture {
	t.WrapU = w
	t.Updates++
	return t
}

// SetWrapV texture wrap for V
func (t *Texture) SetWrapV(w WrapMode) *Texture {
	t.WrapV = w
	t.Updates++
	return t
}

// SetWrapW texture wrap for W
func (t *Texture) SetWrapW(w WrapMode) *Texture {
	t.WrapW = w
	t.Updates++
	return t
}
