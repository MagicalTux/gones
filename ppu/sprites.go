package ppu

func (p *PPU) evaluateSprites() {
	var h int
	if p.getFlag(WideSprites) {
		h = 16
	} else {
		h = 8
	}
	count := 0
	for i := 0; i < 64; i++ {
		y := p.OAM[i*4+0]
		a := p.OAM[i*4+2]
		x := p.OAM[i*4+3]
		row := int(p.scanline) - int(y)
		if row < 0 || row >= h {
			continue
		}
		if count < 8 {
			p.spritePatterns[count] = p.fetchSpritePattern(i, row)
			p.spritePositions[count] = x
			p.spritePriorities[count] = (a >> 5) & 1
			p.spriteIndexes[count] = byte(i)
		}
		count++
	}
	if count > 8 {
		count = 8
		p.stat &= SpriteOverflow
	}
	p.spriteCount = count
}

func (p *PPU) spriteTableBase() uint16 {
	if p.getFlag(AltSprites) {
		return 0x1000
	} else {
		return 0
	}
}

func (p *PPU) fetchSpritePattern(i, row int) uint32 {
	tile := p.OAM[i*4+1]
	attributes := p.OAM[i*4+2]
	address := p.spriteTableBase()
	if !p.getFlag(WideSprites) {
		if attributes&0x80 == 0x80 {
			row = 7 - row
		}
		address += uint16(tile)*16 + uint16(row)
	} else {
		if attributes&0x80 == 0x80 {
			row = 15 - row
		}
		table := tile & 1
		tile &= 0xFE
		if row > 7 {
			tile++
			row -= 8
		}
		address = (uint16(table) << 12) | uint16(tile)<<4 | uint16(row)
	}
	a := (attributes & 3) << 2
	lowTileByte := p.Memory.MemRead(address)
	highTileByte := p.Memory.MemRead(address + 8)
	var data uint32
	for i := 0; i < 8; i++ {
		var p1, p2 byte
		if attributes&0x40 == 0x40 {
			p1 = (lowTileByte & 1) << 0
			p2 = (highTileByte & 1) << 1
			lowTileByte >>= 1
			highTileByte >>= 1
		} else {
			p1 = (lowTileByte & 0x80) >> 7
			p2 = (highTileByte & 0x80) >> 6
			lowTileByte <<= 1
			highTileByte <<= 1
		}
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	return data
}

func (p *PPU) spritePixel() (byte, byte) {
	if !p.getMask(ShowSprites) {
		return 0, 0
	}
	for i := 0; i < p.spriteCount; i++ {
		offset := (int(p.cycle) - 1) - int(p.spritePositions[i])
		if offset < 0 || offset > 7 {
			continue
		}
		offset = 7 - offset
		color := byte((p.spritePatterns[i] >> byte(offset*4)) & 0x0F)
		if color%4 == 0 {
			continue
		}
		return byte(i), color
	}
	return 0, 0
}
