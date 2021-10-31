package widget

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// SplitDirectionType split direction specification
type SplitDirectionType int

// Split directions.
const (
	_ = SplitDirectionType(iota)
	SplitHorizontal
	SplitVertical
)

// Split data.
type Split struct {
	Component
	Direction SplitDirectionType
	Spacing   float32
	Size      float32 // First size
}

// HandleEvent implements event interface and handles update events.
func (s *Split) HandleEvent(e event.Event) {
	if _, ok := e.(gorgeui.EventUpdate); !ok {
		return
	}
	container := s.GetEntities()
	if len(container) == 0 {
		return
	}
	count := len(container)
	_ = count
	for _, c := range container {
		r, ok := c.(gorgeui.RectComponent)
		if !ok {
			return
		}
		rt := r.RectTransform()
		rt.SetAnchor(0, 1)
		rt.SetPivot(.5)
		rt.SetRect(0)
	}
}

// NewSplit returns a new split widget.
func NewSplit() *Split {
	return &Split{
		Component: *NewComponent(),
		Direction: SplitHorizontal,
		Spacing:   0,
	}
}
