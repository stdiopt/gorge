// Package widget contains UI widgets for gorgeui
package widget

import (
	"log"
	"runtime"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// W the widget base interface.
type W interface {
	gorgeui.Entity
	Widget() *Component
}

type (
	renderabler interface {
		Renderable() *gorge.RenderableComponent
	}
	rectTransformer interface {
		RectTransform() *gorgeui.RectComponent
	}
	transformer interface {
		Transform() *gorge.TransformComponent
	}
	widgeter interface {
		Widget() *Component
	}
	graphic interface {
		Transform() *gorge.TransformComponent
		Renderable() *gorge.RenderableComponent
		Colorable() *gorge.ColorableComponent
	}
)

// Widget represents a empty widget with childs
type Widget struct {
	Component
}

// New returns a new empty widget.
func New() *Widget {
	w := &Widget{
		Component: *NewComponent(),
	}
	return w
}

// WidgetCounter debug release stuff
var WidgetCounter int

type gcref struct{ n int }

// Component is the base of all views
type Component struct {
	gcref *gcref
	gorgeui.RectComponent
	gorgeui.ElementComponent

	Name string

	// Data contains optional widget states
	Data map[interface{}]interface{}

	// Initialized, State
	built bool
}

// NewComponent returns a initialized widget component.
func NewComponent() *Component {
	g := &gcref{1}

	// debug object release tracker
	runtime.SetFinalizer(g, func(v interface{}) {
		log.Println("Finalizing UI OBJ:", v)
		WidgetCounter--
	})

	WidgetCounter++
	b := &Component{
		gcref:         g,
		RectComponent: gorgeui.RectIdent(),
		built:         true,
	}
	return b
}

// Widget implements the WidgetComponent for embeding purposes.
func (w *Component) Widget() *Component { return w }

// Attached implements the Attacher interface
func (w *Component) Attached(e gorgeui.Entity) {
	// Sounds bad
	for _, c := range w.GetEntities() {
		if r, ok := c.(renderabler); ok {
			if ui := gorgeui.RootUI(w); ui != nil {
				r.Renderable().CullMask = ui.CullMask
			}
		}
	}
}

// SetName sets widget name for dbg purposes.
func (w *Component) SetName(n string) {
	w.Name = n
}

// GetData returns a widget dynamic property based on k key.
func (w *Component) GetData(k interface{}) interface{} {
	if w.Data == nil {
		return nil
	}
	return w.Data[k]
}

// SetData sets a widget dynamic property based on k key.
func (w *Component) SetData(k, d interface{}) {
	if w.Data == nil {
		w.Data = map[interface{}]interface{}{}
	}
	w.Data[k] = d
}
