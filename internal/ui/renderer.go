package ui

import (
	"RomManager/internal/config"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Renderer struct {
	renderer  *sdl.Renderer
	window    *sdl.Window
	config    *config.Config
	fontCache map[string]*ttf.Font

	textCache map[textKey]*textTexture
}

type textKey struct {
	Text     string
	FontPath string
	Size     int
	Color    config.Color
}

type textTexture struct {
	Texture *sdl.Texture
	Width   int32
	Height  int32
}

func New(renderer *sdl.Renderer, w *sdl.Window, c *config.Config) *Renderer {
	return &Renderer{
		renderer:  renderer,
		window:    w,
		config:    c,
		fontCache: make(map[string]*ttf.Font),
		textCache: make(map[textKey]*textTexture),
	}
}

func (r *Renderer) getFont(path string, size int) (*ttf.Font, error) {
	key := fmt.Sprintf("%s-%d", path, size)
	if font, ok := r.fontCache[key]; ok {
		return font, nil
	}

	font, err := ttf.OpenFont(path, size)
	if err != nil {
		return nil, err
	}

	r.fontCache[key] = font
	return font, nil
}

func (r *Renderer) getTextTexture(text string, fontSize int, color config.Color) (*textTexture, error) {
	key := textKey{
		Text:     text,
		FontPath: r.config.Theme.FontPath,
		Size:     fontSize,
		Color:    color,
	}

	if tt, ok := r.textCache[key]; ok {
		return tt, nil
	}

	font, err := r.getFont(r.config.Theme.FontPath, fontSize)
	if err != nil {
		return nil, err
	}

	surface, err := font.RenderUTF8Solid(text, sdl.Color(color))
	if err != nil {
		return nil, err
	}
	defer surface.Free()

	texture, err := r.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}

	tt := &textTexture{
		Texture: texture,
		Width:   surface.W,
		Height:  surface.H,
	}

	r.textCache[key] = tt
	return tt, nil
}

func (r *Renderer) Destroy() {
	for _, font := range r.fontCache {
		font.Close()
	}
	for _, tt := range r.textCache {
		tt.Texture.Destroy()
	}
}

func (r *Renderer) DrawBox(x int32, y int32, w int32, h int32, color config.Color) {
	r.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	rect := sdl.Rect{X: x, Y: y, W: w, H: h}
	r.renderer.FillRect(&rect)
}

func (r *Renderer) GetWindowSize() (int32, int32) {
	return r.window.GetSize()
}

func (r *Renderer) DrawTextClipped(text string, x int32, y int32, clipRect sdl.Rect, fontSize int, color config.Color) {
	tt, err := r.getTextTexture(text, fontSize, color)
	if err != nil {
		// optionally log once elsewhere; avoid spamming in render loop
		fmt.Println("Failed to get text texture:", err)
		return
	}

	r.renderer.SetClipRect(&clipRect)
	dstRect := sdl.Rect{X: x, Y: y, W: tt.Width, H: tt.Height}
	r.renderer.Copy(tt.Texture, nil, &dstRect)
	r.renderer.SetClipRect(nil)
}

func (r *Renderer) DrawText(text string, x int32, y int32, maxWidth int32, fontSize int, color config.Color, align Align) {
	if maxWidth == 0 {
		windowWidth, _ := r.window.GetSize()
		maxWidth = windowWidth
	}

	tt, err := r.getTextTexture(text, fontSize, color)
	if err != nil {
		fmt.Println("Failed to get text texture:", err)
		return
	}

	textWidth := tt.Width

	if textWidth <= maxWidth {
		// Text fits, draw normally
		var drawX int32
		switch align {
		case AlignCenter:
			drawX = x + (maxWidth-textWidth)/2
		case AlignRight:
			drawX = x + maxWidth - textWidth - 10
		default: // AlignLeft
			drawX = x
		}

		dstRect := sdl.Rect{X: drawX, Y: y, W: tt.Width, H: tt.Height}
		r.renderer.Copy(tt.Texture, nil, &dstRect)
		return
	}

	// Text overflows, scroll it
	scrollSpeed := 35.0  // pixels per second
	pauseDuration := 2.0 // seconds to pause at each end
	extraSpace := int32(20)

	scrollableWidth := textWidth - maxWidth + extraSpace
	if scrollableWidth < 0 {
		scrollableWidth = 0
	}

	totalCycleDuration := (float64(scrollableWidth)/scrollSpeed)*2.0 + (pauseDuration * 2.0)
	if totalCycleDuration <= 0 {
		totalCycleDuration = pauseDuration * 2.0
	}

	timeInCycle := float64(time.Now().UnixMilli()%int64(totalCycleDuration*1000)) / 1000.0

	var scrollOffset int32
	if timeInCycle < pauseDuration {
		scrollOffset = 0
	} else if timeInCycle < pauseDuration+float64(scrollableWidth)/scrollSpeed {
		scrollOffset = int32(scrollSpeed * (timeInCycle - pauseDuration))
	} else if timeInCycle < pauseDuration*2.0+float64(scrollableWidth)/scrollSpeed {
		scrollOffset = scrollableWidth
	} else {
		scrollOffset = scrollableWidth - int32(scrollSpeed*(timeInCycle-(pauseDuration*2.0+float64(scrollableWidth)/scrollSpeed)))
	}

	clipRect := sdl.Rect{X: x, Y: y, W: maxWidth, H: int32(fontSize * 2)}
	r.renderer.SetClipRect(&clipRect)

	dstRect := sdl.Rect{X: x - scrollOffset, Y: y, W: tt.Width, H: tt.Height}
	r.renderer.Copy(tt.Texture, nil, &dstRect)

	r.renderer.SetClipRect(nil)
}

func (r *Renderer) DrawLine(x1 int32, y1 int32, x2 int32, y2 int32, color config.Color) {
	r.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	r.renderer.DrawLine(x1, y1, x2, y2)
}

func (r *Renderer) DrawRect(x int32, pos int32, width int32, height int32, color config.Color) {
	r.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	rect := sdl.Rect{X: x, Y: pos, W: width, H: height}
	r.renderer.FillRect(&rect)
}
