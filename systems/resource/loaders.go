package resource

import "reflect"

var resLoaders = map[loader]LoaderFunc{}

// LoaderFunc is the type of a loader func
type LoaderFunc func(res *Context, v any, names string, opts ...any) error

type loader struct {
	typ reflect.Type
	ext string
}

// Register registers a loader for a type and extension.
func Register(v any, ext string, fn LoaderFunc) {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	resLoaders[loader{typ, ext}] = fn
}

func getLoader(v any, ext string) LoaderFunc {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return resLoaders[loader{typ, ext}]
}
