package notify

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/anim"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/gorgeutil"
	"github.com/stdiopt/gorge/m32"
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

// EventNotify the notification event to be triggered in gorge.
type EventNotify struct {
	Message  string
	Severity Severity
}

func Info(g *gorge.Context, s string) {
	g.Trigger(EventNotify{
		Message:  s,
		Severity: SeverityInfo,
	})
}

// System will handle notification events and create cards.
func System(g *gorge.Context) error {
	u := gorgeutil.FromContext(g)
	uiCam := u.UICamera()
	// uiCam.CullMask |= gorge.CullMaskUIDebug
	// gorgeui.FromContext(g).Debug = gorgeui.DebugRects
	ui := u.UI(uiCam)

	cards := []*card{}
	g.HandleFunc(func(e event.Event) {
		switch e := e.(type) {
		case EventNotify:
			// Should build a card

			w := gorlet.Create(notifyWidget(e))
			w.SetAnchor(1, 1, 1, 1)
			w.SetRect(-21, w.Dim[1], 20, 0)
			ui.Add(w)

			card := &card{
				Widget:  w,
				Timeout: 5,
			}
			card.Enter = anim.New()
			card.Enter.Start()
			{
				ch := anim.NewChannelFuncf32(func(v float32) {
					w.Set("opacity", v)
				})
				ch.SetKey(0, 0)
				ch.SetKey(.5, 1)
				card.Enter.AddChannel(ch)
			}

			card.Exit = anim.New()
			card.Exit.Start()
			{
				const animTime = .3
				ch := anim.NewChannelFuncf32(func(v float32) {
					w.Position[0] = v
				})
				ch.SetKey(0, -w.Dim[0]-1)
				ch.SetKey(animTime, 0)

				opch := anim.NewChannelFuncf32(func(v float32) {
					w.Set("opacity", v)
				})
				opch.SetKey(0, 1)
				opch.SetKey(animTime, 0)

				card.Exit.AddChannel(ch)
				card.Exit.AddChannel(opch)
			}

			cards = append(cards, card)
		case gorge.EventUpdate:
			if len(cards) == 0 {
				return
			}
			curV := float32(0) //-(cards[0].Widget.Dim[1] + 1)
			t := cards
			for i := len(cards) - 1; i >= 0; i-- {
				c := t[i]
				if c.Timeout <= 0 {
					c.Exit.UpdateDelta(e.DeltaTime())
					if c.Exit.State() == anim.StateFinished {
						g.Remove(c.Widget)
						cards = append(cards[:i], cards[i+1:]...)
					}
				}
				c.Enter.UpdateDelta(e.DeltaTime())
				// curV := -float32(i+1) * (5 + 1)
				curV -= c.Widget.Dim[1] + 1

				pos := c.Widget.Position
				pos[1] = curV

				c.Widget.Position = c.Widget.Position.Lerp(pos, e.DeltaTime()*10)
				c.Timeout -= e.DeltaTime()
			}
			if len(t) > len(cards) {
				for i := range t[len(cards):] {
					t[len(cards)+i] = nil
				}
			}
		}
	})
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
	return func(b *gorlet.Builder) {
		color := m32.Color(0, 0, 0, .3)
		switch e.Severity {
		case SeverityInfo:
			color = m32.Color(0, 0, 0, .3)
		case SeverityWarn:
			color = m32.Color(1, 1, 0, .3)
		case SeverityError:
			color = m32.Color(1, 0, 0, .3)
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

		b.Observe("opacity", gorlet.ObsFunc(func(v float32) {
			p.Set("color", m32.Color(
				color[0],
				color[1],
				color[2],
				v*.3,
			))
			l.Set("textColor", m32.Color(1, v))
		}))
	}
}
