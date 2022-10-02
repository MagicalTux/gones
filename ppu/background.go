package ppu

func (p *PPU) fetchTileData() uint32 {
	return uint32(p.tileData >> 32)
}

func (p *PPU) backgroundPixel() byte {
	if !p.getMask(ShowBg) {
		return 0
	}
	data := p.fetchTileData() >> ((7 - p.X) * 4)
	return byte(data & 0x0F)
}
