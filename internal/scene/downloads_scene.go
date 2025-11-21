package scene

import (
	"RomManager/internal/api"
	"RomManager/internal/config"
	"RomManager/internal/downloader"
	"RomManager/internal/input"
	"RomManager/internal/ui"
	"fmt"
	"time"
)

type DownloadsScene struct {
	renderer             *ui.Renderer
	layout               ui.UiElement
	router               SceneManager
	api                  *api.Romm
	downloader           *downloader.Downloader
	jobIDToListItemIndex map[uint]int
	stopChan             chan struct{}
}

func (d *DownloadsScene) HandleInput(action input.Action) {
	d.layout.HandleInput(action)
}

func (d *DownloadsScene) Unload() {
	close(d.stopChan)
}

func (d *DownloadsScene) Draw() {
	width, height := d.renderer.GetWindowSize()
	d.layout.SetSize(width, height)
	d.layout.Draw()
}

func NewDownloadsScene(renderer *ui.Renderer, c *config.Config, router SceneManager, api *api.Romm, d *downloader.Downloader) Scene {
	title := ui.NewTitle("Romm Manager - Downloads", renderer, c)
	titleLayoutElement := &ui.LayoutElement{FullWidth: true, Height: 30, UiElement: title}

	layout := ui.NewVerticalLayout()
	layout.SetPosition(0, 0)

	list := ui.NewList(make([]ui.MenuItem, 0), renderer, c)
	listLayoutElement := &ui.LayoutElement{FullWidth: true, FullHeight: true, UiElement: list}

	jobIDToListItemIndex := make(map[uint]int)
	stopChan := make(chan struct{})
	go func() {
		updateList(d, jobIDToListItemIndex, list)
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				updateList(d, jobIDToListItemIndex, list)
			case <-stopChan:
				return
			}
		}
	}()

	dd := &DownloadsScene{renderer: renderer, layout: layout, router: router, api: api, downloader: d, jobIDToListItemIndex: jobIDToListItemIndex, stopChan: stopChan}

	layout.(*ui.VerticalLayout).AddElement(titleLayoutElement)
	layout.(*ui.VerticalLayout).AddElement(listLayoutElement)

	return dd
}

func updateList(d *downloader.Downloader, jobIDToListItemIndex map[uint]int, list ui.UiElement) {
	jobs, _ := d.GetDownloadJobs()
	for _, job := range jobs {
		progress := 0.0
		if job.Completed != nil && *job.Completed {
			progress = 100.0
		}
		if job.Progress != nil {
			progress = *job.Progress * 100.0
		}

		if _, ok := jobIDToListItemIndex[job.ID]; !ok {
			jobIDToListItemIndex[job.ID] = list.(*ui.List).AddItem(fmt.Sprintf("%06.2f%% %s", progress, job.Name), func() {}, func() {})
		} else {
			list.(*ui.List).UpdateItemText(jobIDToListItemIndex[job.ID], fmt.Sprintf("%06.2f%% %s", progress, job.Name))
		}
	}
}
