package gorlet

import (
	"github.com/stdiopt/gorge/m32"
)

// ListLayout layouter that will rearrange children vertically.
type ListLayout struct {
	// Spacing between elements
	Spacing      float32
	SpacingSides m32.Vec2
	Direction    Direction
}

// Layout implements layouter
func (l *ListLayout) Layout(ent *Entity) {
	cury := float32(0)
	children := ent.Children()
	for _, e := range children {
		r := e.RectTransform()

		h := r.Dim[1]
		// Set a minimal height if 0
		if h == 0 {
			h = 3
		}
		r.SetAnchor(0, 0, 1, 0)
		r.SetRect(l.SpacingSides[0], cury, l.SpacingSides[1], h)
		cury += h + l.Spacing
	}
}

// LayoutList returns a func that will automatically layout children vertically
// based on sizes.
func LayoutList(spacing float32) *ListLayout {
	return &ListLayout{
		Spacing:   spacing,
		Direction: DirectionVertical,
	}
}
