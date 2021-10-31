package widget

import "github.com/stdiopt/gorge/systems/gorgeui"

// EventClick triggers on widget click.
type EventClick struct {
	Widget gorgeui.Entity
}

// EventValueChanged triggers on certain widdgets a value change.
type EventValueChanged struct {
	Value float32
}
