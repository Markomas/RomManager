package ui_render

import (
	"RomManager/internal/config"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Renderer struct {
	renderer *sdl.Renderer
	window   *sdl.Window
	config   *config.Config
}

func (r *Renderer) DrawBox(x int32, y int32, w int32, h int32, color sdl.Color) {
	r.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	rect := sdl.Rect{X: x, Y: y, W: w, H: h}
	r.renderer.FillRect(&rect)
}

func (r *Renderer) GetWindowSize() (int32, int32) {
	return r.window.GetSize()
}

func (r *Renderer) DrawTextClipped(text string, x int, y int, clipRect sdl.Rect, fontSize int, color sdl.Color) {
	font, err := ttf.OpenFont(r.config.Theme.FontPath, fontSize)
	if err != nil {
		fmt.Println("Failed to open font:", err)
		return
	}
	defer font.Close()

	surface, err := font.RenderUTF8Solid(text, color)
	if err != nil {
		fmt.Println("Failed to render text:", err)
		return
	}
	defer surface.Free()

	texture, err := r.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		fmt.Println("Failed to create texture from surface:", err)
		return
	}
	defer texture.Destroy()

	r.renderer.SetClipRect(&clipRect)
	dstRect := sdl.Rect{X: int32(x), Y: int32(y), W: surface.W, H: surface.H}
	r.renderer.Copy(texture, nil, &dstRect)
	r.renderer.SetClipRect(nil)
}

func (r *Renderer) DrawText(text string, x int, y int, maxWidth int, fontSize int, color sdl.Color) {
	font, err := ttf.OpenFont(r.config.Theme.FontPath, fontSize)
	if err != nil {
		fmt.Println("Failed to open font:", err)
		return
	}
	defer font.Close()

	textWidth, _, err := font.SizeUTF8(text)
	if err != nil {
		fmt.Println("Failed to get text size:", err)
		return
	}

	if textWidth <= maxWidth || maxWidth == 0 {
		// Text fits, draw normally
		surface, err := font.RenderUTF8Solid(text, color)
		if err != nil {
			fmt.Println("Failed to render text:", err)
			return
		}
		defer surface.Free()

		texture, err := r.renderer.CreateTextureFromSurface(surface)
		if err != nil {
			fmt.Println("Failed to create texture from surface:", err)
			return
		}
		defer texture.Destroy()

		dstRect := sdl.Rect{X: int32(x), Y: int32(y), W: surface.W, H: surface.H}
		r.renderer.Copy(texture, nil, &dstRect)
	} else {
		// Text overflows, scroll it
		scrollSpeed := 35.0  // pixels per second
		pauseDuration := 2.0 // seconds to pause at each end
		extraSpace := 20     // pixels of space before scrolling back

		scrollableWidth := textWidth - maxWidth + extraSpace
		totalCycleDuration := (float64(scrollableWidth)/scrollSpeed)*2.0 + (pauseDuration * 2.0)

		timeInCycle := float64(time.Now().UnixMilli()%int64(totalCycleDuration*1000)) / 1000.0

		var scrollOffset int
		if timeInCycle < pauseDuration {
			scrollOffset = 0
		} else if timeInCycle < pauseDuration+float64(scrollableWidth)/scrollSpeed {
			scrollOffset = int(scrollSpeed * (timeInCycle - pauseDuration))
		} else if timeInCycle < pauseDuration*2.0+float64(scrollableWidth)/scrollSpeed {
			scrollOffset = scrollableWidth
		} else {
			scrollOffset = scrollableWidth - int(scrollSpeed*(timeInCycle-(pauseDuration*2.0+float64(scrollableWidth)/scrollSpeed)))
		}

		clipRect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(maxWidth), H: int32(fontSize * 2)} // Height can be larger
		r.DrawTextClipped(text, x-scrollOffset, y, clipRect, fontSize, color)
	}
}

func New(renderer *sdl.Renderer, w *sdl.Window, c *config.Config) *Renderer {
	return &Renderer{renderer: renderer, window: w, config: c}
}
