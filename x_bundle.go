// Copyright 2020 Luis Figueiredo
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
	"io"
)

// AssetBundleProvider event
type AssetBundleProvider struct {
	Bundle func() AssetBundler
}

// LoadBundleEvent event for loading assets within the map
type LoadBundleEvent struct {
	Assets map[interface{}]struct{}
}

// Opener contains a method to open a specific asset with no arguments
type Opener interface {
	Open() (io.ReadCloser, error)
}

// AssetBundler methods to load assets
type AssetBundler interface {
	Texture2D(string) *Texture
	Mesh(string) *Mesh
	Font(string, ...FontOptionsFunc) *Font
	Material(...string) *Material
	AddAsset(v interface{})

	Assets() map[interface{}]struct{}
}

// XXX: note
// The purpose of this is to Allow fully unloading assets from GPU
// and that includes primitives
// make primitives return MeshData only
// and document for the user to setup the bundle system beforehand
// Automated stuff should be using a bundle
// Or we can also do a reference walks

// This was on asset pkg, maybe should be back there
// and make the scene use this as in
//	g := gorge.New(...)
//  s := g.Scene(...)
//  s.UseAssets(bundle)
//
//  or even worst
//
//  dependency injection
//  gorge contains a struct/interface with an asset interface
//  gorge requests to the messaging the implementation of Asset Loading
//   after initialization
//  gorge uses that on scene
//

// AssetBundle for loading purposes
// This shall be implemented in asset thing
/*type assetBundle struct {
	assets map[interface{}]struct{}
}

// NewAssetBundle thingy
func newAssetBundle() *assetBundle {
	return &assetBundle{
		assets: map[interface{}]struct{}{},
	}
}
func (b *assetBundle) Assets() map[interface{}]struct{} {
	return b.assets
}

// Texture2D returns a texture
func (b *assetBundle) Texture2D(name string) *Texture {
	t := &Texture{
		TextureLoader: TextureLoaderRef{Name: name},
	}
	b.store(t)
	return t
}

// Texture2DFromData returns a texture from data
// that will be uploaded to gpu
func (b *assetBundle) Texture2DFromData(d *TextureData) *Texture {
	t := &Texture{TextureLoader: d}
	b.store(t)
	return t
}

// Mesh returns a mesh with a loader reference
func (b *assetBundle) Mesh(name string) *Mesh {
	m := &Mesh{
		MeshLoader: MeshLoaderRef{Name: name},
	}
	b.store(m)
	return m
}

// MeshFromData returns a mesh from mesh data that will be uploaded to gpu
func (b *assetBundle) MeshFromData(d *MeshData) *Mesh {
	m := &Mesh{MeshLoader: d}
	b.store(m)
	return m
}

// Font loads a font with options
// XXX: Font is different since we need glyph sizes right away
// So we kind of need asset system/fontLoader here
func (b *assetBundle) Font(fo fontOpener) *Font {
	f, err := NewFont(fo)
	if err != nil {
		// XXX: fix this
		panic(fmt.Sprintf("wrong font: %s", err))
	}
	b.store(f.Texture)
	return f
}

// FontWithOptions loads a font with options
// XXX: Font is different since we need glyph sizes right away
// So we kind of need asset system/fontLoader here
func (b *assetBundle) FontWithOptions(fo fontOpener, opt FontOptions) *Font {
	f, err := FontWithOptions(fo, opt)
	if err != nil {
		// XXX: fix this
		panic(fmt.Sprintf("wrong font: %s", err))
	}
	b.store(f.Texture)
	return f
}

func (b *assetBundle) store(a interface{}) {
	b.assets[a] = struct{}{}
}*/
