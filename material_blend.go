package gorge

// BlendType for material
// TODO: Fix this blending stuff with src and dst for Func maybe
type BlendType uint32

const (
	// BlendOneOneMinusSrcAlpha - gl.ONE, gl.ONE_MINUS_SRC_ALPHA
	BlendOneOneMinusSrcAlpha = BlendType(iota)
	// BlendOneOne - gl.ONE, gl.ONE
	BlendOneOne
	// BlendDisable - disable blending
	BlendDisable
)

func (b BlendType) String() string {
	switch b {
	case BlendOneOneMinusSrcAlpha:
		return "BlendOneOneMinusSrcAlpha"
	case BlendOneOne:
		return "BlendOneOne"
	case BlendDisable:
		return "BlendDisable"
	default:
		return "BlendTypeUnknown"
	}
}
