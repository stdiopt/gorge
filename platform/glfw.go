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

// +build !js, !wasm

package platform

import (
	"log"
	"runtime"
	"time"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gl"
	"github.com/stdiopt/gorge/input"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/renderer"
	"github.com/stdiopt/gorge/resource"

	opengl "github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Type of the platform
const Type = "glfw"

func init() {
	runtime.LockOSThread()
}

// Start will run stuff natively (*nix only maybe)
func Start(opt Options, systems ...gorge.SystemFunc) {

	log.Println("Init GLFW")
	const width, height = 800, 600

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "gorge", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := opengl.Init(); err != nil {
		log.Fatalf("failed to initialize glfw")
	}
	version := opengl.GoStr(opengl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	glfw.WindowHint(glfw.Samples, 4)
	opengl.Enable(opengl.MULTISAMPLE)

	s := glfwSystem{
		glctx:  &gl.Wrapper{},
		window: window,
	}

	resourceLoader := opt.GLFW.Loader
	if resourceLoader == nil {
		resourceLoader = resource.FileLoader{BasePath: "."}
	}

	ggArgs := []gorge.SystemFunc{
		resource.NewSystem(resourceLoader),
		input.System,
		s.Init,
		renderer.System,
	}
	ggArgs = append(ggArgs, systems...)

	// bind stuff together
	g := gorge.New(ggArgs...)
	g.Start()

	///////////////////////////////////////////////////////////////////////////
	// Ticker here
	// Might be moved to elsewhere, JS have requestAnimationFrame which is
	// handled here
	///////////////
	mark := glfw.GetTime()
	for !s.window.ShouldClose() {
		now := glfw.GetTime()
		elapsed := now - mark
		mark = now
		g.UpdateNow(float32(elapsed))

		s.window.SwapBuffers()
		glfw.PollEvents()
		time.Sleep(time.Second / 60)
	}
	g.Close()

}

type glfwSystem struct {
	gorge  *gorge.Gorge
	glctx  *gl.Wrapper
	window *glfw.Window
}

func (s *glfwSystem) Init(g *gorge.Gorge) {
	g.Persist(s.glctx)
	g.Handle(func(gorge.StartEvent) {
		w, h := s.window.GetSize()
		g.Persist(gorge.ResizeEvent{float32(w), float32(h)})
		s.setupEvents()
	})
	s.gorge = g
}

func (s *glfwSystem) setupEvents() {
	s.window.SetSizeCallback(func(w *glfw.Window, width, height int) {
		s.gorge.Persist(gorge.ResizeEvent{float32(width), float32(height)})
	})
	s.window.SetScrollCallback(
		func(w *glfw.Window, xoff, yoff float64) {
			x, y := w.GetCursorPos()
			evt := input.PointerEvent{
				Type: input.MouseWheel,
				Pointers: map[int]input.PointerData{
					0: {
						DeltaZ: -float32(yoff) * 6,
						Pos:    m32.Vec2{float32(x), float32(y)},
					},
				},
			}
			s.gorge.Trigger(evt)

		},
	)
	s.window.SetMouseButtonCallback(
		func(w *glfw.Window, button glfw.MouseButton, a glfw.Action, mog glfw.ModifierKey) {
			if button != glfw.MouseButton1 {
				return
			}
			var typ input.PointerType
			switch a {
			case glfw.Release:
				typ = input.MouseUp
			case glfw.Press:
				typ = input.MouseDown
			default:
				return
			}
			x, y := w.GetCursorPos()
			evt := input.PointerEvent{
				Type: typ,
				Pointers: map[int]input.PointerData{
					0: {Pos: m32.Vec2{float32(x), float32(y)}},
				},
			}
			s.gorge.Trigger(evt)
		},
	)
	s.window.SetCursorPosCallback(func(w *glfw.Window, x, y float64) {
		typ := input.MouseMove
		evt := input.PointerEvent{
			Type: typ,
			Pointers: map[int]input.PointerData{
				0: {Pos: m32.Vec2{float32(x), float32(y)}},
			},
		}
		s.gorge.Trigger(evt)
	})

}
