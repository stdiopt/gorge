package wasm

import "io/fs"

// Options options for Wasm
// TODO: Might change android and ios
type Options struct {
	FS fs.FS
}
