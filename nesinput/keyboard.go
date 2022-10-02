package nesinput

import "github.com/hajimehoshi/ebiten/v2"

type keyboard struct{}

func NewKeyboard() *Generic {
	return &Generic{ButtonDevice: &keyboard{}}
}

func (k *keyboard) Pressed(btn byte) bool {
	switch btn {
	case ButtonA:
		return ebiten.IsKeyPressed(ebiten.KeyZ)
	case ButtonB:
		return ebiten.IsKeyPressed(ebiten.KeyX)
	case ButtonSelect:
		return ebiten.IsKeyPressed(ebiten.KeySpace)
	case ButtonStart:
		return ebiten.IsKeyPressed(ebiten.KeyEnter)
	case ButtonUp:
		return ebiten.IsKeyPressed(ebiten.KeyArrowUp)
	case ButtonDown:
		return ebiten.IsKeyPressed(ebiten.KeyArrowDown)
	case ButtonLeft:
		return ebiten.IsKeyPressed(ebiten.KeyArrowLeft)
	case ButtonRight:
		return ebiten.IsKeyPressed(ebiten.KeyArrowRight)
	default:
		return false
	}
}
