package input

// Input system tied to manager
type Input struct {
	keyManager
	mouseManager
}

// ActionState type
type ActionState int

func (k ActionState) String() string {
	switch k {
	case ActionDown:
		return "ActionDown"
	case ActionHold:
		return "ActionHold"
	case ActionUp:
		return "ActionUp"
	default:
		return "<no state>"
	}
}

// Action states used in keys and buttons.
const (
	_ = ActionState(iota)
	ActionDown
	ActionHold
	ActionUp
)
