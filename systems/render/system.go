package render

import (
	"github.com/stdiopt/gorge"
)

func System(g *gorge.Context) error {
	FromContext(g) // lazy init
	return nil
}
