// Package static contains resources like default shaders, fonts
package static

import (
	"embed"
	"io/fs"
)

// Assets embed static assets
//go:embed src
var Assets embed.FS

var sfs = func() fs.FS {
	sf, err := fs.Sub(Assets, "src")
	if err != nil {
		panic(err)
	}
	return sf
}()

// Data returns the embed data based on name.
func Data(name string) ([]byte, error) {
	return fs.ReadFile(sfs, name)
}

// MustData returns data based on name or panics.
func MustData(name string) []byte {
	data, err := Data(name)
	if err != nil {
		panic(err)
	}
	return data
}
