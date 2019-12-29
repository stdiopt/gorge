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

// Camera thing
type Camera struct {
	Fov         float32
	AspectRatio float32
	Near        float32
	Far         float32

	// Ambient color
	Ambient vec3
	// Or sky box
}

// CameraComponent returns camera component
func (c *Camera) CameraComponent() *Camera { return c }

// Projection returns the projection matrix
func (c Camera) Projection() mat4 {
	return m32.Perspective(c.Fov, c.AspectRatio, c.Near, c.Far)
}

// SetPerspective resets projection matrix to perspective
func (c *Camera) SetPerspective(fov, aspectRatio, near, far float32) *Camera {
	c.Fov = fov
	c.AspectRatio = aspectRatio
	c.Near = near
	c.Far = far
	return c
}

// SetAmbient sets camera ambient/clear color
func (c *Camera) SetAmbient(r, g, b float32) *Camera {
	c.Ambient = vec3{r, g, b}
	return c
}
