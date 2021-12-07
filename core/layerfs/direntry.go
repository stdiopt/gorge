package layerfs

import "io/fs"

type dirEntry struct {
	fileInfo
}

// Type returns the type bits for the entry.
// The type bits are a subset of the usual FileMode bits, those returned by the FileMode.Type method.
func (d dirEntry) Type() fs.FileMode {
	return d.Mode()
}

// Info returns the FileInfo for the file or subdirectory described by the entry.
// The returned FileInfo may be from the time of the original directory read
// or from the time of the call to Info. If the file has been removed or renamed
// since the directory read, Info may return an error satisfying errors.Is(err, ErrNotExist).
// If the entry denotes a symbolic link, Info reports the information about the link itself,
// not the link's target.
func (d dirEntry) Info() (fs.FileInfo, error) {
	return d.fileInfo, nil
}

type entrySet struct {
	set  map[string]struct{}
	list []fs.DirEntry
}

func (e *entrySet) Set(ds ...fs.DirEntry) {
	if e.set == nil {
		e.set = make(map[string]struct{})
	}

	for _, d := range ds {
		if _, ok := e.set[d.Name()]; ok {
			return
		}
		e.set[d.Name()] = struct{}{}
		e.list = append(e.list, d)
	}
}
