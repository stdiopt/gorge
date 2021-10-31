// Package core contains stuff to bootstrap a world
package core

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/core/prop"
)

// State type for gorge.
type State int

func (s State) String() string {
	switch s {
	case StateZero:
		return "zero"
	case StateInitialized:
		return "initialized"
	case StateStarted:
		return "started"
	case StateClosed:
		return "closed"
	}
	return "<undefined>"
}

// gorge States.
const (
	StateZero = iota
	StateInitialized
	StateStarted
	StateClosed
)

type (
	bus   = event.Bus
	props = prop.Props
)

// Core handles event triggering and properties binding for systems.
type Core struct {
	bus
	props
	state State
	done  chan error
	inits []interface{}

	beforeInit []interface{}
}

// New returns a new initialized *Core with system functions that will be
// called upon core.Start.
func New(systems ...interface{}) *Core {
	for _, s := range systems {
		if reflect.TypeOf(s).Kind() != reflect.Func {
			panic(fmt.Errorf("invalid initializer: %T, system initializer should be a func", s))
		}
	}
	return &Core{
		inits: systems,
		done:  make(chan error, 1),
	}
}

// Start initialize the systems.
func (c *Core) Start() error {
	if c.state != StateZero {
		return fmt.Errorf("cannot initialize, current state is: %v", c.state)
	}
	if err := c.init(); err != nil {
		return err
	}
	c.state = StateStarted
	for _, v := range c.beforeInit {
		c.bus.Trigger(v)
	}
	c.beforeInit = nil

	return nil
}

// Wait waits for the core to finish.
func (c *Core) Wait() error {
	if c.done == nil {
		c.done = make(chan error)
	}
	return <-c.done
}

// Run Starts and waits for core to finish.
func (c *Core) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

// Close thuts things down.
func (c *Core) Close(err ...error) {
	if len(err) != 0 {
		c.done <- err[0]
	}
	close(c.done)
	c.done = nil
}

// Trigger triggers an event that will propagate among handlers.
func (c *Core) Trigger(v interface{}) {
	/*if c.state != StateStarted {
		log.Println("Caching trigger")
		c.beforeInit = append(c.beforeInit, v)
		return
	}*/
	c.bus.Trigger(v)
}

// init calls the Initializator funcs passed through init
// if the return of those initializators implements the event.Handler interface
// we add it to the event listeners.
func (c *Core) init() error {
	if c.state != StateZero {
		return errors.New("cannot initialize again")
	}
	c.state = StateInitialized

	c.PutProp(c) // trivial

	rets, err := c.props.BindProps(c.inits...)
	if err != nil {
		return err
	}
	for _, r := range rets {
		if e, ok := r.(event.Handler); ok {
			c.Handle(e)
		}
	}
	return nil
}

// BindProps will call any function passed as parameter with any existing
// registered prop.
func (c *Core) BindProps(inits ...interface{}) error {
	_, err := c.props.BindProps(inits...)
	return err
}
