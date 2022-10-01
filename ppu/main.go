package ppu

import (
	"github.com/MagicalTux/gones/memory"
	"github.com/hajimehoshi/ebiten/v2"
)

// https://www.nesdev.org/wiki/PPU_rendering#Frame_timing_diagram
// The PPU renders 262 scanlines per frame. Each scanline lasts for 341 PPU clock cycles (113.667 CPU clock cycles; 1 CPU cycle = 3 PPU cycles), with each clock cycle producing one pixel.

type PPU struct {
	Memory memory.Master

	ctrl    byte // 0x00 at start
	mask    byte // 0x00 at start
	stat    byte // 0x00 at start
	scroll  byte
	data    byte
	OAM     [256]byte // can be edited by the CPU via DMA
	Palette [32]byte  // color palette

	cycle    uint16
	scanline uint16
	oddframe bool // frame is even/odd (starting at frame 0 which is even)
	frame    uint64

	vblankFlag bool
	vblankNMI  bool

	oamAddr byte   // used to read/write OAM data
	ppuAddr uint16 // addr for PPUDATA

	// https://www.nesdev.org/wiki/PPU_scrolling#PPU_internal_registers
	V uint16 // Current VRAM address (15 bits)
	T uint16 // Temporary VRAM address (15 bits); can also be thought of as the address of the top left onscreen tile.
	X byte   // Fine X scroll (3 bits)
	W bool   // First or second write toggle (1 bit)

	readBuf byte // read buffer for PPUDATA

	front, back *ebiten.Image
	vblank      func()
}

func New() *PPU {
	ppu := &PPU{
		Memory: memory.NewBus(),
		front:  ebiten.NewImage(256, 240),
		back:   ebiten.NewImage(256, 240),
	}

	// https://www.nesdev.org/wiki/PPU_memory_map

	// pattern tables
	// 0x0000~0x2000 is mapped via CHR ROM to the cartridge (mapping is done by the cartridge)

	// nametables
	// 0x2000~0x3f00 is mapped typically to the PPU ram, but the cartridge can clear this mapping to map something else
	// we actually map the whole thing, but reads to values after 0x3f00 will return data from the palette
	ppu.Memory.MapHandler(0x2000, 0x2000, memory.NewRAM(0x800))

	return ppu
}

func (p *PPU) VblankInterrupt(cb func()) {
	p.vblank = cb
}

func (p *PPU) Reset(cnt uint64) {
	// cnt tells us at what point we need to reset, typically 7
	pxls := cnt * 3 // =21

	p.cycle = uint16(pxls)
	p.scanline = 0
	p.frame = 0
}

func (p *PPU) Clock(cnt int) {
	// move clock forward by cnt (1 CPU clock = 3 PPU clock)
	p.cycle += uint16(cnt * 3)

	// each PPU frame is 341*262=89342 PPU clocks long
	p.checkPendingNMI()

	for p.cycle >= 341 {
		p.cycle -= 341
		p.scanline += 1
		if p.scanline >= 241 && !p.vblankFlag {
			p.vblankFlag = true
			p.vblankNMI = true // generate NMI at next available occasion
			p.stat |= VBlankStarted
		}
		if p.scanline >= 261 && p.vblankFlag {
			// clear vblank, sprite, overflow
			p.vblankFlag = false
			p.vblankNMI = false
			// clear SpriteZeroHit & VBlankStarted from p.stat
			p.stat &= ^(SpriteZeroHit | VBlankStarted)
			if p.vblank != nil {
				// trigger vblank interrupt
				p.vblank()
			}
		}
		if p.scanline >= 262 {
			p.scanline -= 262 // 262=0
			p.frame += 1
			p.oddframe = !p.oddframe
		}
		p.checkPendingNMI()
	}
}

func (p *PPU) checkPendingNMI() {
	// only actually send NMI after 3 PPU clocks because it's likely when the CPU would detect it
	// this gives the opportunity for the NMI to not happen if a read on PPUSTATUS happens before the NMI is sent
	if p.scanline == 241 && p.cycle < 2 {
		return
	}
	// check if we have any pending NMI, and send it
	if p.vblankNMI && p.getFlag(GenerateNMI) {
		p.vblankNMI = false
		if p.vblank != nil {
			// trigger vblank interrupt
			p.vblank()
		}
	}
}

func (p *PPU) getFlag(flag byte) bool {
	return p.ctrl&flag == flag
}

func (p *PPU) getMask(m byte) bool {
	return p.mask&m == m
}
