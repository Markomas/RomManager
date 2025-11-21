package scene

import "RomManager/internal/input"

type Scene interface {
	Draw()
	HandleInput(action input.Action)
	Unload()
}

type SceneManager interface {
	AddScene(scene Scene)
	PopScene()
}
