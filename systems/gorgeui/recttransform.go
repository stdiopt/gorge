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
	parent gorge.Matrixer

	Rotation m32.Quat
	Position m32.Vec3
	Scale    m32.Vec3

	Dim m32.Vec2
	// New
	Margin m32.Vec4

	Anchor m32.Vec4 // left, bottom, right, top
	Pivot  m32.Vec2
}

// RectIdent returns a identity rect transform.
func RectIdent() RectComponent {
	// to be able to work with anchors.
	rc := RectComponent{
		Rotation: m32.QIdent(),
		Scale:    m32.Vec3{1, 1, 1},

		Anchor: m32.Vec4{0, 0, 1, 1},
		Dim:    m32.Vec2{0, 0},
		Pivot:  m32.Vec2{.5, .5},
	}
	return rc
}

// NewRectComponent returns a new Rect Transform.
func NewRectComponent() *RectComponent {
	c := RectIdent()
	return &c
}

// RectTransform implements the RectTransform and returns self
func (c *RectComponent) RectTransform() *RectComponent { return c }

// SetParent experiment to relativize anchor.
func (c *RectComponent) SetParent(t gorge.Matrixer) {
	c.parent = t
}

// Parent returns parent sub transform.
func (c *RectComponent) Parent() gorge.Matrixer {
	return c.parent
}

// Mat4 returns 4x4 transform matrix.
func (c *RectComponent) Mat4() m32.Mat4 {
	rect := c.parentRect() // Parent Dim
	anchor := m32.Vec2{
		rect[0] + (rect[2]-rect[0])*c.Anchor[0],
		rect[1] + (rect[3]-rect[1])*c.Anchor[1],
	}
	pivot := m32.Vec2{}
	if m32.FloatEqual(c.Anchor[0], c.Anchor[2]) {
		pivot[0] = -c.Dim[0] * c.Pivot[0]
	}
	if m32.FloatEqual(c.Anchor[1], c.Anchor[3]) {
		pivot[1] = -c.Dim[1] * c.Pivot[1]
	}

	pos := c.Position.
		Add(anchor.Vec3(0)).
		Add(m32.Vec3{c.Margin[0], c.Margin[1], 0})
	m := m32.Translate3D(pos[0], pos[1], pos[2])

	m = m.Mul(c.Rotation.Mat4())
	m = m.Mul(m32.Scale3D(c.Scale[0], c.Scale[1], c.Scale[2]))

	m = m.Mul(m32.Translate3D(pivot[0], pivot[1], 0))

	if c.parent != nil {
		return c.parent.Mat4().Mul(m)
	}
	return m
}

// WorldPosition returns the world position.
func (c *RectComponent) WorldPosition() m32.Vec3 {
	return c.Mat4().Col(3).Vec3()
}

// SetRect with position x,y and dimensions w,h
func (c *RectComponent) SetRect(vs ...float32) {
	v := v4f(vs...)
	c.Position[0] = v[0]
	c.Position[1] = v[1]

	c.Dim[0] = v[2]
	c.Dim[1] = v[3]
}

func (c *RectComponent) SetMargin(vs ...float32) {
	v := v4f(vs...)
	c.Margin[0] = v[0]
	c.Margin[1] = v[1]
	c.Margin[2] = v[2]
	c.Margin[3] = v[3]
}

// SetWidth sets the width.
func (c *RectComponent) SetWidth(w float32) {
	c.Dim[0] = w
}

// SetHeight sets the height.
func (c *RectComponent) SetHeight(w float32) {
	c.Dim[1] = w
}

// SetSize sets the rect dimentions.
func (c *RectComponent) SetSize(v ...float32) {
	c.Dim = m32.V2(v...)
}

// SetAnchor sets anchor.
func (c *RectComponent) SetAnchor(v ...float32) {
	switch len(v) {
	case 1:
		c.Anchor = m32.Vec4{v[0], v[0], v[0], v[0]}
	case 2:
		c.Anchor = m32.Vec4{v[0], v[1], v[0], v[1]}
	case 3:
		c.Anchor = m32.Vec4{v[0], v[1], v[2], v[1]}
	case 4:
		c.Anchor = m32.Vec4(*(*[4]float32)(v))
	default:
		panic("wrong number of params")

	}
}

// SetPivot sets the pivot.
func (c *RectComponent) SetPivot(v ...float32) {
	c.Pivot = v2f(v...)
}

// SetScale sets the scale.
func (c *RectComponent) SetScale(sz ...float32) {
	switch len(sz) {
	case 1:
		c.Scale[0], c.Scale[1], c.Scale[2] = sz[0], sz[0], sz[0]
	case 2, 3:
		copy(c.Scale[:], sz)
	default:
		panic("wrong number of params")
	}
}

// SetPosition sets the position.
func (c *RectComponent) SetPosition(x, y, z float32) {
	c.Position = m32.Vec3{x, y, z}
}

// SetRotation sets the rotation.
func (c *RectComponent) SetRotation(q m32.Quat) {
	c.Rotation = q
}

// Rotate axis
func (c *RectComponent) Rotate(x, y, z float32) {
	c.Rotation = c.Rotation.Mul(m32.QFromAngles(
		x, y, z,
		m32.XYZ,
	))
}

// Translate translates the entity.
func (c *RectComponent) Translate(x, y, z float32) {
	c.Position = c.Position.Add(m32.Vec3{x, y, z})
}

// This should be called Dim which are the dimentions, lefttop will always be 0,0

func (c *RectComponent) CalcSize() m32.Vec2 {
	return c.RelativeSize(c.parentSize())
}

func (c *RectComponent) RelativeSize(parentDim m32.Vec2) m32.Vec2 {
	var right, bottom float32
	// We might discard rect
	right = c.Dim[0]
	bottom = c.Dim[1]
	// If anchor min and max are the same we use pivot
	if c.Anchor[0] != c.Anchor[2] {
		w := parentDim[0] // parentRect[0] is always 0 now
		w -= w*(1-c.Anchor[2]) + w*(c.Anchor[0])
		right = w - c.Dim[0] - c.Position[0]
	}

	if c.Anchor[1] != c.Anchor[3] {
		h := parentDim[1]
		h -= h*(1-c.Anchor[3]) + h*(c.Anchor[1])
		bottom = h - c.Dim[1] - c.Position[1]
	}
	return m32.Vec2{
		right - c.Margin[2] - c.Margin[0],
		bottom - c.Margin[1] - c.Margin[3],
	}
}

func (c *RectComponent) parentSize() m32.Vec2 {
	if p, ok := c.parent.(interface{ CalcSize() m32.Vec2 }); ok {
		return p.CalcSize()
	}
	return m32.Vec2{}
}

// TODO: this will be deprecated

// Rect calculate and returns the rect.
func (c *RectComponent) Rect() m32.Vec4 {
	return c.RelativeRect(c.parentRect())
}

// RelativeRect calculate rect based on other.
func (c *RectComponent) RelativeRect(parentRect m32.Vec4) m32.Vec4 {
	var left, top, right, bottom float32
	// We might discard rect
	left = 0 // parentRect[0] //+ c.Anchor[0]*parentW
	top = 0  // parentRect[1]  //+ c.Anchor[1]*parentH

	right = c.Dim[0]
	bottom = c.Dim[1]
	// If anchor min and max are the same we use pivot
	if c.Anchor[0] != c.Anchor[2] {
		w := parentRect[2] - parentRect[0] // parentRect[0] is always 0 now
		// reduce rect by the relative anchor from both sides
		w -= w*(1-c.Anchor[2]) + w*(c.Anchor[0])
		right = w - c.Dim[0] - c.Position[0]
	}

	if c.Anchor[1] != c.Anchor[3] {
		h := parentRect[3] - parentRect[1]
		h -= h*(1-c.Anchor[3]) + h*(c.Anchor[1])
		// reduce rect by the relative anchor from both sides
		bottom = h - c.Dim[1] - c.Position[1]
	}
	return m32.Vec4{
		left,
		top,
		right - c.Margin[2] - c.Margin[0],
		// bottom - c.Margin[3] - c.Margin[1],
		bottom - c.Margin[1] - c.Margin[3],
	}
}

func (c *RectComponent) parentRect() m32.Vec4 {
	if p, ok := c.parent.(interface{ Rect() m32.Vec4 }); ok {
		return p.Rect()
	}
	return m32.Vec4{}
}
