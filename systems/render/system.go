package render

import (
	"github.com/stdiopt/gorge"
)

func System(g *gorge.Context) {
	FromContext(g) // lazy init
}
