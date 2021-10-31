package resource

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
)

// HTTPFS asset that loads based on a starting url
type HTTPFS struct {
	BaseURL string
}

// Open the asset
func (l HTTPFS) Open(p string) (fs.File, error) {
	url := l.BaseURL + p
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("err: %q %q", url, res.Status)
	}
	return httpFile{res.Body}, nil
}

type httpFile struct {
	io.ReadCloser
}

func (h httpFile) Stat() (fs.FileInfo, error) {
	return nil, os.ErrInvalid
}
