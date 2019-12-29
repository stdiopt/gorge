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

// +build wasm, js

package gophers

import (
	"time"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/input"
	"github.com/stdiopt/gorge/platform/dom"
)

// StatSystem thing
func StatSystem(g *gorge.Gorge) {

	// Dom, wasm specific
	statDom := dom.El("pre",
		dom.Attr{"style": `
			position:fixed;
			padding:10px;
			margin:10px;
			background:rgba(0,0,0,0.8);
			color:white;
			pointer-events:none;
			line-height: 1.5em;
			font-size:0.7em;
			z-index:90;
		`},
	)
	dom.Body.Call("appendChild", statDom)

	// Profiling
	go func() {
		for {
			statStr := statUpdate(g)

			statDom.Set("innerHTML", statStr)

			time.Sleep(time.Second * 3)
		}
	}()

	g.Handle(func(evt input.KeyEvent) {
		if evt.Type == input.KeyUp && evt.Key == "F10" {
			statUpdate(g)
		}
	}).Describe("stat key")

}
