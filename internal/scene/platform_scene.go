package scene

import (
	"RomManager/internal/api"
	"RomManager/internal/config"
	"RomManager/internal/downloader"
	"RomManager/internal/input"
	"RomManager/internal/ui"
	"fmt"
)

type PlatformScene struct {
	renderer *ui.Renderer
	config   *config.Config
	router   SceneManager
	layout   ui.UiElement
	api      *api.Romm
}

func (p *PlatformScene) Unload() {

}

func (p *PlatformScene) Draw() {
	width, height := p.renderer.GetWindowSize()
	p.layout.SetSize(width, height)
	p.layout.Draw()
}

func (p *PlatformScene) HandleInput(action input.Action) {
	p.layout.HandleInput(action)
}

func NewPlatformScene(renderer *ui.Renderer, c *config.Config, router SceneManager, api *api.Romm, d *downloader.Downloader) Scene {
	title := ui.NewTitle("Romm Manager - Platforms", renderer, c)
	titleLayoutElement := &ui.LayoutElement{FullWidth: true, Height: 30, UiElement: title}

	layout := ui.NewVerticalLayout()
	layout.SetPosition(0, 0)
	layout.(*ui.VerticalLayout).AddElement(titleLayoutElement)
	menuItems := []ui.MenuItem{}
	list := ui.NewList(menuItems, renderer, c)

	listLayoutElement := &ui.LayoutElement{FullWidth: true, FullHeight: true, UiElement: list}
	layout.(*ui.VerticalLayout).AddElement(listLayoutElement)

	footer := ui.NewFooter("Loading...", renderer, c)
	footerLayoutElement := &ui.LayoutElement{FullWidth: true, Height: 30, UiElement: footer}
	layout.(*ui.VerticalLayout).AddElement(footerLayoutElement)

	go func() {
		platforms, err := api.GetPlatforms()
		if err != nil {
			footer.(*ui.Footer).SetText(fmt.Sprintf("Error fetching platforms: %v", err))
			return
		}
		for _, platform := range platforms {
			list.(*ui.List).AddItem(
				fmt.Sprintf("%s (%d)", platform.Name, platform.RomCount),
				func() {
					router.AddScene(NewRomsScene(platform.ID, renderer, c, router, api, d))
				},
				func() {
				},
			)
		}
		footer.(*ui.Footer).SetText("")
	}()

	return &PlatformScene{renderer: renderer, config: c, router: router, layout: layout, api: api}
}
