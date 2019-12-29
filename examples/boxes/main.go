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
	"log"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/asset"
	"github.com/stdiopt/gorge/gorgeutils"
	"github.com/stdiopt/gorge/platform"
	"github.com/stdiopt/gorge/primitive"
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
	// Setup the asset system with an http loader
	platform.Start(opt, boxesSystem)
}

func boxesSystem(g *gorge.Gorge) {

	assets := asset.FromECS(g)
	log.Println("asset system:", assets)

	woodTex := assets.Texture2D("wood.png")
	wasmTex := assets.Texture2D("wasm.png")

	gorgeutils.TrackballCamera(g)

	light := gorgeutils.NewLight()
	light.SetPosition(0, 5, 3)

	box := primitive.Cube()
	box.Material.
		SetFloat32("metallic", 0.2).
		SetFloat32("roughness", 0.4).
		SetFloat32("ao", float32(10)).
		SetTexture("albedoMap", woodTex)

	box2 := primitive.Cube()
	box2.Material.
		Set("ao", float32(100)).
		SetTexture("albedoMap", wasmTex)
	box2.
		SetParent(box).
		SetScale(0.2, 0.2, 0.2).
		SetPosition(2, 0, 0)

	// Add entities when start event was fired
	g.Handle(func(gorge.StartEvent) {
		g.AddEntity(light, box, box2)
	})

	// update event at every frameupdate
	g.Handle(func(evt gorge.UpdateEvent) {
		dt := float32(evt)
		box.Rotate(dt/2, dt, 0)
		box2.Rotate(dt*2, 0, dt*2)
	})
}
