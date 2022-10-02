package ppu

// fetch/store pipeline methods

func (p *PPU) fetchNameTableByte() {
	addr := 0x2000 | (p.V & 0x0FFF)
	p.nameTableByte = p.Memory.MemRead(addr)
}

func (p *PPU) fetchAttributeTableByte() {
	v := p.V
	addr := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
	shift := ((v >> 4) & 4) | (v & 2)
	p.attributeTableByte = ((p.Memory.MemRead(addr) >> shift) & 3) << 2
}

func (p *PPU) currentTileAddress() uint16 {
	fineY := (p.V >> 12) & 7
	table := p.ctrl & NameTableMask // 0, 1, 2, 3
	tile := p.nameTableByte
	return (uint16(table) << 12) | (uint16(tile) << 4) | fineY
}

func (p *PPU) fetchLowTileByte() {
	p.lowTileByte = p.Memory.MemRead(p.currentTileAddress())
}

func (p *PPU) fetchHighTileByte() {
	p.highTileByte = p.Memory.MemRead(p.currentTileAddress() + 8)
}

func (p *PPU) storeTileData() {
	var data uint32
	a := p.attributeTableByte

	for i := 0; i < 8; i++ {
		p1 := (p.lowTileByte & 0x80) >> 7
		p2 := (p.highTileByte & 0x80) >> 6
		p.lowTileByte <<= 1
		p.highTileByte <<= 1
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	p.tileData |= uint64(data)
}
