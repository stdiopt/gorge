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

// Text renderable entity
type Text struct {
	gorge.Transform
	gorge.Renderable
	s    string
	mesh *gorge.Mesh
	font *Font

	Min vec2
	Max vec2
	//Width  float32
	//Height float32
}

// New setups the renderable with a font
func New(font *Font) *Text {
	mat := gorge.NewMaterial("pbr")
	mat.DoubleSided = true
	mat.SetTexture("albedoMap", font.Texture)
	mesh := &gorge.Mesh{}
	return &Text{
		Transform: gorge.Transform{},
		Renderable: gorge.Renderable{
			Color:    m32.Vec4{1, 1, 1, 1},
			Material: mat,
			Mesh:     mesh,
		},
		mesh: mesh,
		font: font,
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
			x += float32(t.font.spaceAdv)
			continue
		}
		if ch == '\n' {
			x = 0
			y--
			continue
		}
		g, ok := t.font.glyphs[ch]
		if !ok {
			g = t.font.glyphs['�']
			// default to some odd char
			//panic(fmt.Sprintf("glyph %v not found", ch))
		}
		w := g.size[0]
		h := g.size[1]
		xpos := x + g.bearingH
		ypos := y + h - g.bearingV - 0.6
		verts = append(verts, []float32{
			// Auto letters
			xpos, ypos, 0, g.uv1[0], g.uv1[1],
			xpos + w, ypos, 0, g.uv2[0], g.uv1[1],
			xpos + w, ypos - h, 0, g.uv2[0], g.uv2[1],

			xpos + w, ypos - h, 0, g.uv2[0], g.uv2[1],
			xpos, ypos - h, 0, g.uv1[0], g.uv2[1],
			xpos, ypos, 0, g.uv1[0], g.uv1[1],
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
		x += g.advance
	}

	t.mesh.MeshLoader = &gorge.MeshData{
		Format:   gorge.VertexFormatPT,
		Vertices: verts,
	}
	t.mesh.Updates++
}
