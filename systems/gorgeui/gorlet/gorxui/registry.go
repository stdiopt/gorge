package gorxui

import "github.com/stdiopt/gorge/systems/gorgeui/gorlet"

type buildFunc = func() gorlet.Func

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
	"container":  func() gorlet.Func { return gorlet.Container() },
	"scrollarea": func() gorlet.Func { return gorlet.Scroll() },
	"panel":      func() gorlet.Func { return gorlet.Panel() }, // nolint: gocritic
	"label":      func() gorlet.Func { return gorlet.Label("") },
	"textbutton": func() gorlet.Func { return gorlet.TextButton("", nil) },
	"spinner":    func() gorlet.Func { return gorlet.Spinner("", nil) },
	"slider":     func() gorlet.Func { return gorlet.Slider(0, 1, nil) },
	// containers
	"window":  func() gorlet.Func { return gorlet.Window("") },
	"labeled": func() gorlet.Func { return gorlet.Labeled("") },
}

// Register a Tag func in global registry.
func Register(k string, fn buildFunc) {
	registry[k] = fn
}
