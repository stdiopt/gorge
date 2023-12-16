package text

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"unicode"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/resource"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

func init() {
	resource.Register((*Font)(nil), ".ttf", fontLoader)
}

var commonChars = []rune("`" + `
0123456789µ&
abcdefghijklmnopqrstuvwxyz
ABCDEFGHIJKLMNOPQRSTUVWXYZ
{}()[]|$@?!%/\:;,._-+=<>*"'~#
áéíóúÁÉÍÓÚçÇãÃõÕ
►▼◀▲
`)

func fontLoader(res *resource.Context, v any, name string, opts ...any) error {
	fontOut := v.(*Font)

	opt := FontOptions{}

	for _, o := range opts {
		fontOpt, ok := o.(FontOptionsFunc)
		if !ok {
			return fmt.Errorf("wront options type")
		}
		fontOpt(&opt)
	}

	bg := color.Color(color.RGBA{0, 0, 0, 0})
	fg := color.Color(color.RGBA{255, 255, 255, 255})
	if opt.Background != nil {
		bg = vec4ToColor(*opt.Background)
	}
	if opt.Foreground != nil {
		fg = vec4ToColor(*opt.Foreground)
	}

	resolution := 1024
	if opt.Resolution != 0 {
		resolution = opt.Resolution
	}

	chars := opt.Chars
	if len(chars) == 0 {
		chars = commonChars
	}
	chars = append([]rune{'�'}, chars...) // prepend the interrogation

	// Load font and dependents
	fontData, err := res.LoadBytes(name)
	if err != nil {
		return err
	}
	sff, err := sfnt.Parse(fontData)
	if err != nil {
		return err
	}
	// Ideal font size
	count := glyphCount(chars)
	scale := calcSize(float32(resolution), count)

	var b sfnt.Buffer
	bnd, err := sff.Bounds(&b, fixed.I(int(scale)), font.HintingFull)
	if err != nil {
		return err
	}
	boundW := (bnd.Max.X - bnd.Min.X).Floor()
	boundH := (bnd.Max.Y - bnd.Min.Y).Floor()

	fontScale := float32(.8) // a bit smaller because of weird glyphs
	size := float64(scale * fontScale)

	// Create our image (texture)
	clip := image.Rect(0, 0, resolution, resolution)
	iw := float32(clip.Bounds().Dx())
	ih := float32(clip.Bounds().Dy())

	img := image.NewRGBA(clip)
	draw.Draw(img, img.Bounds(), image.NewUniform(bg), image.Point{}, draw.Src)

	face, err := opentype.NewFace(sff, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}
	// this doesn't mess with Kern func
	defer face.Close() // nolint: errcheck

	dot := fixed.P(0, int(scale))
	glyphs := map[rune]Glyph{}
	ladv := fixed.I(int(scale))
	// Build each char information
	for _, ch := range chars {
		if unicode.IsSpace(ch) {
			continue
		}
		bnd, adv, ok := face.GlyphBounds(ch)
		if !ok {
			return fmt.Errorf("glyph bounds error")
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
			// above can sometimes yield 0 for font smaller than 48pt, 1 is minimum
			if gw == 0 || gh == 0 {
				gw = 1
				gh = 1
			}
		}
		// Center stuff somehow?
		d := dot.Sub(bnd.Min)
		gr, mask, maskp, _, ok := face.Glyph(d, ch)
		if !ok {
			return fmt.Errorf("glyph error")
		}
		// gr = clip.Intersect(gr)
		draw.DrawMask(img, gr, image.NewUniform(fg), image.Point{}, mask, maskp, draw.Over)

		p1x := float32(gr.Min.X)
		p1y := float32(gr.Min.Y)
		p2x := float32(gr.Max.X)
		p2y := float32(gr.Max.Y)
		glyphs[ch] = Glyph{
			Uv1:      gm.Vec2{p1x / iw, p1y / ih},
			Uv2:      gm.Vec2{p2x / iw, p2y / ih},
			Size:     gm.Vec2{float32(gw) / scale, float32(gh) / scale},
			Advance:  float32(adv>>6) / scale,
			BearingV: float32(bnd.Max.Y>>6) / scale,
			BearingH: float32(bnd.Min.X>>6) / scale,
		}

		dot.X += ladv
	}
	adv, _ := face.GlyphAdvance(' ')

	tex := gorge.NewTexture(&gorge.TextureData{
		Source:    fmt.Sprintf("%v", name),
		Format:    gorge.TextureFormatRGBA,
		Width:     img.Bounds().Dx(),
		Height:    img.Bounds().Dy(),
		PixelData: img.Pix,
	})
	tex.FilterMode = gorge.TextureFilterLinear
	tex.ReleaseData(res.Gorge())

	*fontOut = Font{
		SpaceAdv: float32(adv>>6) / scale,
		Size:     float32(size),
		Glyphs:   glyphs,
		Texture:  tex,
		Face:     face,
	}

	return nil
}

func vec4ToColor(v gm.Vec4) color.Color {
	v = v.Mul(255)
	return color.RGBA{uint8(v[0]), uint8(v[1]), uint8(v[2]), uint8(v[3])}
}

// glyph returns the number of printable glyphs and maximum width,height
func glyphCount(s []rune) int {
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
	a := sz * sz
	ia := a / float32(n)
	il := gm.Sqrt(ia)
	nw := gm.Ceil(sz / il)
	nh := gm.Ceil(sz / il)
	l := gm.Min(sz/nw, sz/nh)

	return l
}
