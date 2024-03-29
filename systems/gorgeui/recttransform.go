package gorgeui

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
)

// Unity
// Regular PosX, PosY, PosZ
// Anchor Min Max Vec2 // This affects parenting
// Pivot means the center of the Widget based on our Width and Height

// RectComponent data component based on transform with fields specific for UI
// elements.
type RectComponent struct {
	parent gorge.Matrixer

	Rotation gm.Quat
	Position gm.Vec3
	Scale    gm.Vec3

	Size gm.Vec2
	// New
	Margin gm.Vec4
	Border gm.Vec4

	Anchor gm.Vec4 // left, bottom, right, top
	Pivot  gm.Vec2
}

// RectIdent returns a identity rect transform.
func RectIdent() RectComponent {
	// to be able to work with anchors.
	rc := RectComponent{
		Rotation: gm.QIdent(),
		Scale:    gm.Vec3{1, 1, 1},

		Anchor: gm.Vec4{0, 0, 1, 1},
		Size:   gm.Vec2{0, 0},
		Pivot:  gm.Vec2{.5, .5},
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
func (c *RectComponent) Mat4() gm.Mat4 {
	rect := c.parentRect() // Parent Dim
	anchor := gm.Vec2{
		rect[0] + (rect[2]-rect[0])*c.Anchor[0],
		rect[1] + (rect[3]-rect[1])*c.Anchor[1],
	}
	pivot := gm.Vec2{}
	if gm.FloatEqual(c.Anchor[0], c.Anchor[2]) {
		pivot[0] = -c.Size[0] * c.Pivot[0]
	}
	if gm.FloatEqual(c.Anchor[1], c.Anchor[3]) {
		pivot[1] = -c.Size[1] * c.Pivot[1]
	}

	pos := c.Position.
		Add(anchor.Vec3(0)).
		Add(gm.Vec3{c.Margin[0], c.Margin[1], 0}).
		Add(gm.Vec3{c.Border[0], c.Border[1], 0})

	m := gm.Translate3D(pos[0], pos[1], pos[2])
	m = m.Mul(c.Rotation.Mat4())
	m = m.Mul(gm.Scale3D(c.Scale[0], c.Scale[1], c.Scale[2]))

	m = m.Mul(gm.Translate3D(pivot[0], pivot[1], 0))

	if c.parent != nil {
		return c.parent.Mat4().Mul(m)
	}
	return m
}

// WorldPosition returns the world position.
func (c *RectComponent) WorldPosition() gm.Vec3 {
	return c.Mat4().Col(3).Vec3()
}

// SetRect with position x,y and dimensions w,h
func (c *RectComponent) SetRect(vs ...float32) {
	v := v4f(vs...)
	c.Position[0] = v[0]
	c.Position[1] = v[1]

	c.SetSize(v[2], v[3])
	// c.Size[0] = v[2]
	// c.Size[1] = v[3]
}

func (c *RectComponent) SetBorder(vs ...float32) {
	c.Border = v4f(vs...)
}

func (c *RectComponent) SetBorderv(v gm.Vec4) {
	c.Border = v
}

func (c *RectComponent) SetMargin(vs ...float32) {
	switch len(vs) {
	case 0:
		c.Margin = gm.Vec4{}
	case 1:
		c.Margin = gm.Vec4{vs[0], vs[0], vs[0], vs[0]}
	case 2:
		c.Margin = gm.Vec4{vs[0], vs[1], vs[0], vs[1]}
	case 3:
		c.Margin = gm.Vec4{vs[0], vs[1], vs[2], vs[2]}
	default:
		c.Margin = gm.Vec4{vs[0], vs[1], vs[2], vs[3]}

	}
}

// SetWidth sets the width.
func (c *RectComponent) SetWidth(w float32) {
	c.Size[0] = w
}

// SetHeight sets the height.
func (c *RectComponent) SetHeight(w float32) {
	c.Size[1] = w
}

// SetSize sets the rect dimentions.
func (c *RectComponent) SetSize(v ...float32) {
	c.Size = gm.V2(v...)
}

// SetAnchor sets anchor.
func (c *RectComponent) SetAnchor(v ...float32) {
	switch len(v) {
	case 0:
		c.Anchor = gm.Vec4{}
	case 1:
		c.Anchor = gm.Vec4{v[0], v[0], v[0], v[0]}
	case 2:
		c.Anchor = gm.Vec4{v[0], v[0], v[1], v[1]}
	case 3:
		c.Anchor = gm.Vec4{v[0], v[1], v[2], v[1]}
	default:
		c.Anchor = gm.Vec4{v[0], v[1], v[2], v[3]}

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
	c.Position = gm.Vec3{x, y, z}
}

// Translatev translates the entity.
func (c *RectComponent) Translatev(v gm.Vec3) {
	c.Position = c.Position.Add(v)
}

// SetRotation sets the rotation.
func (c *RectComponent) SetRotation(q gm.Quat) {
	c.Rotation = q
}

// Rotate axis
func (c *RectComponent) Rotate(x, y, z float32) {
	c.Rotation = c.Rotation.Mul(gm.QFromAngles(
		x, y, z,
		gm.XYZ,
	))
}

// Translate translates the entity.
func (c *RectComponent) Translate(x, y, z float32) {
	c.Position = c.Position.Add(gm.Vec3{x, y, z})
}

// This should be called Dim which are the dimentions, lefttop will always be 0,0
// ContentSize
func (c *RectComponent) ContentSize() gm.Vec2 {
	return c.RelativeSize(c.parentSize())
}

func (c *RectComponent) RelativeSize(parentDim gm.Vec2) gm.Vec2 {
	var right, bottom float32
	// We might discard rect
	right = c.Size[0]
	bottom = c.Size[1]
	// If anchor min and max are the same we use pivot
	// We need to claculate Anchor always unless it's 0
	if c.Anchor[0] != c.Anchor[2] {
		w := parentDim[0] // parentRect[0] is always 0 now
		w -= w*(1-c.Anchor[2]) + w*(c.Anchor[0])
		right = w - c.Size[0] - c.Position[0]
	}

	if c.Anchor[1] != c.Anchor[3] {
		h := parentDim[1]
		h -= h*(1-c.Anchor[3]) + h*(c.Anchor[1])
		bottom = h - c.Size[1] - c.Position[1]
	}

	r := gm.Vec2{
		right - c.Margin[2] - c.Margin[0] - c.Border[2] - c.Border[0],
		bottom - c.Margin[1] - c.Margin[3] - c.Border[1] - c.Border[3],
	}
	return r
}

func (c *RectComponent) parentSize() gm.Vec2 {
	if p, ok := c.parent.(interface{ ContentSize() gm.Vec2 }); ok {
		return p.ContentSize()
	}
	return gm.Vec2{}
}

// TODO: this might be deprecated

// Rect calculate and returns the rect.
func (c *RectComponent) Rect() gm.Vec4 {
	return c.RelativeRect(c.parentRect())
}

// RelativeRect calculate rect based on other.
func (c *RectComponent) RelativeRect(parentRect gm.Vec4) gm.Vec4 {
	var left, top, right, bottom float32

	// parentW := parentRect[2] - parentRect[0]
	// parentH := parentRect[3] - parentRect[1]
	// We might discard rect
	left = 0 // c.Anchor[0] * parentW
	top = 0  // c.Anchor[1] * parentH

	right = c.Size[0]
	bottom = c.Size[1]
	// If anchor min and max are the same we use pivot
	if c.Anchor[0] != c.Anchor[2] {
		w := parentRect[2] - parentRect[0] // parentRect[0] is always 0 now
		// reduce rect by the relative anchor from both sides
		w -= w*(1-c.Anchor[2]) + w*(c.Anchor[0])
		right = w - c.Size[0] - c.Position[0]
	}

	if c.Anchor[1] != c.Anchor[3] {
		h := parentRect[3] - parentRect[1]
		h -= h*(1-c.Anchor[3]) + h*(c.Anchor[1])
		// reduce rect by the relative anchor from both sides
		bottom = h - c.Size[1] - c.Position[1]
	}
	return gm.Vec4{
		left,
		top,
		right - c.Margin[2] - c.Margin[0] - c.Border[2] - c.Border[0],
		// bottom - c.Margin[3] - c.Margin[1],
		bottom - c.Margin[1] - c.Margin[3] - c.Border[1] - c.Border[3],
	}
}

func (c *RectComponent) parentRect() gm.Vec4 {
	if p, ok := c.parent.(interface{ Rect() gm.Vec4 }); ok {
		return p.Rect()
	}
	return gm.Vec4{}
}
