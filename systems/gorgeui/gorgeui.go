// Package gorgeui concept gorge UI
package gorgeui

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/text"
)

// DebugFlag flag type to debug UI stuff.
type DebugFlag uint32

// DefaultFont global default font.
var DefaultFont *text.Font

// Debug flags
const (
	DebugRects = 1 << (1 + iota)
	DebugRays
)

// cameraEntity camera composition used in UI system.
type cameraEntity interface {
	gorge.Matrixer
	Transform() *gorge.TransformComponent
	Camera() *gorge.CameraComponent
}

// Entity every entity should implement this interface.
type Entity interface {
	gorge.Matrixer
	RectTransform() *RectComponent
	Element() *ElementComponent
}

// RootUI utility function to retrieve the UI from a widget.
func RootUI(e gorge.Entity) *UI {
	uiEnt, ok := e.(Entity)
	if !ok {
		return nil
	}
	cur, ok := uiEnt.(gorge.ParentGetter)
	if !ok {
		return nil
	}
	for cur != nil {
		if u, ok := cur.(*UI); ok {
			return u
		}
		cur, ok = cur.Parent().(gorge.ParentGetter)
		if !ok {
			return nil
		}
	}
	return nil
}

// EachParent iterates parents
// XXX: Move to gorge?
func EachParent(e Entity, fn func(e Entity) bool) {
	cur := gorge.Entity(e)
	for cur != nil {
		if v, ok := cur.(Entity); ok {
			if !fn(v) {
				return
			}
		}
		p, ok := cur.(gorge.ParentGetter)
		if !ok {
			break
		}
		cur = p.Parent()
	}
}

// HasParent verifies if the parent exists in e hierarchy
func HasParent(e Entity, parent Entity) bool {
	hasParent := false
	EachParent(e, func(e Entity) bool {
		if e == parent {
			hasParent = true
			return false
		}
		return true
	})
	return hasParent
}

func triggerOn(e Entity, v interface{}) bool {
	evt := Event{Entity: e, Value: v}
	if h, ok := e.(Handler); ok {
		// Direct on thing
		h.HandleEvent(evt)
	}

	if h, ok := e.(trigger); ok {
		h.trigger(evt)
	}
	return true
}

// TriggerOn triggers an event on entity and its parents.
func TriggerOn(e Entity, v interface{}) {
	triggerOn(e, v)
}
