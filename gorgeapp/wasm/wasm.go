//go:build (js && ignore) || wasm

// Package wasm provides platform initializations for wasm
package wasm

import (
	"log"
	"strings"
	"syscall/js"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/input"
	"github.com/stdiopt/gorge/systems/render/gl"
	"github.com/stdiopt/gorge/systems/resource"
)

// Start create a premade gorge manager
func Run(opt Options, systems ...interface{}) error {
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
	canvas := El("canvas")

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
	webgl.Call("getExtension", "EXT_color_buffer_float") // possible on chrome

	s := wasmSystem{
		canvas:           canvas,
		glctx:            &gl.Wrapper{Value: webgl},
		CanvasResolution: 1, // TODO: from Opts
	}

	// Default asset loader to .
	resourceFS := opt.FS
	if resourceFS == nil {
		resourceFS = resource.HTTPFS{"."}
	}
	ggArgs := []interface{}{
		func(g *gorge.Context, res *resource.Context) {
			res.AddFS("/", resourceFS)
		},
		s.System,
	}
	ggArgs = append(ggArgs, systems...)

	g := gorge.New(ggArgs...)

	return g.Run()
}

type wasmSystem struct {
	gorge  *gorge.Context
	input  *input.Context
	glctx  *gl.Wrapper
	canvas js.Value

	CanvasResolution float64
	// Width, Height    float64
}

func (s *wasmSystem) HandleEvent(v event.Event) {
	switch v.(type) {
	case gorge.EventStart:
		s.setupEvents()
	case gorge.EventAfterStart:
		var prevFrameTime float64 = js.Global().Get("performance").Call("now").Float() / 1000
		var ticker js.Func
		ticker = js.FuncOf(func(t js.Value, args []js.Value) interface{} {
			s.checkCanvasSize()
			totalTime := args[0].Float() / 1000
			dtSec := totalTime - prevFrameTime

			s.gorge.Update(float32(dtSec))

			prevFrameTime = totalTime
			js.Global().Call("requestAnimationFrame", ticker)
			return nil
		})
		js.Global().Call("requestAnimationFrame", ticker)
	}
}

func (s *wasmSystem) System(g *gorge.Context, ic *input.Context) error {
	s.gorge = g
	s.input = ic
	g.PutProp(s.glctx)
	s.checkCanvasSize()
	g.Handle(s)
	return nil
}

func (s *wasmSystem) checkCanvasSize() {
	size := s.canvas.Call("getBoundingClientRect")
	w := float32(size.Get("width").Float() * s.CanvasResolution)
	h := float32(size.Get("height").Float() * s.CanvasResolution)
	ss := s.gorge.ScreenSize()
	if w != ss[0] || h != ss[1] {
		s.canvas.Set("width", w)
		s.canvas.Set("height", h)

		s.gorge.SetScreenSize(m32.Vec2{float32(w), float32(h)})
	}
}

var evtMap = map[string]input.PointerType{
	"mousedown":   input.MouseDown,
	"mouseup":     input.MouseUp,
	"mousemove":   input.MouseMove,
	"touchstart":  input.PointerDown,
	"touchmove":   input.PointerMove,
	"touchcancel": input.PointerCancel,
	"touchend":    input.PointerEnd,
}

func (s *wasmSystem) setupEvents() {
	ptrEvent := js.FuncOf(s.handleEventPointer)
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

	code := evt.Get("code").String()
	etype := evt.Get("type").String()

	if code == "F12" {
		return nil
	}
	ikey, ok := keyMap[code]
	if !ok {
		log.Println("Key not mapped:", code, ikey)
		js.Global().Get("console").Call("log", evt)
	}

	switch etype {
	case "keydown":
		s.input.SetKeyDown(ikey)
	case "keyup":
		s.input.SetKeyUp(ikey)
	}

	return nil
}

// TODO: pointers are currently relative to window, should be to canvas
func (s *wasmSystem) handleEventPointer(t js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("preventDefault")
	etype := evt.Get("type").String()

	cevt := input.EventPointer{}

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
