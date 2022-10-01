package ppu

const (
	SpriteOverflow byte = 0x20
	SpriteZeroHit  byte = 0x40 // Set when a nonzero pixel of sprite 0 overlaps a nonzero background pixel
	VBlankStarted  byte = 0x80

	NameTableMask   byte = 0x03 // 0, 1, 2 or 3
	LargeIncrements byte = 0x04 // VRAM address increment per CPU read/write of PPUDATA (0: add 1, going across; 1: add 32, going down)
	AltSprites      byte = 0x08 // Sprite pattern table address for 8x8 sprites (0: $0000; 1: $1000; ignored in 8x16 mode)
	AltBackground   byte = 0x10 // Background pattern table address (0: $0000; 1: $1000)
	WideSprites     byte = 0x20 // Sprite size (0: 8x8 pixels; 1: 8x16 pixels – see PPU OAM#Byte 1)
	PpuMaster       byte = 0x40 // PPU master/slave select (0: read backdrop from EXT pins; 1: output color on EXT pins)
	GenerateNMI     byte = 0x80 // Generate an NMI at the start of the vertical blanking interval (0: off; 1: on)

	PPUCTRL   = 0
	PPUMASK   = 1
	PPUSTATUS = 2
	OAMADDR   = 3
	OAMDATA   = 4
	PPUSCROLL = 5
	PPUADDR   = 6
	PPUDATA   = 7
)
