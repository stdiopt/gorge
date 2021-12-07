package gorgeui

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
)

// TODO:

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
	event.Bus

	static   []gorge.Entity
	children []gorge.Entity
	all      []gorge.Entity

	DragEvents     bool
	DisableRaycast bool

	Attached   bool
	LayoutFunc func(e Entity) // LayoutFunc is called when the element state is changed rect etc
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

// SetLayoutFunc sets the layout func.
func (c *ElementComponent) SetLayoutFunc(fn LayoutFunc) {
	c.LayoutFunc = fn
}

// GetEntities implements the gorge.EntityContainer interface.
func (c *ElementComponent) GetEntities() []gorge.Entity {
	return c.all
}

// Children returns the dynamic children of this element.
func (c *ElementComponent) Children() []gorge.Entity {
	return c.children
}

func (c *ElementComponent) add(ents ...gorge.Entity) {
	c.all = append(c.all, ents...)
	c.static = append(c.static, ents...)
}

// AddChildren adds dynamic children to this element.
func (c *ElementComponent) AddChildren(ents ...gorge.Entity) {
	c.all = append(c.all, ents...)
	c.children = append(c.children, ents...)
}

func (c *ElementComponent) removeChildren(ents ...gorge.Entity) {
	for _, e := range ents {
		for i := 0; i < len(c.children); i++ {
			if c.children[i] == e {
				t := c.children
				c.children = append(c.children[:i], c.children[i+1:]...)
				i--
				t[len(t)-1] = nil // remove reference to last so it can be Gc'ed
			}
		}
	}
	for _, e := range ents {
		for i := 0; i < len(c.all); i++ {
			if c.all[i] == e {
				t := c.all
				c.all = append(c.all[:i], c.all[i+1:]...)
				t[len(t)-1] = nil // remove reference to last so it can be Gc'ed
				break
			}
		}
	}
}

// AddChildrenTo adds a children, children will be added to gorge if parent is
// attached.
func AddChildrenTo(parent Entity, ents ...gorge.Entity) {
	for _, cc := range ents {
		if t, ok := cc.(gorge.ParentSetter); ok {
			t.SetParent(parent)
		}
		parent.Element().AddChildren(cc)
	}
	if parent.Element().Attached {
		ui := RootUI(parent)
		ui.Add(ents...)
	}
}

// AddElementTo adds static element.
func AddElementTo(parent Entity, ents ...gorge.Entity) {
	for _, cc := range ents {
		if t, ok := cc.(gorge.ParentSetter); ok {
			t.SetParent(parent)
		}
		parent.Element().add(cc)
	}
	if parent.Element().Attached {
		ui := RootUI(parent)
		ui.Add(ents...)
	}
}
