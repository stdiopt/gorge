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

// LayoutFunc type of func to be attached in UI element to update layout.
type LayoutFunc func(ent Entity)

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
				size = rt.Dim[1]
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

// MultiLayout multiple layout function
func MultiLayout(fns ...LayoutFunc) LayoutFunc {
	return func(ent Entity) {
		for _, fn := range fns {
			fn(ent)
		}
	}
}
