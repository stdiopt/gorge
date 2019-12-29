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

package m32

import "github.com/go-gl/mathgl/mgl32"

//FROM the package above

// RotationOrder is the order in which rotations will be transformed for the
// purposes of AnglesToQuat.
type RotationOrder = mgl32.RotationOrder

// The RotationOrder constants represent a series of rotations along the given
// axes for the use of AnglesToQuat.
const (
	XYX RotationOrder = iota
	XYZ
	XZX
	XZY
	YXY
	YXZ
	YZY
	YZX
	ZYZ
	ZYX
	ZXZ
	ZXY
)

// aliases, the goal is to restrict math 32 to this package
var (
	Translate3D   = mgl32.Translate3D
	Scale3D       = mgl32.Scale3D
	AnglesToQuat  = mgl32.AnglesToQuat
	QuatIdent     = mgl32.QuatIdent
	HomogRotate2D = mgl32.HomogRotate2D
)
