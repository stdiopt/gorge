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

package resource

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/stdiopt/gorge"
)

type textureLoader struct {
	name   string
	loader func() *gorge.TextureData
}

func (l textureLoader) Data() *gorge.TextureData {
	return l.loader()
}

// Texture2D return a texture2D
/*func (a *Asset) Texture2D(name string) *gorge.Texture {
	return gorge.NewTexture(a.texture2DLoader(name)).
		SetFilterMode(gorge.TextureFilterLinear)
}*/

type textureLoaderKey struct {
	name string
}

// Texture2DLoader prepares a TextureLoader for lazy loading
func (m *Manager) texture2DLoader(name string) *textureLoader {
	if m.loaders == nil {
		m.loaders = map[interface{}]interface{}{}
	}

	k := textureLoaderKey{name}
	if loader, ok := m.loaders[k]; ok {
		return loader.(*textureLoader)
	}

	loaderFunc := func() *gorge.TextureData {
		img, err := m.LoadImage(name)
		if err != nil {
			m.Error(err)
			return nil
		}

		width := img.Bounds().Dx()
		height := img.Bounds().Dy()
		var data []byte
		switch img := img.(type) {
		case *image.NRGBA:
			data = img.Pix
		case *image.RGBA:
			data = img.Pix
		case *image.YCbCr:
			b := img.Bounds()
			m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
			draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
			data = m.Pix
		default:
			panic(fmt.Sprintf("unsupported image: %T", img))
		}
		return &gorge.TextureData{
			Source:    name,
			Format:    gorge.TextureFormatRGBA,
			Width:     width,
			Height:    height,
			PixelData: data,
		}
	}
	loader := &textureLoader{
		name:   name,
		loader: loaderFunc,
	}
	m.loaders[k] = loader

	return loader
}
