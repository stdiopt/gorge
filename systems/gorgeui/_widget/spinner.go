package widget

import (
	"fmt"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// EventSpin ?!
type EventSpin float32

// Spinner thingy
type Spinner struct {
	Component
	Label string
	Value float32
	Color m32.Vec4

	panel      *Panel
	labelBg    *Panel
	label      *Label
	valueLabel *Label
}

// HandleEvent handles events.
func (w *Spinner) HandleEvent(ee event.Event) {
	switch e := ee.(type) {
	case gorgeui.EventUpdate:
		if w.valueLabel != nil {
			w.valueLabel.SetText(fmt.Sprintf("%.2f", w.Value))
		}
		if w.labelBg != nil {
			w.labelBg.SetColor(w.Color[0], w.Color[1], w.Color[2], w.Color[3])
		}
	case gorgeui.EventDrag:
		w.Value += e.Delta[0] * 0.01
		w.Trigger(EventSpin(w.Value))
	}
}

// NewSpinner creates a dam spinner.
func NewSpinner(label string, val float32) *Spinner {
	const split = .3

	p := NewPanel()
	p.SetPivot(0)
	p.SetRect(0)
	p.SetAnchor(0, 0, 1, 1)
	// p.SetParent(s)
	lblBg := NewPanel()
	lblBg.SetColor(.2)
	lblBg.SetAnchor(0, 0, split, 1)
	gorgeui.AddChildrenTo(p, lblBg)
	// lp.SetParent(p)

	lbl := NewLabel(label) // needs to update
	lbl.SetColor(1)
	lbl.SetRect(0)
	lbl.SetAnchor(0, 0, 1, 1)
	gorgeui.AddChildrenTo(lblBg, lbl)

	valLbl := NewLabel()
	valLbl.SetColor(1)
	valLbl.SetText("0.0")
	valLbl.SetRect(0)
	valLbl.SetAnchor(split, 0, 1, 1)
	// l.SetParent(p)
	gorgeui.AddChildrenTo(p, valLbl)

	s := &Spinner{
		Component: *NewComponent(),
		Label:     label,
		Value:     val,

		panel:      p,
		labelBg:    lblBg,
		label:      lbl,
		valueLabel: valLbl,
	}
	gorgeui.AddChildrenTo(s, p)
	s.SetRect(1, 1, 10, 5)
	s.SetDragEvents(true)
	return s
}

// NewSpinnerX different way of building
func NewSpinnerX(label string, val float32) *Spinner {
	const split = .4
	s := &Spinner{
		Component: *NewComponent(),
		Label:     label,
		Value:     val,
	}
	Build(s)
	return s
}

// SetColor set spinner color.
func (w *Spinner) SetColor(c m32.Vec4) {
	w.Color = c
}

// SetLabelColor sets the background label color.
func (w *Spinner) SetLabelColor(f ...float32) {
	w.labelBg.SetColor(f...)
	// Updated Layout
}

// Build func
// Using this could be good for building a widget from a struct.
// although widget doesn't have a system.
// we could use the builder to callback this.
func (w *Spinner) Build(b *Builder) {
	const split = .4
	w.SetRect(1, 1, 10, 5)
	w.SetDragEvents(true)
	p := b.BeginPanel()
	p.SetRect(0)
	p.SetAnchor(0, 0, 1, 1)
	{
		p := b.BeginPanel()
		p.SetColor(1, 0, 0)
		p.SetRect(0)
		p.SetAnchor(0, 0, split, 1)
		b.Style().SetColor(1)
		l := b.Label(&w.Label)
		l.SetRect(0, 0, 0, 0)
		l.SetAnchor(0, 0, 1, 1)
		b.EndPanel()
	}
	l := b.Label(func() string { return fmt.Sprintf("%.2f", w.Value) })
	l.SetRect(0)
	l.SetAnchor(split, 0, 1, 1)
	l.SetAlign(AlignCenter)
	b.EndPanel()
}
