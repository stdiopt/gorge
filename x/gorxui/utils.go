package gorxui

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/x/gorlet"
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
		switch strings.ToLower(pp[0]) {
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
	return gorlet.LayoutMulti(layouts...), nil
}

var (
	typOverflow  = reflect.TypeOf(gorlet.Overflow(0))
	typDirection = reflect.TypeOf(gorlet.Direction(0))
	typVec4      = reflect.TypeOf(gm.Vec4{})
)

func parseTyp(typ reflect.Type, s string) (interface{}, error) {
	switch typ {
	case typOverflow:
		switch s {
		case "hidden":
			return gorlet.OverflowHidden, nil
		case "scroll":
			return gorlet.OverflowScroll, nil
		default:
			return gorlet.OverflowVisible, nil
		}
	case typDirection:
		switch s {
		case "horizontal":
			return gorlet.Horizontal, nil
		case "vertical":
			return gorlet.Vertical, nil
		default:
			return nil, fmt.Errorf("Unknown direction %s only horizontal or vertical allowed", s)
		}
	case typVec4:
		f, err := parseFloat32Slice(s)
		if err != nil {
			return nil, err
		}
		return gm.V4(f...), nil

	}
	switch typ.Kind() {
	case reflect.Int:
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, err
		}
		return int(v), nil
	case reflect.Float32:
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, err
		}
		return float32(v), nil
	case reflect.Float64:
		return strconv.ParseFloat(s, 32)
	case reflect.Slice:
		switch typ.Elem().Kind() {
		case reflect.Float32:
			return parseFloat32Slice(s)
		default:
		}
	case reflect.String:
		return s, nil
	case reflect.Bool:
		return strconv.ParseBool(s)
	}
	panic(fmt.Sprintf("unsupported type %v", typ))
	// return nil, fmt.Errorf("unsupported type %s", typ.Kind())
}
