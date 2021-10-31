package event

import (
	"fmt"
	"reflect"
)

// TypedFunc returns an handlerFunc that filters the event by the passed func param type
// if it matches it will call the function with the event.
func TypedFunc(fn interface{}) HandlerFunc {
	typ := fnTyp(fn)
	fnVal := reflect.ValueOf(fn)
	return func(e Event) {
		if reflect.TypeOf(e) == typ {
			fnVal.Call([]reflect.Value{reflect.ValueOf(e)})
		}
	}
}

// Mux returns a new TypeMux.
func Mux(h ...interface{}) *TypeMux {
	m := TypeMux{}
	m.Handle(h...)
	return &m
}

// TypeMux is an event handler that will route events by func param types
// It will call specific funcs for specific event signatures.
type TypeMux struct {
	handlers map[reflect.Type][]HandlerFunc
}

// HandleEvent implements the EventHandler interface.
func (m *TypeMux) HandleEvent(v Event) {
	if m.handlers == nil {
		return
	}
	t := reflect.TypeOf(v)

	handlers := m.handlers[t]
	for _, fn := range handlers {
		fn(v)
	}
}

// Handle accepts functions with a specific argument type based on an event.
func (m *TypeMux) Handle(args ...interface{}) {
	if m.handlers == nil {
		m.handlers = map[reflect.Type][]HandlerFunc{}
	}

	for _, a := range args {
		typ := fnTyp(a)
		fnVal := reflect.ValueOf(a)
		fn := func(v Event) {
			fnVal.Call([]reflect.Value{reflect.ValueOf(v)})
		}

		m.handlers[typ] = append(m.handlers[typ], fn)
	}
}

func fnTyp(fn interface{}) reflect.Type {
	typ := reflect.TypeOf(fn)
	if typ.Kind() != reflect.Func &&
		typ.NumIn() != 1 &&
		typ.NumOut() == 0 {
		panic("wrong type, should be a func with 1 param")
	}
	vtyp := typ.In(0)
	if vtyp.Kind() == reflect.Interface {
		panic(fmt.Sprintf("handler parameter should not be an interface: %v", vtyp))
	}
	return vtyp
}
