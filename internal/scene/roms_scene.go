package scene

import (
	"RomManager/internal/api"
	"RomManager/internal/api/romm"
	"RomManager/internal/config"
	"RomManager/internal/downloader"
	"RomManager/internal/input"
	"RomManager/internal/ui"
	"fmt"
)

type RomsScene struct {
	renderer                   *ui.Renderer
	config                     *config.Config
	router                     SceneManager
	layout                     ui.UiElement
	api                        *api.Romm
	id                         int
	stopChan                   chan struct{}
	previewPanelRom            *romm.Rom
	previewImage               ui.UiElement
	previewDescription         ui.UiElement
	previewLayoutLayoutElement *ui.LayoutElement
	downloader                 *downloader.Downloader
}

func (r *RomsScene) Draw() {
	if r.previewDescription.(*ui.Text).GetText() == "" && r.previewImage.(*ui.Image).GetImagePath() == "" {
		r.previewLayoutLayoutElement.Hidden = true
	} else {
		r.previewLayoutLayoutElement.Hidden = false
	}

	width, height := r.renderer.GetWindowSize()
	r.layout.SetSize(width, height)
	r.layout.Draw()
}

func (r *RomsScene) HandleInput(action input.Action) {
	r.layout.HandleInput(action)
}

func (r *RomsScene) Unload() {
	close(r.stopChan)
}

func NewRomsScene(id int, renderer *ui.Renderer, c *config.Config, router SceneManager, api *api.Romm, d *downloader.Downloader) Scene {
	title := ui.NewTitle("Romm Manager - Roms", renderer, c)
	titleLayoutElement := &ui.LayoutElement{FullWidth: true, Height: 30, UiElement: title}

	layout := ui.NewVerticalLayout()
	layout.SetPosition(0, 0)
	layout.(*ui.VerticalLayout).AddElement(titleLayoutElement)

	layoutListPreview := ui.NewHorizontalLayout()
	layoutListPreviewElement := &ui.LayoutElement{FullWidth: true, FullHeight: true, UiElement: layoutListPreview}

	menuItems := []ui.MenuItem{}
	list := ui.NewList(menuItems, renderer, c)
	listLayoutElement := &ui.LayoutElement{FullWidth: true, FullHeight: true, UiElement: list}
	layoutListPreview.(*ui.HorizontalLayout).AddElement(listLayoutElement)

	previewLayout := ui.NewVerticalLayout()
	previewLayoutLayoutElement := &ui.LayoutElement{FullWidth: true, FullHeight: true, UiElement: previewLayout}
	layoutListPreview.(*ui.HorizontalLayout).AddElement(previewLayoutLayoutElement)

	previewImage := ui.NewImage("", renderer, c)
	previewImageLayoutElement := &ui.LayoutElement{FullWidth: true, FullHeight: true, UiElement: previewImage}
	previewLayout.(*ui.VerticalLayout).AddElement(previewImageLayoutElement)

	previewRomTitleText := ui.NewText("",
		renderer,
		c,
		c.Theme.TextColor,
		19,
		ui.AlignLeft,
		ui.AlignVerticalTop)
	previewTitleLayoutElement := &ui.LayoutElement{FullWidth: true, Height: 30, UiElement: previewRomTitleText}
	previewLayout.(*ui.VerticalLayout).AddElement(previewTitleLayoutElement)

	previewDescription := ui.NewText(
		"",
		renderer,
		c,
		c.Theme.TextColor,
		14,
		ui.AlignLeft,
		ui.AlignVerticalTop,
	)
	previewDescriptionLayoutElement := &ui.LayoutElement{FullWidth: true, FullHeight: true, UiElement: previewDescription}
	previewLayout.(*ui.VerticalLayout).AddElement(previewDescriptionLayoutElement)

	layout.(*ui.VerticalLayout).AddElement(layoutListPreviewElement)

	footer := ui.NewFooter("Loading...", renderer, c)
	footerLayoutElement := &ui.LayoutElement{FullWidth: true, Height: 30, UiElement: footer}
	layout.(*ui.VerticalLayout).AddElement(footerLayoutElement)

	stopChan := make(chan struct{})
	r := &RomsScene{
		renderer:                   renderer,
		config:                     c,
		router:                     router,
		layout:                     layout,
		api:                        api,
		id:                         id,
		stopChan:                   stopChan,
		previewPanelRom:            nil,
		previewImage:               previewImage,
		previewDescription:         previewDescription,
		previewLayoutLayoutElement: previewLayoutLayoutElement,
		downloader:                 d,
	}

	go func() {
		offset := 0
		perPage := 100
		for {
			select {
			case <-r.stopChan:
				footer.(*ui.Footer).SetText("Loading cancelled.")
				return
			default:
			}

			roms, _, err := api.GetRomsByPlatform(r.id, offset, perPage)
			if err != nil {
				footer.(*ui.Footer).SetText(fmt.Sprintf("Error fetching roms: %v", err))
				return
			}
			for _, rom := range roms {
				offset += 1
				list.(*ui.List).AddItem(
					fmt.Sprintf("%s (%s)", rom.FsName, formatBytes(rom.FsSizeBytes)),
					func() {
						d.AddRom(rom)
						footer.(*ui.Footer).SetText(fmt.Sprintf("Scheduled download for %s", rom.Name))
					},
					func() {
						if rom.URLCover != "" {
							imagePath := rom.SsMetadata.ScreenshotUrl
							if imagePath == "" {
								imagePath = rom.URLCover
							}
							previewImage.(*ui.Image).SetImagePath(imagePath)
							previewImageLayoutElement.Hidden = false
						} else {
							previewImageLayoutElement.Hidden = true
						}

						previewDescription.(*ui.Text).SetText(rom.Summary)
						previewRomTitleText.(*ui.Text).SetText(rom.Name)
						r.previewPanelRom = &rom
					},
				)
			}
			if len(roms) == 0 {
				footer.(*ui.Footer).SetText(fmt.Sprintf("Loaded %d roms.", offset))
				break
			}
		}
	}()

	return r
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
