//go:build !js && !android && !mobile

// Package glfw provides initialization for glfw lib
package glfw

import (
	"log"
	"runtime"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/input"
	"github.com/stdiopt/gorge/systems/render/gl"
	"github.com/stdiopt/gorge/systems/resource"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

// Run will run stuff natively (*nix only maybe)
func Run(opt Options, systems ...gorge.InitFunc) error {
	log.Println("Init GLFW")
	const width, height = 800, 600

	if err := glfw.Init(); err != nil {
		log.Fatal("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.StencilBits, 8)

	window, err := glfw.CreateWindow(width, height, "gorge", nil, nil)
	if err != nil {
		return err
	}
	window.MakeContextCurrent()

	// When running opengl "github.com/go-gl/gl/v4.6-core/gl"
	// if err := ogl.Init(); err != nil {
	// 	return err
	// }

	// ogl.Enable(ogl.MULTISAMPLE)
	// ogl.Enable(ogl.PROGRAM_POINT_SIZE)
	// ogl.Disable(opengl.LINE_SMOOTH)

	// This brakes NV
	// ogl.Enable(opengl.POLYGON_SMOOTH)
	// ogl.Hint(opengl.LINE_SMOOTH_HINT, opengl.NICEST)
	// ogl.Hint(opengl.POLYGON_SMOOTH_HINT, opengl.NICEST)

	glw := &gl.Wrapper{}
	gl.Init(glw)

	s := glfwSystem{
		glctx:  glw,
		window: window,
	}

	resourceFS := opt.FS
	if resourceFS == nil {
		resourceFS = RootFS{}
	}

	ggArgs := []gorge.InitFunc{
		func(g *gorge.Context) {
			res := resource.FromContext(g)
			res.AddFS("/", resourceFS)
		},
		s.System,
	}
	ggArgs = append(ggArgs, systems...)

	g := gorge.New(ggArgs...)

	// bind stuff together
	if err := g.Start(); err != nil {
		return err
	}

	triggerPerFrame := float64(1) / 140 // target fps, usually 60
	frameTimeLeft := float64(0)
	lastFrame := float64(0)
	// Ticker here
	mark := glfw.GetTime()
	lastFrame = mark

	for !s.window.ShouldClose() {
		glfw.PollEvents()
		now := glfw.GetTime()
		elapsed := now - mark
		mark = now

		frameTimeLeft -= elapsed
		if frameTimeLeft < 0 {
			g.Update(float32(now - lastFrame))
			s.window.SwapBuffers()

			frameTimeLeft = triggerPerFrame
			lastFrame = now
		}
	}
	g.Close()
	return nil
}

type glfwSystem struct {
	gorge *gorge.Context
	input *input.Context

	glctx  *gl.Wrapper
	window *glfw.Window

	// this should be in a common thing
	cursorRelative bool
}

func (s *glfwSystem) System(g *gorge.Context) {
	s.input = input.FromContext(g)
	s.gorge = g
	s.setupEvents()

	event.Handle(g, func(e gorge.EventCursorRelative) {
		if e { // set cursor to center on turn on
			sx, sy := s.window.GetSize()
			cx, cy := float64(sx/2), float64(sy/2)
			s.window.SetCursorPos(cx, cy)
		}
		s.cursorRelative = bool(e)
	})
	event.Handle(g, func(e gorge.EventCursorHidden) {
		if e {
			s.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
		} else {
			s.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		}
	})
	// Temporary cursor experiment:
	{
		arrowCursor := glfw.CreateStandardCursor(glfw.ArrowCursor)
		handCursor := glfw.CreateStandardCursor(glfw.HandCursor)
		event.Handle(g, func(e gorge.EventCursor) {
			switch gorge.CursorType(e) {
			case gorge.CursorArrow:
				s.window.SetCursor(arrowCursor)
			case gorge.CursorHand:
				s.window.SetCursor(handCursor)
			}
		})
	}
}

func (s *glfwSystem) setupEvents() {
	s.window.SetSizeCallback(func(_ *glfw.Window, width, height int) {
		s.gorge.SetScreenSize(gm.Vec2{float32(width), float32(height)})
	})

	s.window.SetScrollCallback(
		func(_ *glfw.Window, xoff, yoff float64) {
			s.input.SetScrollDelta(gm.Vec2{float32(xoff) * 6, -float32(yoff) * 6})
		},
	)

	s.window.SetMouseButtonCallback(
		func(_ *glfw.Window, button glfw.MouseButton, a glfw.Action, _ glfw.ModifierKey) {
			gbtn, ok := mousebtnMap[button]
			if !ok {
				log.Println("Not mapped", button)
			}
			switch a {
			case glfw.Press, glfw.Repeat:
				s.input.SetMouseButtonState(gbtn, input.ActionDown)
			case glfw.Release:
				s.input.SetMouseButtonState(gbtn, input.ActionUp)
			default:
				return
			}
		},
	)

	// Start in center anyway
	// sx, sy := s.window.GetSize()
	// cx, cy := float64(sx/2), float64(sy/2)
	// s.window.SetCursorPos(cx, cy)
	// s.input.SetCursorPosition(gm.Vec2{float32(cx), float32(cy)})

	s.window.SetCursorPosCallback(func(w *glfw.Window, x, y float64) {
		if !s.cursorRelative {
			s.input.SetCursorPosition(gm.Vec2{float32(x), float32(y)})
			return
		}
		sx, sy := w.GetSize()
		cx, cy := float64(sx/2), float64(sy/2)
		s.input.SetCursorDelta(gm.Vec2{float32(x - cx), float32(y - cy)})
		w.SetCursorPos(cx, cy)
	})
	s.window.SetCharCallback(func(_ *glfw.Window, char rune) {
		s.input.SetKeyChar(char)
	})
	s.window.SetKeyCallback(func(_ *glfw.Window, k glfw.Key, _ int, a glfw.Action, _ glfw.ModifierKey) {
		gkey, ok := keyMap[k]
		if !ok {
			log.Println("Key not mapped:", k, gkey)
		}
		switch a {
		case glfw.Press, glfw.Repeat: // Maybe state hold?
			s.input.SetKeyState(keyMap[k], input.ActionDown)
		case glfw.Release:
			s.input.SetKeyState(keyMap[k], input.ActionUp)
		}
	})
}
