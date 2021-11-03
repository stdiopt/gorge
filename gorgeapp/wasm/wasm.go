//go:build (js && ignore) || wasm

// Package wasm provides platform initializations for wasm
package wasm

import (
	"log"
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
	// Handle platform specific events here.

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

func (s *wasmSystem) setupEvents() {
	keyEvent := js.FuncOf(s.handleKeyEvents)
	mouseEvent := js.FuncOf(s.handleMouseEvents)
	touchEvent := js.FuncOf(s.handleTouchEvents)

	js.Global().Call("addEventListener", "keydown", keyEvent)
	js.Global().Call("addEventListener", "keyup", keyEvent)
	js.Global().Call("addEventListener", "keypress", keyEvent)

	s.canvas.Call("addEventListener", "mousedown", mouseEvent)
	s.canvas.Call("addEventListener", "mouseup", mouseEvent)
	s.canvas.Call("addEventListener", "mousemove", mouseEvent)
	s.canvas.Call("addEventListener", "contextmenu", mouseEvent)
	s.canvas.Call("addEventListener", "wheel", mouseEvent)

	s.canvas.Call("addEventListener", "touchstart", touchEvent)
	s.canvas.Call("addEventListener", "touchmove", touchEvent)
	s.canvas.Call("addEventListener", "touchcancel", touchEvent)
	s.canvas.Call("addEventListener", "touchend", touchEvent)
}

func (s *wasmSystem) handleKeyEvents(t js.Value, args []js.Value) interface{} {
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
		s.input.SetKeyState(ikey, input.ActionDown)
	case "keyup":
		s.input.SetKeyState(ikey, input.ActionUp)
	}

	return nil
}

func (s *wasmSystem) handleMouseEvents(t js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("preventDefault")
	etype := evt.Get("type").String()

	switch etype {
	case "wheel":
		// Maybe grab X too but this is about mousewheel
		s.input.SetScrollDelta(float32(evt.Get("deltaY").Float()))
	case "mousemove":
		s.input.SetCursorPosition(
			float32(evt.Get("pageX").Float()*s.CanvasResolution),
			float32(evt.Get("pageY").Float()*s.CanvasResolution),
		)
	case "contextmenu":
		s.input.SetMouseButtonState(input.MouseRight, input.ActionDown)
	case "mousedown":
		btn := evt.Get("button").Int()
		gbtn, ok := mousebtnMap[btn]
		if !ok {
			log.Println("Mouse button not mapped:", gbtn)
		}
		s.input.SetMouseButtonState(gbtn, input.ActionDown)
	case "mouseup":
		btn := evt.Get("button").Int()
		gbtn, ok := mousebtnMap[btn]
		if !ok {
			log.Println("Mouse button not mapped:", gbtn)
		}
		s.input.SetMouseButtonState(gbtn, input.ActionUp)
	}

	return nil
}

func (s *wasmSystem) handleTouchEvents(t js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("preventDefault")
	etype := evt.Get("type").String()

	var gtyp input.PointerType
	switch etype {
	case "touchstart":
		gtyp = input.PointerDown
	case "touchmove":
		gtyp = input.PointerMove
	case "touchcancel":
		gtyp = input.PointerCancel
	case "touchend":
		gtyp = input.PointerEnd
	}
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
	s.gorge.Trigger(input.EventPointer{
		Type:     gtyp,
		Pointers: pts,
	})
	return nil
}
