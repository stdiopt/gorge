// Package prop handles binding properties to functions, dependency injection etc
package prop

import (
	"fmt"
	"log"
	"reflect"
)

// ErrMissing represents an error for a dependency miss.
type ErrMissing struct {
	prop string
}

func (e ErrMissing) Error() string {
	return fmt.Sprintf("resource '%s' missing", e.prop)
}

type propFetcher interface {
	prop() interface{}
}

// Props contains type based props.
type Props struct {
	props map[reflect.Type]propFetcher
}

// PutProp put properties into the registry
func (p *Props) PutProp(vs ...interface{}) {
	if p.props == nil {
		p.props = map[reflect.Type]propFetcher{}
	}
	for i, v := range vs {
		if v == nil {
			log.Printf("Warning param %d is nil", i)
			continue
		}
		v := v // intentional shadow
		var pp propFetcher

		typ := reflect.TypeOf(v)
		// Specific constructor case
		if typ.Kind() == reflect.Func {
			if typ.NumIn() != 0 || typ.NumOut() != 1 {
				panic("wrong resource type")
			}
			fnVal := reflect.ValueOf(v)

			typ = typ.Out(0)
			pp = funcFetcher(func() interface{} {
				return fnVal.Call([]reflect.Value{})[0].Interface()
			})
		} else {
			pp = valueFetcher{v}
		}
		p.props[typ] = pp
	}
}

// BindProps call the inits.. funcs with the properties on the container
// order is not guaranteed since an initialization might be delayed for
// dependency solving, it will put initialization returns as a property and
// return a slice with those
func (p *Props) BindProps(inits ...interface{}) ([]interface{}, error) {
	var rets []interface{}
	missing := make([]interface{}, len(inits))

	for len(inits) > 0 {
		errs := []error{}
		missing = missing[:0]

		for i, fni := range inits {
			fnRets, err := p.bindProps(fni)
			if _, ok := err.(ErrMissing); ok {
				missing = append(missing, fni)
				errs = append(errs, err)
				continue
			}
			if err != nil {
				return nil, err
			}
			rets = append(rets, fnRets...)

			if len(missing) > 0 {
				missing = append(missing, inits[i+1:]...)
				break
			}
		}
		// Did not solve much, so we are stalled
		if len(inits) == len(missing) {
			return nil, fmt.Errorf("%v", errs)
		}
		inits = missing
	}
	return rets, nil
}

// BindProps will solve the function parameters and inject the properties as
// arguments, any func returns will be set in props
func (p *Props) bindProps(fn interface{}) ([]interface{}, error) {
	fnVal := reflect.ValueOf(fn)
	fnTyp := reflect.TypeOf(fn)
	if fnTyp.Kind() != reflect.Func {
		panic("bind: should be a func")
	}
	args := make([]reflect.Value, fnTyp.NumIn())
	for i := 0; i < fnTyp.NumIn(); i++ {
		argTyp := fnTyp.In(i)
		pp, ok := p.getProp(argTyp)
		if !ok {
			return nil, ErrMissing{fmt.Sprint(argTyp)}
		}
		args[i] = reflect.ValueOf(pp.prop())
	}
	res := fnVal.Call(args)
	// Deal with returns
	var rets []interface{}
	for _, r := range res {
		v := r.Interface()
		if v == nil || r.IsZero() {
			continue
		}
		// If one of the returns is an error we return nil and that error
		// note that we already nil/zero checked in the previous condition
		if err, ok := v.(error); ok {
			return nil, err
		}
		rets = append(rets, v)
	}
	if len(rets) > 0 {
		// Put any return type in the deps stuff
		p.PutProp(rets...)
	}

	return rets, nil
}

func (p *Props) getProp(t reflect.Type) (propFetcher, bool) {
	if p.props == nil {
		return nil, false
	}

	// Special case
	fn := p.props[t]
	if fn == nil {
		return nil, false
	}

	return fn, true
}

// DumpProps dump props to default logger.
func (p *Props) DumpProps() {
	log.Println("Registered props:")
	for k, v := range p.props {
		if v, ok := v.(valueFetcher); ok {
			log.Printf("\t%v: %v", k, reflect.TypeOf(v))
		}
	}
}

type funcFetcher func() interface{}

func (fn funcFetcher) prop() interface{} { return fn() }

type valueFetcher struct {
	val interface{}
}

func (v valueFetcher) prop() interface{} { return v.val }
