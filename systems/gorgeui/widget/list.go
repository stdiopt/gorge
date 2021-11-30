package widget

import (
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type (
	// ListDirection list direction description.
	ListDirection int
	// ListSize list size description.
	ListSize int
)

// List directions
const (
	_ = ListDirection(iota)
	ListHorizontal
	ListVertical
)

// List size specifications.
const (
	_ = ListSize(iota)
	ListSizeFixed
	ListSizeStretch
	ListSizeOriginal
)

// List state for list widget.
type List struct {
	widget    W
	Direction ListDirection
	SizeMode  ListSize

	Padding  m32.Vec2
	Spacing  float32
	ElemSize float32
}

// SetElemSize sets the element size.
func (l *List) SetElemSize(v float32) {
	l.ElemSize = v
}

// SetSpacing the space between child elements.
func (l *List) SetSpacing(v float32) {
	l.Spacing = v
}

// SetSizeMode sets the way to resize the elements.
func (l *List) SetSizeMode(v ListSize) {
	l.SizeMode = v
}

// SetDirection sets list direction vertical or horizontal.
func (l *List) SetDirection(v ListDirection) {
	l.Direction = v
}

// SetPadding sets list padding
func (l *List) SetPadding(v ...float32) {
	switch len(v) {
	case 1:
		l.Padding = m32.Vec2{v[0], v[0]}
	case 2:
		l.Padding = m32.Vec2{v[0], v[1]}
	default:
		l.Padding = m32.Vec2{}

	}
}

// HandleEvent implements event.Handler
func (l *List) HandleEvent(e gorgeui.Event) {
	_, ok := e.Value.(gorgeui.EventUpdate)
	if !ok {
		return
	}
	container := l.widget.Widget().GetEntities()
	cur := l.Padding[0]
	for i, c := range container {
		r, ok := c.(rectTransformer)
		if !ok {
			continue
		}
		rt := r.RectTransform()
		if l.Direction == ListVertical {
			rt.SetAnchor(0, 0, 1, 0)
			rt.SetPivot(0, 0)
			var size float32
			switch l.SizeMode {
			case ListSizeFixed:
				size = l.ElemSize
			case ListSizeOriginal:
				size = rt.Dim[1]
			case ListSizeStretch:
				pr := l.widget.RectTransform().Rect()
				fc := float32(len(container))
				size = (pr[3] - pr[1]) / fc
				size += l.Spacing * (fc - 1) / fc
			}
			rt.SetRect(l.Padding[0], cur, l.Padding[1], size)
			cur += size + l.Spacing
			continue
		}
		rt.SetAnchor(0, 0, 0, 1)
		rt.SetPivot(0, .5)

		var size float32
		switch l.SizeMode {
		case ListSizeFixed:
			size = l.ElemSize
		case ListSizeOriginal:
			size = rt.Dim[0]
		case ListSizeStretch:
			pr := l.widget.RectTransform().Rect()
			size = (pr[2] - pr[0]) / float32(len(container))
			if i < len(container)-1 {
				size -= l.Spacing
			}
		}
		rt.SetRect(cur, l.Padding[0], size, l.Padding[1])
		cur += size + l.Spacing
	}
}

// ListController Attaches a layout handler to the widget which will resize
// children accordingly.
func ListController(w W) *List {
	l := &List{
		widget:    w,
		Direction: ListVertical,
		SizeMode:  ListSizeFixed,
		Spacing:   2,
	}
	w.Widget().Handle(l)
	// w.Widget().SetData(listKey, l)
	// w.Widget().HandleFunc(func(e event.Event) {
	// })
	return l
}

/*func NewList(opts ...gorgeui.EntityFunc) *List {
	l := &List{
		WidgetComponent: *NewComponent(),
		Direction:       ListVertical,
		SizeMode:        ListSizeFixed,
		Spacing:         2,
	}

	gorgeui.ApplyTo(l,
		gorgeui.WithHandleFunc(func(e event.Event) {
			_, ok := e.(gorgeui.EventUpdate)
			if !ok {
				return
			}
			container := l.Container.GetEntities()
			var cur float32
			for i, c := range container {
				r, ok := c.(rectTransformer)
				if !ok {
					continue
				}
				rt := r.RectTransform()
				if l.Direction == ListVertical {
					rt.SetAnchor(0, 1, 1, 1)
					rt.SetPivot(.5, 1)
					var size float32
					switch l.SizeMode {
					case ListSizeFixed:
						size = l.ElemSize
					case ListSizeOriginal:
						size = rt.Dim[1]
					case ListSizeStretch:
						pr := l.Rect()
						fc := float32(len(container))
						size = (pr[3] - pr[1]) / fc
						size -= l.Spacing * (fc - 1) / fc
					}
					rt.SetRect(0, cur, 0, size)
					cur -= size + l.Spacing
				} else {
					rt.SetAnchor(0, 0, 0, 1)
					rt.SetPivot(0, .5)

					var size float32
					switch l.SizeMode {
					case ListSizeFixed:
						size = l.ElemSize
					case ListSizeOriginal:
						size = rt.Dim[0]
					case ListSizeStretch:
						pr := l.Rect()
						size = (pr[2] - pr[0]) / float32(len(container))
						if i < len(container)-1 {
							size -= l.Spacing
						}
					}
					rt.SetRect(cur, 0, size, 0)
					cur += size + l.Spacing
				}
			}
		}),
		gorgeui.FuncGroup(opts...),
	)

	return l
}*/
