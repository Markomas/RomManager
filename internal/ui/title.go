package ui

import (
	"RomManager/internal/config"
	"RomManager/internal/input"
)

type Title struct {
	text     string
	x        int32
	y        int32
	width    int32
	height   int32
	renderer *Renderer
	config   *config.Config
}

func (t *Title) HandleInput(action input.Action) {}

func (t *Title) SetSize(width, height int32) {
	t.width = width
	t.height = height
}

func (t *Title) SetPosition(x, y int32) {
	t.y = y
	t.x = x
}

func (t *Title) Draw() {
	t.renderer.DrawBox(t.x, t.y, t.width, t.height, t.config.Theme.TitleBackgroundColor)
	t.renderer.DrawText(t.text, t.x, t.y, t.width, t.config.Theme.TitleFontSize, t.config.Theme.TitleColor, AlignCenter)
	t.renderer.DrawLine(t.x, t.y+t.height-1, t.width, t.y+t.height-1, t.config.Theme.TitleLineColor)
}

func NewTitle(text string, renderer *Renderer, c *config.Config) UiElement {
	return &Title{text: text, renderer: renderer, config: c, x: 0, y: 0, width: 0, height: 0}
}
