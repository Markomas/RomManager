package ui

import (
	"RomManager/internal/config"
	"RomManager/internal/input"
)

type MenuItem struct {
	Label    string
	Action   func()
	OnChange func()
}

type List struct {
	renderer      *Renderer
	config        *config.Config
	height        int32
	width         int32
	y             int32
	x             int32
	items         []MenuItem
	selectedIndex int
	padding       int32
	topIndex      int
}

func (l *List) Draw() {
	l.renderer.DrawRect(l.x, l.y, l.width, l.height, l.config.Theme.ListBackgroundColor)
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
		overflowScroll := false
		if itemIndex == l.selectedIndex {
			bgColor := l.config.Theme.ListSelectedTextBackgroundColor
			l.renderer.DrawRect(l.x, yPos, l.width, itemHeight, bgColor)
			textColor = l.config.Theme.ListSelectedTextColor
			overflowScroll = true
		}

		l.renderer.DrawText(item.Label, l.x+l.padding, yPos+l.padding, l.width, l.config.Theme.ListTextFontSize, textColor, AlignLeft, overflowScroll)
	}
}

func (l *List) HandleInput(action input.Action) {
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
		if l.selectedIndex < len(l.items) && l.items[l.selectedIndex].OnChange != nil {
			l.items[l.selectedIndex].OnChange()
		}

		break
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
		if l.selectedIndex < len(l.items) && l.items[l.selectedIndex].OnChange != nil {
			l.items[l.selectedIndex].OnChange()
		}
		break
	case input.ActionUp:
		if l.selectedIndex > 0 {
			l.selectedIndex--
			if l.selectedIndex < l.topIndex {
				l.topIndex = l.selectedIndex
			}
		}
		if l.selectedIndex < len(l.items) && l.items[l.selectedIndex].OnChange != nil {
			l.items[l.selectedIndex].OnChange()
		}
		l.items[l.selectedIndex].OnChange()
		break
	case input.ActionDown:
		if l.selectedIndex < len(l.items)-1 {
			l.selectedIndex++
			if l.selectedIndex >= l.topIndex+visibleItems {
				l.topIndex = l.selectedIndex - visibleItems + 1
			}
		}
		if l.selectedIndex < len(l.items) && l.items[l.selectedIndex].OnChange != nil {
			l.items[l.selectedIndex].OnChange()
		}
		break
	case input.ActionSelect:
		if l.selectedIndex >= 0 && l.selectedIndex < len(l.items) {
			l.items[l.selectedIndex].Action()
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

func (l *List) AddItem(label string, action func(), onChange func()) int {
	l.items = append(l.items, MenuItem{Label: label, Action: action, OnChange: onChange})
	if l.selectedIndex < len(l.items) && l.items[l.selectedIndex].OnChange != nil {
		l.items[l.selectedIndex].OnChange()
	}
	return len(l.items) - 1
}

func (l *List) ClearItems() {
	l.items = make([]MenuItem, 0)
}

func (l *List) UpdateItemText(i int, text string) {
	if i >= 0 && i < len(l.items) {
		l.items[i].Label = text
	}
}

func NewList(items []MenuItem, renderer *Renderer, c *config.Config) UiElement {
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
