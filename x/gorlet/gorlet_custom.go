package gorlet

import (
	"fmt"
	"reflect"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/systems/gorgeui"
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
	Widget[*WCustom]
	Layout    Layouter
	observers map[string]*Observer
}

func Custom(fn BuildFunc) *WCustom {
	w := Build(&WCustom{})
	b := &B{root: &curEntity{entity: w}}
	fn(b)
	if b.clientArea != nil {
		w.SetClientArea(b.clientArea)
	}

	event.Handle(w, func(gorgeui.EventUpdate) {
		if w.Layout != nil {
			w.Layout.Layout(w)
		}
	})

	return w
}

func (w *WCustom) SetLayout(li ...Layouter) *WCustom {
	w.Layout = LayoutMulti(li...)
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
