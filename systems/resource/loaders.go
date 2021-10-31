package resource

import "reflect"

var resLoaders = map[loader]LoaderFunc{}

// LoaderFunc is the type of a loader func
type LoaderFunc func(res *Context, v interface{}, names string, opts ...interface{}) error

type loader struct {
	typ reflect.Type
	ext string
}

// Register registers a loader for a type and extension.
func Register(v interface{}, ext string, fn LoaderFunc) {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	resLoaders[loader{typ, ext}] = fn
}

func getLoader(v interface{}, ext string) LoaderFunc {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return resLoaders[loader{typ, ext}]
}
