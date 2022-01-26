package gorlet

import (
	"github.com/stdiopt/gorge/math/gm"
)

func Rect(v ...float32) gm.Vec4 {
	switch len(v) {
	case 0:
		return gm.Vec4{0, 0, 0, 0}
	case 1:
		return gm.Vec4{v[0], v[0], v[0], v[0]}
	case 2:
		return gm.Vec4{v[0], v[0], v[1], v[1]}
	case 3:
		return gm.Vec4{v[0], v[1], v[2], v[2]}
	default:
		return gm.Vec4{v[0], v[1], v[2], v[3]}
	}
}

func Anchor(v ...float32) gm.Vec4 {
	switch len(v) {
	case 1:
		return gm.Vec4{v[0], v[0], v[0], v[0]}
	case 2:
		return gm.Vec4{v[0], v[1], v[0], v[1]}
	case 3:
		return gm.Vec4{v[0], v[1], v[2], v[1]}
	default:
		return gm.Vec4{v[0], v[1], v[2], v[3]}
	}
}

func Margin(v ...float32) gm.Vec4 {
	switch len(v) {
	case 1:
		return gm.Vec4{v[0], v[0], v[0], v[0]}
	case 2:
		return gm.Vec4{v[0], v[1], v[0], v[1]}
	case 3:
		return gm.Vec4{v[0], v[1], v[2], v[2]}
	default:
		return gm.Vec4{v[0], v[1], v[2], v[3]}
	}
}
