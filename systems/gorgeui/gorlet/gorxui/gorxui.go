// Package gorxui experimental XML to gorlet serializer.
package gorxui

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/systems/gorgeui/gorlet"
	"github.com/stdiopt/gorge/text"
)

type XUI struct {
	registry Registry
}

func New() *XUI {
	return &XUI{}
}

func (x *XUI) Define(k string, fn buildFunc) {
	if x.registry == nil {
		x.registry = make(Registry)
	}
	x.registry[k] = fn
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
	fns := []func(*gorlet.Builder){}

	push := func(fn func(*gorlet.Builder)) {
		fns = append(fns, fn)
	}

	rr := registry.merge(x.registry)

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
			push(func(b *gorlet.Builder) {
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
			var e *gorlet.Entity
			fn, ok := rr[strings.ToLower(tok.Name.Local)]
			if !ok {
				fn = func() gorlet.Func { return gorlet.Container() }
			}
			// log.Println("Start:", tok)
			// b.UseLayout(gorlet.LayoutList(.5))
			push(func(b *gorlet.Builder) {
				e = b.Begin(fn())
				for _, a := range tok.Attr {
					_ = setProp(b.Root(), e, a) // errcheck
				}
			})
		case xml.EndElement:
			push(func(b *gorlet.Builder) {
				b.End()
			})
		}
	}
	bf := func(b *gorlet.Builder) {
		for _, fn := range fns {
			fn(b)
		}
	}
	return bf, nil
}

// EventAction action triggered on stuff.
type EventAction struct {
	Action string
	Entity *gorlet.Entity
}

func setProp(root, e *gorlet.Entity, a xml.Attr) error {
	if a.Name.Space == "a" {
		switch a.Name.Local {
		case "click":
			event.Handle(e, func(gorgeui.EventPointerUp) {
				gorge.Trigger(root, EventAction{a.Value, e})
			})
		default:
			event.Handle(e, func(evt EventAction) {
				if evt.Action != a.Name.Local {
					return
				}

				gorge.Trigger(root, EventAction{a.Value, e})
			})
			log.Println("Unknown action", a.Name.Local)
		}
		return nil
	}
	if a.Name.Space == "p" {
		log.Printf("Observing %v in %v", a.Value, a.Name.Local)
		root.Observe(a.Value, e.PropSetter(a.Name.Local))
		return nil
	}
	switch a.Name.Local {
	case "layout":
		parts := strings.SplitN(a.Value, " ", 2)
		switch parts[0] {
		case "list", "vlist":
			var spacing float32
			if len(parts) > 1 {
				s, err := strconv.ParseFloat(parts[1], 32)
				if err != nil {
					return err
				}
				spacing = float32(s)
			}
			e.SetLayout(gorlet.LayoutList(spacing))
		case "flex":
			f, err := parseFlexProp(parts[1])
			if err != nil {
				return err
			}
			e.SetLayout(f)
		case "autoHeight":
			var spacing float32
			if len(parts) > 1 {
				s, err := strconv.ParseFloat(parts[1], 32)
				if err != nil {
					return err
				}
				spacing = float32(s)
			}
			e.SetLayout(gorlet.AutoHeight(spacing))
		default:
			return fmt.Errorf("layout %q not implemented", a.Value)
		}
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
		e.Dim[0] = float32(v)
	case "height":
		v, err := strconv.ParseFloat(a.Value, 32)
		if err != nil {
			return err
		}
		e.Dim[1] = float32(v)
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

	case "color", "textColor", "handlerColor":
		p, err := parseFloat32Slice(a.Value)
		if err != nil {
			return err
		}
		e.Set(a.Name.Local, gm.Color(p...))
	default:
		fn := e.PropSetter(a.Name.Local)
		typ := reflect.TypeOf(fn).In(0)
		v, err := parseTyp(typ, a.Value)
		if err != nil {
			return err
		}
		e.Set(a.Name.Local, v)
	}
	return nil
}
