package event

import "github.com/stdiopt/gorge/core/setlist"

type Event any

type Handler interface {
	HandleEvent(Event)
}
type HandlerFunc[T any] func(T)

// catch all event
type handleEvent struct {
	// To make it hashable
	fn HandlerFunc[Event]
}

func (h *handleEvent) HandleEvent(e Event) { (h.fn)(e) }

type Buser interface {
	bus() *Bus
}

type Bus struct {
	listeners []any // slices

	handlers setlist.SetList[Handler]
}

func (b *Bus) bus() *Bus { return b }

func (b *Bus) AddHandler(h Handler) {
	b.handlers.Add(h)
}

func (b *Bus) Remove(h Handler) {
	b.handlers.Remove(h)
}

func Trigger[T any](bb Buser, v T) {
	b := bb.bus()
	var t []HandlerFunc[T]
	for _, l := range b.listeners {
		if tt, ok := l.([]HandlerFunc[T]); ok {
			t = tt
			break
		}
	}
	for _, fn := range t {
		fn(v)
	}
	for _, h := range b.handlers.Items() {
		h.HandleEvent(v)
	}
}

func Handle[T any](bb Buser, fn HandlerFunc[T]) {
	b := bb.bus()

	if fn, ok := any(fn).(HandlerFunc[Event]); ok {
		b.handlers.Add(&handleEvent{fn: fn})
		return
	}

	i, l := search[T](bb)
	if i != -1 {
		b.listeners[i] = append(l, fn)
	}
	b.listeners = append(b.listeners, []HandlerFunc[T]{fn})
}

func search[T any](bb Buser) (int, []HandlerFunc[T]) {
	for i, l := range bb.bus().listeners {
		if tt, ok := l.([]HandlerFunc[T]); ok {
			return i, tt
		}
	}
	return -1, nil
}
