package widget

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// AutoWidth sets a auto width resize handler that resizes based on children.
func AutoWidth(w W, extra float32) event.HandlerFunc {
	return func(e event.Event) {
		_, ok := e.(gorgeui.EventUpdate)
		if !ok {
			return
		}

		dim := m32.Vec2{}
		container := w.Element().GetEntities()
		for _, c := range container {
			rt, ok := c.(rectTransformer)
			if !ok {
				continue
			}
			r := rt.RectTransform()

			left := r.Position[0] // Left most
			right := left + r.Dim[0]
			dim[0] = m32.Max(right+extra, dim[0])

		}
		w.RectTransform().Dim[0] = dim[0]
	}
}

// AutoHeight sets an auto height resizer that resizes based on children.
func AutoHeight(w W, extra float32) event.HandlerFunc {
	return func(e event.Event) {
		_, ok := e.(gorgeui.EventUpdate)
		if !ok {
			return
		}

		dim := m32.Vec2{}
		container := w.Element().GetEntities()
		for _, c := range container {
			rt, ok := c.(rectTransformer)
			if !ok {
				continue
			}
			r := rt.RectTransform()

			top := r.Position[1]
			bottom := top + r.Dim[1]
			dim[1] = m32.Max(bottom+extra, dim[1])

		}
		w.RectTransform().Dim[1] = dim[1]
	}
}

// ResizeToContentHandler will resize a rect to children content.
type ResizeToContentHandler struct {
	w     W
	extra float32
}

// HandleEvent handles gorge events.
func (h ResizeToContentHandler) HandleEvent(e event.Event) {
	_, ok := e.(gorgeui.EventUpdate)
	if !ok {
		return
	}

	dim := m32.Vec2{}
	container := h.w.Element().GetEntities()
	for _, c := range container {
		rt, ok := c.(rectTransformer)
		if !ok {
			continue
		}
		r := rt.RectTransform()
		// find Max and Min
		// left := r.Position[0] // Left most
		right := r.Dim[0]

		// top := r.Position[1]
		bottom := r.Dim[1]

		dim[0] = m32.Max(right+h.extra, dim[0])
		dim[1] = m32.Max(bottom+h.extra, dim[1])

	}
	h.w.RectTransform().Dim[0] = dim[0]
	h.w.RectTransform().Dim[1] = dim[1]
}

// ResizeToContent only works if anchors are relative?
func ResizeToContent(w W, extra float32) *ResizeToContentHandler {
	return &ResizeToContentHandler{w, extra}
}
