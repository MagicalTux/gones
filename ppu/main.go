package ppu

import (
	"image"
	"sync"

	"github.com/MagicalTux/gones/memory"
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

	// variables used during rendering
	nameTableByte      byte
	attributeTableByte byte
	lowTileByte        byte
	highTileByte       byte
	tileData           uint64
	nameTableMemory    memory.RAM

	// sprites
	spriteCount      int
	spritePatterns   [8]uint32
	spritePositions  [8]byte
	spritePriorities [8]byte
	spriteIndexes    [8]byte

	front, back     *image.RGBA
	frontLk         sync.Mutex
	VBlankInterrupt func()

	sync chan *image.RGBA
}

func New() *PPU {
	ppu := &PPU{
		Memory:          memory.NewBus(),
		front:           image.NewRGBA(image.Rect(0, 0, 256, 240)),
		back:            image.NewRGBA(image.Rect(0, 0, 256, 240)),
		sync:            make(chan *image.RGBA),
		nameTableMemory: memory.NewRAM(0x800), // NEW standard 2kB PPU work ram
	}

	// https://www.nesdev.org/wiki/PPU_memory_map

	// pattern tables
	// 0x0000~0x2000 is mapped via CHR ROM to the cartridge (mapping is done by the cartridge)

	// nametables
	// 0x2000~0x3f00 is mapped typically to the PPU ram, but the cartridge can clear this mapping to map something else
	// we actually map the whole thing, but reads to values after 0x3f00 will return data from the palette
	// NOTE: see mirroring.go for ways this can be re-mapped if required
	ppu.Memory.MapHandler(0x2000, 0x2000, ppu.nameTableMemory)

	return ppu
}

func (p *PPU) Reset(cnt uint64) {
	// cnt tells us at what point we need to reset, typically 7
	pxls := cnt * 3 // =21

	p.cycle = uint16(pxls)
	p.scanline = 0
	p.frame = 0
}

func (p *PPU) Clock(cnt uint64) uint64 {
	p.checkPendingNMI()

	if cnt == 0 {
		// should not happen
		return 0
	}

	// read some status stuff
	renderEnabled := p.getMask(ShowBg) || p.getMask(ShowSprites)

	// each PPU frame is 341*262=89342 PPU clocks long

	for xrun := uint64(0); xrun < cnt; xrun += 1 {
		p.cycle += 1

		if p.cycle == 341 {
			// increase scanline
			p.cycle = 0
			p.scanline += 1

			if p.scanline == 262 {
				p.scanline = 0
				p.frame += 1
				p.oddframe = p.frame&1 == 1 // !p.oddframe
			}
		}

		if renderEnabled {
			p.triggerRender()
		}

		// generate a position uint32 identifier to easily test various positions on the frame.
		// Because we loop on all pixels all legal values are guaranteed to happen within a frame
		// this allow very fine tuning of events happening at specific pixels such as NMI, etc
		posId := uint32(p.scanline)<<16 | uint32(p.cycle)

		// See: https://www.nesdev.org/w/images/default/d/d1/Ntsc_timing.png

		switch posId {
		case 0x00f10001: // scanline=241 cycle=1
			p.vblankFlag = true
			p.stat |= VBlankStarted
			p.Flip() // perform double buffer flip
		case 0x00f10003: // scanline=241 cycle=3
			p.vblankNMI = true // generate NMI at next available occasion (slightly delayed compared to flag)
		case 0x01050001: // scanline=261 cycle=1
			// clear vblank, sprite, overflow
			p.vblankFlag = false
			p.vblankNMI = false
			// clear SpriteZeroHit & VBlankStarted from p.stat
			p.stat &= ^(SpriteZeroHit | SpriteOverflow | VBlankStarted)
		case 0x01050153: // scanline=261 cycle=340
			if renderEnabled && p.oddframe {
				// when rendering, frames after an odd frame will skip their first 0,0 pixel, act as if it was done just now
				p.cycle = 0
				p.scanline = 0
				p.frame += 1
				p.oddframe = p.frame&1 == 1 // !p.oddframe
				posId = 0
			}
		}
	}

	return cnt
}

func (p *PPU) Flip() {
	p.frontLk.Lock()
	p.front, p.back = p.back, p.front
	p.frontLk.Unlock()
}

func (p *PPU) Front(cb func(*image.RGBA)) {
	p.frontLk.Lock()
	defer p.frontLk.Unlock()

	cb(p.front)
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
		if f := p.VBlankInterrupt; f != nil {
			// trigger vblank interrupt
			f()
		}
	}
}

func (p *PPU) getFlag(flag byte) bool {
	return p.ctrl&flag == flag
}

func (p *PPU) getMask(m byte) bool {
	return p.mask&m == m
}
