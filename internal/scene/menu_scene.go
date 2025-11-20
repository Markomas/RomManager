package scene

import (
	"RomManager/internal/config"
	"RomManager/internal/input"
	"RomManager/internal/ui"
	"fmt"
)

type MenuScene struct {
	renderer          *ui.Renderer
	layout            ui.UiElement
	previewVisible    bool
	box4LayoutElement *ui.LayoutElement
}

func (m *MenuScene) Draw() {
	width, height := m.renderer.GetWindowSize()
	m.box4LayoutElement.Hidden = !m.previewVisible
	m.layout.SetSize(width, height)
	m.layout.Draw()
}

func (m *MenuScene) HandleInput(action input.Action) {
	m.layout.HandleInput(action)
	m.previewVisible = true
}

func NewMenuScene(renderer *ui.Renderer, c *config.Config) Scene {
	title := ui.NewTitle("Romm Manager", renderer, c)
	titleLayoutElement := &ui.LayoutElement{FullWidth: true, Height: 30, UiElement: title}

	items := []string{"File", "Edit", "View", "Help"}
	//add 100 random items
	for i := 0; i < 100; i++ {
		items = append(items, fmt.Sprintf("Item %d", i))
	}

	box2 := ui.NewList(items, renderer, c)
	box2LayoutElement := &ui.LayoutElement{FullWidth: true, FullHeight: true, UiElement: box2}

	box3 := ui.NewBox(c.Theme.BackgroundColor, renderer)
	box3LayoutElement := &ui.LayoutElement{FullWidth: true, Height: 30, UiElement: box3}

	previewVisible := false
	box4 := ui.NewBox(config.Color{122, 122, 122, 122}, renderer)
	box4LayoutElement := &ui.LayoutElement{Hidden: !previewVisible, Width: 300, FullHeight: true, UiElement: box4}

	horizontalLayout := ui.NewHorizontalLayout()
	horizontalLayout.(*ui.HorizontalLayout).AddElement(box2LayoutElement)
	horizontalLayout.(*ui.HorizontalLayout).AddElement(box4LayoutElement)

	horizontalLayoutElement := &ui.LayoutElement{FullWidth: true, FullHeight: true, UiElement: horizontalLayout}

	layout := ui.NewVerticalLayout()
	layout.SetPosition(0, 0)

	layout.(*ui.VerticalLayout).AddElement(titleLayoutElement)
	layout.(*ui.VerticalLayout).AddElement(horizontalLayoutElement)
	layout.(*ui.VerticalLayout).AddElement(box3LayoutElement)

	return &MenuScene{renderer: renderer, layout: layout, box4LayoutElement: box4LayoutElement, previewVisible: previewVisible}
}
