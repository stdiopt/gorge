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
	"github.com/stdiopt/gorge/examples/gophers"
	"github.com/stdiopt/gorge/examples/gophers/assets"
	"github.com/stdiopt/gorge/platform"
	"github.com/stdiopt/gorge/renderer"
	"github.com/stdiopt/gorge/resource"
)

func main() {
	log.SetFlags(0)

	renderer.ExperimentalSkybox = true

	loader := resource.EmbedLoader{Data: assets.Data}
	opt := platform.Options{
		Wasm: platform.WasmOptions{Loader: loader},
		GLFW: platform.GLFWOptions{Loader: loader},
	}

	platform.Start(opt, gophers.System, errorReporter)

}

func errorReporter(g *gorge.Gorge) {
	g.Handle(func(e gorge.ErrorEvent) {
		log.Printf("\033[01;31m%v\033[0m", e.Err)
	})
}
