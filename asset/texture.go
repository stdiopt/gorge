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

package asset

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/stdiopt/gorge"
)

type textureLoader func() *gorge.TextureData

func (l textureLoader) Data() *gorge.TextureData { return l() }

// Texture2D prepares a texture2D for loading
func (s *System) Texture2D(name string) *gorge.Texture {

	loader := textureLoader(func() *gorge.TextureData {
		img, err := s.LoadImage(name)
		if err != nil {
			s.Error(err)
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
	})

	t := &gorge.Texture{
		TextureLoader: loader,
		FilterMode:    gorge.TextureFilterLinear,
	}

	s.gorge.Trigger(AddEvent{Asset: t})
	return t
}
