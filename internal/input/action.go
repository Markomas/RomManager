package input

type Action int

const (
	ActionNone Action = iota
	ActionQuit
	ActionUp
	ActionDown
	ActionLeft
	ActionRight
	ActionSelect
	ActionBack
	ActionJumpDown
	ActionJumpUp
)

func (a Action) String() string {
	switch a {
	case ActionNone:
		return "None"
	case ActionQuit:
		return "Quit"
	case ActionUp:
		return "Up"
	case ActionDown:
		return "Down"
	case ActionLeft:
		return "Left"
	case ActionRight:
		return "Right"
	case ActionSelect:
		return "Select"
	case ActionBack:
		return "Back"
	case ActionJumpDown:
		return "Jump Down"
	case ActionJumpUp:
		return "Jump Up"
	default:
		return "Unknown"
	}
}
