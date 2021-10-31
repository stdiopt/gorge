package text

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"golang.org/x/image/font"
)

// Glyph char information
type Glyph struct {
	Uv1  m32.Vec2
	Uv2  m32.Vec2
	Size m32.Vec2

	Advance  float32
	BearingH float32
	BearingV float32
}

// Font is like a texture with extra information about glyphs
type Font struct {
	face     font.Face
	Glyphs   map[rune]Glyph
	SpaceAdv float32
	// Size is the internally rendered size
	Size float32
	*gorge.Texture
}

func (f *Font) getGlyph(ch rune) Glyph {
	g, ok := f.Glyphs[ch]
	if !ok {
		return f.Glyphs['�'] // Special one
	}
	return g
}

// FontOptions options for font
type FontOptions struct {
	Resolution int
	Chars      string
	Background *m32.Vec4
	Foreground *m32.Vec4
}

// FontOptionsFunc func to manipulate font options.
type FontOptionsFunc func(o *FontOptions)

// FontResolution sets the font texture resolution option.
func FontResolution(n int) FontOptionsFunc {
	return func(opt *FontOptions) {
		opt.Resolution = n
	}
}

// FontBackground sets the font texture background option.
func FontBackground(c m32.Vec4) FontOptionsFunc {
	return func(opt *FontOptions) {
		opt.Background = &c
	}
}

// FontForeground sets the font texture foreground option.
func FontForeground(c m32.Vec4) FontOptionsFunc {
	return func(opt *FontOptions) {
		opt.Foreground = &c
	}
}
