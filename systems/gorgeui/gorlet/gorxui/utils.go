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

func parseTyp(typ reflect.Type, s string) (interface{}, error) {
	switch typ.Kind() {
	case reflect.String:
		return s, nil
	case reflect.Bool:
		return strconv.ParseBool(s)
	}
	panic(fmt.Sprintf("unsupported type %v", typ))
	// return nil, fmt.Errorf("unsupported type %s", typ.Kind())
}
