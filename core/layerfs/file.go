package layerfs

import (
	"io"
	"io/fs"
)

type file struct {
	fileInfo fileInfo
}

func (f file) Stat() (fs.FileInfo, error) {
	return f.fileInfo, nil
}

func (f file) Read(buf []byte) (int, error) {
	return 0, io.EOF
}

func (f file) Close() error {
	return nil
}
