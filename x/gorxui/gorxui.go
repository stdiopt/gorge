// Package gorxui experimental XML to gorlet serializer.
package gorxui

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/text"
	"github.com/stdiopt/gorge/x/gorlet"
)

type XUI struct {
	registry Registry
}

func New() *XUI {
	return &XUI{}
}

func Create(c any) *gorlet.WCustom {
	return New().Create(c)
}

func (x *XUI) Create(c any) *gorlet.WCustom {
	switch c := c.(type) {
	case string:
		return x.MustFromString(c)
	case io.Reader:
		comp, err := x.read(c)
		if err != nil {
			panic(err)
		}
		return comp
	case buildFunc:
		return gorlet.Custom(func(b *gorlet.B) {
			b.Add(c())
		})
	case gorlet.BuildFunc:
		return gorlet.Custom(c)
	default:
		panic(fmt.Errorf("gorxui: unknown type %T", c))
	}
}

func (x *XUI) Define(k string, c any) {
	if x.registry == nil {
		x.registry = make(Registry)
	}
	switch c := c.(type) {
	case buildFunc:
		x.registry[strings.ToLower(k)] = c
	case gorlet.BuildFunc:
		x.registry[strings.ToLower(k)] = func() gorlet.Entity {
			return gorlet.Custom(c)
		}
	}

	// comp := x.Create(c)
	// fn := func() gorlet.Entity { return comp }
	// x.registry[strings.ToLower(k)] = fn
}

func (x *XUI) Get(k string) (buildFunc, bool) {
	k = strings.ToLower(k)
	if fn, ok := x.registry[k]; ok {
		return fn, true
	}
	if fn, ok := registry[k]; ok {
		return fn, true
	}
	return nil, false
}

// FromString parses a xml string into a gorlet.Entity.
func (x *XUI) FromString(s string) (*gorlet.WCustom, error) {
	return x.read(strings.NewReader(s))
}

// MustFromString parses xml and panic if error.
func (x *XUI) MustFromString(s string) *gorlet.WCustom {
	fn, err := x.read(strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	return fn
}

func (x *XUI) read(rd io.Reader) (*gorlet.WCustom, error) {
	var target string

	type tmpl struct {
		attr  []xml.Attr
		funcs []gorlet.BuildFunc
	}
	defs := map[string]*tmpl{}
	// main
	fns := []gorlet.BuildFunc{}

	push := func(fn func(*gorlet.B)) {
		if target == "" {
			fns = append(fns, fn)
			return
		}
		t, ok := defs[target]
		if !ok {
			t := &tmpl{}
			defs[target] = t
		}
		t.funcs = append(t.funcs, fn)
	}

	xm := xml.NewDecoder(rd)
	for {
		tok, err := xm.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("err:", err)
			// do Something else
			break
		}
		switch tok := tok.(type) {
		case xml.CharData:
			t := strings.TrimSpace(string(tok))
			if t == "" {
				continue
			}
			push(func(b *gorlet.B) {
				b.Label(t).
					SetAutoSize(true).
					SetTextAlign(text.AlignStart, text.AlignCenter).
					SetAnchor(0)
			})

		case xml.StartElement:
			// var e *gorlet.Entity
			t := strings.ToLower(tok.Name.Local)
			switch t {
			case "clientarea":
				push(func(b *gorlet.B) {
					b.ClientArea()
				})
				continue
			case "template":
				if target != "" {
					return nil, fmt.Errorf("gorxui: template can't be nested")
				}
				var name string
				for _, a := range tok.Attr {
					if a.Name.Local == "name" {
						name = strings.ToLower(a.Value)
						break
					}
				}
				if name == "" {
					return nil, fmt.Errorf("gorxui: template must have name attribute")
				}
				if _, ok := defs[name]; ok {
					return nil, fmt.Errorf("gorxui: template %q already defined", name)
				}

				defs[name] = &tmpl{
					attr: tok.Attr,
				}
				target = name
				continue
			default:
			}

			fn, ok := x.Get(tok.Name.Local)
			if !ok {
				fn = func() gorlet.Entity { return gorlet.Container() }
			}
			attr := append([]xml.Attr{}, tok.Attr...)
			push(func(b *gorlet.B) {
				e := fn()
				b.Begin(e)
				for _, a := range attr {
					setProp(b.Root().(*gorlet.WCustom), e, a) // errcheck
				}
			})
		case xml.EndElement:
			if tok.Name.Local == "template" {
				t := defs[target]
				bf := func(b *gorlet.B) {
					for _, fn := range t.funcs {
						fn(b)
					}
				}
				x.Define(target, func() gorlet.Entity {
					e := gorlet.Custom(bf)
					for _, a := range t.attr {
						setProp(e, e, a) // errcheck
					}
					return e
				})
				target = ""
				continue
			}
			push(func(b *gorlet.B) {
				b.End()
			})
		}
	}
	bf := gorlet.Custom(func(b *gorlet.B) {
		for _, fn := range fns {
			fn(b)
		}
	})
	return bf, nil
}

// EventAction action triggered on stuff.
type EventAction struct {
	Action string
	Orig   event.Event
	Entity gorlet.Entity
}

func setProp(root *gorlet.WCustom, e gorlet.Entity, a xml.Attr) error {
	switch a.Name.Space {
	case "a":
		switch a.Name.Local {
		case "click":
			log.Printf("Attaching handler click: %v to %v", a.Value, e)
			event.Handle(e, func(evt gorlet.EventClick) {
				log.Println("Clicked")
				event.Trigger(root, EventAction{
					a.Value,
					evt,
					e,
				})
			})
		case "changed":
			event.Handle(e, func(evt gorlet.EventValueChanged) {
				event.Trigger(root, EventAction{
					a.Value,
					evt,
					e,
				})
			})

		default:
		}
		event.Handle(e, func(evt EventAction) {
			if evt.Action != a.Name.Local {
				return
			}
			event.Trigger(root, EventAction{
				a.Value,
				evt.Orig,
				e,
			})
		})
		return nil
	case "p":
		if c, ok := e.(*gorlet.WCustom); ok {
			if o := c.Observer(a.Name.Local); o != nil {
				root.ObserveType(a.Value, o.Type, func(v any) {
					c.Set(a.Name.Local, v)
				})
			}
		}
		vmethod := reflect.ValueOf(e).MethodByName(fmt.Sprintf("Set%s", strings.Title(a.Name.Local)))
		if vmethod.IsValid() {
			root.Observe(a.Value, vmethod.Interface())
		}
		return nil
	case "":
		// Do nothing move next
	default:
		panic(fmt.Errorf("%v unknown namespace", a.Name.Space))
	}

	switch a.Name.Local {
	case "id":
		e.SetID(a.Value)
	case "layout":
		l, err := parseLayout(a.Value)
		if err != nil {
			return err
		}
		if c, ok := e.(*gorlet.WCustom); ok {
			c.SetLayout(l)
		}
	default:
		type customSetter interface {
			Observer()
			Set(string, any) *gorlet.WCustom
		}
		if c, ok := e.(*gorlet.WCustom); ok {
			if o := c.Observer(a.Name.Local); o != nil {
				v, err := parseTyp(o.Type.In(0), a.Value)
				if err != nil {
					return err
				}
				c.Set(a.Name.Local, v)
			}
		}
		// a.Name.Local
		mval := reflect.ValueOf(e)
		etyp := mval.Type()
		for i := 0; i < etyp.NumMethod(); i++ {
			f := etyp.Method(i)
			if !strings.HasPrefix(f.Name, "Set") {
				continue
			}
			if !strings.EqualFold(f.Name[3:], a.Name.Local) {
				continue
			}

			// Solve multiple params into a slice on v
			v, err := parseTyp(f.Type.In(1), a.Value)
			if err != nil {
				return err
			}
			// What if variadic?
			pv := reflect.ValueOf(v)
			args := []reflect.Value{}
			if f.Type.IsVariadic() {
				for i := 0; i < pv.Len(); i++ {
					args = append(args, pv.Index(i))
				}
			} else {
				args = append(args, pv)
			}
			mval.MethodByName(f.Name).Call(args)
			return nil

		}
	}

	return nil
}
