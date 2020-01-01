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
	// TGA image support
	//_ "github.com/ftrvxmtrx/tga"
)

// Image interface
type Image interface {
	Image() (image.Image, error)
}

// ImageFunc returns an image
type ImageFunc func() (image.Image, error)

// Image to satisfy interface
func (f ImageFunc) Image() (image.Image, error) {
	return f()
}

// Image returns an image loader from the manager
func (m *Manager) Image(name string) Image {
	return ImageFunc(func() (image.Image, error) {
		rd, err := m.Loader.Open(name)
		if err != nil {
			return nil, fmt.Errorf("error getting image: %w", err)
		}

		img, _, err := image.Decode(rd)
		if err != nil {
			return nil, fmt.Errorf("error deconding image: %w", err)
		}
		return img, nil
	})
}
