package ui

import (
	"RomManager/internal/config"
	"RomManager/internal/input"
)

type Footer struct {
	text     string
	renderer *Renderer
	x        int32
	y        int32
	width    int32
	height   int32
	config   *config.Config
}

func (f *Footer) Draw() {
	f.renderer.DrawBox(f.x, f.y, f.width, f.height, f.config.Theme.FooterBackgroundColor)
	f.renderer.DrawText(f.text, f.x+5, f.y+5, 0, f.config.Theme.FooterTextFontSize, f.config.Theme.FooterTextColor, AlignLeft, true)
	f.renderer.DrawLine(f.x, f.y+1, f.width, f.y+1, f.config.Theme.FooterLineColor)
}

func (f *Footer) HandleInput(action input.Action) {

}

func (f *Footer) SetSize(width, height int32) {
	f.width = width
	f.height = height
}

func (f *Footer) SetPosition(x, y int32) {
	f.x = x
	f.y = y
}

func (f *Footer) SetText(text string) {
	f.text = text
}

func NewFooter(text string, renderer *Renderer, c *config.Config) UiElement {
	return &Footer{
		text: text, renderer: renderer, config: c, x: 0, y: 0, width: 0, height: 0,
	}
}
