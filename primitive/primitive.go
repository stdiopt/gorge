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

package primitive

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
)

// Cool aliases
type (
	vec2 = m32.Vec2
	vec3 = m32.Vec3
	vec4 = m32.Vec4
	mat4 = m32.Mat4
)

// MeshEntity thing
type MeshEntity struct {
	gorge.Transform
	gorge.Renderable
}

// Primitive generator
/*type Primitive struct {
	scene *gorge.Scene
}

// Context returns a primitive creater based on scene
func Context(s *gorge.Scene) *Primitive {
	return &Primitive{s}
}

// Cube creates a cube entity
func (p *Primitive) Cube() *MeshEntity {
	// This is needed to be binded to scene
	mesh := p.scene.Assets().MeshFromData(Cube())
	mat := gorge.NewMaterial("")

	return &MeshEntity{
		*gorge.NewTransform(),
		gorge.Renderable{
			Color:    vec4{1, 1, 1, 1},
			Mesh:     mesh,
			Material: mat,
		},
	}
}

// Plane returns a plane entity
func (p *Primitive) Plane() *MeshEntity {
	mesh := p.scene.Assets().MeshFromData(Plane())
	mat := gorge.NewMaterial("")
	return &MeshEntity{
		*gorge.NewTransform(),
		gorge.Renderable{
			Color:    vec4{1, 1, 1, 1},
			Mesh:     mesh,
			Material: mat,
		},
	}
}

// Poly returns a polygon with n sides
func (p *Primitive) Poly(n int) *MeshEntity {
	mesh := p.scene.Assets().MeshFromData(Poly(n))
	mat := gorge.NewMaterial("")
	mat.DrawType = gorge.DrawTriangleFan

	return &MeshEntity{
		*gorge.NewTransform(),
		gorge.Renderable{
			Color:    vec4{1, 1, 1, 1},
			Mesh:     mesh,
			Material: mat,
		},
	}
}*/
