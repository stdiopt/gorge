package gorlet

import (
	"fmt"
	"reflect"
)

// Make dynamic observers here

type observerFunc = func(any)

type Observer struct {
	Type  reflect.Type // Func type
	Funcs []observerFunc
}

func (o *Observer) Call(v any) {
	for _, fn := range o.Funcs {
		fn(v)
	}
}

type WCustom struct {
	Widget[WCustom]
	observers map[string]*Observer
}

func Custom(fn BuildFunc) *WCustom {
	w := Build(&WCustom{})
	b := &B{root: &curEntity{entity: w}}
	fn(b)
	if b.clientArea != nil {
		w.SetClientArea(b.clientArea)
	}

	return w
}

func (w *WCustom) String() string {
	return fmt.Sprintf("%s obs:%d", w.Widget.String(), len(w.observers))
}

func (w *WCustom) Set(k string, v any) *WCustom {
	w.set(k, v)
	return w
}

// Observer returns the observer
func (w *WCustom) Observer(k string) *Observer {
	return w.observers[k]
}

func (w *WCustom) Observe(k string, ifn any) {
	fnVal := reflect.ValueOf(ifn)
	typ := fnVal.Type().In(0)
	// tTyp := reflect.TypeOf(ifn)
	args := make([]reflect.Value, 1)
	fn := func(v any) {
		args[0] = reflect.ValueOf(v)
		if args[0].Type() != typ {
			if reflect.TypeOf(v).ConvertibleTo(typ) {
				panic(fmt.Sprintf("Can't convert prop %T(%v) to %v", v, v, typ))
			}
			args[0] = args[0].Convert(typ)
		}
		fnVal.Call(args)
	}
	w.ObserveType(k, fnVal.Type(), fn)
}

func (w *WCustom) set(k string, v any) {
	o, ok := w.observers[k]
	if !ok {
		return
	}
	o.Call(v)
}

func (w *WCustom) ObserveType(k string, typ reflect.Type, fn func(any)) {
	if w.observers == nil {
		w.observers = make(map[string]*Observer)
	}

	o, ok := w.observers[k]
	if !ok {
		o = &Observer{
			Type:  typ,
			Funcs: make([]observerFunc, 0),
		}
		w.observers[k] = o
	}

	if typ != o.Type {
		panic("wrong observed func")
	}

	o.Funcs = append(o.Funcs, fn)
}

/*
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
	if typ.Kind() != reflect.Func {
		panic(fmt.Sprintf("Can't observe %v", typ))
	}
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
}*/
