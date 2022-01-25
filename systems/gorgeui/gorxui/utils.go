package gorxui

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/stdiopt/gorge/systems/gorgeui/gorlet"
)

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
func parseFlexProp(param string) (*gorlet.FlexLayout, error) {
	parts := strings.Split(param, " ")
	flex := gorlet.FlexLayout{}
	for _, p := range parts {
		kv := strings.Split(p, "=")
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
		}
	}
	return &flex, nil
}

func parseLayout(param string) (gorlet.Layouter, error) {
	layouts := []gorlet.Layouter{}
	parts := strings.Split(param, ";")
	for _, p := range parts {
		pp := strings.Split(p, ":")
		switch pp[0] {
		case "flex", "hflex":
			sz, err := parseFloat32Slice(pp[1])
			if err != nil {
				return nil, err
			}

			flex := gorlet.FlexLayout{
				Direction: gorlet.Horizontal,
			}
			flex.SetSizes(sz...)
			layouts = append(layouts, flex)
		case "vflex":
			sz, err := parseFloat32Slice(pp[1])
			if err != nil {
				return nil, err
			}
			flex := gorlet.FlexLayout{}
			flex.SetSizes(sz...)
			layouts = append(layouts, flex)
		case "autoheight":
			layouts = append(layouts, gorlet.AutoHeight(0))
		case "list", "vlist":
			var spacing float32
			if len(pp) > 1 {
				s := strings.TrimSpace(pp[1])
				f, err := strconv.ParseFloat(s, 32)
				if err != nil {
					return nil, err
				}
				spacing = float32(f)
			}
			layouts = append(layouts, gorlet.LayoutList(spacing))
		}
	}
	return gorlet.MultiLayout(layouts...), nil
}

var typDirection = reflect.TypeOf(gorlet.Direction(0))

func parseTyp(typ reflect.Type, s string) (interface{}, error) {
	switch typ {
	case typDirection:
		switch s {
		case "horizontal":
			return gorlet.Horizontal, nil
		case "vertical":
			return gorlet.Vertical, nil
		default:
			return nil, fmt.Errorf("Unknown direction %s only horizontal or vertical allowed", s)
		}
	}
	switch typ.Kind() {
	case reflect.Float32:
		return strconv.ParseFloat(s, 32)
	case reflect.Float64:
		return strconv.ParseFloat(s, 32)
	case reflect.Slice:
		switch typ.Elem().Kind() {
		case reflect.Float32:
			return parseFloat32Slice(s)
		}
	case reflect.String:
		return s, nil
	case reflect.Bool:
		return strconv.ParseBool(s)
	}
	panic(fmt.Sprintf("unsupported type %v", typ))
	// return nil, fmt.Errorf("unsupported type %s", typ.Kind())
}
