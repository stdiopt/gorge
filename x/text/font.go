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

package text

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"io/ioutil"
	"unicode"

	"github.com/golang/freetype/truetype"
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Might miss some chars
const commonChars = `
0123456789µ
abcdefghijklmnopqrstuvwxyz
ABCDEFGHIJKLMNOPQRSTUVWXYZ
{}()[]|$@?%/\:;,._-+=<>*"'~
`

type (
	vec2 = m32.Vec2
)

// Glyph char information
type Glyph struct {
	uv1  vec2
	uv2  vec2
	size vec2

	advance  float32
	bearingH float32
	bearingV float32
}

// Font contains glyphs and texture
type Font struct {
	glyphs   map[rune]Glyph
	spaceAdv float32
	Texture  *gorge.Texture
}

// FontOptions options for font
type FontOptions struct {
	Resolution int
	Chars      string
	Background *m32.Vec4
	Foreground *m32.Vec4
}

type fontOpener interface {
	Open() (io.ReadCloser, error)
}

// NewFont returns a new font
func NewFont(f fontOpener) (*Font, error) {
	return FontWithOptions(f, FontOptions{})
}

// NewFontWithOptions get a font
func FontWithOptions(f fontOpener, opt FontOptions) (*Font, error) {
	bg := color.Color(color.RGBA{0, 0, 0, 0})
	fg := color.Color(color.RGBA{255, 255, 255, 255})
	if opt.Background != nil {
		bg = vec4ToColor(*opt.Background)
	}
	if opt.Foreground != nil {
		fg = vec4ToColor(*opt.Foreground)
	}

	res := 1024
	if opt.Resolution != 0 {
		res = opt.Resolution
	}

	chars := opt.Chars
	if chars == "" {
		chars = commonChars
	}
	chars = "�" + chars // prepend the interrogation

	b := &fontBuilder{
		res:   res,
		bg:    bg,
		fg:    fg,
		chars: chars,
		font:  f,
	}
	return b.build()
}

// Maintain the state
type fontBuilder struct {
	res   int
	chars string
	bg    color.Color
	fg    color.Color
	font  fontOpener
}

// Two similar steps
// - the Font info
// - the texture generator
func (b fontBuilder) image() (image.Image, error) {
	// Maybe this
	f, err := b.font.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fontData, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	ttf, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	// Ideal font size
	count := glyphCount(b.chars)
	scale := calcSize(float32(b.res), count)
	fontScale := float32(0.8)

	// Font info

	// Create our image (texture)
	clip := image.Rect(0, 0, b.res, b.res)
	img := image.NewRGBA(clip)
	draw.Draw(img, img.Bounds(), image.NewUniform(b.bg), image.ZP, draw.Src)

	ttfFace := truetype.NewFace(ttf, &truetype.Options{
		Size:    float64(scale * fontScale), // a bit smaller because of weird glyphs
		DPI:     72,
		Hinting: font.HintingFull,
	})
	defer ttfFace.Close()

	dot := fixed.P(0, int(scale))
	ladv := fixed.I(int(scale))
	// Each char
	for _, ch := range b.chars {
		if unicode.IsSpace(ch) {
			continue
		}
		bnd, _, ok := ttfFace.GlyphBounds(ch)
		if !ok {
			return nil, fmt.Errorf("glyph bounds error")
		}
		if (dot.X + ladv) > fixed.I(clip.Dx()) {
			dot.X = 0
			dot.Y += fixed.I(int(scale))
		}
		// Center stuff somehow?
		d := dot.Sub(bnd.Min)
		gr, mask, maskp, _, ok := ttfFace.Glyph(d, ch)
		if !ok {
			return nil, fmt.Errorf("glyph error")
		}
		//gr = clip.Intersect(gr)
		draw.DrawMask(img, gr, image.NewUniform(b.fg), image.ZP, mask, maskp, draw.Over)
		dot.X += ladv
	}

	return img, nil
}
func (b fontBuilder) build() (*Font, error) {
	f, err := b.font.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fontData, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	ttf, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	// Ideal font size
	count := glyphCount(b.chars)
	scale := calcSize(float32(b.res), count)

	bnd := ttf.Bounds(fixed.I(int(scale)))
	boundW := (bnd.Max.X - bnd.Min.X).Floor()
	boundH := (bnd.Max.Y - bnd.Min.Y).Floor()
	fontScale := float32(0.8)

	clip := image.Rect(0, 0, b.res, b.res)
	iw := float32(clip.Bounds().Dx())
	ih := float32(clip.Bounds().Dy())

	ttfFace := truetype.NewFace(ttf, &truetype.Options{
		Size:    float64(scale * fontScale), // a bit smaller because of weird glyphs
		DPI:     72,
		Hinting: font.HintingFull,
	})
	defer ttfFace.Close()

	dot := fixed.P(0, int(scale))
	glyphs := map[rune]Glyph{}
	ladv := fixed.I(int(scale))
	// Each char
	for _, ch := range b.chars {
		if unicode.IsSpace(ch) {
			continue
		}
		bnd, adv, ok := ttfFace.GlyphBounds(ch)
		if !ok {
			return nil, fmt.Errorf("glyph bounds error")
		}
		if (dot.X + ladv) > fixed.I(clip.Dx()) {
			dot.X = 0
			dot.Y += fixed.I(int(scale))
		}
		gw := int((bnd.Max.X - bnd.Min.X) >> 6)
		gh := int((bnd.Max.Y - bnd.Min.Y) >> 6)

		if gw == 0 || gh == 0 {
			gw = boundW
			gh = boundH
			//above can sometimes yield 0 for font smaller than 48pt, 1 is minimum
			if gw == 0 || gh == 0 {
				gw = 1
				gh = 1
			}
		}

		// Center stuff somehow?
		d := dot.Sub(bnd.Min)
		gr, _, _, _, ok := ttfFace.Glyph(d, ch)
		if !ok {
			return nil, fmt.Errorf("glyph error")
		}

		p1x := float32(gr.Min.X)
		p1y := float32(gr.Min.Y)
		p2x := float32(gr.Max.X)
		p2y := float32(gr.Max.Y)
		glyphs[ch] = Glyph{
			// uv.V inverted
			uv1:      vec2{p1x / iw, p1y / ih},
			uv2:      vec2{p2x / iw, p2y / ih},
			size:     vec2{float32(gw) / scale, float32(gh) / scale},
			advance:  float32(adv>>6) / scale,
			bearingV: float32(bnd.Max.Y>>6) / scale,
			bearingH: float32(bnd.Min.X>>6) / scale,
		}

		dot.X += ladv
	}

	adv, _ := ttfFace.GlyphAdvance(' ')

	font := &Font{
		spaceAdv: float32(adv>>6) / scale,
		glyphs:   glyphs,
		Texture: &gorge.Texture{
			Name:          fmt.Sprintf("%#v", b.font),
			TextureLoader: &texGen{b},
			FilterMode:    gorge.TextureFilterLinear,
		},
	}

	return font, nil
}

func vec4ToColor(v m32.Vec4) color.Color {
	v = v.Mul(255)
	return color.RGBA{uint8(v[0]), uint8(v[1]), uint8(v[2]), uint8(v[3])}
}

// glyph returns the number of printable glyphs and maximum width,height
func glyphCount(s string) int {
	count := 0
	for _, ch := range s {
		if unicode.IsSpace(ch) {
			continue
		}
		count++
	}
	return count
}

// given a max square size we calculate the optimal subsize
func calcSize(sz float32, n int) float32 {
	a := float32(sz * sz)
	ia := a / float32(n)
	il := m32.Sqrt(ia)
	nw := m32.Ceil(sz / il)
	nh := m32.Ceil(sz / il)
	l := m32.Min(sz/nw, sz/nh)

	return l
}

type texGen struct {
	builder fontBuilder
}

// Data will call underlying func to return texture
func (g texGen) Data() *gorge.TextureData {
	img, err := g.builder.image()
	if err != nil {
		panic("uh oh")
	}
	rgb := img.(*image.RGBA)
	return &gorge.TextureData{
		Source:    fmt.Sprintf("%#v", g.builder.font),
		Format:    gorge.TextureFormatRGBA,
		Width:     rgb.Bounds().Dx(),
		Height:    rgb.Bounds().Dy(),
		PixelData: rgb.Pix,
	}
}
