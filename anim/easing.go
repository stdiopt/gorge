package anim

// TODO: Add the rest of known functions
// And even maybe bezier paths and stuff

// OutBounce easing
func OutBounce(t float32) float32 {
	switch {
	case t < 4/11.0:
		return (121 * t * t) / 16.0
	case t < 8/11.0:
		return (363 / 40.0 * t * t) - (99 / 10.0 * t) + 17/5.0
	case t < 9/10.0:
		return (4356 / 361.0 * t * t) - (35442 / 1805.0 * t) + 16061/1805.0
	default:
		return (54 / 5.0 * t * t) - (513 / 25.0 * t) + 268/25.0
	}
}

// InQuad takes a value from [0,1] and do Quad easing
func InQuad(t float32) float32 {
	return t * t
}

// OutQuad takes a value from [0,1] and do Quad easing
func OutQuad(t float32) float32 {
	return -t * (t - 2)
}

// InOutQuad takes a value from [0,1] and do easing
func InOutQuad(t float32) float32 {
	if t < 0.5 {
		return 2 * t * t
	}
	t = 2*t - 1
	return -0.5 * (t*(t-2) - 1)
}

// Step easing.
func Step(t float32) float32 {
	if t < .5 {
		return 0
	}
	return 1
}
