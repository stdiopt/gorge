package gorlet

import (
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// GridLayout creates a grid layouter that organizes children in a grid
type GridLayout struct {
	Cols    int
	Rows    int
	Spacing float32
}

// Layout implements layouter interface.
func (l *GridLayout) Layout(ent Entity) {
	sw := 1 / float32(l.Cols)
	sh := 1 / float32(l.Rows)

	children := ent.GetEntities()
	for i, e := range children {
		rr, ok := e.(gorgeui.RectTransformer)
		if !ok {
			continue
		}
		rt := rr.RectTransform()

		cw := float32(i%l.Cols) / float32(l.Cols)
		ch := float32(i/l.Cols) / float32(l.Rows)
		rt.SetAnchor(cw, ch, cw+sw, ch+sh)
		// s := gm.Vec4{l.Spacing / 2, l.Spacing / 2, l.Spacing / 2, l.Spacing / 2}
		s := gm.Vec4{l.Spacing / 2, l.Spacing / 2, l.Spacing / 2, l.Spacing / 2}
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

		rt.SetRect(s[:]...)
	}
}

func (l *GridLayout) SetSpacing(spacing float32) {
	l.Spacing = spacing
}

func (l *GridLayout) SetDimensions(cols, rows int) {
	l.Cols = cols
	l.Rows = rows
}

// LayoutGrid creates a grid probably only needs one dimension.
func LayoutGrid(cols, rows int, spacing float32) *GridLayout {
	return &GridLayout{
		Cols:    cols,
		Rows:    rows,
		Spacing: spacing,
	}
}
