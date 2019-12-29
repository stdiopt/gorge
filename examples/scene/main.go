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
	"github.com/stdiopt/gorge/x/scene"
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
	platform.Start(opt, scene.System, sceneStuff)
}

func sceneStuff(g *gorge.Gorge) {
	sm := scene.ManagerFromGorge(g)
	s := sm.New(scene1)
	s.Init()

	g.Handle(func(gorge.StartEvent) {
		log.Println("Will load scene")
		sm.Load(s)
	})

}

func scene1(s *scene.Scene) {

	gorgeutils.TrackballCamera(s)

	light := gorgeutils.NewLight()
	light.SetPosition(0, 10, 0)
	box := primitive.Cube()
	box.SetPosition(-1, 0, 0)
	// Will be added to scene
	s.AddEntity(
		box,
		light,
	)

	box2 := primitive.Cube()
	box2.SetPosition(1, 0, 0)

	var triggerTime = float32(10)
	s.Handle(func(e gorge.UpdateEvent) {
		box.Rotate(float32(e)*2, 0, 0)
		box2.Rotate(0, float32(e)*4, 0)
		triggerTime -= float32(e)

		if triggerTime < 0 {
			s.AddEntity(box2)
			triggerTime = 1000000
		}

	})

}
