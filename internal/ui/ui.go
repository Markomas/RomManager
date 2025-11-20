package ui

import "RomManager/internal/input"

type UiElement interface {
	Draw()
	HandleInput(action input.Action)
	SetSize(width, height int32)
	SetPosition(x, y int32)
}
