package main

const (
	ButtonA = iota
	ButtonB
	ButtonSelect
	ButtonStart
	ButtonUp
	ButtonDown
	ButtonLeft
	ButtonRight
)

type Controller struct {
	Buttons [8]bool
	index   byte
	OUT0    bool
}

func (c *Controller) Read() byte {
	if c.OUT0 {
		// if OUT0 is set, only return state of button 0 (ButtonA)
		if c.Buttons[ButtonA] {
			return 1
		} else {
			return 0
		}
	}

	var val byte
	if c.index < 8 && c.Buttons[c.index] {
		val = 1
	}
	c.index += 1
	return val
}

func (c *Controller) Write(value byte) {
	c.OUT0 = value&1 == 1
	if c.OUT0 {
		c.index = 0
	}
}
