// Package gorge contains mostly data only components
package gorge

import (
	"errors"
	"log"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/internal/logger"
	"github.com/stdiopt/gorge/math/gm"
)

// ErrAlreadyStarted is returned when Start is called more than once
var ErrAlreadyStarted = errors.New("already started")

func init() {
	logger.Global()
}

// InitFunc type of function to initialize gorge.
type InitFunc func(*Context) error

type eventBus = event.Bus

// Gorge main state manager and message bus
type Gorge struct {
	eventBus
	contexts []any
	// tcore
	// screenSize since this is shared between places
	// Maybe create Device/Display so we can even use multiple displays
	screenSize gm.Vec2
	inits      []InitFunc

	fnch chan syncFunc
	done chan error
}

// New create a new manager with default systems
func New(inits ...InitFunc) *Gorge {
	return &Gorge{
		inits: inits,
		fnch:  make(chan syncFunc, 64),
	}
}

// SetScreenSize used by the gorgeapp to set the current screensize
// Might be changed in the future if we use multiple Display devices.
func (g *Gorge) SetScreenSize(s gm.Vec2) {
	g.screenSize = s
}

// ScreenSize returns the previously set screensize
func (g *Gorge) ScreenSize() gm.Vec2 {
	return g.screenSize
}

// Start the systems
// nolint: errcheck
func (g *Gorge) Start() error {
	if g.done != nil {
		return ErrAlreadyStarted
	}
	g.done = make(chan error)
	// Call every init func
	c := &Context{g}
	for _, fn := range g.inits {
		fn(c)
	}
	/*if err := g.tcore.Start(); err != nil {
		return err
	}*/

	Trigger(g, EventStart{})
	Trigger(g, EventAfterStart{})
	return nil
}

// Wait waits for execution to finish.
func (g *Gorge) Wait() error {
	if g.done == nil {
		return nil
	}
	return <-g.done
}

// Close closes the running instance.
func (g *Gorge) Close() {
	close(g.done)
	g.done = nil
}

// CloseWithError closes with an error which will be sent on Wait() call
func (g *Gorge) CloseWithError(err error) {
	g.done <- err
	close(g.done)
}

// Run will initialize gorge, Start and wait.
func (g *Gorge) Run() error {
	if err := g.Start(); err != nil {
		return err
	}
	return g.Wait()
}

type syncFunc struct {
	fn   func()
	done chan struct{}
}

// RunInMain schedule a func to be run on main loop
// It will wait for the function to return
func (g *Gorge) RunInMain(fn func()) {
	sf := syncFunc{
		fn:   fn,
		done: make(chan struct{}),
	}
	g.fnch <- sf
	// Wait for func to finish
	<-sf.done
}

// Update just updates stuff right away
// nolint: errcheck
func (g *Gorge) Update(dt float32) {
	// How much does this costs?
	// Calls any schedule func
	select {
	case sf := <-g.fnch:
		sf.fn()
		close(sf.done)
	default:
	}
	Trigger(g, EventPreUpdate(dt))
	Trigger(g, EventUpdate(dt))
	Trigger(g, EventPostUpdate(dt))
	Trigger(g, EventRender(dt))
}

// Add adds an entity
// nolint: errcheck
func (g *Gorge) Add(ents ...Entity) {
	for _, e := range ents {
		EachEntity(e, func(e Entity) {
			Trigger(g, EventAddEntity{e})
		})
	}
}

// Remove an entity
// nolint: errcheck
func (g *Gorge) Remove(ents ...Entity) {
	for _, e := range ents {
		EachEntity(e, func(e Entity) {
			Trigger(g, EventRemoveEntity{e})
		})
	}
}

// TriggerInMain does not trigger synchronous per say but does trigger in main loop
// this is useful for GL related operations that depends on the specific thread it's running
// since we don't control much threads

//func (g *Gorge) TriggerInMain(v any) {
//	g.RunInMain(func() { g.Trigger(v) })
//}

// Error persists an error in the event system
// nolint: errcheck
func (g *Gorge) Error(err error) {
	log.Printf("[error] %v", err)
	Trigger(g, EventError{err})
}

// Handlers helpers

// HandleUpdate adds a listener that filters events and calls fn if it is the
// EventUpdate.
func (g *Gorge) HandleUpdate(fn func(float32)) {
	event.Handle(g, func(e EventUpdate) {
		fn(float32(e))
	})
}

// HandleError registers a function that filters events and calls fn if event
// is the EventError.
func (g *Gorge) HandleError(fn func(err error)) {
	event.Handle(g, func(e EventError) {
		fn(e.Err)
	})
}
