package nesinput

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type gamepad struct {
	id ebiten.GamepadID
}

func NewGamepad(id ebiten.GamepadID) (*Generic, error) {
	if !ebiten.IsStandardGamepadLayoutAvailable(id) {
		return nil, fmt.Errorf("no layout available for gamepad %s", ebiten.GamepadName(id))
	}
	return &Generic{ButtonDevice: &gamepad{id}}, nil
}

func (c *gamepad) Pressed(btn byte) bool {
	switch btn {
	case ButtonA:
		return ebiten.IsStandardGamepadButtonPressed(c.id, ebiten.StandardGamepadButtonRightRight)
	case ButtonB:
		return ebiten.IsStandardGamepadButtonPressed(c.id, ebiten.StandardGamepadButtonRightBottom)
	case ButtonSelect:
		return ebiten.IsStandardGamepadButtonPressed(c.id, ebiten.StandardGamepadButtonCenterLeft)
	case ButtonStart:
		return ebiten.IsStandardGamepadButtonPressed(c.id, ebiten.StandardGamepadButtonCenterRight)
	case ButtonUp:
		return ebiten.IsStandardGamepadButtonPressed(c.id, ebiten.StandardGamepadButtonLeftTop)
	case ButtonDown:
		return ebiten.IsStandardGamepadButtonPressed(c.id, ebiten.StandardGamepadButtonLeftBottom)
	case ButtonLeft:
		return ebiten.IsStandardGamepadButtonPressed(c.id, ebiten.StandardGamepadButtonLeftLeft)
	case ButtonRight:
		return ebiten.IsStandardGamepadButtonPressed(c.id, ebiten.StandardGamepadButtonLeftRight)
	default:
		return false
	}
}
