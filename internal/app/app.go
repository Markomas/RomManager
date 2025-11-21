package app

import (
	"RomManager/internal/api"
	"RomManager/internal/config"
	"RomManager/internal/downloader"
	"RomManager/internal/input"
	"RomManager/internal/router"
	"RomManager/internal/scene"
	"RomManager/internal/ui"
	"fmt"
	"runtime"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type App struct {
	config      *config.Config
	r           *sdl.Renderer
	w           *sdl.Window
	running     bool
	uiRender    *ui.Renderer
	frameCount  uint64
	lastFPS     float64
	lastTime    uint64
	inputMapper *input.Mapper
	sceneRouter *router.Router
}

func New(c *config.Config) (*App, error) {
	fmt.Printf("OS: %s, Arch: %s\n", runtime.GOOS, runtime.GOARCH)
	var err error
	if err = sdl.Init(sdl.INIT_VIDEO | sdl.INIT_GAMECONTROLLER | sdl.INIT_EVENTS); err != nil {
		return nil, err
	}

	if err = img.Init(img.INIT_PNG | img.INIT_JPG); err != nil {
		return nil, err
	}

	if err = ttf.Init(); err != nil {
		return nil, err
	}

	displayMode, err := sdl.GetDesktopDisplayMode(0)
	if err != nil {
		return nil, err
	}

	var w *sdl.Window
	var r *sdl.Renderer

	if runtime.GOARCH == "amd64" {
		w, r, err = sdl.CreateWindowAndRenderer(720, 720, sdl.WINDOW_SHOWN|sdl.WINDOW_ALLOW_HIGHDPI|sdl.WINDOW_RESIZABLE)
	} else {
		w, r, err = sdl.CreateWindowAndRenderer(displayMode.W, displayMode.H, sdl.WINDOW_FULLSCREEN_DESKTOP|sdl.WINDOW_ALLOW_HIGHDPI)
	}

	d := downloader.NewDownloader(c)

	if err != nil {
		return nil, err
	}

	uiRender := ui.New(r, w, c)
	api := api.New(c)
	sceneRouter := router.New()
	menuScene := scene.NewMenuScene(uiRender, c, sceneRouter, api, d)
	sceneRouter.AddScene(menuScene)

	inputMapper := input.New()

	d.Run()

	return &App{config: c, w: w, r: r, uiRender: uiRender, sceneRouter: sceneRouter, inputMapper: inputMapper}, nil
}

func (a *App) Run() {
	a.running = true
	a.lastTime = sdl.GetTicks64()
	for a.running {
		a.handleEvents()
		a.update()
		a.render()

		a.frameCount++
		currentTime := sdl.GetTicks64()
		if currentTime > a.lastTime+1000 {
			a.lastFPS = float64(a.frameCount) * 1000.0 / float64(currentTime-a.lastTime)
			a.lastTime = currentTime
			a.frameCount = 0
		}

		if a.config.System.MaxFPS > 0 {
			frameTime := sdl.GetTicks64() - currentTime
			targetFrameTime := uint64(1000.0 / float64(a.config.System.MaxFPS))
			if frameTime < targetFrameTime {
				sdl.Delay(uint32(targetFrameTime - frameTime))
			}
		}

	}
}

func (a *App) handleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			a.running = false
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_ESCAPE && t.State == sdl.PRESSED {
				a.running = false
			}
		}
		a.inputMapper.ConsumeSDLEvent(event)
	}
}

func (a *App) update() {
	a.inputMapper.ProcessHeldActions()
	for _, action := range a.inputMapper.DrainActions() {
		a.sceneRouter.HandleInput(action)
	}
}

func (a *App) render() {
	a.r.Clear()

	a.sceneRouter.DrawCurrentScene()
	if a.config.System.ShowFPS {
		a.uiRender.DrawText(fmt.Sprintf("FPS: %.0f", a.lastFPS), 10, 10, 0, 16, config.Color{R: 255, G: 255, B: 255, A: 255}, ui.AlignRight, false)
	}
	a.r.Present()
}

func (a *App) Destroy() {
	a.r.Destroy()
	a.w.Destroy()
	ttf.Quit()
	img.Quit()
	sdl.Quit()
}
