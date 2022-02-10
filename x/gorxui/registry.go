package gorxui

import "github.com/stdiopt/gorge/x/gorlet"

type buildFunc = func() gorlet.Entity

type Registry map[string]buildFunc

func (r Registry) merge(r2 Registry) Registry {
	ret := Registry{}
	for k, v := range r {
		ret[k] = v
	}
	for k, v := range r2 {
		ret[k] = v
	}
	return ret
}

var registry = Registry{
	"container":   func() gorlet.Entity { return gorlet.Container() },
	"panel":       func() gorlet.Entity { return gorlet.Panel() }, // nolint: gocritic
	"label":       func() gorlet.Entity { return gorlet.Label("") },
	"textbutton":  func() gorlet.Entity { return gorlet.TextButton("") },
	"spinner":     func() gorlet.Entity { return gorlet.Spinner("") },
	"spinnervec3": func() gorlet.Entity { return gorlet.SpinnerVec3() },
	"slider":      func() gorlet.Entity { return gorlet.Slider(0, 1) },
	"colorpicker": func() gorlet.Entity { return gorlet.ColorPicker() },
	// containers
	//"window":  func() gorlet.Func { return gorlet.Window("") },
	//"labeled": func() gorlet.Func { return gorlet.Labeled("") },
	"list": func() gorlet.Entity { return gorlet.List() },
	"flex": func() gorlet.Entity { return gorlet.Flex() },
	"grid": func() gorlet.Entity { return gorlet.Grid() },
}

// Register a Tag func in global registry.
func Register(k string, fn buildFunc) {
	registry[k] = fn
}
