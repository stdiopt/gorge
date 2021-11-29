package gorgeui

import (
	"github.com/stdiopt/gorge"
)

type (
	// Attacher interface used to automatically add an entity in the UI hierarchy.
	Attacher interface{ Attached(e Entity) }
	// Detacher interface used to automatically remove an entity in the UI hierarchy.
	Detacher interface{ Detached(e Entity) }

	elementer interface{ Element() *ElementComponent }
)

// ElementComponent is a base widget for UI things
// it contains some state
type ElementComponent struct {
	eventBus
	gorge.Container
	DragEvents     bool
	DisableRaycast bool

	Attached bool
}

// Element implements the Element.
func (c *ElementComponent) Element() *ElementComponent { return c }

// SetDragEvents enable or disable drag events for this element.
func (c *ElementComponent) SetDragEvents(b bool) {
	c.DragEvents = b
}

// SetDisableRaycast disable or enable ray casting on this element.
func (c *ElementComponent) SetDisableRaycast(b bool) {
	c.DisableRaycast = b
}

// AddChildrenTo adds a children, children will be added to gorge if parent is
// attached.
func AddChildrenTo(parent Entity, ents ...gorge.Entity) {
	for _, cc := range ents {
		if t, ok := cc.(gorge.ParentSetter); ok {
			t.SetParent(parent)
		}
		parent.Element().Add(cc)
	}
	if parent.Element().Attached {
		ui := RootUI(parent)
		ui.Add(ents...)
	}
}
