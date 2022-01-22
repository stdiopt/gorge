package gorlet

import (
	"fmt"
	"reflect"
)

type ObserverFunc = func(any)

func Ptr[T any](p *T) func(T) {
	return func(v T) {
		*p = v
	}
}

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
