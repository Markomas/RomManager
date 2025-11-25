package ui

import (
	"RomManager/internal/config"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Renderer struct {
	renderer         *sdl.Renderer
	window           *sdl.Window
	config           *config.Config
	fontCache        map[string]*ttf.Font
	textCache        map[textKey]*textTexture
	imageCache       map[string]*sdl.Texture
	scrollTracker    map[textKey]time.Time
	imageDownloadJob *ImageDownloadJob
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
		renderer:         renderer,
		window:           w,
		config:           c,
		fontCache:        make(map[string]*ttf.Font),
		textCache:        make(map[textKey]*textTexture),
		imageCache:       make(map[string]*sdl.Texture),
		scrollTracker:    make(map[textKey]time.Time),
		imageDownloadJob: NewImageDownloadJob(c),
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

	s, err := font.RenderUTF8Solid(text, sdl.Color(color))
	if err != nil {
		return nil, err
	}
	defer s.Free()

	texture, err := r.renderer.CreateTextureFromSurface(s)
	if err != nil {
		return nil, err
	}

	tt := &textTexture{
		Texture: texture,
		Width:   s.W,
		Height:  s.H,
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
	for _, texture := range r.imageCache {
		texture.Destroy()
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

func (r *Renderer) DrawText(text string, x int32, y int32, maxWidth int32, fontSize int, color config.Color, align Align, scrollOnOverflow bool) {
	if len(text) == 0 {
		return
	}

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

	key := textKey{
		Text:     text,
		FontPath: r.config.Theme.FontPath,
		Size:     fontSize,
		Color:    color,
	}

	var timeInCycle float64
	if scrollOnOverflow {
		startTime, ok := r.scrollTracker[key]
		if !ok {
			startTime = time.Now()
			r.scrollTracker[key] = startTime
		}
		timeInCycle = float64(time.Since(startTime).Milliseconds()%int64(totalCycleDuration*1000)) / 1000.0
	} else {
		delete(r.scrollTracker, key)
	}

	var scrollOffset int32
	if !scrollOnOverflow || timeInCycle < pauseDuration {
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

func (r *Renderer) DrawImage(path string, x int32, y int32, width int32, height int32) {
	if path == "" {
		return
	}

	textureCacheKey := path
	imageLoadPath := path

	// 1. Check if texture is already in cache
	texture, ok := r.imageCache[textureCacheKey]
	if !ok {
		// It's not in cache, we might need to load it
		isURL := strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")

		if isURL {
			if !r.imageDownloadJob.CheckIfFileIsDownloaded(path) {
				go r.imageDownloadJob.AddDownloadJob(path)
				r.imageDownloadJob.Start()
				return
			}
			imageLoadPath = r.imageDownloadJob.UrlToLocalPath(path)
		}

		// 5. Load from disk (either local file or downloaded file)
		surface, err := img.Load(imageLoadPath)
		if err != nil {
			log.Printf("Failed to load image %s: %v", imageLoadPath, err)
			return
		}
		defer surface.Free()

		newTexture, err := r.renderer.CreateTextureFromSurface(surface)
		if err != nil {
			log.Printf("Failed to create texture for %s: %v", imageLoadPath, err)
			return
		}
		r.imageCache[textureCacheKey] = newTexture
		texture = newTexture
	}

	// 6. Calculate dimensions and draw
	_, _, originalWidth, originalHeight, err := texture.Query()
	if err != nil {
		log.Printf("Failed to query texture for %s: %v", path, err)
		return
	}

	// Calculate new dimensions to maintain aspect ratio
	aspectRatio := float32(originalWidth) / float32(originalHeight)
	newWidth := width
	newHeight := int32(float32(newWidth) / aspectRatio)

	if newHeight > height {
		newHeight = height
		newWidth = int32(float32(newHeight) * aspectRatio)
	}

	// Center the image
	newX := x + (width-newWidth)/2
	newY := y + (height-newHeight)/2

	dstRect := sdl.Rect{X: newX, Y: newY, W: newWidth, H: newHeight}
	r.renderer.Copy(texture, nil, &dstRect)
}

func (r *Renderer) DrawTextBox(text string, x int32, y int32, width int32, height int32, size int, color config.Color, verticalAlign Align, horizontalAlign Align) {
	if text == "" || width <= 0 || height <= 0 {
		return
	}

	font, err := r.getFont(r.config.Theme.FontPath, size)
	if err != nil {
		log.Printf("Failed to get font: %v", err)
		return
	}

	lineHeight := int32(font.Height())
	if lineHeight == 0 {
		return
	}

	words := strings.Split(text, " ")
	var lines []string
	currentLine := ""

	for _, word := range words {
		// Handle multiple words in one loop, separated by newlines in original text
		subWords := strings.Split(word, "\n")
		for i, subWord := range subWords {
			testLine := currentLine
			if len(currentLine) > 0 {
				testLine += " "
			}
			testLine += subWord

			lineWidth, _, _ := font.SizeUTF8(testLine)

			if int32(lineWidth) > width {
				if len(currentLine) > 0 {
					lines = append(lines, currentLine)
					currentLine = subWord
				} else { // Word itself is longer than width
					lines = append(lines, subWord)
					currentLine = ""
				}
			} else {
				currentLine = testLine
			}

			// Handle explicit newlines
			if i < len(subWords)-1 {
				lines = append(lines, currentLine)
				currentLine = ""
			}
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	maxLines := int(height / lineHeight)
	if maxLines <= 0 {
		return
	}

	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lastLineIndex := len(lines) - 1
		lastLine := lines[lastLineIndex]

		ellipsis := "..."
		for len(lastLine) > 0 {
			lineWidth, _, _ := font.SizeUTF8(lastLine + ellipsis)
			if int32(lineWidth) <= width {
				break
			}
			// Remove last character
			if len(lastLine) > 0 {
				lastLine = lastLine[:len(lastLine)-1]
			}
		}
		lines[lastLineIndex] = lastLine + ellipsis
	}

	totalTextHeight := int32(len(lines)) * lineHeight
	var startY int32
	switch verticalAlign {
	case AlignVerticalTop:
		startY = y
	case AlignVerticalCenter:
		startY = y + (height-totalTextHeight)/2
	case AlignVerticalBottom:
		startY = y + height - totalTextHeight
	default:
		startY = y
	}

	clipRect := sdl.Rect{X: x, Y: y, W: width, H: height}
	r.renderer.SetClipRect(&clipRect)
	defer r.renderer.SetClipRect(nil)

	for i, line := range lines {
		tt, err := r.getTextTexture(line, size, color)
		if err != nil {
			continue
		}

		var startX int32
		switch horizontalAlign {
		case AlignLeft:
			startX = x
		case AlignCenter:
			startX = x + (width-tt.Width)/2
		case AlignRight:
			startX = x + width - tt.Width
		default:
			startX = x
		}

		dstRect := sdl.Rect{X: startX, Y: startY + int32(i)*lineHeight, W: tt.Width, H: tt.Height}
		r.renderer.Copy(tt.Texture, nil, &dstRect)
	}
}
