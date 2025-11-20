package scene

import (
	"RomManager/internal/input"
	"RomManager/internal/ui"
	"RomManager/internal/ui_render"

	"github.com/veandco/go-sdl2/sdl"
)

type MenuScene struct {
	renderer *ui_render.Renderer
	layout   ui.UiElement
}

func (m *MenuScene) Draw() {
	width, height := m.renderer.GetWindowSize()
	m.layout.SetSize(width, height)
	m.layout.Draw()
}

func (m *MenuScene) HandleInput(action input.Action) {

}

func NewMenuScene(renderer *ui_render.Renderer) Scene {
	box := ui.NewBox(sdl.Color{R: 255, G: 0, B: 0, A: 255}, renderer)
	boxLayoutElement := ui.LayoutElement{FullWidth: true, Height: 30, UiElement: box}

	box2 := ui.NewBox(sdl.Color{R: 0, G: 255, B: 0, A: 255}, renderer)
	box2LayoutElement := ui.LayoutElement{FullWidth: true, FullHeight: false, UiElement: box2}

	box3 := ui.NewBox(sdl.Color{R: 255, G: 0, B: 255, A: 255}, renderer)
	box3LayoutElement := ui.LayoutElement{FullWidth: true, Height: 30, UiElement: box3}

	layout := ui.NewVerticalLayout()
	layout.SetPosition(0, 0)

	layout.(*ui.VerticalLayout).AddElement(boxLayoutElement)
	layout.(*ui.VerticalLayout).AddElement(box2LayoutElement)
	layout.(*ui.VerticalLayout).AddElement(box3LayoutElement)

	return &MenuScene{renderer: renderer, layout: layout}
}
