package ui

import (
	"RomManager/internal/config"
	"RomManager/internal/input"
)

type Text struct {
	text       string
	x          int32
	y          int32
	width      int32
	height     int32
	renderer   *Renderer
	config     *config.Config
	color      config.Color
	fontSize   int
	horizontal Align
	vertical   Align
}

func NewText(text string, renderer *Renderer, c *config.Config, color config.Color, fontSize int, horizontal Align, vertical Align) UiElement {
	return &Text{
		text:       text,
		renderer:   renderer,
		config:     c,
		color:      color,
		fontSize:   fontSize,
		horizontal: horizontal,
		vertical:   vertical,
		x:          0,
		y:          0,
		width:      0,
		height:     0,
	}
}

func (t *Text) SetSize(width, height int32) {
	t.width = width
	t.height = height
}

func (t *Text) SetPosition(x, y int32) {
	t.x = x
	t.y = y
}

func (t *Text) Draw() {
	t.renderer.DrawBox(t.x, t.y, t.width, t.height, t.config.Theme.TextBackgroundColor)
	padding := int32(t.config.Theme.TextPadding)
	t.renderer.DrawTextBox(t.text, t.x+padding, t.y+padding, t.width-padding, t.height-padding, t.fontSize, t.color, t.vertical, t.horizontal)
}

func (t *Text) HandleInput(action input.Action) {
	// Text element is not interactive
}

func (t *Text) SetText(text string) {
	t.text = text
}

func (t *Text) GetText() string {
	return t.text
}
