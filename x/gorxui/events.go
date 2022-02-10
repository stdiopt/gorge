package gorxui

import (
	"fmt"
	"reflect"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/x/gorlet"
)

type Actions map[string]any

func On(e gorlet.Entity, a Actions) {
	event.Handle(e, a.handler())
}

func (a Actions) handler() func(e EventAction) {
	value := map[string]func(v any){}
	click := map[string]func(){}

	for k, v := range a {
		// Should be either func() or func(T)
		switch v := v.(type) {
		case func():
			click[k] = v
		default:
			value[k] = makeFunc(v)
		}
	}
	return func(e EventAction) {
		switch evt := e.Orig.(type) {
		case gorlet.EventValueChanged:
			if fn, ok := value[e.Action]; ok {
				fn(evt.Value)
			}
		case gorlet.EventClick:
			if fn, ok := click[e.Action]; ok {
				fn()
			}
		}
	}
}

func (a Actions) For(e gorlet.Entity) {
	event.Handle(e, a.handler())
}

func makeFunc(ifn any) func(any) {
	if fn, ok := ifn.(func(any)); ok {
		return fn
	}
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

	return func(v any) {
		arg := reflect.ValueOf(v)
		if aTyp := arg.Type(); aTyp != inTyp {
			if !arg.CanConvert(inTyp) {
				panic(fmt.Sprintf("Can't convert prop %v(%v) to %v", aTyp, v, inTyp))
			}
			arg = arg.Convert(inTyp)
		}
		fnVal.Call([]reflect.Value{arg})
	}
}
