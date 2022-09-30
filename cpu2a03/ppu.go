package cpu2a03

import (
	"log"
	"unsafe"
)

// https://www.nesdev.org/wiki/PPU_rendering#Frame_timing_diagram
// The PPU renders 262 scanlines per frame. Each scanline lasts for 341 PPU clock cycles (113.667 CPU clock cycles; 1 CPU cycle = 3 PPU cycles), with each clock cycle producing one pixel.

type PPU struct {
	cpu *Cpu2A03

	ctrl    byte // 0x00 at start
	mask    byte // 0x00 at start
	stat    byte // 0x00 at start
	oamaddr byte
	scroll  byte
	addr    byte
	data    byte
}

func (p *PPU) Clock(cnt int) {
	// move clock forward by cnt (1 CPU clock = 3 PPU clock)
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
