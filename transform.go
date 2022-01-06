package gorge

import (
	"github.com/stdiopt/gorge/m32"
)

// TransformBuilds tracks the numbers of transforms for performance purposes
var (
	TransformBuilds    = 0
	TransformBuildSave = 0
)

type (
	// ParentGetter interface for a parent getter.
	ParentGetter interface{ Parent() Matrixer }
	// ParentSetter interface that Sets a parent.
	ParentSetter interface{ SetParent(Matrixer) }
	// Transformer interface for the transform component implementer.
	Transformer interface {
		Mat4() m32.Mat4
		Transform() *TransformComponent
	}
	// Matrixer interface to mat4
	Matrixer interface {
		Mat4() m32.Mat4
	}
)

// TransformComponent component

type affine struct {
	Position m32.Vec3
	Rotation m32.Quat
	Scale    m32.Vec3
}

// TransformComponent Thing
type TransformComponent struct {
	affine
	parent Matrixer

	cached         affine
	cachedWorldMat m32.Mat4

	updates      int
	parentUpdate int
}

// TransformIdent returns a transform identity copy.
func TransformIdent() TransformComponent {
	return TransformComponent{
		affine: affine{
			Rotation: m32.QIdent(),
			Scale:    m32.Vec3{1, 1, 1},
		},
	}
}

// NewTransformComponent returns an initialized transform component
func NewTransformComponent() *TransformComponent {
	return &TransformComponent{
		affine: affine{
			Rotation: m32.QIdent(),
			Scale:    m32.Vec3{1, 1, 1},
		},
	}
}

// Transform component
func (c *TransformComponent) Transform() *TransformComponent { return c }

// Updated returns true if the transform or relative parents were updated.
func (c *TransformComponent) Updated() bool {
	// Should force update if it is 0
	if c.updates == 0 {
		return true
	}
	if c.cached != c.affine {
		return true
	}
	switch v := c.parent.(type) {
	case interface{ Transform() *TransformComponent }:
		t := v.Transform()
		if t.Updated() {
			return true
		}
		if c.parentUpdate != t.Updates() {
			return true
		}
	case Matrixer:
		return true
	}
	return false
}

// Updates return the current number of udpates.
func (c *TransformComponent) Updates() int {
	return c.updates
}

// Parent returns the current parent.
func (c *TransformComponent) Parent() Matrixer {
	return c.parent
}

// Mat4 returns the World m32.Mat4 from the transformations
func (c *TransformComponent) Mat4() m32.Mat4 {
	if c == nil {
		return m32.M4Ident()
	}
	if !c.Updated() {
		TransformBuildSave++
		return c.cachedWorldMat
	}

	TransformBuilds++

	// TODO: {lpf} Could have a local matrix cache as well instead of recalculating

	c.cachedWorldMat = m32.Translate3D(c.Position[0], c.Position[1], c.Position[2])
	c.cachedWorldMat = c.cachedWorldMat.Mul(c.Rotation.Mat4())
	c.cachedWorldMat = c.cachedWorldMat.Mul(m32.Scale3D(c.Scale[0], c.Scale[1], c.Scale[2]))
	if c.parent != nil {
		// t := c.parent.Transform()
		c.cachedWorldMat = c.parent.Mat4().Mul(c.cachedWorldMat)

		if t, ok := c.parent.(interface{ Updates() int }); ok {
			c.parentUpdate = t.Updates()
		}
	}
	c.cached = c.affine
	c.updates++
	return c.cachedWorldMat
}

// SetMat4Decompose decomposes a 4x4 into position, rotation and scale.
// https://answers.unity.com/questions/402280/how-to-decompose-a-trs-matrix.html
func (c *TransformComponent) SetMat4Decompose(m m32.Mat4) {
	c.Position = m.Col(3).Vec3()
	c.Rotation = m.Quat()
	// Scale might not work if it has negative scales they say?
	c.Scale = m32.Vec3{
		m.Col(0).Len(),
		m.Col(1).Len(),
		m.Col(2).Len(),
	}
}

// SetParent of the transform
func (c *TransformComponent) SetParent(p Matrixer) {
	c.parent = p
}

// Set full transform
func (c *TransformComponent) Set(position m32.Vec3, euler m32.Vec3, scale m32.Vec3) {
	c.SetPositionv(position)
	c.SetEulerv(euler)
	c.SetScalev(scale)
}

// SetPositionv sets the current position on the world with a vector
func (c *TransformComponent) SetPositionv(pos m32.Vec3) {
	c.Position = pos
}

// SetEulerv sets the euler angles as a vector
func (c *TransformComponent) SetEulerv(angles m32.Vec3) {
	c.Rotation = m32.QFromAngles(
		angles[2], angles[1], angles[0],
		m32.ZYX,
	)
}

// SetPosition sets the current position on the world
func (c *TransformComponent) SetPosition(x, y, z float32) {
	c.SetPositionv(m32.Vec3{x, y, z})
}

// SetRotation set a quaternion
func (c *TransformComponent) SetRotation(v m32.Quat) {
	c.Rotation = v
}

// SetEuler convenient func
func (c *TransformComponent) SetEuler(x, y, z float32) {
	c.SetEulerv(m32.Vec3{x, y, z})
}

// SetScale will set scale
// 1 argument, will set all axis
// 2 arguments, will set only x and y and z to 1
// 3 arguments, will set all
func (c *TransformComponent) SetScale(sz ...float32) {
	switch len(sz) {
	case 1:
		c.Scale[0], c.Scale[1], c.Scale[2] = sz[0], sz[0], sz[0]
	case 2, 3:
		copy(c.Scale[:], sz)
	default:
		panic("wrong number of params")
	}
}

// SetScalev just sets the scale
func (c *TransformComponent) SetScalev(scale m32.Vec3) {
	c.Scale = scale
}

// LookAt resets the local rotation to lookAt
// if 1 param is used, we will Use default m32.Up() +Y vector
func (c *TransformComponent) LookAt(target Matrixer, v ...m32.Vec3) {
	up := m32.Vec3{0, 1, 0}
	if len(v) > 1 {
		up = v[0]
	}
	pos := target.Mat4().Col(3)

	dir := c.Position.Sub(pos.Vec3()).Normalize()
	c.SetRotation(m32.QLookAt(dir, up))
}

// LookAtPosition resets the local rotation to lookAt
// if 1 param is used, we will Use default m32.Up() +Y vector
func (c *TransformComponent) LookAtPosition(target m32.Vec3, v ...m32.Vec3) {
	up := m32.Vec3{0, 1, 0}
	if len(v) > 1 {
		up = v[0]
	}

	dir := c.Position.Sub(target).Normalize()
	c.SetRotation(m32.QLookAt(dir, up))
}

// LookDir looks at direction
func (c *TransformComponent) LookDir(dir m32.Vec3, v ...m32.Vec3) {
	up := m32.Vec3{0, 1, 0}
	if len(v) > 1 {
		up = v[0]
	}
	c.SetRotation(m32.QLookAt(dir, up))
}

// //////////////////////////
// Relative operations

// Translate the thing
func (c *TransformComponent) Translate(x, y, z float32) {
	c.Position = c.Position.Add(m32.Vec3{x, y, z})
}

// Translatev translate by vector
func (c *TransformComponent) Translatev(axis m32.Vec3) {
	c.Position = c.Position.Add(axis)
}

// Rotate axis
func (c *TransformComponent) Rotate(x, y, z float32) {
	c.Rotation = c.Rotation.Mul(m32.QFromAngles(
		x, y, z,
		m32.XYZ,
	))
}

// Rotatev axis by vector
func (c *TransformComponent) Rotatev(angles m32.Vec3) {
	c.Rotation = c.Rotation.Mul(m32.QFromAngles(
		angles[0], angles[1], angles[2],
		m32.XYZ,
	))
}

// WorldPosition returns world position
func (c *TransformComponent) WorldPosition() m32.Vec3 {
	return c.Mat4().Col(3).Vec3()
}

// WorldRotation returns world rotation
func (c *TransformComponent) WorldRotation() m32.Quat {
	return c.Mat4().Quat()
}

// Left returns World left of the transform
func (c *TransformComponent) Left() m32.Vec3 {
	return c.Mat4().MulV4(m32.Vec4{-1, 0, 0, 0}).Vec3()
}

// Right returns World right of the transform
func (c *TransformComponent) Right() m32.Vec3 {
	return c.Mat4().MulV4(m32.Vec4{1, 0, 0, 0}).Vec3()
}

// Up returns World up of the transform
func (c *TransformComponent) Up() m32.Vec3 {
	return c.Mat4().MulV4(m32.Vec4{0, 1, 0, 0}).Vec3()
}

// Down returns World Fown of the transform
func (c *TransformComponent) Down() m32.Vec3 {
	return c.Mat4().MulV4(m32.Vec4{0, -1, 0, 0}).Vec3()
}

// Forward returns World Forward vector of the transform
func (c *TransformComponent) Forward() m32.Vec3 {
	return c.Mat4().MulV4(m32.Vec4{0, 0, -1, 0}).Vec3()
}

// Backward returns backward vector of the transform
func (c *TransformComponent) Backward() m32.Vec3 {
	return c.Mat4().MulV4(m32.Vec4{0, 0, 1, 0}).Vec3()
}

// Inv return the inverse matrix of the transform
func (c *TransformComponent) Inv() m32.Mat4 {
	return c.Mat4().Inv()
}
