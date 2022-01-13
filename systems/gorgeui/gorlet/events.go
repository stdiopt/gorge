package gorlet

type (
	// EventValueChanged triggered when an input value is changed.
	EventValueChanged struct{ Value any }
	// EventClick is trigger when an event is clicked.
	EventClick struct{} // need more info
)
