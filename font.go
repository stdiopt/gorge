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
	"github.com/stdiopt/gorge/m32"
)

// Glyph char information
type Glyph struct {
	Uv1  vec2
	Uv2  vec2
	Size vec2

	Advance  float32
	BearingH float32
	BearingV float32
}

// Font is like a texture with extra information about glyphs
type Font struct {
	Glyphs   map[rune]Glyph
	SpaceAdv float32
	*Texture
}

// FontOptions options for font
type FontOptions struct {
	Resolution int
	Chars      string
	Background *m32.Vec4
	Foreground *m32.Vec4
}

// FontOptionsFunc func to manipulate font options
type FontOptionsFunc func(o *FontOptions)

// FontResolution option
func FontResolution(n int) FontOptionsFunc {
	return func(opt *FontOptions) {
		opt.Resolution = n
	}
}
