package notify

import (
	"fmt"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gorgeutil"
)

// Context notification context.
type Context struct {
	system *system
	gorge  *gorge.Context
}

// Info sends a info notification.
func (c *Context) Info(s string) {
	c.gorge.Trigger(EventNotify{
		Message:  s,
		Severity: SeverityInfo,
	})
}

// Warn sends a warn notification.
func (c *Context) Warn(s string) {
	c.gorge.Trigger(EventNotify{
		Message:  s,
		Severity: SeverityWarn,
	})
}

// Error sends a error notification.
func (c *Context) Error(s string) {
	c.gorge.Trigger(EventNotify{
		Message:  s,
		Severity: SeverityError,
	})
}

// Infof sames as info but accepts a fmt format string
func (c *Context) Infof(f string, args ...interface{}) {
	c.Info(fmt.Sprintf(f, args...))
}

// Warnf sames as warn but accepts a fmt format string
func (c *Context) Warnf(f string, args ...interface{}) {
	c.Warn(fmt.Sprintf(f, args...))
}

// Errorf sames as error but accepts a fmt format string
func (c *Context) Errorf(f string, args ...interface{}) {
	c.Error(fmt.Sprintf(f, args...))
}

// FromContext returns a notify context from the given gorge context
func FromContext(g *gorge.Context) *Context {
	if ctx, ok := gorge.GetSystem(g, ctxKey).(*Context); ok {
		return ctx
	}
	u := gorgeutil.FromContext(g)

	uiCam := u.UICamera()
	ui := u.UI(uiCam)

	s := &system{
		uiCam: uiCam,
		ui:    ui,
	}
	g.Handle(s)
	ctx := &Context{
		gorge:  g,
		system: s,
	}
	gorge.AddSystem(g, ctxKey, ctx)
	return ctx
}
