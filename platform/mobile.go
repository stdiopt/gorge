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

// +build android

package platform

import (
	"log"
	"time"

	"github.com/stdiopt/gorge"
	gorgegl "github.com/stdiopt/gorge/gl"
	"github.com/stdiopt/gorge/input"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/renderer"
	"github.com/stdiopt/gorge/resource"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

// Type of the platform
const Type = "mobile"

// Mobile stars a mobile up
func Mobile(path string, systems ...gorge.SystemFunc) {

	resourceLoader := resource.FileLoader{BasePath: path}

	ggArgs := []gorge.SystemFunc{
		resource.NewSystem(resourceLoader),
		input.System,
	}
	ggArgs = append(ggArgs, systems...)

	g := gorge.New(ggArgs...)

	// Initialize stuff behind

	var glctx gl.Context3
	// Run in background?
	app.Main(func(a app.App) {
		mark := time.Now()
		for e := range a.Events() {

			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					var ok bool
					glctx, ok = e.DrawContext.(gl.Context3)
					if !ok {
						log.Println("gles3 not supported?")
						_, ok = e.DrawContext.(gl.Context)
						if !ok {
							log.Println("gles2 not supported?")
							return
						}
						log.Println("gles2 supported - but we waant 3 anyway :shrug:")
						g.Close()
						return
					}

					gw := gorgegl.Wrapper{glctx}
					g.Persist(gw)
					renderer.System(gw, g)

					g.Start()
				case lifecycle.CrossOff:
					g.Close()
				}

			case paint.Event:
				if glctx == nil {
					continue
				}

				now := time.Now()
				sub := float32(now.Sub(mark)) / 1000000000
				g.UpdateNow(sub)

				a.Publish()
				a.Send(paint.Event{}) // keep animating
				mark = now
			case mouse.Event:
				log.Println("Mouse event:", e.Button)
				switch e.Button {
				case mouse.ButtonWheelUp:
					g.Trigger(input.PointerEvent{
						Type: input.MouseWheel,
						Pointers: map[int]input.PointerData{
							0: {
								DeltaZ: -1,
								Pos:    m32.Vec2{e.X, e.Y},
							},
						},
					})
				case mouse.ButtonWheelDown:
					g.Trigger(input.PointerEvent{
						Type: input.MouseWheel,
						Pointers: map[int]input.PointerData{
							0: {
								DeltaZ: 1,
								Pos:    m32.Vec2{e.X, e.Y},
							},
						},
					})
				}
			case touch.Event:
				log.Println("TouchEvent:", e)
				var typ input.PointerType = input.MouseMove
				switch e.Type {
				case touch.TypeBegin:
					typ = input.MouseDown
				case touch.TypeEnd:
					typ = input.MouseUp
				}

				g.Trigger(input.PointerEvent{
					Type: typ,
					Pointers: map[int]input.PointerData{
						0: {Pos: m32.Vec2{e.X, e.Y}},
					},
				})

			case size.Event:
				b := e.Bounds()
				g.Persist(gorge.ResizeEvent{float32(b.Dx()), float32(b.Dy())})
			}

		}
	})
}
