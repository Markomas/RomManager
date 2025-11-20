package scene

import "RomManager/internal/input"

type Scene interface {
	Draw()
	HandleInput(action input.Action)
}
