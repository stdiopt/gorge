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

// ShaderLoader loads shader data for a material
type ShaderLoader interface {
	Data() *ShaderData
}

// ShaderData contains shaders sources
type ShaderData struct {
	VertSrc string
	FragSrc string
}

// Data convinient to return memory data
func (d ShaderData) Data() *ShaderData { return &d }

// Material the material
type Material struct {
	loader ShaderLoader
	// Primitive stuff
	Depth       bool
	DoubleSided bool
	// Program name (not being used yet)
	Program string
	//Color   vec4

	// Texture loaders instead?
	Textures map[string]*Texture
	props    map[string]interface{}
}

// NewMaterial returns a initialized Material
func NewMaterial(l ShaderLoader) *Material {
	return &Material{
		loader: l,
		Depth:  true,
		//Color:    vec4{0.5, 0.5, 0.5, 1},
	}
}

// Loader returns the shader loader for this material
func (m *Material) Loader() ShaderLoader {
	return m.loader
}

// SetTexture uniform thing for specific name
func (m *Material) SetTexture(k string, t *Texture) *Material {
	if m.Textures == nil {
		m.Textures = map[string]*Texture{}
	}
	m.Textures[k] = t
	return m
}

// Set properties by name
func (m *Material) Set(name string, value interface{}) *Material {
	if m.props == nil {
		m.props = map[string]interface{}{}
	}
	if f, ok := value.(float64); ok {
		m.props[name] = float32(f)
		return m
	}
	m.props[name] = value

	return m
}

// Get return named property
func (m *Material) Get(name string) interface{} {
	if m.props == nil {
		return nil
	}
	return m.props[name]
}

// Props returns the properties of this material
func (m *Material) Props() map[string]interface{} {
	return m.props
}

// SetFloat32 XXX testing sets a float32
func (m *Material) SetFloat32(name string, v float32) *Material {
	return m.Set(name, v)
}
