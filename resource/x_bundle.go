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
	"github.com/stdiopt/gorge"
)

// TODO: Store the loader and not the model?

// BundleFunc returns a func that creates bundles
func BundleFunc(m *Manager) func() gorge.AssetBundler {
	return func() gorge.AssetBundler {
		return newBundle(m)
	}
}

// This implements gorge.AssetBundler
type bundle struct {
	manager *Manager
	assets  map[interface{}]struct{}
}

// NewAssetBundle thingy
func newBundle(m *Manager) *bundle {
	return &bundle{
		manager: m,
		assets:  map[interface{}]struct{}{},
	}
}
func (b *bundle) Assets() map[interface{}]struct{} {
	return b.assets
}

func (b *bundle) AddAsset(v interface{}) {
	b.store(v)
}

// Texture2D returns a texture
func (b *bundle) Texture2D(name string) *gorge.Texture {
	loader := b.manager.texture2DLoader(name)

	t := gorge.NewTexture(loader)
	b.store(t)
	return t
}

// Mesh returns a mesh with a loader reference
func (b *bundle) Mesh(name string) *gorge.Mesh {
	loader := b.manager.meshLoader(name)
	mesh := gorge.NewMesh(loader)
	b.store(mesh)
	return mesh
}

// MeshFromData returns a mesh from mesh data that will be uploaded to gpu
/*func (b *bundle) MeshFromData(d *gorge.MeshData) *gorge.Mesh {
	mesh := gorge.NewMesh(d)
	b.store(mesh)
	return mesh
}*/

// Texture2DFromData returns a texture from data
// that will be uploaded to gpu
/*func (b *bundle) Texture2DFromData(d *gorge.TextureData) *gorge.Texture {
	texture := gorge.NewTexture(d)
	b.store(texture)
	return texture
}*/

// Material loads a shader and returns the material
func (b *bundle) Material(name ...string) *gorge.Material {
	loader := b.manager.shaderLoader(name...)
	material := gorge.NewMaterial(loader)
	b.store(material)
	return material
}

// Font loads a font and returns the text.Font
func (b *bundle) Font(name string, optFn ...gorge.FontOptionsFunc) *gorge.Font {
	font := b.manager.Font(name, optFn...)
	b.store(font.Texture)
	return font
}

func (b *bundle) store(a interface{}) {
	b.assets[a] = struct{}{}
}
