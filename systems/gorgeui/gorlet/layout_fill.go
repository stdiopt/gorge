package gorlet

// Vertical only
type FillLayout struct {
	Spacing float32
}

func (f FillLayout) Layout(e *Entity) {
	children := e.Children()
	var childHeight float32
	var fill int

	for _, c := range children {
		c.SetAnchor(0, 0, 1, 0)
		if c.fill {
			fill++
			continue
		}
		childHeight += c.ContentSize()[1]
	}

	msz := e.ContentSize()

	spacingSum := f.Spacing * float32(len(children)-1)
	left := msz[1] - spacingSum - childHeight
	off := float32(0)
	for _, c := range children {
		sz := c.ContentSize()
		if c.fill {
			sz[1] = left / float32(fill)
			c.SetHeight(sz[1])
		}
		c.Position[1] = off
		off += sz[1] + f.Spacing
	}
}

func LayoutFill(spacing float32) FillLayout {
	return FillLayout{Spacing: spacing}
}
