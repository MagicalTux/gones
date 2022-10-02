package ppu

// NTSC Timing Helper Functions

func (p *PPU) incrementX() {
	// increment hori(v)
	// if coarse X == 31
	if p.V&0x001F == 31 {
		// coarse X = 0
		p.V &= 0xFFE0
		// switch horizontal nametable
		p.V ^= 0x0400
	} else {
		// increment coarse X
		p.V++
	}
}

func (p *PPU) incrementY() {
	// increment vert(v)
	// if fine Y < 7
	if p.V&0x7000 != 0x7000 {
		// increment fine Y
		p.V += 0x1000
	} else {
		// fine Y = 0
		p.V &= 0x8FFF
		// let y = coarse Y
		y := (p.V & 0x03E0) >> 5
		if y == 29 {
			// coarse Y = 0
			y = 0
			// switch vertical nametable
			p.V ^= 0x0800
		} else if y == 31 {
			// coarse Y = 0, nametable not switched
			y = 0
		} else {
			// increment coarse Y
			y++
		}
		// put coarse Y back into v
		p.V = (p.V & 0xFC1F) | (y << 5)
	}
}

func (p *PPU) copyX() {
	// hori(v) = hori(t)
	// v: .....F.. ...EDCBA = t: .....F.. ...EDCBA
	p.V = (p.V & 0xFBE0) | (p.T & 0x041F)
}

func (p *PPU) copyY() {
	// vert(v) = vert(t)
	// v: .IHGF.ED CBA..... = t: .IHGF.ED CBA.....
	p.V = (p.V & 0x841F) | (p.T & 0x7BE0)
}
