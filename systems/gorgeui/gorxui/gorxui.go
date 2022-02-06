// Package gorxui experimental XML to gorlet serializer.
package gorxui

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui/gorlet"
	"github.com/stdiopt/gorge/text"
)

type XUI struct {
	registry Registry
}

func New() *XUI {
	return &XUI{}
}

func Create(c any) gorlet.Func {
	return New().Create(c)
}

func (x *XUI) Create(c any) gorlet.Func {
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
		return c()
	default:
		panic(fmt.Errorf("gorxui: unknown type %T", c))
	}
}

func (x *XUI) Define(k string, c any) {
	if x.registry == nil {
		x.registry = make(Registry)
	}
	comp := x.Create(c)
	fn := func() gorlet.Func { return comp }
	x.registry[strings.ToLower(k)] = fn
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
func (x *XUI) FromString(s string) (gorlet.Func, error) {
	return x.read(strings.NewReader(s))
}

// MustFromString parses xml and panic if error.
func (x *XUI) MustFromString(s string) gorlet.Func {
	fn, err := x.read(strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	return fn
}

func (x *XUI) read(rd io.Reader) (gorlet.Func, error) {
	var target string
	defs := map[string][]gorlet.Func{}
	// main
	fns := []gorlet.Func{}

	push := func(fn func(*gorlet.B)) {
		if target == "" {
			fns = append(fns, fn)
			return
		}
		defs[target] = append(defs[target], fn)
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
				// b.Save()
				// add a label right here
				b.Use("autoSize", true)
				b.Use("textAlign", gorlet.TextAlign(text.AlignStart, text.AlignCenter))
				b.Use("text", b.Prop("text", t))
				b.UseAnchor(0)
				b.Label("")
				// b.Restore()
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
				target = name
				continue
			default:
			}
			fn, ok := x.Get(tok.Name.Local)
			if !ok {
				fn = func() gorlet.Func { return gorlet.Container() }
			}

			isFirst := len(fns) == 0
			if target != "" {
				isFirst = len(defs[target]) == 0
			}
			// This uses set root which might not be ideal for certain situations
			if isFirst { // this uses SetRoot
				push(func(b *gorlet.B) {
					e := b.SetRoot(fn())
					for _, a := range tok.Attr {
						setProp(b.Root(), e, a) // errcheck
					}
				})
				continue
			}
			push(func(b *gorlet.B) {
				e := b.Begin(fn())
				for _, a := range tok.Attr {
					setProp(b.Root(), e, a) // errcheck
				}
			})
		case xml.EndElement:
			if tok.Name.Local == "template" {
				fns := defs[target]
				bf := func(b *gorlet.B) {
					for _, fn := range fns {
						fn(b)
					}
				}
				x.Define(target, func() gorlet.Func { return bf })
				target = ""
				continue
			}
			push(func(b *gorlet.B) {
				b.End()
			})
		}
	}
	bf := func(b *gorlet.B) {
		for _, fn := range fns {
			fn(b)
		}
	}
	return bf, nil
}

// EventAction action triggered on stuff.
type EventAction struct {
	Action string
	Orig   event.Event
	Entity *gorlet.Entity
}

func setProp(root, e *gorlet.Entity, a xml.Attr) error {
	switch a.Name.Space {
	case "a":
		switch a.Name.Local {
		case "click":
			event.Handle(e, func(evt gorlet.EventClick) {
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
		root.ObserveTo(a.Value, e, a.Name.Local)
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
		e.SetLayout(l)
	case "margin":
		p, err := parseFloat32Slice(a.Value)
		if err != nil {
			return err
		}
		e.SetMargin(p...)
	case "rect":
		p, err := parseFloat32Slice(a.Value)
		if err != nil {
			return err
		}
		e.SetRect(p...)
	case "width":
		v, err := strconv.ParseFloat(a.Value, 32)
		if err != nil {
			return err
		}
		e.Size[0] = float32(v)
	case "height":
		v, err := strconv.ParseFloat(a.Value, 32)
		if err != nil {
			return err
		}
		e.Size[1] = float32(v)
	case "anchor":
		p, err := parseFloat32Slice(a.Value)
		if err != nil {
			return err
		}
		e.SetAnchor(p...)
	case "pivot":
		p, err := parseFloat32Slice(a.Value)
		if err != nil {
			return err
		}
		e.SetAnchor(p...)

	case "color", "textColor", "handlerColor", "borderColor":
		p, err := parseFloat32Slice(a.Value)
		if err != nil {
			return err
		}
		e.Set(a.Name.Local, gm.Color(p...))
	default:
		o := e.Observer(a.Name.Local)
		if o == nil {
			return nil
		}
		v, err := parseTyp(o.Type, a.Value)
		if err != nil {
			return err
		}
		e.Set(a.Name.Local, v)
	}
	return nil
}
