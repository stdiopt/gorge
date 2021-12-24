// Package gorgeapp initializes gorge with default systems for specific platform
package gorgeapp

import (
	"io/fs"
	"log"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gorgeapp/glfw"
	"github.com/stdiopt/gorge/gorgeapp/wasm"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/systems/input"
	"github.com/stdiopt/gorge/systems/render"
	"github.com/stdiopt/gorge/systems/render/renderpl"
	"github.com/stdiopt/gorge/systems/resource"
)

type (
	// WasmOptions options specific for wasm.
	WasmOptions = wasm.Options
	// GLFWOptions glfw options.
	GLFWOptions = glfw.Options
)

// AppFunc func to be used in options.
type AppFunc func(p *App)

// App the bootstrapper.
type App struct {
	inits []gorge.InitFunc

	wasmOptions wasm.Options
	glfwOptions glfw.Options
}

// New creates a new App.
func New(inits ...gorge.InitFunc) *App {
	// Defaults should be on each platforms code but not deep enough so it wont
	// be hard to find
	log.Println("Initializing platform:", Type)
	a := &App{
		// Defaults
		wasmOptions: wasm.Options{
			FS: resource.HTTPFS{BaseURL: ""},
		},
		glfwOptions: glfw.Options{
			FS: resource.FileFS{BasePath: "."},
		},
	}

	// define default rendering pipeline
	defInits := []gorge.InitFunc{
		func(g *gorge.Context) error {
			resource.FromContext(g)
			input.FromContext(g)
			render.FromContext(g)
			renderpl.Default(g)
			// This initializes some global fonts
			gorgeui.FromContext(g)
			return nil
			// gorgeutil.FromContext(g)
		},
		// resource.System,
		// input.System,
		// Disable audio system for android for now, since oto conflicts symbols because of
		// x/mobile/app so it's being added in other platforms
		// audio.System,
		// render.System,
		// renderpl.Default,
		// gorgeui.System,
		// gorgeutil.System,
	}
	a.inits = append(defInits, inits...) // nolint
	return a
}

// Options calls the params appfuncs.
func (a *App) Options(opt ...AppFunc) {
	for _, ofn := range opt {
		ofn(a)
	}
}

// this needs to be here since the other files are build restricted.

// GLFWOpt sets glfw options.
func GLFWOpt(o GLFWOptions) AppFunc {
	return func(a *App) {
		a.glfwOptions = o
	}
}

// GLFWSourcer sets the GLFWSourcer
func GLFWSourcer(s fs.FS) AppFunc {
	return func(a *App) {
		a.glfwOptions.FS = s
	}
}

// WasmOpt sets the wasm options.
func WasmOpt(o WasmOptions) AppFunc {
	return func(a *App) {
		a.wasmOptions = o
	}
}
