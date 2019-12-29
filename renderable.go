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

// Renderable contains info for renderer
// material and mesh
type Renderable struct {
	Name     string
	Color    vec4
	Material *Material
	Mesh     *Mesh
}

// RenderableComponent to satisfy component
func (r *Renderable) RenderableComponent() *Renderable { return r }

// SetMaterial sets the material
func (r *Renderable) SetMaterial(m *Material) *Renderable {
	r.Material = m
	return r
}

// SetMesh sets the mesh
func (r *Renderable) SetMesh(m *Mesh) *Renderable {
	r.Mesh = m
	return r
}

// SetColor sets the color
func (r *Renderable) SetColor(c vec4) *Renderable {
	r.Color = c
	return r
}

// NewRenderable returns a new mesh
func NewRenderable(n string, mesh *Mesh, mat *Material) *Renderable {
	return &Renderable{
		Name:     n,
		Color:    vec4{0.5, 0.5, 0.5, 1},
		Mesh:     mesh,
		Material: mat,
	}
}
