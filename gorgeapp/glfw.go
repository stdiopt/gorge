//go:build !js && !android && !mobile

package gorgeapp

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gorgeapp/glfw"
	"github.com/stdiopt/gorge/systems/audio"
)

// Type shows the plataform it is being run on
const Type = "glfw"

// Run the glfw app
func (a *App) Run() error {
	inits := append([]gorge.InitFunc{
		func(g *gorge.Context) {
			audio.FromContext(g)
		},
	},
		a.inits...,
	)
	return glfw.Run(a.glfwOptions, inits...)
}
