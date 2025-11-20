package ui

import (
	"RomManager/internal/config"
	"RomManager/internal/input"
)

type List struct {
	renderer      *Renderer
	config        *config.Config
	height        int32
	width         int32
	y             int32
	x             int32
	items         []string
	selectedIndex int
	padding       int32
}

func (l *List) Draw() {
	itemHeight := int32(l.config.Theme.ListTextFontSize) + l.padding*2
	for i, item := range l.items {
		yPos := l.y + int32(i)*itemHeight
		if yPos+itemHeight > l.y+l.height {
			break // Don't draw items that are outside the list's bounds
		}

		textColor := l.config.Theme.TextColor
		if i == l.selectedIndex {
			bgColor := l.config.Theme.ListSelectedTextBackgroundColor
			l.renderer.DrawRect(l.x, yPos, l.width, itemHeight, bgColor)
			textColor = l.config.Theme.ListSelectedTextColor
		}

		l.renderer.DrawText(item, l.x+l.padding, yPos+l.padding, 0, l.config.Theme.ListTextFontSize, textColor, AlignLeft)
	}
}

func (l *List) HandleInput(action input.Action) {
	switch action {
	case input.ActionUp:
		if l.selectedIndex > 0 {
			l.selectedIndex--
		}
	case input.ActionDown:
		if l.selectedIndex < len(l.items)-1 {
			l.selectedIndex++
		}
	}
}

func (l *List) SetSize(width, height int32) {
	l.width = width
	l.height = height
}

func (l *List) SetPosition(x, y int32) {
	l.x = x
	l.y = y
}

func NewList(items []string, renderer *Renderer, c *config.Config) UiElement {
	return &List{
		renderer:      renderer,
		config:        c,
		x:             0,
		y:             0,
		width:         0,
		height:        0,
		items:         items,
		selectedIndex: 0,
		padding:       5,
	}
}
