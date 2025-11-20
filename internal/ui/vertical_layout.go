package ui

import "RomManager/internal/input"

type VerticalLayout struct {
	elements []LayoutElement
	width    int32
	height   int32
	x        int32
	y        int32
}

type LayoutElement struct {
	FullWidth  bool
	Width      int32
	FullHeight bool
	Height     int32
	UiElement  UiElement
}

func (l *VerticalLayout) Draw() {
	fixedHeightSum := int32(0)
	fullHeightCount := 0
	for _, element := range l.elements {
		if element.FullHeight {
			fullHeightCount++
		} else {
			fixedHeightSum += element.Height
		}
	}

	remainingHeight := l.height - fixedHeightSum
	var heightPerFullHeightElement int32
	if fullHeightCount > 0 {
		heightPerFullHeightElement = remainingHeight / int32(fullHeightCount)
	}

	currentY := l.y
	for _, element := range l.elements {
		width := element.Width
		if element.FullWidth {
			width = l.width
		}

		height := element.Height
		if element.FullHeight {
			height = heightPerFullHeightElement
		}

		element.UiElement.SetPosition(l.x, currentY)
		element.UiElement.SetSize(width, height)
		element.UiElement.Draw()
		currentY += height
	}
}

func (l *VerticalLayout) HandleInput(action input.Action) {

}

func (l *VerticalLayout) SetSize(width, height int32) {
	l.width = width
	l.height = height
}

func (l *VerticalLayout) SetPosition(x, y int32) {
	l.x = x
	l.y = y
}

func (l *VerticalLayout) AddElement(element LayoutElement) {
	l.elements = append(l.elements, element)
}

func NewVerticalLayout() UiElement {
	elements := make([]LayoutElement, 0)
	return &VerticalLayout{elements: elements, width: 0, height: 0, x: 0, y: 0}
}
