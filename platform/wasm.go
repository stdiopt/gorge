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

// +build js, wasm

package platform

import (
	"strings"
	"syscall/js"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gl"
	"github.com/stdiopt/gorge/input"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/platform/dom"
	"github.com/stdiopt/gorge/renderer"
	"github.com/stdiopt/gorge/resource"
)

// Type of the platform
const Type = "wasm"

// Alias from dom
var (
	Document = dom.Document
	Body     = dom.Body
	El       = dom.El
)

type (
	//Attr attributes
	Attr = dom.Attr
	//Text dom text
	Text = dom.Text
)

// Start create a premade gorge manager
func Start(opt Options, systems ...gorge.SystemFunc) {
	Document.Get("head").Set("innerHTML", `
	<meta name="mobile-web-app-capable" content="yes">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title> gorge - go reduced game engine</title>
	<style>
		* {box-sizing: border-box;}
		body{ height: 100vh; margin:0; padding:0; }
		canvas { position:fixed; top: 0px; height:100%; width: 100%; }
		#fs-btn { z-index:10; position:fixed; top:5px;right:5px; }
	</style>
	`)
	fullScreenBtn := El("button", Attr{"id": "fs-btn"}, Text("fullscreen"))
	fullScreenBtn.Call("addEventListener", "click", js.FuncOf(func(t js.Value, args []js.Value) interface{} {
		Body.Call("requestFullscreen")
		return nil
	}))
	canvas := dom.El("canvas")

	Body.Call("appendChild", fullScreenBtn)
	Body.Call("appendChild", canvas)

	// Get gl Context from WebGL thingy
	ctxOpt := map[string]interface{}{
		"preserveDrawingBuffer": true,
		"antialias":             true,
	}
	webgl := canvas.Call("getContext", "webgl2", ctxOpt)
	js.Global().Get("console").Call("log", webgl.Call("getSupportedExtensions"))
	webgl.Call("getExtension", "EXT_texture_filter_anisotropic")

	s := wasmSystem{
		canvas:           canvas,
		glctx:            &gl.Wrapper{Value: webgl},
		CanvasResolution: 1, // TODO: from Opts
	}

	// Default asset loader to .
	resourceLoader := opt.Wasm.Loader
	if resourceLoader == nil {
		resourceLoader = resource.HTTPLoader{"."}
	}
	ggArgs := []gorge.SystemFunc{
		resource.NewSystem(resourceLoader),
		input.System,
		s.init,
		renderer.System, // will load glctx and asset
	}
	ggArgs = append(ggArgs, systems...)

	// bind stuff together
	g := gorge.New(ggArgs...)
	g.Run()

}

type wasmSystem struct {
	gorge  *gorge.Gorge
	glctx  *gl.Wrapper
	canvas js.Value

	CanvasResolution float64
	Width, Height    float64
}

func (s *wasmSystem) init(g *gorge.Gorge) {
	g.Persist(s.glctx)
	g.Handle(func(gorge.StartEvent) {
		s.checkCanvasSize()
		s.setupEvents()
	})
	// XXX: The looper happens after start, since it can have race condition on start
	// go wasm is weird
	g.Handle(func(gorge.AfterStartEvent) {
		var prevFrameTime float64
		var ticker js.Func
		ticker = js.FuncOf(func(t js.Value, args []js.Value) interface{} {
			s.checkCanvasSize()
			dt := args[0].Float()
			dtSec := (dt - prevFrameTime) / 1000

			g.UpdateNow(float32(dtSec))

			prevFrameTime = dt
			js.Global().Call("requestAnimationFrame", ticker)
			return nil
		})

		js.Global().Call("requestAnimationFrame", ticker)
	})
	s.gorge = g
}

func (s *wasmSystem) checkCanvasSize() {
	size := s.canvas.Call("getBoundingClientRect")
	w := size.Get("width").Float() * s.CanvasResolution
	h := size.Get("height").Float() * s.CanvasResolution
	if w != s.Width || h != s.Height {
		s.canvas.Set("width", w)
		s.canvas.Set("height", h)
		s.Width, s.Height = w, h
		s.gorge.Persist(gorge.ResizeEvent{float32(w), float32(h)})
	}
}

var (
	evtMap = map[string]input.PointerType{
		"mousedown":   input.MouseDown,
		"mouseup":     input.MouseUp,
		"mousemove":   input.MouseMove,
		"touchstart":  input.PointerDown,
		"touchmove":   input.PointerMove,
		"touchcancel": input.PointerCancel,
		"touchend":    input.PointerEnd,
	}
)

func (s *wasmSystem) setupEvents() {

	ptrEvent := js.FuncOf(s.handlePointerEvent)
	keyEvent := js.FuncOf(s.handleKeyEvent)

	for k := range evtMap {
		s.canvas.Call("addEventListener", k, ptrEvent)
	}
	s.canvas.Call("addEventListener", "wheel", ptrEvent)
	js.Global().Call("addEventListener", "keydown", keyEvent)
	js.Global().Call("addEventListener", "keyup", keyEvent)
	js.Global().Call("addEventListener", "keypress", keyEvent)

}

func (s *wasmSystem) handleKeyEvent(t js.Value, args []js.Value) interface{} {
	evt := args[0]

	key := evt.Get("key").String()
	etype := evt.Get("type").String()

	if key == "F12" {
		return nil
	}
	evt.Call("preventDefault")

	switch etype {
	case "keydown":
		input.FromECS(s.gorge).KeyDown(key)
	case "keyup":
		input.FromECS(s.gorge).KeyUp(key)
	}
	//s.manager.Trigger(input.KeyEvent{keyEvtMap[etype], code})

	return nil
}

// TODO: pointers are currently relative to window, should be to canvas
func (s *wasmSystem) handlePointerEvent(t js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("preventDefault")
	etype := evt.Get("type").String()

	cevt := input.PointerEvent{}

	switch {
	case strings.HasPrefix(etype, "wheel"):
		cevt.Type = input.MouseWheel
		cevt.Pointers = map[int]input.PointerData{
			0: {
				DeltaZ: float32(evt.Get("deltaY").Float()),
				Pos: m32.Vec2{
					float32(evt.Get("pageX").Float() * s.CanvasResolution),
					float32(evt.Get("pageY").Float() * s.CanvasResolution),
				},
			},
		}
	case strings.HasPrefix(etype, "mouse"):
		cevt.Type = evtMap[etype]
		cevt.Pointers = map[int]input.PointerData{
			0: {
				Pos: m32.Vec2{
					float32(evt.Get("pageX").Float() * s.CanvasResolution),
					float32(evt.Get("pageY").Float() * s.CanvasResolution),
				},
			},
		}
	case strings.HasPrefix(etype, "touch"):
		cevt.Type = evtMap[etype]
		pts := map[int]input.PointerData{}

		touches := evt.Get("changedTouches")
		for i := 0; i < touches.Length(); i++ {
			t := touches.Index(i)
			id := t.Get("identifier").Int()
			pts[id] = input.PointerData{
				Pos: m32.Vec2{
					float32(t.Get("pageX").Float() * s.CanvasResolution),
					float32(t.Get("pageY").Float() * s.CanvasResolution),
				},
			}
		}
		touches = evt.Get("touches")
		for i := 0; i < touches.Length(); i++ {
			t := touches.Index(i)
			id := t.Get("identifier").Int()
			pts[id] = input.PointerData{
				Pos: m32.Vec2{
					float32(t.Get("pageX").Float() * s.CanvasResolution),
					float32(t.Get("pageY").Float() * s.CanvasResolution),
				},
			}
		}
		cevt.Pointers = pts
	}

	s.gorge.Trigger(cevt)

	return nil
}
