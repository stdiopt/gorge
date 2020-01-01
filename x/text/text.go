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

package text

import (
	"fmt"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
)

type vec2 = m32.Vec2

// Text renderable entity
type Text struct {
	gorge.Transform
	gorge.Renderable
	s        string
	mesh     *gorge.Mesh
	meshData *gorge.MeshData
	font     *gorge.Font

	Min vec2
	Max vec2
	//Width  float32
	//Height float32
}

// New setups the renderable with a font
func New(font *gorge.Font) *Text {
	mat := gorge.NewMaterial(nil)
	mat.SetTexture("albedoMap", font.Texture)
	meshData := &gorge.MeshData{Name: "text"}
	mesh := gorge.NewMesh(meshData)
	return &Text{
		Transform: gorge.Transform{},
		Renderable: gorge.Renderable{
			Color:    m32.Vec4{1, 1, 1, 1},
			Material: mat,
			Mesh:     mesh,
		},
		mesh:     mesh,
		meshData: meshData,
		font:     font,
	}
}

// SetTextf sets formated text
func (t *Text) SetTextf(f string, args ...interface{}) {
	t.SetText(fmt.Sprintf(f, args...))
}

// SetText updates the underlying vertices with the new text
// TODO: reduce float allocations by editing floats
func (t *Text) SetText(s string) {
	verts := []float32{}

	var x, y float32
	for _, ch := range s {
		if ch == ' ' {
			x += float32(t.font.SpaceAdv)
			continue
		}
		if ch == '\n' {
			x = 0
			y--
			continue
		}
		g, ok := t.font.Glyphs[ch]
		if !ok {
			g = t.font.Glyphs['�']
			// default to some odd char
			//panic(fmt.Sprintf("glyph %v not found", ch))
		}
		w := g.Size[0]
		h := g.Size[1]
		xpos := x + g.BearingH
		ypos := y + h - g.BearingV - 0.6
		verts = append(verts, []float32{
			// Auto letters
			xpos, ypos, 0, g.Uv1[0], g.Uv1[1],
			xpos + w, ypos - h, 0, g.Uv2[0], g.Uv2[1],
			xpos + w, ypos, 0, g.Uv2[0], g.Uv1[1],

			xpos + w, ypos - h, 0, g.Uv2[0], g.Uv2[1],
			xpos, ypos, 0, g.Uv1[0], g.Uv1[1],
			xpos, ypos - h, 0, g.Uv1[0], g.Uv2[1],
		}...)
		if xpos < t.Min[0] {
			t.Min[0] = xpos
		}
		if xpos+w > t.Max[0] {
			t.Max[0] = xpos + w
		}
		if ypos-h < t.Min[1] {
			t.Min[1] = ypos - h
		}
		if ypos > t.Max[1] {
			t.Max[1] = ypos
		}
		x += g.Advance
	}

	t.meshData.Format = gorge.VertexFormatPT
	t.meshData.Vertices = verts
	t.meshData.Updates++
}
