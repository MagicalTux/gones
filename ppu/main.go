package ppu

import (
	"log"
	"unsafe"

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
	oamaddr byte
	scroll  byte
	addr    byte
	data    byte
	oam     [256]byte

	cycle    uint16
	scanline uint16
	oddframe bool // frame is even/odd (starting at frame 0 which is even)
	frame    uint64

	vblankFlag bool
	vblankNMI  bool

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
	ppu.Memory.MapHandler(0x2000, 0x1f00, memory.NewRAM(0x800))

	// $3F00-$3F1F	contains Palette RAM indexes - this needs to be separate
	ppu.Memory.MapHandler(0x3f00, 0x100, memory.NewRAM(32))

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
	}
	p.checkPendingNMI()
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

func (p *PPU) MemRead(offset uint16) byte {
	// only care about first 3 bits (&0x7)

	switch offset & 7 {
	case 2: // PPU status
		stat := p.stat
		p.stat &= ^VBlankStarted // always clear VBlankStarted when reading PPU STATUS
		if p.scanline == 241 && p.cycle == 0 {
			// special case, hide vblank flag and don't send the NMI
			p.vblankNMI = false
			stat &= ^VBlankStarted
		}
		if p.scanline == 241 && p.cycle <= 2 {
			// if we are within 2 PPU clocks of setting p.vblankFlag we should inhibit the NMI
			p.vblankNMI = false
		}
		return stat
	default:
		log.Printf("Unhandled PPU read: $%04x", offset)
	}
	return 0
}

func (p *PPU) MemWrite(offset uint16, val byte) byte {
	// only care about first 3 bits (&0x7)
	switch offset & 7 {
	default:
		log.Printf("Unhandled PPU write: $%04x = $%02x", offset, val)
	}
	return 0
}

func (p *PPU) Ptr() uintptr {
	return uintptr(unsafe.Pointer(p))
}

func (p *PPU) String() string {
	return "PPU"
}
