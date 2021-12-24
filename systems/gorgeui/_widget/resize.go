package widget

// AutoWidth sets a auto width resize handler that resizes based on children.
/*func AutoWidth(extra float32) gorgeui.HandlerFunc {
	return func(e gorgeui.Event) {
		_, ok := e.Value.(gorgeui.EventUpdate)
		if !ok {
			return
		}
		w, ok := e.Entity.(W)
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
func AutoHeight(extra float32) gorgeui.HandlerFunc {
	return func(e gorgeui.Event) {
		_, ok := e.Value.(gorgeui.EventUpdate)
		if !ok {
			return
		}
		w, ok := e.Entity.(W)
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
	extra float32
}

// HandleEvent handles gorge events.
func (h ResizeToContentHandler) HandleEvent(e gorgeui.Event) {
	_, ok := e.Value.(gorgeui.EventUpdate)
	if !ok {
		return
	}
	w, ok := e.Entity.(W)
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
		// find Max and Min
		// left := r.Position[0] // Left most
		right := r.Dim[0]

		// top := r.Position[1]
		bottom := r.Dim[1]

		dim[0] = m32.Max(right+h.extra, dim[0])
		dim[1] = m32.Max(bottom+h.extra, dim[1])

	}
	w.RectTransform().Dim[0] = dim[0]
	w.RectTransform().Dim[1] = dim[1]
}

// ResizeToContent only works if anchors are relative?
func ResizeToContent(extra float32) *ResizeToContentHandler {
	return &ResizeToContentHandler{extra}
}
*/
