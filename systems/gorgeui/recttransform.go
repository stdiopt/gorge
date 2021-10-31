package gorgeui

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
)

// Unity
// Regular PosX, PosY, PosZ
// Anchor Min Max Vec2 // This affects parenting
// Pivot means the center of the Widget based on our Width and Height

// RectComponent data component based on transform with fields specific for UI
// elements.
type RectComponent struct {
	gorge.TransformComponent

	Dim    m32.Vec2
	Anchor m32.Vec4 // left, bottom, right, top

	Pivot m32.Vec2
}

// RectIdent returns a identity rect transform.
func RectIdent() RectComponent {
	return RectComponent{
		TransformComponent: gorge.TransformIdent(),
		Anchor:             m32.Vec4{0, 0, 1, 1},
		Dim:                m32.Vec2{0, 0},
		Pivot:              m32.Vec2{.5, .5},
	}
}

// NewRectComponent returns a new Rect Transform.
func NewRectComponent() *RectComponent {
	r := RectIdent()
	return &r
}

// RectTransform implements the RectTransform and returns self
func (r *RectComponent) RectTransform() *RectComponent { return r }

// SetRect sets the rect.
func (r *RectComponent) SetRect(v ...float32) {
	v4 := v4f(v...)
	r.Position[0] = v4[0]
	r.Position[1] = v4[1]

	r.Dim[0] = v4[2]
	r.Dim[1] = v4[3]
}

// SetAnchor sets anchor.
func (r *RectComponent) SetAnchor(v ...float32) {
	r.Anchor = v4f(v...)
}

// SetPivot sets the pivot.
func (r *RectComponent) SetPivot(v ...float32) {
	r.Pivot = v2f(v...)
}

// Rect calculate and returns the rect.
func (r *RectComponent) Rect() m32.Vec4 {
	parentRect := m32.Vec4{}
	if p := r.Parent(); p != nil {
		if c, ok := p.(interface{ Rect() m32.Vec4 }); ok {
			parentRect = c.Rect()
		}
	}
	return r.RelativeRect(parentRect)
}

// RelativeRect calculate rect based on other.
func (r *RectComponent) RelativeRect(parentRect m32.Vec4) m32.Vec4 {
	parentW := parentRect[2] - parentRect[0]
	parentH := parentRect[3] - parentRect[1]

	var left, top, right, bottom float32

	left = parentRect[0] + r.Anchor[0]*parentW
	bottom = parentRect[1] + r.Anchor[1]*parentH

	// If anchor min and max are the same we use pivot
	if r.Anchor[0] == r.Anchor[2] {
		w := r.Dim[0]
		left -= r.Pivot[0] * w
		right = left + w
	} else {
		right = -r.Dim[0] + parentRect[0] + r.Anchor[2]*parentW - r.Position[0]
	}

	if r.Anchor[1] == r.Anchor[3] {
		h := r.Dim[1]
		bottom -= r.Pivot[1] * h
		top = bottom + h
	} else {
		top = -r.Dim[1] + parentRect[1] + r.Anchor[3]*parentH - r.Position[1]
	}

	return m32.Vec4{left, bottom, right, top}
}

// ApplyTo will copy this rect transform to the specific entity
func (r RectComponent) ApplyTo(e Entity) {
	// copy Rect transform
	*e.RectTransform() = r
}
