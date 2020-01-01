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

package renderer

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gl"
)

// TODO: complex data struct here with sets and all of that?
// State management
// type texture struct {
//   gl.Texture
// }
// func( t *texture) SetWrapS(gl.REPEAT) {
//   If texture not binded in manager, bind it
//   gl.TexParameteri etc...
// }
//

type textureManager struct {
	g gl.Context3
	//assets *asset.System

	// Use name to loose dependency?
	textures map[interface{}]*texture

	gray    *texture
	invalid *texture
	normal  *texture
	white   *texture
}

func newTextureManager(g gl.Context3) *textureManager {
	// Prepare some weird textures here too

	tm := &textureManager{
		g: g,
		//assets:   assets,
		textures: map[interface{}]*texture{},
	}

	grayTex := gorge.NewTexture(&gorge.TextureData{
		Width: 1, Height: 1,
		PixelData: []byte{127, 127, 127, 255},
	})
	invalidTex := gorge.NewTexture(&gorge.TextureData{
		Width: 1, Height: 1,
		PixelData: []byte{255, 0, 255, 255},
	})
	tm.gray = tm.Get(grayTex)
	tm.invalid = tm.Get(invalidTex)

	//tm.gray = tm.Get(gray)
	//tm.gray = tm.Create2D(1, 1)
	//tm.gray.SetImage2D([]byte{127, 127, 127, 255})

	tm.normal = tm.Create2D(1, 1)
	tm.normal.SetImage2D([]byte{127, 127, 255, 255})

	tm.white = tm.Create2D(1, 1)
	tm.white.SetImage2D([]byte{255, 255, 255, 255})

	return tm
}

func (tm *textureManager) Get(t *gorge.Texture) *texture {
	return tm.get(t)
}

func (tm *textureManager) get(t *gorge.Texture) *texture {
	k := t
	if tex, ok := tm.textures[k]; ok {
		tex.update()
		return tex
	}

	tex := &texture{
		g:           tm.g,
		texture:     t,
		id:          tm.g.CreateTexture(),
		updates:     -1,
		dataUpdates: -1,
	}
	tex.update()

	tm.textures[k] = tex
	return tex
}

// Experiments
// CreateCubeMap size width equal to height per side
func (tm *textureManager) CreateCubeMap(size int) *texture {
	g := tm.g

	tid := g.CreateTexture()
	g.BindTexture(gl.TEXTURE_CUBE_MAP, tid)

	// Default texture params
	g.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	g.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	g.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	g.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	g.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	/*for i := uint32(0); i < 6; i++ {
		g.TexImage2D(
			gl.TEXTURE_CUBE_MAP_POSITIVE_X+i,
			0,
			gl.RGB,
			width, height, gl.RGB, gl.UNSIGNED_BYTE, nil)
	}*/
	tex := &texture{
		g:      tm.g,
		id:     tid,
		width:  size,
		height: size,
	}
	return tex
}

// Local state tracking
type texture struct {
	manager *textureManager
	g       gl.Context3
	// TODO: Ref or bundle ref here?
	id gl.Texture

	width, height int
	texture       *gorge.Texture
	updates       int
	dataUpdates   int
}

func (t *texture) update() {
	t.g.BindTexture(gl.TEXTURE_2D, t.id)

	t.updateData()
	t.updateParam()

}

// Assume texture is binded
func (t *texture) updateData() {
	if t.texture.DataUpdates == t.dataUpdates {
		return
	}
	t.dataUpdates = t.texture.DataUpdates

	g := t.g

	texData := t.texture.Loader().Data()
	if texData == nil {
		// Update a pink image
		g.TexImage2D(gl.TEXTURE_2D, 0,
			gl.RGBA, 1, 1,
			gl.RGBA, gl.UNSIGNED_BYTE, []byte{255, 0, 255, 255},
		)
		return
	}
	// Set the rest
	g.TexImage2D(gl.TEXTURE_2D, 0,
		gl.RGBA, texData.Width, texData.Height,
		gl.RGBA, gl.UNSIGNED_BYTE, texData.PixelData,
	)
	if texData.Width == texData.Height && texData.Width > 1 {
		g.GenerateMipmap(gl.TEXTURE_2D)
	}
}

func (t *texture) updateParam() {
	// No updates needed
	if t.texture.Updates == t.updates {
		return
	}
	g := t.g
	t.updates = t.texture.Updates

	if t.texture.WrapU != 0 {
		g.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, int(t.texture.WrapU))
	}
	if t.texture.WrapV != 0 {
		g.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, int(t.texture.WrapV))
	}
	if t.texture.WrapW != 0 {
		g.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, int(t.texture.WrapW))
	}

	switch t.texture.FilterMode {
	case gorge.TextureFilterPoint:
		g.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST_MIPMAP_NEAREST)
		g.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	//case gorge.TextureFilterLinear:
	default:
		g.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		g.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	}

	g.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAX_ANISOTROPY_EXT, 16)
}

// Create2D Creates a 2D texture in the gpu
func (tm *textureManager) Create2D(width, height int) *texture {
	g := tm.g

	tid := g.CreateTexture()
	// Setup tex stuff from gorge.Texture
	g.BindTexture(gl.TEXTURE_2D, tid)
	g.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	g.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	g.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	g.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	g.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAX_ANISOTROPY_EXT, 16)

	return &texture{
		g:      tm.g,
		id:     tid,
		width:  width,
		height: height,
	}
}

// Set clamp and what nots?
func (t *texture) SetImage2D(data []byte) {
	t.g.BindTexture(gl.TEXTURE_2D, t.id)
	// Carefull texture needs to be binded
	t.g.TexImage2D(
		gl.TEXTURE_2D, 0,
		gl.RGBA, t.width, t.height,
		gl.RGBA, gl.UNSIGNED_BYTE, data,
	)
	t.g.GenerateMipmap(gl.TEXTURE_2D)
}

func (t *texture) SetImageCube(i int, data []byte) {
	// Carefull texture needs to be binded
	t.g.TexImage2D(
		gl.Enum(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i),
		0, gl.RGBA, t.width, t.height,
		gl.RGBA, gl.UNSIGNED_BYTE, data,
	)
}
