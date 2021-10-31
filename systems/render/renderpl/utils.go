package renderpl

import "github.com/stdiopt/gorge/systems/render"

type cameraSorter []render.Camera

// Len is the number of elements in the collection.
func (c cameraSorter) Len() int { return len(c) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (c cameraSorter) Less(i int, j int) bool {
	return c[i].Camera().Order < c[j].Camera().Order
}

// Swap swaps the elements with indexes i and j.
func (c cameraSorter) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}
