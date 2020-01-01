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

package main

import (
	"fmt"
	"log"
	"math"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gorgeutils"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/platform"
	"github.com/stdiopt/gorge/primitive"
	"github.com/stdiopt/gorge/resource"
	"github.com/stdiopt/gorge/x/text"
)

type (
	vec2 = m32.Vec2
	vec3 = m32.Vec3
	vec4 = m32.Vec4
)

func main() {

	opt := platform.Options{
		Wasm: platform.WasmOptions{
			Loader: resource.HTTPLoader{BaseURL: "../assets"},
		},
		GLFW: platform.GLFWOptions{
			Loader: resource.FileLoader{BasePath: "/assets"},
		},
	}
	platform.Start(opt, func(g *gorge.Gorge) {
		s := g.Scene(textSystem)
		g.StartScene(s)
	})
}

/*func loadFont(s *gorge.Scene, name string) *gorge.Font {

	font :=
	// fontTexture
	font, err := text.NewFont(assets.Asset(rsrc))
	if err != nil {
		log.Println("error loading font:", err)
	}
	return font
}*/

func textSystem(s *gorge.Scene) {

	light := gorgeutils.NewLight()
	light.SetPosition(0, 0, -1)
	gym := gorgeutils.NewGimbal()
	gym.SetParent(light)

	// Add entity
	gorgeutils.TrackballCamera(s)

	assets := s.Assets()
	fonts := []*gorge.Font{
		assets.Font("fonts/open-sans.ttf"),
		assets.Font("fonts/mashanzheng.ttf"),
		assets.Font("fonts/inria.ttf"),
	}

	// Preload textures :/

	textParent := gorge.NewTransform()
	var planes []*primitive.MeshEntity
	var texts []*text.Text
	for i, f := range fonts {
		p := primitive.Plane() // new plane I guess
		p.Material = assets.Material("shaders/unlit")
		p.Material.SetTexture("albedoMap", f.Texture)
		p.SetPosition(float32(-3+i*3), 0, 0.1).
			Rotate(-math.Pi/2, 0, 0)
		planes = append(planes, p)

		t := text.New(f)
		t.SetParent(textParent).
			SetPosition(0, -0.5+float32(i)*1.5, 1).
			SetScale(.4)
		t.Material.
			SetFloat32("ao", 1).
			SetFloat32("roughness", 0.1).
			SetFloat32("metallic", 0.9)

		texts = append(texts, t)
	}
	s.AddEntity(gym.Entities...)
	s.AddEntity(light)

	for _, p := range planes {
		s.AddEntity(p)
	}
	for _, t := range texts {
		s.AddEntity(t)
	}

	log.Println("Adding entities")

	total := float32(0)
	s.Handle(func(evt gorge.UpdateEvent) {
		total += float32(evt)
		s := fmt.Sprintf("together sep The weird BROWN fox quick jumped the gray fox: %.2f", total)
		for i, t := range texts {
			t.SetText(s)
			size := t.Max.Sub(t.Min)
			t.SetPosition(-(size[0]/2)*t.Scale[0], float32(i)*0.4, 1)

		}
	})

}
