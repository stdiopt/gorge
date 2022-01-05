package render

import (
	"runtime"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/render/gl"
)

type textureManager struct {
	gorge      *gorge.Context
	texInvalid *Texture
	texWhite   *Texture

	count int
}

func newTextureManager(g *gorge.Context) *textureManager {
	m := &textureManager{
		gorge: g,
	}

	m.texInvalid = m.New(&gorge.TextureData{
		Source: "texturemanager.invalid",
		Width:  1, Height: 1,
		PixelData: []byte{255, 0, 255, 255},
	})

	m.texWhite = m.New(&gorge.TextureData{
		Source: "texturemanager.white",
		Width:  1, Height: 1,
		PixelData: []byte{255, 255, 255, 255},
	})

	return m
}

func (m *textureManager) New(r gorge.TextureResource) *Texture {
	t := &Texture{
		manager: m,
		updates: -1,
	}
	runtime.SetFinalizer(t, func(t *Texture) {
		m.gorge.RunInMain(func() {
			m.destroy(t)
		})
	})

	if d, ok := r.(*gorge.TextureData); ok {
		t.upload(d)
		return t
	}

	// Temporary texture
	t.upload(&gorge.TextureData{
		Format: gorge.TextureFormatRGBA,
		Width:  1, Height: 1,
		PixelData: []byte{0, 0, 0, 50},
	})

	return t
}

func (m *textureManager) destroy(t *Texture) {
	t.destroy()
}

func (m *textureManager) Bind(tex *gorge.Texture) {
	t := m.Get(tex)
	t.bind(tex)
}

func (m *textureManager) GetByRef(r gorge.TextureResource) *Texture {
	if r == nil {
		return m.texWhite
	}
	t, ok := gorge.GetGPU(r).(*Texture)
	if !ok {
		t = m.New(r)
		gorge.SetGPU(r, t)
		return t
	}
	if d, ok := r.(*gorge.TextureData); ok {
		t.update(d)
	}
	return t
}

func (m *textureManager) Get(tex *gorge.Texture) *Texture {
	if tex == nil {
		return m.texWhite
	}

	return m.GetByRef(tex.Resource)
}

func (m *textureManager) Update(r *gorge.TextureData) {
	t, ok := gorge.GetGPU(r).(*Texture)
	if !ok {
		t = m.New(r)
		gorge.SetGPU(r, t)
	}
	// Force an update
	t.updates--
	t.update(r)
}

// Texture is a opengl texture controller
type Texture struct {
	manager *textureManager // TODO: avoid putting manager here
	ID      gl.Texture
	Type    gl.Enum // gl.TEXTURE_2D, gl.TEXTURE_CUBE etc.. defaults to 2D
	// width, height int
	// Might not be needed?
	mipmap bool
	// updates indicates updates for dynamic TextureData
	updates int
}

func (t *Texture) destroy() {
	gl.DeleteTexture(t.ID)
	t.manager.count--
}

func (t *Texture) upload(data *gorge.TextureData) {
	if !gl.IsValid(t.ID) { // should assume valid?
		t.ID = gl.CreateTexture()
		t.manager.count++
	}
	t.Type = gl.TEXTURE_2D
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
	if data == nil || len(data.PixelData) == 0 {
		// Upload a pink image
		gl.TexImage2D(gl.TEXTURE_2D, 0,
			gl.RGBA, 1, 1,
			gl.RGBA, gl.UNSIGNED_BYTE, []byte{255, 0, 255, 255},
		)
		return
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)

	iformat, format, dt := TextureFormat(data.Format)
	// Set the rest
	gl.TexImage2D(gl.TEXTURE_2D, 0,
		iformat, data.Width, data.Height,
		format, dt, data.PixelData,
	)

	// Might need to recheck this for dynamic textures
	// Check if power of 2
	if /*data.Width == data.Height && */ data.Width > 1 {
		t.mipmap = true
		gl.GenerateMipmap(gl.TEXTURE_2D)
	}

	t.updates = data.Updates
}

// We should only update right on Get
func (t *Texture) update(data *gorge.TextureData) bool {
	if t.updates == data.Updates {
		return false
	}
	t.upload(data)
	return true
}

// Need better binding or update verifier
// This updates should happen only once when we update the Texture thingy
// since we are not tracking if the thing is updated we just push here
// might as well as add a local verifier
// should be called setup based on gorge tex
func (t *Texture) bind(gt *gorge.Texture) {
	gl.BindTexture(gl.TEXTURE_2D, t.ID)

	if gt == nil {
		return
	}
	// Only if updated? but if we reuse a texture data might change things
	// for this texture
	wu, wv, ww := gt.GetWrap()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, TextureWrap(wu))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, TextureWrap(wv))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, TextureWrap(ww))

	fm := gt.GetFilterMode()
	switch fm {
	case gorge.TextureFilterPoint:
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST_MIPMAP_NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	// case gorge.TextureFilterLinear:
	default:
		// if t.mipmap {
		// 	t.gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		// } else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		//}
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	}
	// Not all systems supports this
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAX_ANISOTROPY_EXT, 16)
}
