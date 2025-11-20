package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %s\n", err)
		os.Exit(1)
	}
	defer sdl.Quit()

	displayMode, err := sdl.GetDesktopDisplayMode(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get display mode: %s\n", err)
		os.Exit(1)
	}

	window, err := sdl.CreateWindow(
		"Red Square Fullscreen",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		displayMode.W,
		displayMode.H,
		sdl.WINDOW_FULLSCREEN_DESKTOP,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		os.Exit(1)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		os.Exit(1)
	}
	defer renderer.Destroy()

	renderer.SetDrawColor(255, 0, 0, 255) // Set draw color to red
	renderer.Clear()
	renderer.Present()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			}
		}
	}
}
