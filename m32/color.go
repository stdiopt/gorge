package m32

// Color utils

// Color returns a vec4 with the colors value
// if the number of params is:
// 0 - returns a zero Vec4
// 1 - returns a Vec4 with all 3 elements set to param 0 and alpha 1
// 2 - returns a Vec4 with all 3 elements set to param 0 and alpha to param 1
// 3 - returns a Vec4 with all 3 elements set to params and alpha 1
// 4 - returns a Vec4 with all 4 elements set to params.
func Color(v ...float32) Vec4 {
	switch len(v) {
	case 0:
		return Vec4{}
	case 1:
		return Vec4{v[0], v[0], v[0], 1}
	case 2:
		return Vec4{v[0], v[0], v[0], v[1]}
	case 3:
		return Vec4{v[0], v[1], v[2], 1}
	default:
		return Vec4{v[0], v[1], v[2], v[3]}
	}
}

// HSV2RGB Hue Saturation Value to RGB
// values are in the range 0-1 and will be clamped to that range.
func HSV2RGB(h, s, v float32) Vec3 {
	h = Mod(h, 1)
	s = Clamp(s, 0, 1)
	v = Clamp(v, 0, 1)

	i := int(h * 6)
	f := h*6 - float32(i)
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)
	switch i % 6 {
	case 0:
		return Vec3{v, t, p}
	case 1:
		return Vec3{q, v, p}
	case 2:
		return Vec3{p, v, t}
	case 3:
		return Vec3{p, q, v}
	case 4:
		return Vec3{t, p, v}
	case 5:
		return Vec3{v, p, q}
	default:
		return Vec3{}
	}
}
