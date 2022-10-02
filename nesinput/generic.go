package nesinput

type ButtonDevice interface {
	Pressed(key byte) bool
}

type Generic struct {
	ButtonDevice
	index byte
	OUT0  bool
}

func (c *Generic) Read() byte {
	if c.OUT0 {
		// if OUT0 is set, only return state of button 0 (ButtonA)
		if c.Pressed(ButtonA) {
			return 1
		} else {
			return 0
		}
	}

	var val byte
	if c.index < 8 && c.Pressed(c.index) {
		val = 1
	}
	c.index += 1
	return val
}

func (c *Generic) Write(value byte) {
	c.OUT0 = value&1 == 1
	if c.OUT0 {
		c.index = 0
	}
}
