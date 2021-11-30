package gorgeui

import (
	"fmt"
	"reflect"
	"sync"
)

// Event base event for the thing.
type Event struct {
	Entity          Entity
	Value           interface{}
	stopPropagation bool
}

// HandlerFunc type of func to handle an event
type HandlerFunc func(Event)

// HandleEvent implements the Handler interface.
func (fn HandlerFunc) HandleEvent(e Event) { fn(e) }

// Handler event Handler interface.
type Handler interface {
	HandleEvent(Event)
}

type trigger interface {
	trigger(Event)
}

// eventBus eventBus
type eventBus struct {
	mu          sync.Mutex
	listeners   []Handler
	comparables map[Handler]struct{}
}

// HandleFunc adds a func based listener.
func (b *eventBus) HandleFunc(fn HandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handle(fn)
}

// Handle adds a Listener that handles events.
func (b *eventBus) Handle(h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handle(h)
}

func (b *eventBus) handle(h Handler) {
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

func (b *eventBus) remove(h Handler) {
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

func (b *eventBus) trigger(e Event) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, h := range b.listeners {
		h.HandleEvent(e)
	}
}

type listenerSlice struct {
	listeners []Handler
}
