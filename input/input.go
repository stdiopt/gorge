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

// Package input normalizes inputs from systems
package input

import "github.com/stdiopt/gorge/m32"

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

// Consts for key handlers
const (
	KeyDown = iota
	KeyUp
)

// PointerData common
type PointerData struct {
	DeltaZ float32 // for Wheel
	Pos    m32.Vec2
}

// PointerEvent on canvas
type PointerEvent struct {
	Type     PointerType
	Pointers map[int]PointerData
}

// KeyEvent thing
type KeyEvent struct {
	Type int
	Key  string // temp, it should be code
}
