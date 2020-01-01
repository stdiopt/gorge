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

// PreUpdateEvent type
type PreUpdateEvent float32

// UpdateEvent type
type UpdateEvent float32

// PostUpdateEvent type
type PostUpdateEvent float32

// RenderEvent happens after pre,update and post update events
type RenderEvent float32

// ResizeEvent ...
type ResizeEvent m32.Vec2

// EntitiesAddEvent is triggered when entities are added
type EntitiesAddEvent []Entity

// EntitiesRemoveEvent is triggered when entities are destroyed
type EntitiesRemoveEvent []Entity

// StartEvent fired when things starts
type StartEvent struct{}

// AfterStartEvent to attach stuff (wasm request animation frame workaround)
type AfterStartEvent struct{}

// DestroyEvent is called when system is shutting down
type DestroyEvent struct{}

// ErrorEvent contains an error
type ErrorEvent struct{ Err error }

// WarnEvent contains a warning
type WarnEvent string
