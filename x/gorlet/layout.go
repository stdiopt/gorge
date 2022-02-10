package gorlet

import (
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// Direction for certain types of layouts
type Direction int

// Layout direction
const (
	Horizontal Direction = iota
	Vertical
)

// Layouter interface for layouting entities.
type Layouter interface {
	Layout(Entity)
}

// LayoutFunc type of func to be attached in UI element to update layout.
type LayoutFunc func(ent Entity)

// Layout implements layouter interface.
func (fn LayoutFunc) Layout(ent Entity) {
	fn(ent)
}

// MultiLayout multiple layout function
func MultiLayout(ls ...Layouter) LayoutFunc {
	return func(ent Entity) {
		for _, l := range ls {
			l.Layout(ent)
		}
	}
}

// AutoHeight be resize to content.
func AutoHeight(spacing float32) LayoutFunc {
	return func(ent Entity) {
		mainrt := ent.RectTransform()
		// Anchor Y should be Dimentional.
		mainrt.Anchor[3] = mainrt.Anchor[1]
		children := ent.GetEntities()

		dim := gm.Vec2{}
		for _, e := range children {
			rr, ok := e.(gorgeui.RectTransformer)
			if !ok {
				continue
			}
			rt := rr.RectTransform()
			rect := rt.Rect()
			h := rect[3] - rect[1] + (rt.Margin[1] + rt.Margin[3] + rt.Border[1] + rt.Border[3])

			top := rt.Position[1]
			bottom := top + h
			dim[1] = gm.Max(bottom+spacing, dim[1])

		}
		rt := ent.RectTransform()
		rt.Size[1] = dim[1] + rt.Margin[1] + rt.Margin[3] + rt.Border[1] + rt.Border[3]
	}
}

/*
func ContentSize() LayoutFunc {
	return func(ent *Entity) {
		ent.Size = gm.Vec2{}
		ent.Anchor[2] = ent.Anchor[0]
		ent.Anchor[3] = ent.Anchor[1]

		b := ent.CalcMax()

		ent.Size[0] = b[0]
		ent.Size[1] = b[1]
	}
}

func ContentSizeY() LayoutFunc {
	return func(ent *Entity) {
		ent.Size = gm.Vec2{}
		ent.Anchor[2] = ent.Anchor[0]
		// ent.Anchor[3] = ent.Anchor[1]

		b := ent.CalcMax()

		// ent.Dim[0] = b[0]
		ent.Size[1] = b[1]
	}
}

// Placement

// Vertical placement
/*func Vertical(spacing gm.Vec4, dim gm.Vec2) PlacementFunc {
	var pos gm.Vec2
	return func(w *Entity) {
		w.SetAnchor(0, 0, 1, 0)
		w.SetRect(spacing[0], spacing[1]+pos[1], spacing[2], dim[1])
		w.SetPivot(0)
		pos[1] += dim[1] + spacing[3]
	}
}*/
