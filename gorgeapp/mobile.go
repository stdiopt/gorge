//go:build android || mobile

package gorgeapp

import (
	"github.com/stdiopt/gorge/gorgeapp/mobile"
)

const Type = "mobile"

func (a *App) Run() error {
	return mobile.Run(a.mobileOptions, a.inits...)
}
