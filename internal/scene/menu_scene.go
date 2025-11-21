package scene

import (
	"RomManager/internal/api"
	"RomManager/internal/config"
	"RomManager/internal/downloader"
	"RomManager/internal/input"
	"RomManager/internal/ui"
	"fmt"
)

type MenuScene struct {
	renderer   *ui.Renderer
	layout     ui.UiElement
	router     SceneManager
	api        *api.Romm
	downloader *downloader.Downloader
}

func (m *MenuScene) Unload() {

}

func (m *MenuScene) Draw() {
	width, height := m.renderer.GetWindowSize()
	m.layout.SetSize(width, height)
	m.layout.Draw()
}

func (m *MenuScene) HandleInput(action input.Action) {
	m.layout.HandleInput(action)
}

func NewMenuScene(renderer *ui.Renderer, c *config.Config, router SceneManager, api *api.Romm, d *downloader.Downloader) Scene {
	title := ui.NewTitle("Romm Manager", renderer, c)
	titleLayoutElement := &ui.LayoutElement{FullWidth: true, Height: 30, UiElement: title}

	footer := ui.NewBox(c.Theme.BackgroundColor, renderer)
	footerLayoutElement := &ui.LayoutElement{FullWidth: true, Height: 30, UiElement: footer}

	layout := ui.NewVerticalLayout()
	layout.SetPosition(0, 0)
	m := &MenuScene{renderer: renderer, layout: layout, router: router, api: api, downloader: d}

	layout.(*ui.VerticalLayout).AddElement(titleLayoutElement)
	menuItems := []ui.MenuItem{
		{"Platforms", func() { router.AddScene(NewPlatformScene(renderer, c, m.router, api, d)) }, func() {}},
		{"Downloads", func() { router.AddScene(NewDownloadsScene(renderer, c, m.router, api, d)) }, func() {}},
	}
	list := ui.NewList(menuItems, renderer, c)

	listLayoutElement := &ui.LayoutElement{FullWidth: true, FullHeight: true, UiElement: list}
	layout.(*ui.VerticalLayout).AddElement(listLayoutElement)
	layout.(*ui.VerticalLayout).AddElement(footerLayoutElement)

	return m
}

func (m *MenuScene) SelectItem(i int) {
	fmt.Printf("Selected item: %d\n", i)
}
