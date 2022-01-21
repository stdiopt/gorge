package gorlet

import "github.com/stdiopt/gorge/math/gm"

// GridLayout creates a grid layouter that organizes children in a grid
type GridLayout struct {
	Cols    int
	Rows    int
	Spacing float32
}

// Layout implements layouter interface.
func (l *GridLayout) Layout(e *Entity) {
	sw := 1 / float32(l.Cols)
	sh := 1 / float32(l.Rows)
	for i, e := range e.Children() {
		cw := float32(i%l.Cols) / float32(l.Cols)
		ch := float32(i/l.Cols) / float32(l.Rows)
		e.SetAnchor(cw, ch, cw+sw, ch+sh)
		// s := gm.Vec4{l.Spacing / 2, l.Spacing / 2, l.Spacing / 2, l.Spacing / 2}
		s := gm.Vec4{0, 0, l.Spacing / 2, l.Spacing / 2}
		if cw == 0 {
			s[0] = 0
		} else if cw+sw == 1 {
			s[2] = 0
		}
		if ch == 0 {
			s[1] = 0
		} else if ch+sh == 1 {
			s[3] = 0
		}

		e.SetRect(s[:]...)
	}
}

// LayoutGrid creates a grid probably only needs one dimension.
func LayoutGrid(cols, rows int, spacing float32) *GridLayout {
	return &GridLayout{
		Cols:    cols,
		Rows:    rows,
		Spacing: spacing,
	}
}
