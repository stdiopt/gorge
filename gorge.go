// Package gorge contains mostly data only components
package gorge

import (
	"log"

	"github.com/stdiopt/gorge/core"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/internal/logger"
	"github.com/stdiopt/gorge/m32"
)

func init() {
	logger.Global()
}

type tcore = core.Core

// Gorge main state manager and message bus
type Gorge struct {
	tcore
	// screenSize since this is shared between places
	// Maybe create Device/Dysplay so we can even use multiple displays
	screenSize m32.Vec2

	fnch chan syncFunc
}

// New create a new manager with default systems
func New(systems ...interface{}) *Gorge {
	g := &Gorge{
		tcore: *core.New(systems...),
		fnch:  make(chan syncFunc, 64),
	}
	g.PutProp(func() *Context {
		return &Context{g}
	})
	return g
}

// SetScreenSize used by the gorgeapp to set the current screensize
// Might be changed in the future if we use multiple Display devices.
func (g *Gorge) SetScreenSize(s m32.Vec2) {
	g.screenSize = s
}

// ScreenSize returns the previously set screensize
func (g *Gorge) ScreenSize() m32.Vec2 {
	return g.screenSize
}

// Start the systems
// nolint: errcheck
func (g *Gorge) Start() error {
	if err := g.tcore.Start(); err != nil {
		return err
	}

	g.Trigger(EventStart{})
	g.Trigger(EventAfterStart{})
	return nil
}

// Run will initialize gorge, Start and wait.
func (g *Gorge) Run() error {
	if err := g.Start(); err != nil {
		return err
	}
	return g.Wait()
}

type syncFunc struct {
	fn func()
	ch chan struct{}
}

// RunInMain schedule a func to be run on main loop
// It will wait for the function to return
func (g *Gorge) RunInMain(fn func()) {
	sf := syncFunc{
		fn: fn,
		ch: make(chan struct{}),
	}
	g.fnch <- sf
	// Wait for func to finish
	<-sf.ch
}

// Update just updates stuff right away
// nolint: errcheck
func (g *Gorge) Update(dt float32) {
	// How much does this costs?
	// Calls any schedule func
	select {
	case sf := <-g.fnch:
		sf.fn()
		close(sf.ch)
	default:
	}
	g.Trigger(EventPreUpdate(dt))
	g.Trigger(EventUpdate(dt))
	g.Trigger(EventPostUpdate(dt))
	g.Trigger(EventRender(dt))
}

// Add adds an entity
// nolint: errcheck
func (g *Gorge) Add(ents ...Entity) {
	for _, e := range ents {
		EachEntity(e, func(e Entity) {
			g.Trigger(EventAddEntity{e})
		})
	}
}

// Remove an entity
// nolint: errcheck
func (g *Gorge) Remove(ents ...Entity) {
	for _, e := range ents {
		EachEntity(e, func(e Entity) {
			g.Trigger(EventRemoveEntity{e})
		})
	}
}

// TriggerOnUpdate does not trigger synchronous per say but does trigger in main loop
// this is useful for GL related operations that depends on the specific thread it's running
// since we don't control much threads
func (g *Gorge) TriggerOnUpdate(v interface{}) {
	g.RunInMain(func() { g.Trigger(v) })
}

// Error persists an error in the event system
// nolint: errcheck
func (g *Gorge) Error(err error) {
	log.Printf("[error] %v", err)
	g.Trigger(EventError{err})
}

// UpdateResource triggers an EventResourceUpdate event
// so the systems can garsp the resource (put textures in gpu, etc...).
func (g *Gorge) UpdateResource(r ResourceRef) {
	g.TriggerOnUpdate(EventResourceUpdate{
		Resource: r,
	})
}

// Handlers helpers

// HandleUpdate adds a listener that filters events and calls fn if it is the
// EventUpdate.
func (g *Gorge) HandleUpdate(fn func(EventUpdate)) {
	g.HandleFunc(func(e event.Event) {
		if e, ok := e.(EventUpdate); ok {
			fn(e)
		}
	})
}

// HandleError registers a function that filters events and calls fn if event
// is the EventError.
func (g *Gorge) HandleError(fn func(err error)) {
	g.HandleFunc(func(v event.Event) {
		if e, ok := v.(EventError); ok {
			fn(e.Err)
		}
	})
}
