package router

import (
	"RomManager/internal/input"
	"RomManager/internal/scene"
	"os"
)

type SceneManager interface {
	AddScene(scene scene.Scene)
	PopScene()
}

func New() *Router {
	return &Router{
		scenes: make([]scene.Scene, 0),
	}
}

type Router struct {
	scenes []scene.Scene
}

func (r *Router) AddScene(scene scene.Scene) {
	r.scenes = append(r.scenes, scene)
}

func (r *Router) HandleInput(action input.Action) {
	if action == input.ActionBack {
		r.PopScene()
		if len(r.scenes) <= 0 {
			os.Exit(0)
		}
	}
	if len(r.scenes) > 0 {
		currentScene := r.scenes[len(r.scenes)-1]
		currentScene.HandleInput(action)
	}
}

func (r *Router) DrawCurrentScene() {
	if len(r.scenes) > 0 {
		currentScene := r.scenes[len(r.scenes)-1]
		currentScene.Draw()
	}
}

func (r *Router) PopScene() {
	if len(r.scenes) > 0 {
		currentScene := r.scenes[len(r.scenes)-1]
		currentScene.Unload()
		r.scenes = r.scenes[:len(r.scenes)-1]
	}
}
