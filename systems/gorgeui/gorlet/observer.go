package gorlet

import (
	"fmt"
	"reflect"
)

type ObserverFunc = func(any)

type Observer struct {
	Type  reflect.Type // First type
	Funcs []ObserverFunc
}

func (o *Observer) Call(v interface{}) {
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

func (o *observers) observe(k string, ifn any) {
	if ifn == nil { // this will delete
		o.observeWithType(k, nil, nil)
		return
	}

	fn, typ := makeObserverFunc(ifn)

	o.observeWithType(k, typ, fn)
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
	if o.observers == nil {
		o.observers = make(map[string]*Observer)
	}
}

func (o *observers) observer(k string) *Observer {
	return o.observers[k]
}

/*
// ObsFunc creates a typed observer func from reflection.
func ObsFunc(fn any) ObserverFunc {
	if fn, ok := fn.(ObserverFunc); ok {
		return fn
	}
	fnVal := reflect.ValueOf(fn)
	inTyp := fnVal.Type().In(0)

	return func(v any) {
		arg := reflect.ValueOf(v)
		if aTyp := arg.Type(); aTyp != inTyp {
			if !arg.CanConvert(inTyp) {
				panic(fmt.Sprintf("Can't convert prop %v(%v) to %v", aTyp, v, inTyp))
			}
			arg = arg.Convert(inTyp)
		}
		// Type check somewhere
		fnVal.Call([]reflect.Value{arg})
	}
}
*/
func makeObserverFunc(ifn any) (ObserverFunc, reflect.Type) {
	fnVal := reflect.ValueOf(ifn)
	fnTyp := fnVal.Type()

	if fnTyp.Kind() != reflect.Func {
		panic("not a function")
	}
	if fnTyp.NumIn() != 1 {
		panic("function must have one input parameter")
	}
	if fnTyp.NumOut() != 0 {
		panic("function must have no output parameters")
	}
	inTyp := fnTyp.In(0)

	fn := func(v any) {
		arg := reflect.ValueOf(v)
		if aTyp := arg.Type(); aTyp != inTyp {
			if !arg.CanConvert(inTyp) {
				panic(fmt.Sprintf("Can't convert prop %v(%v) to %v", aTyp, v, inTyp))
			}
			arg = arg.Convert(inTyp)
		}
		fnVal.Call([]reflect.Value{arg})
	}
	return fn, inTyp
}

func Ptr[T any](p *T) func(T) {
	return func(v T) {
		*p = v
	}
}

/*
type ObserverFunc = func(any)


var (
	typAny    = reflect.TypeOf((*any)(nil)).Elem()
	typString = reflect.TypeOf("")
)

func makeObserver(fn any) ObserverFunc {
	switch fn := fn.(type) {
	case func(any):
		return fn
	}

	fnVal := reflect.ValueOf(fn)
	fnTyp := fnVal.Type()
	if fnTyp.Kind() != reflect.Func {
		panic("not a function")
	}
	if fnTyp.NumIn() != 1 {
		panic("not a function with one argument")
	}
	if fnTyp.NumOut() != 0 {
		panic("not a function with no return values")
	}
	inTyp := fnTyp.In(0)

	return func(v any) {
		arg := reflect.ValueOf(v)
		if aTyp := arg.Type(); aTyp != inTyp {
			if !arg.CanConvert(inTyp) {
				panic(fmt.Sprintf("Can't convert prop %v(%v) to %v", aTyp, v, inTyp))
			}
			arg = arg.Convert(inTyp)
		}
		// Type check somewhere
		fnVal.Call([]reflect.Value{arg})
	}
}

/*
type observer interface {
	Observe(string, func(interface{}))
}

func Observe[T any](o observer, name string, fn func(T)) {
	// ObsFunc because we still use reflection to avoid int to float conversions
	// which panic on interface{} until we find something better
	o.Observe(name, ObsFunc(fn))
}
*/
