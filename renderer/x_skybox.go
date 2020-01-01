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
	"fmt"
	"image"
	"image/draw"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gl"
	"github.com/stdiopt/gorge/resource"
)

// PrepareSkybox prepares textures and skybox cube
func (rs *Renderer) PrepareSkybox() {
	g := rs.g

	var assets *resource.Manager
	rs.gorge.Query(func(m *resource.Manager) { assets = m })

	///////////////////////////////////////////////////////////////////////////
	// XXX: PLAYING AROUND HERE
	///////////////////////////////////////////////////////////////////////////

	srcs := []string{
		"skybox/right.jpg",
		"skybox/left.jpg",
		"skybox/top.jpg",
		"skybox/bottom.jpg",
		"skybox/front.jpg",
		"skybox/back.jpg",
	}
	images := [6]image.Image{}

	var err error
	for i, s := range srcs {
		images[i], err = assets.LoadImage(s)
		if err != nil {
			panic(err)
		}
	}
	width, _ := images[0].Bounds().Dx(), images[0].Bounds().Dy()

	tex := rs.textures.CreateCubeMap(width)
	for i, img := range images {

		switch img := img.(type) {
		case *image.RGBA:
			tex.SetImageCube(i, img.Pix)
		case *image.NRGBA:
			tex.SetImageCube(i, img.Pix)
		case *image.YCbCr:
			b := img.Bounds()
			m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
			draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
			tex.SetImageCube(i, m.Pix)
		default:
			panic(fmt.Sprintf("image is wrong: %T", img))
		}
	}

	mat := gorge.NewMaterial(skyboxShader)

	rs.skyboxTex = tex
	rs.skyboxShader = rs.shaders.Get(mat)

	rs.skyboxVAO = g.CreateVertexArray()
	g.BindVertexArray(rs.skyboxVAO)

	VBO := g.CreateBuffer()
	g.BindBuffer(gl.ARRAY_BUFFER, VBO)

	g.BufferDataX(gl.ARRAY_BUFFER, skyboxVert, gl.STATIC_DRAW)

	if a, ok := rs.skyboxShader.Attrib("aPosition"); ok {
		// Prepare a VAO and VBO
		g.EnableVertexAttribArray(a)
		g.VertexAttribPointer(a, 3, gl.FLOAT, false, 0, 0)
	}

}

// Skybox experimental
func (rs *Renderer) Skybox(ri *renderInfo) {
	g := rs.g

	g.ClearColor(ri.ambient[0], ri.ambient[1], ri.ambient[2], 1)
	g.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	g.Enable(gl.CULL_FACE)
	g.FrontFace(gl.CW)
	g.DepthMask(false)

	rs.skyboxShader.bind()
	rs.skyboxShader.Set("projection", ri.projection)
	rs.skyboxShader.Set("view", ri.view)

	g.BindVertexArray(rs.skyboxVAO)
	g.ActiveTexture(gl.TEXTURE0)
	g.BindTexture(gl.TEXTURE_CUBE_MAP, rs.skyboxTex.id)
	rs.skyboxShader.Set("skybox", 0)

	g.DrawArrays(gl.TRIANGLES, 0, 36)

	g.DepthMask(true)
	g.FrontFace(gl.CCW)

}

var skyboxVert = []float32{
	// positions
	// Back
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,

	// Left
	-1.0, -1.0, 1.0,
	-1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0,

	// Right
	1.0, -1.0, -1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, -1.0, -1.0,

	// Front
	-1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, -1.0, 1.0,
	-1.0, -1.0, 1.0,

	// Up
	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,

	// Down
	1.0, -1.0, -1.0, // 5
	-1.0, -1.0, -1.0, // 4
	1.0, -1.0, 1.0, // 3
	1.0, -1.0, 1.0, // 2
	-1.0, -1.0, -1.0, // 1
	-1.0, -1.0, 1.0, // 0
}

var skyboxShader = &gorge.ShaderData{
	VertSrc: `#version 300 es
layout (location = 0) in vec3 aPosition;

out vec3 TexCoords;

uniform mat4 projection;
uniform mat4 view;

void main()
{
	TexCoords = aPosition;
	mat4 lview = mat4(mat3(view));
	vec4 pos =  projection * lview * vec4(aPosition, 1.0);
	gl_Position = pos.xyww;
}`,
	FragSrc: `#version 300 es
precision highp float;

out vec4 FragColor;

in vec3 TexCoords;

uniform samplerCube skybox;

void main() {
	FragColor = texture(skybox, TexCoords);
}`,
}
