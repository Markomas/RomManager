package ui

import "RomManager/internal/input"

type UiElement interface {
	Draw()
	HandleInput(action input.Action)
	SetSize(width, height int32)
	SetPosition(x, y int32)
}

type LayoutElement struct {
	FullWidth  bool
	Width      int32
	FullHeight bool
	Height     int32
	UiElement  UiElement
}
