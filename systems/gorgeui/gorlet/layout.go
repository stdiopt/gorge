package gorlet

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

// Layouter interface for layouting entities.
type Layouter interface {
	Layout(*Entity)
}

// LayoutFunc type of func to be attached in UI element to update layout.
type LayoutFunc func(ent *Entity)

// Layout implements layouter interface.
func (fn LayoutFunc) Layout(ent *Entity) {
	fn(ent)
}

// MultiLayout multiple layout function
func MultiLayout(ls ...Layouter) LayoutFunc {
	return func(ent *Entity) {
		for _, l := range ls {
			l.Layout(ent)
		}
	}
}

// AutoHeight be resize to content.
func AutoHeight(spacing float32) LayoutFunc {
	return func(ent *Entity) {
		// Anchor Y should be Dimentional.
		ent.Anchor[3] = ent.Anchor[1]
		children := ent.Children()

		dim := m32.Vec2{}
		for _, c := range children {
			rect := c.Rect()
			h := rect[3] - rect[1] + (c.Margin[1] + c.Margin[3])

			top := c.Position[1]
			bottom := top + h
			dim[1] = m32.Max(bottom+spacing, dim[1])

		}
		rt := ent.RectTransform()
		rt.Dim[1] = dim[1] + rt.Margin[1] + rt.Margin[3]
	}
}

func ContentSize() LayoutFunc {
	return func(ent *Entity) {
		ent.Dim = m32.Vec2{}
		ent.Anchor[2] = ent.Anchor[0]
		ent.Anchor[3] = ent.Anchor[1]

		b := ent.CalcBounds()

		ent.Dim[0] = b[2]
		ent.Dim[1] = b[3]
	}
}

// Placement

// Vertical placement
/*func Vertical(spacing m32.Vec4, dim m32.Vec2) PlacementFunc {
	var pos m32.Vec2
	return func(w *Entity) {
		w.SetAnchor(0, 0, 1, 0)
		w.SetRect(spacing[0], spacing[1]+pos[1], spacing[2], dim[1])
		w.SetPivot(0)
		pos[1] += dim[1] + spacing[3]
	}
}*/
