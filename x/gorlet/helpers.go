package gorlet

import (
	"github.com/stdiopt/gorge"
)

func Find[T any](e Entity, fn func(v *T)) {
	if e, ok := any(e).(*T); ok {
		fn(e)
	}
	c, ok := e.(gorge.EntityContainer)
	if !ok {
		return
	}
	for _, e := range c.GetEntities() {
		ge, ok := e.(Entity)
		if !ok {
			continue
		}
		Find(ge, fn)
	}
}
