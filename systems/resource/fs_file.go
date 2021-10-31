package resource

import (
	"io/fs"
	"os"
	"path/filepath"
)

// FileFS file system implementation
// opened files will be prefixed with BasePath
type FileFS struct {
	BasePath string
}

// Open File from filesystem prefixed by BasePath.
func (l FileFS) Open(p string) (fs.File, error) {
	path := filepath.Join(l.BasePath, p)
	return os.Open(path)
}
