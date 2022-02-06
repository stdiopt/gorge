package notify

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/anim"
	"github.com/stdiopt/gorge/gorgeutil"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/systems/gorgeui/gorlet"
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
	w := gorlet.Create(notifyWidget(e))
	w.SetAnchor(1, 1, 1, 1)
	w.SetRect(-21, w.Size[1], 20, 0)
	s.ui.Add(w)

	card := &card{
		Widget:  w,
		Timeout: 5,
	}
	card.Enter = anim.New()
	card.Enter.Start()
	{
		ch := anim.AddChannel(card.Enter, anim.Float32)
		ch.On(func(v float32) { w.Set("opacity", v) })
		ch.SetKey(0, 0)
		ch.SetKey(.5, 1)
	}

	card.Exit = anim.New()
	card.Exit.Start()
	{
		const animTime = .3
		ch := anim.AddChannel(card.Exit, anim.Float32)
		ch.On(func(v float32) {
			w.Position[0] = v
		})
		ch.SetKey(0, -w.Size[0]-1)
		ch.SetKey(animTime*2, 0)

		opch := anim.AddChannel(card.Exit, anim.Float32)
		opch.On(func(v float32) {
			w.Set("opacity", v)
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
	Enter *anim.Animation
	Exit  *anim.Animation

	Widget  *gorlet.Entity
	Timeout float32
}

func notifyWidget(e EventNotify) gorlet.Func {
	return func(b *gorlet.B) {
		color := gm.Color(0, 0, 0, .3)
		switch e.Severity {
		case SeverityInfo:
			color = gm.Color(0, 0, 0, .3)
		case SeverityWarn:
			color = gm.Color(1, 1, 0, .3)
		case SeverityError:
			color = gm.Color(1, 0, 0, .3)
		}
		b.Root().SetLayout(gorlet.AutoHeight(0))
		// Initial
		b.Root().SetAnchor(1)
		b.Root().SetRect(0, 0, 20, 0)

		b.Use("color", color)
		b.UseAnchor(0, 0, 1, 0)
		p := b.BeginPanel(gorlet.AutoHeight(2))
		b.UseAnchor(0, .5, 1, .5)
		b.UsePivot(.5)
		b.UseRect(1, 0, 1, 0)
		b.Use("autoSize", true)
		l := b.Label(e.Message)
		// TODO: This needs to be FIXED since we need to recaculate
		// stuff internally before showing in the screen
		b.EndPanel()

		gorlet.Observe(b, "opacity", func(v float32) {
			p.Set("color", gm.Color(
				color[0],
				color[1],
				color[2],
				v*.3,
			))
			l.Set("textColor", gm.Color(1, v))
		})
	}
}
