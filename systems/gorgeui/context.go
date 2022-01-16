package gorgeui

import (
	"log"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/resource"
	"github.com/stdiopt/gorge/text"
)

var ctxKey = struct{ string }{"gorgeui"}

// Context to be used in gorge systems.
type Context struct {
	*system
}

// FromContext retrieve gorgeui context from gorge
func FromContext(g *gorge.Context) *Context {
	if ctx, ok := gorge.GetContext[*Context](g); ok {
		return ctx
	}

	rc := resource.FromContext(g)

	log.Println("Initializing system")
	dbg := newDebugLines()
	dbg.SetQueue(200)
	dbg.SetCullMask(gorge.CullMaskUIDebug)
	g.Add(dbg)

	DefaultFont = &text.Font{}
	if err := rc.Load(DefaultFont, "_gorge/fonts/font.ttf"); err != nil {
		// what to do here, send error to gorge??!
		log.Println("error loading font:", err)
		return nil
	}

	s := &system{
		gorge: g,
		font:  DefaultFont,
		dbg:   dbg,
	}
	s.setupEvents(g)

	return gorge.AddContext(g, &Context{s})
}

// New returns a new UI
func (c Context) New(cam cameraEntity) *UI {
	ui := New()
	ui.SetCamera(cam)
	c.gorge.Add(ui)
	return ui
}
