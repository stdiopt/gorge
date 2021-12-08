package gorlet

import (
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// FlexVertical automatically layout children vertically based on sizes.
func FlexVertical(sizes ...float32) gorgeui.LayoutFunc {
	return Flex(gorgeui.DirectionVertical, sizes...)
}

// FlexHorizontal automatically layout children horizontally based on sizes.
func FlexHorizontal(sizes ...float32) gorgeui.LayoutFunc {
	return Flex(gorgeui.DirectionHorizontal, sizes...)
}

// TODO: {lpf} Transform this to am interface{} as we might want to change size
// params live.

// Flex layout
func Flex(dir gorgeui.Direction, sizes ...float32) gorgeui.LayoutFunc {
	var sum float32
	smaller := sizes[0]
	for _, f := range sizes {
		sum += f
		smaller = m32.Min(smaller, f)
	}
	// spacing := m32.Vec4{1, 1, 1, 1}
	spacing := m32.Vec4{.3, .3, .3, .3}
	return func(ent gorgeui.Entity) {
		children := ent.Element().Children()
		esum := sum // effective sum
		if d := len(children) - len(sizes); d > 0 {
			esum = sum + float32(d)*smaller
			// log.Println("Sum is:", sum+float32(d)*smaller)
		}
		var start float32
		for i, e := range children {
			rt, ok := e.(interface{ RectTransform() *gorgeui.RectComponent })
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
			case gorgeui.DirectionHorizontal:
				r.SetAnchor(start, 0, end, 1)
			case gorgeui.DirectionVertical:
				r.SetAnchor(0, start, 1, end)
			}
			r.SetRect(spacing[0], spacing[1], spacing[2], spacing[3])
			start = end
		}
	}
}
