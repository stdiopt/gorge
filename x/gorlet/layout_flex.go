package gorlet

import (
	"github.com/stdiopt/gorge/math/gm"
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

func (f *FlexLayout) String() string {
	return "FlexLayout{sizeS: %v]"
}

// Layout implements layouter interface.
func (l FlexLayout) Layout(ent Entity) {
	children := ent.GetEntities()
	esum := l.sum // effective sum
	if d := len(children) - len(l.sizes); d > 0 {
		esum = l.sum + float32(d)*l.smaller
	}
	var start float32
	for i, e := range children {
		rr, ok := e.(gorgeui.RectTransformer)
		if !ok {
			continue
		}
		rt := rr.RectTransform()
		sz := l.smaller
		if i < len(l.sizes) {
			sz = l.sizes[i]
		}

		end := start + sz/esum
		switch l.Direction {
		case Horizontal:
			spacing := float32(0)
			if i != 0 {
				spacing = l.Spacing
			}
			rt.SetRect(spacing, 0, 0, 0)
			rt.SetAnchor(start, 0, end, 1)
		case Vertical:
			spacing := float32(0)
			if i != 0 {
				spacing = l.Spacing
			}
			rt.SetRect(0, spacing, 0, 0)
			rt.SetAnchor(0, start, 1, end)
		}
		start = end
	}
}

// SetSizes sets new sizes and recalculates flex.
func (l *FlexLayout) SetSizes(sizes ...float32) {
	if len(sizes) == 0 {
		return
	}
	l.sizes = sizes
	l.sum = 0
	l.smaller = sizes[0]
	for _, f := range sizes {
		l.sum += f
		l.smaller = gm.Min(l.smaller, f)
	}
}

func (l *FlexLayout) SetSpacing(spacing float32) {
	l.Spacing = spacing
}

func (l *FlexLayout) SetDirection(d Direction) {
	l.Direction = d
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
	return layoutFlex(Vertical, sizes...)
}

// LayoutFlexHorizontal automatically layout children horizontally based on sizes.
func LayoutFlexHorizontal(sizes ...float32) *FlexLayout {
	return layoutFlex(Horizontal, sizes...)
}
