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

package gorgeutils

import (
	"math"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/primitive"
)

// MeshEntity thing
type MeshEntity struct {
	gorge.Transform
	gorge.Renderable
}

// Gimbal Compost object
type Gimbal struct {
	Entities []gorge.Entity
	*gorge.Transform
}

// NewGimbal creates entities on manager
func NewGimbal() *Gimbal {
	// Parent thing
	root := gorge.NewTransform()

	line := gorge.NewMesh(&gorge.MeshData{
		Name:   "gorgeutils.Gimbal",
		Format: gorge.VertexFormatP,
		Vertices: []float32{
			0, 0, 0,
			0, 0, 1,
		},
		Indices: []uint32{},
	})
	line.DrawType = gorge.DrawLines

	rot90 := float32(math.Pi / 2)

	objs := []struct {
		axis vec3
		rot  vec3
	}{
		{axis: vec3{0, 0, 1}, rot: vec3{}},
		{axis: vec3{0, 1, 0}, rot: vec3{-rot90, 0, 0}},
		{axis: vec3{1, 0, 0}, rot: vec3{0, rot90, 0}},
	}

	g := &Gimbal{
		Entities:  []gorge.Entity{},
		Transform: root,
	}

	for _, o := range objs {
		color := o.axis.Vec4(1)

		mat := gorge.NewMaterial(nil)

		l := &MeshEntity{
			*gorge.NewTransform(),
			gorge.Renderable{
				Mesh:     line,
				Material: mat,
				Color:    color,
			},
		}
		l.TransformComponent().Rotatev(o.rot).SetParent(root)

		g.Entities = append(g.Entities, l)

	}
	for _, o := range objs {
		color := o.axis.Vec4(1)

		b := primitive.Cube()
		b.SetPositionv(o.axis).
			SetScale(0.08).
			SetParent(root)
		b.Renderable.Color = color

		g.Entities = append(g.Entities, b)
	}
	return g
}
