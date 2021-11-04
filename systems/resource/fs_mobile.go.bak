//go:build !js && !wasm

package resource

import (
	"io"
	"io/fs"
	"os"

	"github.com/stdiopt/gorge/gorgeapp/mobile/asset"
)

// MobileFS loads files using gomobile
type MobileFS struct{}

// Open stuff
func (l MobileFS) Open(p string) (fs.File, error) {
	r, err := asset.Open(p)
	if err != nil {
		return nil, err
	}
	return mobileFile{r}, nil
}

type mobileFile struct {
	io.ReadCloser
}

func (f mobileFile) Stat() (fs.FileInfo, error) {
	return nil, os.ErrInvalid
}
