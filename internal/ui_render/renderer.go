package ui_render

import (
	"RomManager/internal/config"
	"fmt"

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

func (r *Renderer) DrawText(text string, x int, y int, maxWidth int, fontSize int, color sdl.Color) {
	font, err := ttf.OpenFont(r.config.Theme.FontPath, fontSize)
	if err != nil {
		fmt.Println("Failed to open font:", err)
		return
	}
	defer font.Close()

	textToRender := text
	textWidth, _, err := font.SizeUTF8(textToRender)
	if err != nil {
		fmt.Println("Failed to get text size:", err)
		return
	}

	if textWidth > maxWidth {
		ellipsis := "..."
		for len(textToRender) > 0 {
			runes := []rune(textToRender)
			textToRender = string(runes[:len(runes)-1])
			currentWidth, _, _ := font.SizeUTF8(textToRender + ellipsis)
			if currentWidth <= maxWidth {
				textToRender += ellipsis
				break
			}
		}
	}

	if textToRender == "" {
		return
	}

	surface, err := font.RenderUTF8Solid(textToRender, color)
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

	var dstRect sdl.Rect
	dstRect = sdl.Rect{X: int32(x), Y: int32(y), W: surface.W, H: surface.H}

	r.renderer.Copy(texture, nil, &dstRect)
}

func New(renderer *sdl.Renderer, w *sdl.Window, c *config.Config) *Renderer {
	return &Renderer{renderer: renderer, window: w, config: c}
}
