package renderpl

import "github.com/stdiopt/gorge/systems/render"

type cameraSorter []render.HCamera

// Len is the number of elements in the collection.
func (s cameraSorter) Len() int { return len(s) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (s cameraSorter) Less(i int, j int) bool {
	return s[i].Camera.Camera().Order < s[j].Camera.Camera().Order
}

// Swap swaps the elements with indexes i and j.
func (s cameraSorter) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

type renderableGroupSorter []*render.RenderableGroup

func (s renderableGroupSorter) Len() int { return len(s) }

func (s renderableGroupSorter) Less(i int, j int) bool {
	return s[i].Renderable().Order < s[j].Renderable().Order
}

func (s renderableGroupSorter) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}
