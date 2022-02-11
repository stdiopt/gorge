package gorgeapp

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/input"
	"github.com/stdiopt/gorge/systems/resource"
)

type (
	gorgeContext    = gorge.Context
	inputContext    = input.Context
	resourceContext = resource.Context
)

type AppContextFunc func(*Context)

type Context struct {
	*gorgeContext
	*inputContext
	*resourceContext
}
