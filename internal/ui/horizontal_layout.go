package ui

import "RomManager/internal/input"

type HorizontalLayout struct {
	elements []*LayoutElement
	width    int32
	height   int32
	x        int32
	y        int32
}

func (l *HorizontalLayout) Draw() {
	fixedWidthSum := int32(0)
	fullWidthCount := 0
	for _, element := range l.elements {
		if element.Hidden {
			continue
		}
		if element.FullWidth {
			fullWidthCount++
		} else {
			fixedWidthSum += element.Width
		}
	}

	remainingWidth := l.width - fixedWidthSum
	var widthPerFullWidthElement int32
	if fullWidthCount > 0 {
		widthPerFullWidthElement = remainingWidth / int32(fullWidthCount)
	}

	currentX := l.x
	for _, element := range l.elements {
		if element.Hidden {
			continue
		}
		width := element.Width
		if element.FullWidth {
			width = widthPerFullWidthElement
		}

		height := element.Height
		if element.FullHeight {
			height = l.height
		}

		element.UiElement.SetPosition(currentX, l.y)
		element.UiElement.SetSize(width, height)
		element.UiElement.Draw()
		currentX += width
	}
}

func (l *HorizontalLayout) HandleInput(action input.Action) {
	for _, element := range l.elements {
		element.UiElement.HandleInput(action)
	}
}

func (l *HorizontalLayout) SetSize(width, height int32) {
	l.width = width
	l.height = height
}

func (l *HorizontalLayout) SetPosition(x, y int32) {
	l.x = x
	l.y = y
}

func (l *HorizontalLayout) AddElement(element *LayoutElement) {
	l.elements = append(l.elements, element)
}

func NewHorizontalLayout() UiElement {
	elements := make([]*LayoutElement, 0)
	return &HorizontalLayout{elements: elements, width: 0, height: 0, x: 0, y: 0}
}
