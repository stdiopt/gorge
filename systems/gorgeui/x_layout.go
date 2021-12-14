package gorgeui

import (
	"github.com/stdiopt/gorge/m32"
)

// Direction for certain types of layouts
type Direction int

// Layout direction
const (
	DirectionHorizontal Direction = iota
	DirectionVertical
)

// Layouter interface that can be used to layout children.
type Layouter interface {
	Layout(ent Entity)
}

// LayoutFunc type of func to be attached in UI element to update layout.
type LayoutFunc func(ent Entity)

// Layout implements layouter interface.
func (fn LayoutFunc) Layout(ent Entity) {
	fn(ent)
}

// AutoHeight be resize to content.
func AutoHeight(spacing float32) LayoutFunc {
	return func(ent Entity) {
		el := ent.Element()

		dim := m32.Vec2{}
		children := el.Children()
		for _, c := range children {
			rt, ok := c.(rectTransformer)
			if !ok {
				continue
			}
			r := rt.RectTransform()

			top := r.Position[1]
			bottom := top + r.Dim[1]
			dim[1] = m32.Max(bottom+spacing, dim[1])

		}
		ent.RectTransform().Dim[1] = dim[1]
	}
}

// ListVertical returns a layoutfunc that will arrange children vertically.
func ListVertical(spacing float32) LayoutFunc {
	return List(DirectionVertical, spacing)
}

// List automatically layout children vertically.
func List(dir Direction, spacing float32) LayoutFunc {
	padding := m32.Vec4{spacing, spacing, spacing, spacing}
	return func(ent Entity) {
		children := ent.Element().Children()
		cur := padding[0]
		for _, c := range children {
			r, ok := c.(rectTransformer)
			if !ok {
				continue
			}
			rt := r.RectTransform()
			var size float32
			switch dir {
			case DirectionVertical:
				rt.SetAnchor(0, 0, 1, 0)
				rt.SetPivot(0, 0)
				rect := rt.Rect()
				size = rect[3] - rect[0]
				rt.SetRect(padding[0], cur, padding[1], size)
			case DirectionHorizontal:
				rt.SetAnchor(0, 0, 0, 1)
				rt.SetPivot(0, 0)
				size = rt.Dim[1]
				rt.SetRect(padding[0], cur, padding[1], size)
			}
			cur += size + spacing
			continue
		}
	}
}

// FlexVertical automatically layout children vertically based on sizes.
func FlexVertical(sizes ...float32) LayoutFunc {
	return Flex(DirectionVertical, sizes...)
}

// FlexHorizontal automatically layout children horizontally based on sizes.
func FlexHorizontal(sizes ...float32) LayoutFunc {
	return Flex(DirectionHorizontal, sizes...)
}

// TODO: {lpf} Transform this to am interface{} as we might want to change size
// params live.

// Flex layout
func Flex(dir Direction, sizes ...float32) LayoutFunc {
	var sum float32
	smaller := sizes[0]
	for _, f := range sizes {
		sum += f
		smaller = m32.Min(smaller, f)
	}
	// spacing := m32.Vec4{1, 1, 1, 1}
	spacing := m32.Vec4{0}
	return func(ent Entity) {
		children := ent.Element().Children()
		esum := sum // effective sum
		if d := len(children) - len(sizes); d > 0 {
			esum = sum + float32(d)*smaller
			// log.Println("Sum is:", sum+float32(d)*smaller)
		}
		var start float32
		for i, e := range children {
			rt, ok := e.(rectTransformer)
			if !ok {
				continue
			}
			r := rt.RectTransform()
			sz := smaller
			if i < len(sizes) {
				sz = sizes[i]
			}

			end := start + sz/esum
			// log.Println("Size:", sz)
			switch dir {
			case DirectionHorizontal:
				r.SetAnchor(start, 0, end, 1)
			case DirectionVertical:
				r.SetAnchor(0, start, 1, end)
			}
			r.SetRect(spacing[0], spacing[1], spacing[2], spacing[3])
			start = end
		}
	}
}

// GridLayout creates a grid layouter that organizes children in a grid
type GridLayout struct {
	Cols    int
	Rows    int
	Spacing float32
}

// Layout implements layouter interface.
func (l *GridLayout) Layout(e Entity) {
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

// MultiLayout multiple layout function
func MultiLayout(ls ...Layouter) LayoutFunc {
	return func(ent Entity) {
		for _, l := range ls {
			l.Layout(ent)
		}
	}
}
