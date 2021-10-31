// Package event simpler event
package event

import (
	"fmt"
	"reflect"
	"sync"
)

// Event is the raw event interface which accepts anything.
type Event interface{}

// HandlerFunc is the handler type for an event handler.
type HandlerFunc func(Event)

// HandleEvent implements the event.Handler interface.
func (fn HandlerFunc) HandleEvent(v Event) { fn(v) }

// Handler interface for event handler.
type Handler interface {
	HandleEvent(v Event)
}

// Trigger interface for even triggers chaining.
type Trigger interface {
	Trigger(v Event)
}

// Bus it's the event bus that holds listeners and triggers events.
type Bus struct {
	mu sync.Mutex

	listeners []Handler

	// provides a way to check if the comparable handler was registered already
	comparables map[Handler]struct{}
	pool        *sync.Pool
}

// HandleFunc adds a func based listener.
func (b *Bus) HandleFunc(fn func(v Event)) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handle(HandlerFunc(fn))
}

// Handle adds a Listener that handles events.
func (b *Bus) Handle(h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handle(h)
}

// RemoveHandler remove an handler if the handler is a comparable type (i.e.
// not a func).
func (b *Bus) RemoveHandler(h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.remove(h)
}

// Trigger broadcast a message across all listeners
// Pool reason:
//   Although it's a bit extra for triggers which handles updates it's also
//   SAFE: We work with a copied slice instead of the original since an handler
//   can manipulate the original slice by calling Handler or triggering even
//   other event This way we avoid alocating slice copies everytime we call
//   trigger by reusing previous slices from the pool, the only drawBack will
//   be that these slices will get stuck in the pool
func (b *Bus) Trigger(v Event) {
	if b.pool == nil {
		b.pool = &sync.Pool{
			New: func() interface{} {
				return &listenerSlice{
					listeners: []Handler{},
				}
			},
		}
	}
	p := b.pool.Get().(*listenerSlice)
	// Copy listeners so we can manipulate the original in some Handle call by
	// adding handlers or triggering new things, also this will avoid
	// alocations of new listener slice
	b.mu.Lock()
	p.listeners = append(p.listeners[:0], b.listeners...)
	b.mu.Unlock()

	for i, h := range p.listeners {
		h.HandleEvent(v)
		p.listeners[i] = nil // deref the handler from the slice copy
	}

	b.pool.Put(p)
}

func (b *Bus) handle(h Handler) {
	b.listeners = append(b.listeners, h)

	typ := reflect.TypeOf(h)
	if !typ.Comparable() {
		return
	}
	if b.comparables == nil {
		b.comparables = map[Handler]struct{}{}
	}
	if _, ok := b.comparables[h]; ok {
		panic(fmt.Errorf("event handler already registered: %T", h))
	}
	b.comparables[h] = struct{}{}
}

func (b *Bus) remove(h Handler) {
	typ := reflect.TypeOf(h)
	// if h is not a comparable we can't do much
	if !typ.Comparable() {
		return
	}

	if b.comparables == nil {
		return
	}
	delete(b.comparables, h)

	for i, hh := range b.listeners {
		if h == hh {
			t := b.listeners
			b.listeners = append(b.listeners[:i], b.listeners[i+1:]...)
			t[len(t)-1] = nil // remove last one as it was copied
			return
		}
	}
}

type listenerSlice struct {
	listeners []Handler
}
