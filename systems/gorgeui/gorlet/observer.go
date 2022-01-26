package gorlet

import (
	"fmt"
	"reflect"
)

type ObserverFunc = func(any)

type Observer struct {
	Type  reflect.Type // Func type
	Funcs []ObserverFunc
}

func (o *Observer) Call(v any) {
	for _, fn := range o.Funcs {
		fn(v)
	}
}

type observers struct {
	observers map[string]*Observer
}

func (o *observers) set(k string, v any) bool {
	oo := o.observers[k]
	if oo == nil {
		return false
	}
	oo.Call(v)
	return true
}

func (o *observers) observeWithType(k string, typ reflect.Type, fn ObserverFunc) {
	if fn == nil {
		if o.observers == nil {
			return
		}
		delete(o.observers, k)
		return
	}
	if o.observers == nil {
		o.observers = make(map[string]*Observer)
	}

	op, ok := o.observers[k]
	if !ok {
		op = &Observer{
			Type: typ,
		}
		o.observers[k] = op
	}
	if op.Type != typ {
		panic(fmt.Sprintf("type mismatch: %s != %s, observer already registered with different type", op.Type, typ))
	}
	op.Funcs = append(op.Funcs, fn)
}

func (o *observers) observer(k string) *Observer {
	return o.observers[k]
}

type observer interface {
	observeWithType(string, reflect.Type, ObserverFunc)
}

func Observe[T any](o observer, k string, tfn func(T)) {
	fn, typ := makeObserverFuncGen(tfn)
	o.observeWithType(k, typ, fn)
}

func Ptr[T any](p *T) func(T) {
	return func(v T) {
		*p = v
	}
}

func makeObserverFuncGen[T any](tfn func(T)) (func(any), reflect.Type) {
	tTyp := reflect.TypeOf(*new(T))

	fn := func(v any) {
		if v, ok := v.(T); ok {
			tfn(v)
			return
		}
		if !reflect.TypeOf(v).ConvertibleTo(tTyp) {
			panic(fmt.Sprintf("Can't convert prop %T(%v) to %v", v, v, tTyp))
		}
		tfn(reflect.ValueOf(v).Convert(tTyp).Interface().(T))
	}
	return fn, reflect.TypeOf(tfn)
}
