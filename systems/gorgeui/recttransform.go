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
	parent   gorge.Transformer
	Rotation m32.Quat
	Position m32.Vec3
	Scale    m32.Vec3

	Dim    m32.Vec2
	Anchor m32.Vec4 // left, bottom, right, top

	Pivot m32.Vec2

	// Anchor relative rect
	t1 gorge.TransformComponent // parent offset transform
	t2 gorge.TransformComponent // local transform
}

// RectIdent returns a identity rect transform.
func RectIdent() RectComponent {
	// to be able to work with anchors.
	rc := RectComponent{
		t1:       gorge.TransformIdent(),
		t2:       gorge.TransformIdent(),
		Rotation: m32.QIdent(),
		Scale:    m32.Vec3{1, 1, 1},

		Anchor: m32.Vec4{0, 0, 1, 1},
		Dim:    m32.Vec2{0, 0},
		Pivot:  m32.Vec2{.5, .5},
	}
	rc.Transform() // rebuild transform
	// rc.TransformComponent.SetParent(rc.container)
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
func (c *RectComponent) SetParent(t gorge.Transformer) {
	c.parent = t
}

// Parent returns parent sub transform.
func (c *RectComponent) Parent() gorge.Transformer {
	return c.parent
}

// Transform calculates transform and return it.
func (c *RectComponent) Transform() *gorge.TransformComponent {
	// This is heavy.
	c.t1 = gorge.TransformIdent()
	c.t1.SetParent(c.parent)
	c.t1.Rotation = c.Rotation
	c.t1.Scale = c.Scale
	// This should be parent based on parent Rect position
	rect := c.parentRect() // Parent Dim

	anchor := m32.Vec2{
		(rect[2] - rect[0]) * c.Anchor[0],
		(rect[3] - rect[1]) * c.Anchor[1],
	}
	c.t1.Position = c.Position.Add(anchor.Vec3(0))

	c.t2 = gorge.TransformIdent()
	c.t2.SetParent(&c.t1)
	// T2 position  should only be valid if we don't have the anchor? else should be always .5?
	pivot := m32.Vec2{
		-c.Dim[0] * c.Pivot[0],
		-c.Dim[1] * c.Pivot[1],
	}
	if c.Anchor[0] != c.Anchor[2] {
		anchor[0] -= pivot[0]
		// c.t2.Position[0] = pivot[0]
	} else {
		c.t2.Position[0] = pivot[0]
	}
	// c.t2.Position[0] = pivot[0]
	if c.Anchor[1] != c.Anchor[3] {
		anchor[1] -= pivot[1]
		// c.t2.Position[1] = pivot[1]
	} else {
		c.t2.Position[1] = pivot[1]
	}

	// This DIM might differ from the parent rect.
	// the t2 position should be based on pivot
	// but if anchor is that thing it should be based on other stuff?

	// Returns the child most transform?
	return &c.t2
}

// Mat4 returns 4x4 transform matrix.
func (c *RectComponent) Mat4() m32.Mat4 {
	return c.Transform().Mat4()
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

// Rect calculate and returns the rect.
func (c *RectComponent) Rect() m32.Vec4 {
	return c.RelativeRect(c.parentRect())
}

// RelativeRect calculate rect based on other.
func (c *RectComponent) RelativeRect(parentRect m32.Vec4) m32.Vec4 {
	var left, top, right, bottom float32
	// We might discard rect
	left = 0 // c.Position[0] // parentRect[0] //+ c.Anchor[0]*parentW
	top = 0  // c.Position[1]  // parentRect[1]  //+ c.Anchor[1]*parentH

	right = c.Dim[0]
	bottom = c.Dim[1]
	// If anchor min and max are the same we use pivot
	if c.Anchor[0] != c.Anchor[2] {
		w := parentRect[2] // parentRect[0] is always 0 now
		// reduce rect by the relative anchor from both sides
		w -= w*(1-c.Anchor[2]) + w*(c.Anchor[0])
		right = w - c.Dim[0] - c.Position[0]
	}

	if c.Anchor[1] != c.Anchor[3] {
		w := parentRect[3]
		w -= w*(1-c.Anchor[3]) + w*(c.Anchor[1])
		// reduce rect by the relative anchor from both sides
		bottom = w - c.Dim[1] - c.Position[1]
	}
	return m32.Vec4{left, top, right, bottom}
}

func (c *RectComponent) parentRect() m32.Vec4 {
	parentRect := m32.Vec4{}
	if p := c.parent; p != nil {
		if c, ok := p.(interface{ Rect() m32.Vec4 }); ok {
			parentRect = c.Rect()
		}
	}
	return parentRect
}
