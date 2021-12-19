package gorlet

import (
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// FlexLayout will redimension children based on sizes.
type FlexLayout struct {
	Spacing   float32
	Direction Direction

	sizes   []float32
	smaller float32
	sum     float32
}

// Layout implements layouter interface.
func (l FlexLayout) Layout(ent gorgeui.Entity) {
	children := ent.Element().Children()
	esum := l.sum // effective sum
	if d := len(children) - len(l.sizes); d > 0 {
		esum = l.sum + float32(d)*l.smaller
	}
	// log.Println("Smaller:", l.smaller)
	// log.Println("Sizes:", l.sizes)
	var start float32
	for i, e := range children {
		rt, ok := e.(interface{ RectTransform() *gorgeui.RectComponent })
		if !ok {
			continue
		}
		r := rt.RectTransform()
		sz := l.smaller
		if i < len(l.sizes) {
			sz = l.sizes[i]
		}

		end := start + sz/esum
		switch l.Direction {
		case DirectionHorizontal:
			r.SetAnchor(start, 0, end, 1)
		case DirectionVertical:
			r.SetAnchor(0, start, 1, end)
		}
		// Do we need to set rect here?
		// r.SetRect(l.Spacing)
		start = end
	}
}

// SetSizes sets new sizes and recalculates flex.
func (l *FlexLayout) SetSizes(sizes ...float32) {
	l.sizes = sizes
	l.sum = 0
	l.smaller = sizes[0]
	for _, f := range sizes {
		l.sum += f
		l.smaller = m32.Min(l.smaller, f)
	}
}

// layoutFlex layout
func layoutFlex(dir Direction, sizes ...float32) *FlexLayout {
	l := &FlexLayout{
		Spacing:   .3,
		Direction: dir,
	}
	l.SetSizes(sizes...)
	return l
}

// LayoutFlexVertical automatically layout children vertically based on sizes.
func LayoutFlexVertical(sizes ...float32) *FlexLayout {
	return layoutFlex(DirectionVertical, sizes...)
}

// LayoutFlexHorizontal automatically layout children horizontally based on sizes.
func LayoutFlexHorizontal(sizes ...float32) *FlexLayout {
	return layoutFlex(DirectionHorizontal, sizes...)
}
