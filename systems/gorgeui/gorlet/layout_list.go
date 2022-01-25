package gorlet

// ListLayout layouter that will rearrange children vertically.
type ListLayout struct {
	// Spacing between elements
	Spacing   float32
	Direction Direction
}

// Layout implements layouter
func (l *ListLayout) Layout(ent *Entity) {
	cury := float32(0)
	children := ent.Children()
	for _, e := range children {
		rt := e.RectTransform()
		rt.SetAnchor(0, 0, 1, 0)

		r := rt.Rect()

		h := r[3] - r[1] + rt.Margin[1] + rt.Margin[3]
		rt.Position[1] = cury
		cury += h + l.Spacing
	}
}

// LayoutList returns a func that will automatically layout children vertically
// based on sizes.
func LayoutList(spacing float32) *ListLayout {
	return &ListLayout{
		Spacing:   spacing,
		Direction: Vertical,
	}
}
