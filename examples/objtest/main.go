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
	"math"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/asset"
	"github.com/stdiopt/gorge/gorgeutils"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/platform"
	"github.com/stdiopt/gorge/primitive"
	"github.com/stdiopt/gorge/renderer"
)

func main() {
	renderer.ExperimentalSkybox = true
	opt := platform.Options{
		Wasm: platform.WasmOptions{
			AssetLoader: asset.HTTPLoader{BaseURL: "../assets"},
		},
		GLFW: platform.GLFWOptions{
			AssetLoader: asset.FileLoader{BasePath: "/assets"},
		},
	}
	// Setup the asset system with an http loader
	platform.Start(opt, objSystem)
}

func objSystem(g *gorge.Gorge) {
	assets := asset.FromECS(g)
	g.Handle(func(err gorge.ErrorEvent) {
		log.Println("Something errored:", err)
	})

	tex := assets.Texture2D("dog/dog.jpg")

	gorgeutils.TrackballCamera(g)

	light := gorgeutils.NewLight()
	light.SetPosition(0, 20, 0)
	light.Color = m32.Vec3{10, 10, 10}

	//box := primitive.Cube("box1")
	//box.Transform.
	//SetScale(6, 0.2, 4).
	//SetPosition(0, -1, 0)

	//box.Renderable.Material = gorge.NewMaterial("").
	//Set("ao", float32(10))

	// Load here instead
	objData := assets.Mesh("dog/dog.obj").Data()

	objMesh := &gorge.Mesh{MeshLoader: objData}

	// MeshEntity a basic renderable entity
	r1 := &primitive.MeshEntity{
		Transform: *gorge.NewTransform(),
		Renderable: *gorge.NewRenderable("",
			objMesh,
			gorge.NewMaterial("").
				SetFloat32("ao", 10).
				SetTexture("albedoMap", tex),
		),
	}
	r1.SetScale(0.1).
		SetPosition(0, -0.8, 0).
		SetEuler(math.Pi/2, 0, 0)

	r2 := &primitive.MeshEntity{
		Transform: *gorge.NewTransform(),
		Renderable: *gorge.NewRenderable(
			"",
			objMesh,
			gorge.NewMaterial("reflect"),
		),
	}
	r2.Mesh = objMesh
	r2.SetScale(0.1).
		SetPosition(-4, -0.8, 0).
		SetEuler(math.Pi/2, 0, 0)

	r3 := &primitive.MeshEntity{
		Transform: *gorge.NewTransform(),
		Renderable: *gorge.NewRenderable("",
			objMesh,
			gorge.NewMaterial("refract"),
		),
	}
	r3.SetScale(0.1).
		SetPosition(4, -0.8, 0).
		SetEuler(math.Pi/2, 0, 0)

	g.Handle(func(gorge.StartEvent) {
		g.AddEntity(light)
		//m.AddEntity(box)
		g.AddEntity(r1, r2, r3)

	})
}

// Create a mesher
