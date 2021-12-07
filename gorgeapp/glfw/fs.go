package glfw

import (
	"io/fs"
	"os"
)

// RootFS wraps os.Open in a fs.FS.
type RootFS struct{}

// Open opens the named file.
//
// When Open returns an error, it should be of type *PathError
// with the Op field set to "open", the Path field set to name,
// and the Err field describing the problem.
//
// Open should reject attempts to open names that do not satisfy
// ValidPath(name), returning a *PathError with Err set to
// ErrInvalid or ErrNotExist.
func (r RootFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}
