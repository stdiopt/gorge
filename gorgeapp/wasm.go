//go:build (js && ignore) || wasm

package gorgeapp

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gorgeapp/wasm"
	"github.com/stdiopt/gorge/systems/audio"
)

const Type = "wasm"

func (a *App) Run() error {
	inits := append([]gorge.InitFunc{audio.System}, a.inits...)
	return wasm.Run(a.wasmOptions, inits...)
}
