package resource

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"log"
	"unsafe"

	xdraw "golang.org/x/image/draw"

	// Image loaders
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	// External support for .hdr files
	"github.com/mdouchement/hdr"

	// Import HDR rgbe driver
	_ "github.com/mdouchement/hdr/codec/rgbe"

	"github.com/stdiopt/gorge"
)

const maxaddr = 0x7FFFFFFF

func init() {
	exts := []string{
		".jpg", ".jpeg",
		".png",
		".gif",
		".hdr",
	}

	for _, ext := range exts {
		// Remove texture from here not valid, as we will only load Data instead of binding to texture directly
		Register((*gorge.Texture)(nil), ext, textureLoader)
		Register((*gorge.TextureData)(nil), ext, textureDataLoader)
	}
}

func textureLoader(res *Context, v interface{}, name string, _ ...interface{}) error {
	tex, ok := v.(*gorge.Texture)
	if !ok {
		return fmt.Errorf("unable to load data into: %T", v)
	}

	var texData gorge.TextureData
	if err := textureDataLoader(res, &texData, name); err != nil {
		return err
	}
	tex.Resource = &texData

	return nil
}

func textureDataLoader(res *Context, v interface{}, name string, _ ...interface{}) error {
	texData := v.(*gorge.TextureData)

	rd, err := res.Open(name)
	if err != nil {
		return fmt.Errorf("[resource] error opening image: %w", err)
	}

	td, err := ReadTexture(rd)
	if err != nil {
		return err
	}
	*texData = *td
	texData.Source = name

	return nil
}

// ReadTexture reads an image from the reader and returns a textureData.
func ReadTexture(rd io.Reader) (*gorge.TextureData, error) {
	img, _, err := image.Decode(rd)
	if err != nil {
		return nil, fmt.Errorf("[resource] error decoding image: %w", err)
	}

	return TextureDataFromImage(img)
}

// TextureDataFromImage converts a go image.Image to gorge.TextureData.
func TextureDataFromImage(im image.Image) (*gorge.TextureData, error) {
	dim := im.Bounds()
	var format gorge.TextureFormat
	var pixData []byte
	switch im := im.(type) {
	case *hdr.RGB:
		sz := len(im.Pix) * 4
		byteData := (*(*[maxaddr]byte)(unsafe.Pointer(&im.Pix[0])))[:sz:sz]

		format = gorge.TextureFormatRGB32F
		pixData = append([]byte{}, byteData...)
	case *image.NRGBA:
		if dim = getPowerOf2Dim(dim); dim != im.Bounds() {
			orig := im
			im = image.NewNRGBA(dim)
			convImg(im, orig)
		}

		format = gorge.TextureFormatRGBA
		pixData = im.Pix
	case *image.RGBA:
		if dim = getPowerOf2Dim(dim); dim != im.Bounds() {
			orig := im
			im = image.NewRGBA(dim)
			convImg(im, orig)
		}

		format = gorge.TextureFormatRGBA
		pixData = im.Pix
	case *image.YCbCr, *image.RGBA64, *image.Paletted: // We convert these for now
		dim = getPowerOf2Dim(dim)

		dimg := image.NewRGBA(dim)
		convImg(dimg, im)

		format = gorge.TextureFormatRGBA
		pixData = dimg.Pix
	case *image.Alpha:
		if dim = getPowerOf2Dim(dim); dim != im.Bounds() {
			orig := im
			im = image.NewAlpha(dim)
			convImg(im, orig)
		}

		format = gorge.TextureFormatGray
		pixData = im.Pix
	case *image.Gray:
		if dim = getPowerOf2Dim(dim); dim != im.Bounds() {
			orig := im
			im = image.NewGray(dim)
			convImg(im, orig)
		}

		format = gorge.TextureFormatGray
		pixData = im.Pix
	case *image.Gray16:
		dim = getPowerOf2Dim(dim)

		dimg := image.NewGray(dim)
		convImg(dimg, im)

		format = gorge.TextureFormatGray
		pixData = dimg.Pix
	default:
		return nil, fmt.Errorf("[resource] unsupported image: %T", im)
	}

	texData := gorge.TextureData{
		Format:    format,
		Width:     dim.Dx(),
		Height:    dim.Dy(),
		PixelData: pixData,
	}
	return &texData, nil
}

// TextureFromImage converts a go image.Image to gorge.Texture.
func TextureFromImage(im image.Image) (*gorge.Texture, error) {
	texData, err := TextureDataFromImage(im)
	if err != nil {
		return nil, err
	}
	tex := gorge.NewTexture(texData)
	return tex, nil
}

func getPowerOf2Dim(src image.Rectangle) image.Rectangle {
	w := src.Bounds().Dx()
	h := src.Bounds().Dy()
	if isPowerOf2(w) && isPowerOf2(h) {
		return src
	}
	if w == h { // swquare
		s := nearestPowerOf2(w)
		return image.Rect(0, 0, s, s)
	}
	w = nearestPowerOf2(makeEven(w))
	h = nearestPowerOf2(makeEven(h))

	return image.Rect(0, 0, w, h)
}

func convImg(dst draw.Image, src image.Image) {
	if dst.Bounds() == src.Bounds() {
		log.Printf("Converting image: %T:%v %T:%v", src, src.Bounds(), dst, dst.Bounds())
		draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)
		return
	}
	log.Printf("Rescaling image: %T:%v -> %T%v", src, src.Bounds(), dst, dst.Bounds())
	xdraw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), xdraw.Over, nil)
}
