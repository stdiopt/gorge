package gorge

// Core events are events handled by gorgeapp that will be handled by the platform.
type (
	// EventCursorRelative turns true/false relative cursor mode uppin triggering event.
	EventCursorRelative bool
	// EventCursorHidden turns true/false hidden cursor uppon triggering event.
	EventCursorHidden bool

	EventCursor CursorType
)

type CursorType int

const (
	CursorArrow CursorType = iota
	CursorHand
)
