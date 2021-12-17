package gorlet

import (
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// GridLayout creates a grid layouter that organizes children in a grid
type GridLayout struct {
	Cols    int
	Rows    int
	Spacing float32
}

// Layout implements layouter interface.
func (l *GridLayout) Layout(e gorgeui.Entity) {
	sw := 1 / float32(l.Cols)
	sh := 1 / float32(l.Rows)
	for i, c := range e.Element().Children() {
		cui, ok := c.(Entity)
		if !ok {
			continue
		}
		cw := float32(i%l.Cols) / float32(l.Cols)
		ch := float32(i/l.Cols) / float32(l.Rows)
		cui.RectTransform().SetAnchor(cw, ch, cw+sw, ch+sh)
		// s := m32.Vec4{l.Spacing / 2, l.Spacing / 2, l.Spacing / 2, l.Spacing / 2}
		s := m32.Vec4{0, 0, l.Spacing / 2, l.Spacing / 2}
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

		cui.RectTransform().SetRect(s[:]...)
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
