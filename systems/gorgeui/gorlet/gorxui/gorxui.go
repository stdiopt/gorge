// Package gorxui experimental XML to gorlet serializer.
package gorxui

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/systems/gorgeui/gorlet"
	"github.com/stdiopt/gorge/text"
)

type buildFunc = gorlet.Func

var registry = map[string]func() buildFunc{
	"panel":      func() buildFunc { return gorlet.Panel() }, // nolint: gocritic
	"label":      func() buildFunc { return gorlet.Label("") },
	"textbutton": func() buildFunc { return gorlet.TextButton("", nil) },
	"spinner":    func() buildFunc { return gorlet.Spinner("", nil) },
	"slider":     func() buildFunc { return gorlet.Slider(0, 1, nil) },
	// containers
	"window":  func() buildFunc { return gorlet.Window("") },
	"labeled": func() buildFunc { return gorlet.Labeled("") },
}

// Register a Tag func.
func Register(k string, fn func() buildFunc) {
	registry[k] = fn
}

// FromString parses a xml string into a gorlet.Entity.
func FromString(s string) (gorlet.Func, error) {
	return read(strings.NewReader(s))
}

// MustFromString parses xml and panic if error.
func MustFromString(s string) gorlet.Func {
	fn, err := read(strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	return fn
}

func read(r io.Reader) (gorlet.Func, error) {
	fns := []func(*gorlet.Builder){}

	push := func(fn func(*gorlet.Builder)) {
		fns = append(fns, fn)
	}

	x := xml.NewDecoder(r)
	for {
		tok, err := x.Token()
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
			fn, ok := registry[strings.ToLower(tok.Name.Local)]
			if !ok {
				fn = func() buildFunc { return gorlet.Container }
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

func parseFloat32Slice(str string) ([]float32, error) {
	var ret []float32
	sp := strings.Split(str, ",")
	for _, s := range sp {
		f, err := strconv.ParseFloat(strings.TrimSpace(s), 32)
		if err != nil {
			return nil, err
		}
		ret = append(ret, float32(f))
	}
	return ret, nil
}

// Property parser? what if json?
func flexProp(param string) (*gorlet.FlexLayout, error) {
	props := strings.Split(param, ";")

	flex := gorlet.FlexLayout{}
	for _, p := range props {
		kv := strings.Split(p, ":")
		switch kv[0] {
		case "spacing":
			sz, err := parseFloat32Slice(kv[1])
			if err != nil {
				return nil, err
			}
			flex.Spacing = sz[0]
		case "sizes":
			sz, err := parseFloat32Slice(kv[1])
			if err != nil {
				return nil, err
			}
			flex.SetSizes(sz...)
		case "dir":
			switch kv[1] {
			case "v":
				flex.Direction = gorlet.DirectionVertical
			case "h":
				flex.Direction = gorlet.DirectionHorizontal
			}
		}
	}
	return &flex, nil
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
			f, err := flexProp(parts[1])
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
		e.Set(a.Name.Local, m32.Color(p...))
	default:
		e.Set(a.Name.Local, a.Value)
	}
	return nil
}
