package ppu

func (p *PPU) triggerRender() {
	preLine := p.scanline == 261
	visibleLine := p.scanline < 240
	renderLine := preLine || visibleLine
	preFetchCycle := p.cycle >= 321 && p.cycle <= 336
	visibleCycle := p.cycle >= 1 && p.cycle <= 256
	fetchCycle := preFetchCycle || visibleCycle

	if visibleLine && visibleCycle {
		p.renderPixel()
	}
	if renderLine && fetchCycle {
		p.tileData <<= 4
		// see https://www.nesdev.org/w/images/default/d/d1/Ntsc_timing.png
		switch p.cycle & 7 {
		case 1:
			p.fetchNameTableByte()
		case 3:
			p.fetchAttributeTableByte()
		case 5:
			p.fetchLowTileByte()
		case 7:
			p.fetchHighTileByte()
		case 0:
			p.storeTileData()
		}
	}
	if preLine && p.cycle >= 280 && p.cycle <= 304 {
		p.copyY()
	}
	if renderLine {
		if fetchCycle && p.cycle%8 == 0 {
			p.incrementX()
		}
		if p.cycle == 256 {
			p.incrementY()
		}
		if p.cycle == 257 {
			p.copyX()
		}
	}
	if p.cycle == 257 {
		if visibleLine {
			p.evaluateSprites()
		} else {
			p.spriteCount = 0
		}
	}
}

func (p *PPU) renderPixel() {
	x := p.cycle - 1
	y := p.scanline
	background := p.backgroundPixel()
	i, sprite := p.spritePixel()
	if x < 8 && !p.getMask(ShowLeftBg) {
		background = 0
	}
	if x < 8 && !p.getMask(ShowLeftSprites) {
		sprite = 0
	}
	b := background%4 != 0
	s := sprite%4 != 0
	var color byte
	if !b && !s {
		// no background and no sprite → transparent
		color = 0
	} else if !b && s {
		// sprite only
		color = sprite | 0x10
	} else if b && !s {
		// background only
		color = background
	} else {
		// both background & sprite → collision
		if p.spriteIndexes[i] == 0 && x < 255 {
			p.stat &= SpriteZeroHit
		}
		// check sprite priority vs background
		if p.spritePriorities[i] == 0 {
			color = sprite | 0x10
		} else {
			color = background
		}
	}
	c := Palette[p.readPalette(uint16(color))%64]
	p.back.Set(int(x), int(y), c)
}

func (p *PPU) readPalette(address uint16) byte {
	if address >= 16 && address%4 == 0 {
		address -= 16
	}
	return p.Palette[address]
}
