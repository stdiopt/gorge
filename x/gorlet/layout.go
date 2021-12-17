package gorlet

import "github.com/stdiopt/gorge/systems/gorgeui"

// Direction for certain types of layouts
type Direction int

// Layout direction
const (
	DirectionHorizontal Direction = iota
	DirectionVertical
)

type Layouter interface {
	Layout(gorgeui.Entity)
}

// LayoutFunc type of func to be attached in UI element to update layout.
type LayoutFunc func(ent gorgeui.Entity)

// Layout implements layouter interface.
func (fn LayoutFunc) Layout(ent gorgeui.Entity) {
	fn(ent)
}

// MultiLayout multiple layout function
func MultiLayout(ls ...Layouter) LayoutFunc {
	return func(ent gorgeui.Entity) {
		for _, l := range ls {
			l.Layout(ent)
		}
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
