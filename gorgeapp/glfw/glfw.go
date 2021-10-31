//go:build !js && !android && !mobile

// Package glfw provides initialization for glfw lib
package glfw

import (
	"log"
	"os"
	"runtime"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/input"
	"github.com/stdiopt/gorge/systems/render/gl"
	"github.com/stdiopt/gorge/systems/resource"

	// ogl "github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

// Run will run stuff natively (*nix only maybe)
func Run(opt Options, systems ...interface{}) error {
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

	window, err := glfw.CreateWindow(width, height, "gorge", nil, nil)
	if err != nil {
		return err
	}
	window.MakeContextCurrent()

	glw := &gl.Wrapper{}

	// When running opengl "github.com/go-gl/gl/v4.6-core/gl"
	/*if err := opengl.Init(); err != nil {
		return err
	}*/

	// ogl.Enable(opengl.MULTISAMPLE)
	// ogl.Enable(opengl.PROGRAM_POINT_SIZE)
	// ogl.Disable(opengl.LINE_SMOOTH)

	// This brakes NV
	// ogl.Enable(opengl.POLYGON_SMOOTH)
	// ogl.Hint(opengl.LINE_SMOOTH_HINT, opengl.NICEST)
	// ogl.Hint(opengl.POLYGON_SMOOTH_HINT, opengl.NICEST)

	s := glfwSystem{
		glctx:  glw,
		window: window,
	}

	resourceFS := opt.FS
	if resourceFS == nil {
		resourceFS = os.DirFS(".")
	}

	ggArgs := []interface{}{
		func(g *gorge.Context, res *resource.Context) {
			res.AddFS("/", resourceFS)
			g.PutProp(glw)
		},
		s.System,
	}
	ggArgs = append(ggArgs, systems...)

	g := gorge.New(ggArgs...)
	// bind stuff together
	if err := g.Start(); err != nil {
		return err
	}

	triggerPerFrame := float64(1) / 60 // target fps, usually 60
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
}

func (s *glfwSystem) System(g *gorge.Context, ic *input.Context) error {
	s.input = ic
	s.gorge = g
	s.setupEvents()
	return nil
}

func (s *glfwSystem) setupEvents() {
	s.window.SetSizeCallback(func(_ *glfw.Window, width, height int) {
		s.gorge.SetScreenSize(m32.Vec2{float32(width), float32(height)})
	})

	s.window.SetScrollCallback(
		func(w *glfw.Window, _, yoff float64) {
			x, y := w.GetCursorPos()
			evt := input.EventPointer{
				Type: input.MouseWheel,
				Pointers: map[int]input.PointerData{
					0: {
						DeltaZ: -float32(yoff) * 6,
						Pos:    m32.Vec2{float32(x), float32(y)},
					},
				},
			}
			s.gorge.Trigger(evt) // nolint: errcheck
		},
	)

	s.window.SetMouseButtonCallback(
		func(w *glfw.Window, button glfw.MouseButton, a glfw.Action, _ glfw.ModifierKey) {
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
			evt := input.EventPointer{
				Type:   typ,
				Button: int(button),
				Pointers: map[int]input.PointerData{
					0: {Pos: m32.Vec2{float32(x), float32(y)}},
				},
			}
			s.gorge.Trigger(evt) // nolint: errcheck
		},
	)

	s.window.SetCursorPosCallback(func(_ *glfw.Window, x, y float64) {
		typ := input.MouseMove
		evt := input.EventPointer{
			Type: typ,
			Pointers: map[int]input.PointerData{
				0: {Pos: m32.Vec2{float32(x), float32(y)}},
			},
		}
		s.gorge.Trigger(evt) // nolint: errcheck
	})

	s.window.SetKeyCallback(func(_ *glfw.Window, k glfw.Key, _ int, a glfw.Action, _ glfw.ModifierKey) {
		keyFn := s.input.SetKeyDown
		if a == glfw.Release {
			keyFn = s.input.SetKeyUp
		}
		keyFn(keyMap[k])
	})
}
