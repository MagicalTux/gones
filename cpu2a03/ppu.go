package cpu2a03

import (
	"log"
	"unsafe"
)

type PPU struct {
	cpu *Cpu2A03
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