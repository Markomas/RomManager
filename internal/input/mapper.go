package input

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	repeatDelay    = 500 // ms
	repeatInterval = 10  // ms
)

type heldKey struct {
	action        Action
	pressTime     uint32
	nextRepeat    uint32
	isInitialSent bool
}

type Mapper struct {
	queue       []Action
	heldKeys    map[sdl.Keycode]*heldKey
	heldButtons map[uint8]*heldKey
}

func New() *Mapper {
	return &Mapper{
		queue:       make([]Action, 0, 16),
		heldKeys:    make(map[sdl.Keycode]*heldKey),
		heldButtons: make(map[uint8]*heldKey),
	}
}

func (m *Mapper) push(a Action) {
	m.queue = append(m.queue, a)
}

func (m *Mapper) DrainActions() []Action {
	out := m.queue
	m.queue = m.queue[:0]
	return out
}

func (m *Mapper) ProcessHeldActions() {
	now := sdl.GetTicks()
	for _, k := range m.heldKeys {
		if !k.isInitialSent {
			m.push(k.action)
			k.isInitialSent = true
			k.nextRepeat = now + repeatDelay
		} else if now >= k.nextRepeat {
			m.push(k.action)
			k.nextRepeat = now + repeatInterval
		}
	}
	for _, b := range m.heldButtons {
		if !b.isInitialSent {
			m.push(b.action)
			b.isInitialSent = true
			b.nextRepeat = now + repeatDelay
		} else if now >= b.nextRepeat {
			m.push(b.action)
			b.nextRepeat = now + repeatInterval
		}
	}
}

func (m *Mapper) ConsumeSDLEvent(ev sdl.Event) bool {
	switch e := ev.(type) {
	case *sdl.KeyboardEvent:
		var action Action
		switch e.Keysym.Sym {
		case sdl.K_UP:
			action = ActionUp
		case sdl.K_DOWN:
			action = ActionDown
		case sdl.K_LEFT:
			action = ActionLeft
		case sdl.K_RIGHT:
			action = ActionRight
		case sdl.K_RETURN, sdl.K_SPACE:
			action = ActionSelect
		case sdl.K_BACKSPACE:
			action = ActionBack
		case sdl.K_ESCAPE:
			action = ActionQuit
		case sdl.K_PAGEUP:
			action = ActionJumpUp
		case sdl.K_PAGEDOWN:
			action = ActionJumpDown
		default:
			return false
		}

		if e.State == sdl.PRESSED {
			if _, exists := m.heldKeys[e.Keysym.Sym]; !exists {
				now := sdl.GetTicks()
				m.heldKeys[e.Keysym.Sym] = &heldKey{
					action:    action,
					pressTime: now,
				}
			}
		} else if e.State == sdl.RELEASED {
			delete(m.heldKeys, e.Keysym.Sym)
		}
		return true

	case *sdl.ControllerDeviceEvent:
		if e.Type == sdl.CONTROLLERDEVICEADDED {
			fmt.Println("Controller connected:", e.Which)
			sdl.GameControllerOpen(int(e.Which))
		}
		return true

	case *sdl.ControllerButtonEvent:
		var action Action
		switch e.Button {
		case sdl.CONTROLLER_BUTTON_DPAD_UP:
			action = ActionUp
		case sdl.CONTROLLER_BUTTON_DPAD_DOWN:
			action = ActionDown
		case sdl.CONTROLLER_BUTTON_DPAD_LEFT:
			action = ActionLeft
		case sdl.CONTROLLER_BUTTON_DPAD_RIGHT:
			action = ActionRight
		case sdl.CONTROLLER_BUTTON_A:
			action = ActionSelect
		case sdl.CONTROLLER_BUTTON_B:
			action = ActionBack
		case sdl.CONTROLLER_BUTTON_GUIDE:
			action = ActionQuit
		case sdl.CONTROLLER_BUTTON_RIGHTSHOULDER:
			action = ActionJumpDown
		case sdl.CONTROLLER_BUTTON_LEFTSHOULDER:
			action = ActionJumpUp
		default:
			return false
		}

		if e.State == sdl.PRESSED {
			if _, exists := m.heldButtons[e.Button]; !exists {
				now := sdl.GetTicks()
				m.heldButtons[e.Button] = &heldKey{
					action:    action,
					pressTime: now,
				}
			}
		} else if e.State == sdl.RELEASED {
			delete(m.heldButtons, e.Button)
		}
		return true
	}
	return false
}
