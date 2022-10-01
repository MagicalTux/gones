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

	x, y uint16

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
	// 0x0000~0x2000 is mapped via CHR ROM to the cartridge, typically
	// nametables
	// 0x2000~0x3f00 is mapped typically to the ram
	//ppu.Memory.MapHandler(0x2000, 0x1f00, memory.NewRAM(0x1000))

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

	p.x = 0
	p.y = uint16(pxls)
}

func (p *PPU) Clock(cnt int) {
	// move clock forward by cnt (1 CPU clock = 3 PPU clock)
	p.y += uint16(cnt * 3)
}

func (p *PPU) MemRead(offset uint16) byte {
	// only care about first 3 bits (&0x7)
	log.Printf("PPU read: $%04x", offset)

	switch offset & 7 {
	case 2: // PPU status
		return 0x80 // means "vblank", games will check if this bit is set or not and wait
	}
	return 0
}

func (p *PPU) MemWrite(offset uint16, val byte) byte {
	// only care about first 3 bits (&0x7)
	log.Printf("PPU write: $%04x = %d", offset, val)
	switch offset & 7 {
	case 0: // PPU control
		// TODO
	}
	return 0
}

func (p *PPU) Ptr() uintptr {
	return uintptr(unsafe.Pointer(p))
}

func (p *PPU) String() string {
	return "PPU"
}
