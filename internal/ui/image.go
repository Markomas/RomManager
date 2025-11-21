package ui

import (
	"RomManager/internal/config"
	"RomManager/internal/input"
)

type Image struct {
	renderer  *Renderer
	x         int32
	y         int32
	width     int32
	height    int32
	config    *config.Config
	imagePath string
}

func (i *Image) Draw() {
	i.renderer.DrawBox(i.x, i.y, i.width, i.height, i.config.Theme.ImageBackgroundColor)
	i.renderer.DrawImage(i.imagePath, i.x, i.y, i.width, i.height)
}

func (i *Image) HandleInput(action input.Action) {

}

func (i *Image) SetSize(width, height int32) {
	i.width = width
	i.height = height
}

func (i *Image) SetPosition(x, y int32) {
	i.x = x
	i.y = y
}

func (i *Image) SetImagePath(path string) {
	i.imagePath = path
}

func (i *Image) GetImagePath() string {
	return i.imagePath
}

func NewImage(imagePath string, renderer *Renderer, c *config.Config) UiElement {
	return &Image{
		renderer:  renderer,
		imagePath: imagePath,
		config:    c,
		x:         0,
		y:         0,
		width:     0,
		height:    0,
	}
}
