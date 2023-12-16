package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type WQuad struct {
	Widget[WQuad]

	ent *gEntity
}

func Quad(c ...gorge.Entity) *WQuad {
	return Build(&WQuad{}).Add(c...)
}

func (w *WQuad) Build(b *B) {
	w.ent = newEntity("quad", quadMesh())
	w.Add(w.ent)
	event.Handle(w, func(gorgeui.EventUpdate) {
		sz := w.ContentSize()
		w.ent.SetScale(sz[0], sz[1])
	})
}

func (w *WQuad) SetColor(vs ...float32) {
	c := gm.Color(vs...)
	w.ent.Colorable().SetColor(c[0], c[1], c[2], c[3])
}
