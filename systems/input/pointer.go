package input

// Mouse still triggers Pointer events but this would be mainly used for touch
// tbf more pointer devices are available so this can still be valid

// PointerType pointer event type
type PointerType int

func (p PointerType) String() string {
	switch p {
	case MouseDown:
		return "MouseDown"
	case MouseUp:
		return "MouseUp"
	case MouseMove:
		return "MouseMove"
	case MouseWheel:
		return "MouseWheel"
	case PointerDown:
		return "PointerDown"
	case PointerMove:
		return "PointerMove"
	case PointerEnd:
		return "PointerEnd"
	case PointerCancel:
		return "PointerCancel"
	default:
		return "<invalid>"
	}
}

// Pointer comments
const (
	_ = PointerType(iota)
	MouseDown
	MouseUp
	MouseMove
	MouseWheel
	PointerDown
	PointerMove
	PointerEnd
	PointerCancel
)
