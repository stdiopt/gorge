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

package gorgeutils

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/input"
	"github.com/stdiopt/gorge/m32"
)

type (
	vec2 = m32.Vec2
	vec3 = m32.Vec3
)

// CameraRig thing
type CameraRig struct {
	*gorge.Transform
	Entity *Camera
}

type msg interface {
	Handle(interface{}) *gorge.Handler
	AddEntity(...gorge.Entity)
}

// TrackballCamera attaches events and all to make a trackball
func TrackballCamera(g msg) *CameraRig {
	cameraRig := &CameraRig{
		Transform: gorge.NewTransform(),
		Entity:    NewCamera(),
	}

	g.Handle(func(evt gorge.StartEvent) {
		g.AddEntity(cameraRig.Entity)
	})
	g.Handle(func(evt gorge.ResizeEvent) {
		sz := vec2(evt)
		cameraRig.Entity.Camera.AspectRatio = sz[0] / sz[1]
	})
	cameraRig.Rotate(0.7, 0, 0)
	cameraRig.Entity.Camera.
		SetAmbient(0.4, 0.4, 0.4)

	cameraRig.Entity.Transform.
		SetParent(cameraRig.Transform).
		SetEuler(0, 0, 0).
		SetPosition(0, 0, -10)

	var lastP vec2
	var camRot vec2 = vec2{-0.7, 0}
	var dragging = false
	g.Handle(func(evt input.PointerEvent) {
		delta := vec2(evt.Pointers[0].Pos).Sub(lastP)
		lastP = vec2(evt.Pointers[0].Pos)
		if evt.Type == input.MouseWheel {
			dist := cameraRig.Entity.Transform.WorldPosition().Len()
			multiplier := dist * 0.005
			cameraRig.Entity.Transform.Translate(0, 0, -evt.Pointers[0].DeltaZ*multiplier)
		}
		if evt.Type == input.MouseDown {
			dragging = true
		}
		if evt.Type == input.MouseUp {
			dragging = false
		}

		// If dragging or pointer move
		if dragging || evt.Type == input.PointerMove {
			if len(evt.Pointers) == 1 {
				scale := float32(0.005)
				camRot = camRot.Add(
					vec2{-delta[1], delta[0]}.Mul(scale),
				)
				cameraRig.SetRotation(m32.QuatEuler(camRot[0], camRot[1], 0))
			}
			/*
				if len(evt.Pointers) == 2 {
					v := vec2(evt.Pointers[0].Pos).Sub(vec2(evt.Pointers[1].Pos))
					curPinch := v.Len()
					if !pinching {
						lastPinch = curPinch
						pinching = true
					}
					deltaPinch := curPinch - lastPinch
					lastPinch = curPinch
					s.camera.Transform().Translate(0, 0, deltaPinch*0.1)
				}*/
		}

	})
	return cameraRig
}
