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
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gorgeutils"
	"github.com/stdiopt/gorge/platform"
	"github.com/stdiopt/gorge/primitive"
)

func main() {
	opt := platform.Options{}
	platform.Start(opt, func(g *gorge.Gorge) {
		gorgeutils.TrackballCamera(g)

		light := gorgeutils.NewLight()
		light.SetPosition(0, 10, -4)

		cube := primitive.Cube()

		g.Handle(func(gorge.StartEvent) {
			g.AddEntity(light)
			g.AddEntity(cube)
		})
		g.Handle(func(dt gorge.UpdateEvent) {
			cube.Rotate(0, 1*float32(dt), 0)
		})
	})
}
