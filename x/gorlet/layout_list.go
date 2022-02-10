package gorlet

import (
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// ListLayout layouter that will rearrange children vertically.
type ListLayout struct {
	// Spacing between elements
	Spacing   float32
	Direction Direction
}

// Layout implements layouter
func (l *ListLayout) Layout(ent Entity) {
	cur := float32(0)

	children := ent.GetEntities()
	for _, e := range children {
		rr, ok := e.(gorgeui.RectTransformer)
		if !ok {
			continue
		}

		rt := rr.RectTransform()
		r := rt.Rect()

		switch l.Direction {
		case Vertical:
			rt.SetAnchor(0, 0, 1, 0)
			d := r[3] - r[1] + rt.Margin[1] + rt.Margin[3] + rt.Border[1] + rt.Border[3]
			rt.Position[1] = cur
			cur += d + l.Spacing
		case Horizontal:
			rt.SetAnchor(0, 0, 0, 1)
			d := r[2] - r[0] + rt.Margin[0] + rt.Margin[2] + rt.Border[0] + rt.Border[2]
			rt.Position[0] = cur
			cur += d + l.Spacing
		}
	}
}

func (l *ListLayout) SetSpacing(v float32) {
	l.Spacing = v
}

func (l *ListLayout) SetDirection(v Direction) {
	l.Direction = v
}

// LayoutList returns a func that will automatically layout children vertically
// based on sizes.
func LayoutList(spacing float32) *ListLayout {
	return &ListLayout{
		Spacing:   spacing,
		Direction: Vertical,
	}
}
