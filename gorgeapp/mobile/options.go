package mobile

import "io/fs"

// Options options for Mobile
// TODO: Might change android and ios
type Options struct {
	FS fs.FS
}
