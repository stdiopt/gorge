package resource

import "github.com/stdiopt/gorge"

func System(g *gorge.Context) error {
	FromContext(g)
	return nil
}
