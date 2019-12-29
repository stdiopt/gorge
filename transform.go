// Copyright 2019 Luis Figueiredo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gorge

import (
	"github.com/stdiopt/gorge/m32"
)

var (
	//TransformBuilds tracks the numbers of transforms for performance purposes
	TransformBuilds = 0
)

// TransformComponent component
type transformComponent interface {
	TransformComponent() *Transform
}

type affine struct {
	Position vec3
	Rotation quat
	Scale    vec3
}

// Transform Thing
type Transform struct {
	parent transformComponent
	affine

	mat mat4

	shadow affine
}

// TransformComponent component
func (c *Transform) TransformComponent() *Transform { return c }

func (c *Transform) update() {
}

// Mat4 returns the mat4 from the transformations
// Multiply world * local
func (c *Transform) Mat4() mat4 {
	if c.shadow != c.affine {
		TransformBuilds++

		m := m32.Translate3D(c.Position[0], c.Position[1], c.Position[2])
		m = m.Mul4(c.Rotation.Mat4())
		m = m.Mul4(m32.Scale3D(c.Scale[0], c.Scale[1], c.Scale[2]))
		c.mat = m
		c.shadow = c.affine
	}
	if c.parent != nil {
		return c.parent.TransformComponent().Mat4().Mul4(c.mat)
	}
	return c.mat
}

// SetParent of the transform
func (c *Transform) SetParent(p transformComponent) *Transform {
	c.parent = p
	return c
}

// Set full transform
func (c *Transform) Set(position vec3, euler vec3, scale vec3) *Transform {
	c.SetPositionv(position)
	c.SetEulerv(euler)
	c.SetScalev(scale)
	return c
}

// SetPositionv sets the current position on the world with a vector
func (c *Transform) SetPositionv(pos vec3) *Transform {
	c.Position = pos
	return c
}

// SetEulerv sets the euler angles as a vector
func (c *Transform) SetEulerv(angles vec3) *Transform {
	c.Rotation = m32.AnglesToQuat(
		angles[2], angles[1], angles[0],
		m32.ZYX,
	)
	return c
}

// SetPosition sets the current position on the world
func (c *Transform) SetPosition(x, y, z float32) *Transform {
	return c.SetPositionv(vec3{x, y, z})
}

// SetRotation set a quaternion
func (c *Transform) SetRotation(v quat) *Transform {
	c.Rotation = v
	return c
}

// SetEuler convenient func
func (c *Transform) SetEuler(x, y, z float32) *Transform {
	return c.SetEulerv(vec3{x, y, z})
}

// SetScale will set scale
// 1 argument, will set all axis
// 2 arguments, will set only x and y and z to 1
// 3 arguments, will set all
func (c *Transform) SetScale(sz ...float32) *Transform {
	switch len(sz) {
	case 1:
		c.Scale[0], c.Scale[1], c.Scale[2] = sz[0], sz[0], sz[0]
	case 2, 3:
		copy(c.Scale[:], sz)
	default:
		panic("wrong number of params")
	}
	return c
}

// SetScalev just sets the scale
func (c *Transform) SetScalev(scale vec3) *Transform {
	c.Scale = scale
	return c
}

// LookAt resets the local rotation to lookAt
func (c *Transform) LookAt(target, up vec3) *Transform {
	dir := target.Sub(c.Position).Normalize()
	return c.SetRotation(m32.QuatLookAt(dir, up))
}

////////////////////////////
// Relative operations
////

// Translate the thing
func (c *Transform) Translate(x, y, z float32) *Transform {
	c.Position = c.Position.Add(vec3{x, y, z})
	return c
}

// Translatev translate by vector
func (c *Transform) Translatev(axis vec3) *Transform {
	c.Position = c.Position.Add(axis)
	return c
}

// Rotate axis
func (c *Transform) Rotate(x, y, z float32) *Transform {
	c.Rotation = c.Rotation.Mul(m32.AnglesToQuat(
		x, y, z,
		m32.XYZ,
	))
	return c
}

// Rotatev axis by vector
func (c *Transform) Rotatev(angles vec3) *Transform {
	c.Rotation = c.Rotation.Mul(m32.AnglesToQuat(
		angles[0], angles[1], angles[2],
		m32.XYZ,
	))
	return c
}

// WorldPosition returns world position
func (c *Transform) WorldPosition() vec3 {
	return c.Mat4().Col(3).Vec3()
}

// Left returns left of the transform
func (c *Transform) Left() vec3 {
	return c.Mat4().Mul4x1(vec4{-1, 0, 0, 0}).Vec3()
}

// Right returns right of the transform
func (c *Transform) Right() vec3 {
	return c.Mat4().Mul4x1(vec4{1, 0, 0, 0}).Vec3()
}

// Up returns up of the transform
func (c *Transform) Up() vec3 {
	return c.Mat4().Mul4x1(vec4{0, 1, 0, 0}).Vec3()
}

// Down returns down of the transform
func (c *Transform) Down() vec3 {
	return c.Mat4().Mul4x1(vec4{0, -1, 0, 0}).Vec3()
}

// Forward returns forward vector of the transform
func (c *Transform) Forward() vec3 {
	return c.Mat4().Mul4x1(m32.Forward().Vec4(0)).Vec3()
}

// Backward returns backward vector of the transform
func (c *Transform) Backward() vec3 {
	return c.Mat4().Mul4x1(vec4{0, 0, -1, 0}).Vec3()
}

// Inv return the inverse matrix of the transform
func (c *Transform) Inv() mat4 {
	return c.Mat4().Inv()
}

// Position returns the position
//func (c Transform) Position() vec3 {
//return c.position
//}

//// Rotation returns the rotation
//func (c Transform) Rotation() quat {
//return c.rotation
//}

//// Scale returns the scale
//func (c Transform) Scale() vec3 {
//return c.size
//}

// NewTransform returns an initialized transform component
func NewTransform() *Transform {
	return &Transform{
		affine: affine{
			Rotation: m32.QuatIdent(),
			Scale:    vec3{1, 1, 1},
		},
	}
}
