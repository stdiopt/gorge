package layerfs

import (
	"io/fs"
	"os"
	"strings"
)

// FS is a list of FS that can be used to mount an overlay filesystem.
type FS []fs.FS

// Mount file system with a specific prefix.
func (f *FS) Mount(prefix string, ffs fs.FS) {
	prefix = strings.Trim(prefix, "/")
	if prefix == "" {
		*f = append(*f, ffs)
		return
	}
	*f = append(*f, Prefix(prefix, ffs))
}

// ReadDir reads the named directory from last the fs to first, the latest dir
// entry prevails.
func (f FS) ReadDir(name string) ([]fs.DirEntry, error) {
	// This check was removed to be able to use real FS, although
	/*if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: os.ErrInvalid}
	}
	name = strings.Trim(name, "/")
	*/
	entries := entrySet{}
	for i := len(f) - 1; i >= 0; i-- {
		ee, err := fs.ReadDir(f[i], name)
		if err != nil {
			continue
		}
		// found = true
		entries.Set(ee...)
	}
	return entries.list, nil
}

// Open opens the named file from the last fs to first, the latest file
// prevails.
func (f FS) Open(name string) (fs.File, error) {
	/*
		if !fs.ValidPath(name) {
			return nil, &fs.PathError{Op: "open", Path: name, Err: os.ErrInvalid}
		}
		name = strings.Trim(name, "/")
	*/
	for i := len(f) - 1; i >= 0; i-- {
		fl, err := f[i].Open(name)
		if err == nil { // inverse err check
			return fl, nil
		}
	}
	return nil, &fs.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}
