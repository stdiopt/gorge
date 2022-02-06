package notify

import (
	"fmt"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/anim"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/gorgeutil"
)

// Context notification context.
type Context struct {
	system *system
	gorge  *gorge.Context
}

// Info sends a info notification.
func (c *Context) Info(s string) {
	event.Trigger(c.gorge, EventNotify{
		Message:  s,
		Severity: SeverityInfo,
	})
}

// Warn sends a warn notification.
func (c *Context) Warn(s string) {
	event.Trigger(c.gorge, EventNotify{
		Message:  s,
		Severity: SeverityWarn,
	})
}

// Error sends a error notification.
func (c *Context) Error(s string) {
	event.Trigger(c.gorge, EventNotify{
		Message:  s,
		Severity: SeverityError,
	})
}

// Infof sames as info but accepts a fmt format string
func (c *Context) Infof(f string, args ...any) {
	c.Info(fmt.Sprintf(f, args...))
}

// Warnf sames as warn but accepts a fmt format string
func (c *Context) Warnf(f string, args ...any) {
	c.Warn(fmt.Sprintf(f, args...))
}

// Errorf sames as error but accepts a fmt format string
func (c *Context) Errorf(f string, args ...any) {
	c.Error(fmt.Sprintf(f, args...))
}

// FromContext returns a notify context from the given gorge context
func FromContext(g *gorge.Context) *Context {
	if ctx, ok := gorge.GetContext[*Context](g); ok {
		return ctx
	}
	u := gorgeutil.FromContext(g)

	uiCam := u.UICamera()
	ui := u.UI(uiCam)

	s := &system{
		uiCam: uiCam,
		ui:    ui,
	}
	event.Handle(g, s.createNotification)
	event.Handle(g, func(e gorge.EventUpdate) {
		if len(s.cards) == 0 {
			return
		}
		curV := float32(0) //-(cards[0].Widget.Dim[1] + 1)
		t := s.cards
		for i := len(s.cards) - 1; i >= 0; i-- {
			c := t[i]
			if c.Timeout <= 0 {
				c.Exit.UpdateDelta(e.DeltaTime())
				if c.Exit.State() == anim.StateFinished {
					s.ui.Remove(c.Widget)
					s.cards = append(s.cards[:i], s.cards[i+1:]...)
				}
			} else {
				c.Enter.UpdateDelta(e.DeltaTime())
			}
			// curV := -float32(i+1) * (5 + 1)
			curV -= c.Widget.Size[1] + 1

			pos := c.Widget.Position
			pos[1] = curV

			c.Widget.Position = c.Widget.Position.Lerp(pos, e.DeltaTime()*10)
			c.Timeout -= e.DeltaTime()
		}
		if len(t) > len(s.cards) {
			for i := range t[len(s.cards):] {
				t[len(s.cards)+i] = nil
			}
		}
	})

	ctx := &Context{
		gorge:  g,
		system: s,
	}
	gorge.SetContext(g, ctx)
	return ctx
}
