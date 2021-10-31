package widget

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/m32/ray"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// Slider state for slider widget
type Slider struct {
	Component
	Value        float32
	Color        m32.Vec4
	Padding      m32.Vec4
	HandlerColor m32.Vec4

	// Track   W
	// Handler W

	// Back *Panel

	entity   *Panel
	track    *Widget
	handler  W
	dragging bool
}

// HandleEvent handles gorgeui events.
func (s *Slider) HandleEvent(e event.Event) {
	switch e := e.(type) {
	case gorgeui.EventUpdate:
		if s.handler != nil {
			s.handler.Widget().SetPivot(.5)
			s.handler.Widget().SetAnchor(s.Value, 0, s.Value, 1)
		}
		hs := s.HandlerSize()
		s.track.SetRect(hs/2, 0, hs/2, 0)
		// Wonrg
	case gorgeui.EventPointerUp:
		if s.dragging {
			return
		}
		res := e.RayResult
		r := s.track.Rect()
		fullw := r[2] - r[0]

		wp := s.WorldPosition()
		val := (res.Position[0] - (wp[0] + r[0])) / fullw // Ray in thing position
		val -= s.HandlerSize() / fullw / 2
		val = m32.Clamp(val, 0, 1)
		if val != s.Value {
			s.Value = val
			s.Trigger(EventValueChanged{val})
		}
	case gorgeui.EventDrag:
		s.dragging = true
		rect := s.track.Rect()
		fullw := rect[2] - rect[0]

		m := s.track.Mat4()
		v0 := m.MulV4(m32.Vec4{rect[0], rect[1], 0, 1}).Vec3()
		v1 := m.MulV4(m32.Vec4{rect[2], rect[1], 0, 1}).Vec3() // right
		v2 := m.MulV4(m32.Vec4{rect[0], rect[3], 0, 1}).Vec3() // up)

		ui := gorgeui.RootUI(s)
		r := ray.FromScreen(ui.ScreenSize(), ui.Camera, e.Position)
		res := ray.IntersectRect(r, v0, v1, v2)

		wp := s.WorldPosition()
		val := (res.Position[0] - (wp[0] + rect[0])) / fullw // Ray in thing position
		val -= s.HandlerSize() / fullw / 2

		val = m32.Clamp(val, 0, 1)

		if val != s.Value {
			s.Value = val
			s.Trigger(EventValueChanged{val})
		}
	case gorgeui.EventDragEnd:
		s.dragging = false
	}
}

// NewSlider returns a new slider widget.
func NewSlider() *Slider {
	// entity := QuadEntity()
	// var sliderw *widget.W

	sliderTrack := New()
	sliderTrack.SetAnchor(0, 1)
	sliderTrack.SetPivot(.5)
	sliderTrack.SetRect(0)
	entity := NewPanel()
	entity.SetAnchor(0, 1)

	gorgeui.AddChildrenTo(entity, sliderTrack)

	s := &Slider{
		Component:    *NewComponent(),
		Value:        .5,
		Color:        m32.Vec4{.5, .5, .5, 1},
		HandlerColor: m32.Vec4{1, .5, .5, 1},

		track:  sliderTrack,
		entity: entity,
	}
	s.DragEvents = true

	gorgeui.AddChildrenTo(s, entity)

	// MainRect // it's a panel
	//   +- Track // It's the track where the thing runs
	//      +- Handler // Handler belongs to track
	// Having the panel on track
	// would be harder to configure unless we have things like
	//  widget.SliderTrackRect? widget.SliderTrackAnchor? widget.SliderTrackPivot?

	return s
}

// HandlerSize returns the relative handlersize.
func (s *Slider) HandlerSize() float32 {
	if s.handler == nil {
		return 1
	}

	r := s.handler.Widget().Rect()
	w := r[2] - r[0]

	if w != w {
		return 1
	}
	return w
}

// SetHandler sets the handler widget.
func (s *Slider) SetHandler(h W) {
	s.handler = h
	gorgeui.AddChildrenTo(s.track, h)
}

// SetValue updates the slider value.
func (s *Slider) SetValue(v float32) {
	s.Value = v
}

// SetColor sets the main slider color.
func (s *Slider) SetColor(v ...float32) {
	s.entity.Color = v4Color(v...)
}
