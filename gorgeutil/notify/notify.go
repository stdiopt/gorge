package notify

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/anim"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/gorgeutil"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/x/gorlet"
)

// Severity will affect the color of the notification.
type Severity int

// Severity levels
const (
	SeverityInfo Severity = iota
	SeverityWarn
	SeverityError
)

var ctxKey = struct{ string }{"notify"}

// EventNotify the notification event to be triggered in gorge.
type EventNotify struct {
	Message  string
	Severity Severity
}

type system struct {
	uiCam *gorgeutil.Camera
	ui    *gorgeui.UI
	cards []*card
}

func (s *system) createNotification(e EventNotify) {
	card := gorlet.Build(&card{
		Timeout: 5,
		event:   e,
	})
	card.SetAnchor(1, 1, 1, 1)
	card.SetPivot(1, 0)
	card.SetRect(0, card.Size[1], 20, 0)
	s.ui.Add(card)

	// Move to build or something
	card.Enter = anim.New()
	card.Enter.Start()
	{
		ch := anim.AddChannel(card.Enter, anim.Float32)
		ch.On(func(v float32) { card.SetOpacity(v) })
		ch.SetKey(0, 0)
		ch.SetKey(.5, 1)
	}

	card.Exit = anim.New()
	card.Exit.Start()
	{
		const animTime = .3
		ch := anim.AddChannel(card.Exit, anim.Float32)
		ch.On(func(v float32) {
			card.Position[0] = v
		})
		ch.SetKey(0, 0)
		ch.SetKey(animTime*2, card.Size[0])

		opch := anim.AddChannel(card.Exit, anim.Float32)
		opch.On(func(v float32) {
			card.SetOpacity(v)
		})
		opch.SetKey(0, 1)
		opch.SetKey(animTime, 0)
	}

	s.cards = append(s.cards, card)
}

// System  initializes notification system in gorge.
func System(g *gorge.Context) error {
	FromContext(g)
	return nil
}

// card is a card that can be shown on the screen
type card struct {
	gorlet.Widget[card]
	Enter   *anim.Animation
	Exit    *anim.Animation
	Timeout float32

	// Switch to not notify
	event EventNotify
	color gm.Vec4

	pane *gorlet.WPane
	lbl  *gorlet.WLabel
}

func (w *card) Build(b *gorlet.B) {
	w.SetAnchor(0)
	w.SetRect(0, 0, 20, 0)

	w.color = gm.Color(0, 0, 0, .3)
	switch w.event.Severity {
	case SeverityInfo:
		w.color = gm.Color(0, 0, 0, .3)
	case SeverityWarn:
		w.color = gm.Color(1, 1, 0, .3)
	case SeverityError:
		w.color = gm.Color(1, 0, 0, .3)
	}
	w.pane = b.Pane().SetColor(w.color[:]...)
	w.lbl = b.Label(w.event.Message).
		SetRect(1, 0, 1, 0).
		SetAnchor(0, .5, 1, .5).
		SetPivot(.5).
		SetAutoSize(true)
	// a := widget.AutoHeight(0)
	event.Handle(w, func(gorgeui.EventUpdate) {
		w.Size[1] = w.lbl.ContentSize()[1] + 2
	})
}

func (w *card) SetOpacity(v float32) *card {
	w.pane.SetColor(w.color[0], w.color[1], w.color[2], v*.3)
	w.lbl.SetColor(1, v)
	return w
}
