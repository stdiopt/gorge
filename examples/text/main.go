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
	"github.com/stdiopt/gorge/asset"
	"github.com/stdiopt/gorge/gorgeutils"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/platform"
	"github.com/stdiopt/gorge/primitive"
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
			AssetLoader: asset.HTTPLoader{BaseURL: "../assets"},
		},
		GLFW: platform.GLFWOptions{
			AssetLoader: asset.FileLoader{BasePath: "/assets"},
		},
	}
	platform.Start(opt, textSystem)
}

func loadFont(assets *asset.System, rsrc string) *text.Font {
	// fontTexture
	font, err := text.NewFont(assets.Asset(rsrc))
	if err != nil {
		log.Println("error loading font:", err)
	}
	return font
}

func textSystem(g *gorge.Gorge) {
	assets := asset.FromECS(g)

	light := gorgeutils.NewLight()
	light.SetPosition(0, 0, -1)
	gym := primitive.NewGimbal()
	gym.SetParent(light)

	// Add entity
	gorgeutils.TrackballCamera(g)

	fonts := []*text.Font{
		loadFont(assets, "fonts/open-sans.ttf"),
		loadFont(assets, "fonts/mashanzheng.ttf"),
		loadFont(assets, "fonts/inria.ttf"),
	}

	// Preload textures :/

	textParent := gorge.NewTransform()
	var planes []*primitive.MeshEntity
	var texts []*text.Text
	for i, f := range fonts {
		g.Trigger(asset.AddEvent{Asset: f.Texture})

		p := primitive.Plane() // new plane I guess
		p.Material.Name = "unlit"
		p.Material.SetTexture("albedoMap", f.Texture)
		p.SetPosition(float32(-3+i*3), 0, 0.1).
			Rotate(-math.Pi/2, 0, 0)
		planes = append(planes, p)

		t := text.New(f)
		t.SetParent(textParent).
			SetPosition(0, -0.5+float32(i)*1.5, 1).
			SetScale(.2)
		t.Material.
			SetFloat32("ao", 10).
			SetFloat32("roughness", 0.1).
			SetFloat32("metallic", 0.9)

		texts = append(texts, t)
	}

	log.Println("Adding entities")
	g.Handle(func(gorge.StartEvent) {
		g.AddEntity(gym.Entities)
		g.AddEntity(light)
		for _, p := range planes {
			log.Println("Plane texture:", p.Material.Textures["albedoMap"])
			g.AddEntity(p)
		}
		for _, t := range texts {
			g.AddEntity(t)
		}
	})

	total := float32(0)
	g.Handle(func(evt gorge.UpdateEvent) {
		total += float32(evt)
		s := fmt.Sprintf("together sep The weird BROWN fox quick jumped the gray fox: %.2f", total)
		for i, t := range texts {
			t.SetText(s)
			t.TransformComponent().SetPosition(0, float32(i), 1)

		}
	})

}
