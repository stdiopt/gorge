package event

type Event any

type Handler interface {
	HandleEvent(Event)
}
type HandlerFunc[T any] func(T)

type Buser interface {
	bus() *Bus
}

type Bus struct {
	listeners  []any // slices
	handlers   []Handler
	handlerSet map[Handler]struct{}
}

func (b *Bus) bus() *Bus { return b }

func (b *Bus) AddHandler(h Handler) {
	if b.handlerSet == nil {
		b.handlerSet = make(map[Handler]struct{})
	}
	if _, ok := b.handlerSet[h]; ok {
		return
	}
	b.handlers = append(b.handlers, h)
	b.handlerSet[h] = struct{}{}
}

func (b *Bus) Remove(h Handler) {
	if b.handlerSet == nil {
		return
	}
	delete(b.handlerSet, h)
	for i, v := range b.handlers {
		if v == h {
			t := b.handlers
			b.handlers = append(b.handlers[:i], b.handlers[i+1:]...)
			t[len(t)-1] = nil
			return
		}
	}
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
	for _, h := range b.handlers {
		h.HandleEvent(v)
	}
}

func HandleFunc[T any](bb Buser, h HandlerFunc[T]) {
	b := bb.bus()

	i, l := search[T](bb)
	if i != -1 {
		b.listeners[i] = append(l, h)
	}
	b.listeners = append(b.listeners, []HandlerFunc[T]{h})
}

func search[T any](bb Buser) (int, []HandlerFunc[T]) {
	for i, l := range bb.bus().listeners {
		if tt, ok := l.([]HandlerFunc[T]); ok {
			return i, tt
		}
	}
	return -1, nil
}
