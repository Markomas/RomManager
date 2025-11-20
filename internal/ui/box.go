package ui

import (
	"RomManager/internal/input"
	"RomManager/internal/ui_render"

	"github.com/veandco/go-sdl2/sdl"
)

type Box struct {
	renderer *ui_render.Renderer
	x        int32
	y        int32
	width    int32
	height   int32
	color    sdl.Color
}

func (b *Box) SetSize(width, height int32) {
	b.width = width
	b.height = height
}

func (b *Box) SetPosition(x, y int32) {
	b.x = x
	b.y = y
}

func (b *Box) Draw() {
	b.renderer.DrawBox(b.x, b.y, b.width, b.height, b.color)
}

func (b *Box) HandleInput(action input.Action) {
}

func NewBox(color sdl.Color, renderer *ui_render.Renderer) UiElement {
	x, y, width, height := int32(0), int32(0), int32(0), int32(0)
	return &Box{renderer: renderer, x: x, y: y, width: width, height: height, color: color}
}
