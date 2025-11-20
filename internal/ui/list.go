package ui

import (
	"RomManager/internal/config"
	"RomManager/internal/input"
	"fmt"
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
	topIndex      int
}

func (l *List) Draw() {
	itemHeight := int32(l.config.Theme.ListTextFontSize) + l.padding*2
	if itemHeight <= 0 {
		return
	}
	for i := 0; ; i++ {
		itemIndex := l.topIndex + i
		if itemIndex >= len(l.items) {
			break
		}

		yPos := l.y + int32(i)*itemHeight
		if yPos+itemHeight > l.y+l.height {
			break // Don't draw items that are outside the list's bounds
		}

		item := l.items[itemIndex]
		textColor := l.config.Theme.TextColor
		if itemIndex == l.selectedIndex {
			bgColor := l.config.Theme.ListSelectedTextBackgroundColor
			l.renderer.DrawRect(l.x, yPos, l.width, itemHeight, bgColor)
			textColor = l.config.Theme.ListSelectedTextColor
		}

		l.renderer.DrawText(item, l.x+l.padding, yPos+l.padding, 0, l.config.Theme.ListTextFontSize, textColor, AlignLeft)
	}
}

func (l *List) HandleInput(action input.Action) {
	fmt.Printf("Handling input: %v\n", action)
	itemHeight := int32(l.config.Theme.ListTextFontSize) + l.padding*2
	if itemHeight <= 0 {
		return
	}
	visibleItems := int(l.height / itemHeight)

	switch action {
	case input.ActionJumpUp:
		l.selectedIndex -= visibleItems
		if l.selectedIndex < 0 {
			l.selectedIndex = 0
		}
		l.topIndex = l.selectedIndex
	case input.ActionJumpDown:
		l.selectedIndex += visibleItems
		if l.selectedIndex >= len(l.items) {
			l.selectedIndex = len(l.items) - 1
		}
		l.topIndex = l.selectedIndex - visibleItems + 1
		if l.topIndex > len(l.items)-visibleItems {
			l.topIndex = len(l.items) - visibleItems
		}
		if l.topIndex < 0 {
			l.topIndex = 0
		}
	case input.ActionUp:
		if l.selectedIndex > 0 {
			l.selectedIndex--
			if l.selectedIndex < l.topIndex {
				l.topIndex = l.selectedIndex
			}
		}
	case input.ActionDown:
		if l.selectedIndex < len(l.items)-1 {
			l.selectedIndex++
			if l.selectedIndex >= l.topIndex+visibleItems {
				l.topIndex = l.selectedIndex - visibleItems + 1
			}
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
		topIndex:      0,
	}
}
