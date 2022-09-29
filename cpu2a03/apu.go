package cpu2a03

import (
	"log"
	"unsafe"
)

type APU struct {
	cpu *Cpu2A03
}

func (p *APU) MemRead(offset uint16) byte {
	log.Printf("APU read: $%04x", offset)
	return 0
}

func (p *APU) MemWrite(offset uint16, val byte) byte {
	log.Printf("APU write: $%04x = %d", offset, val)
	return 0
}

func (p *APU) Ptr() uintptr {
	return uintptr(unsafe.Pointer(p))
}

func (p *APU) String() string {
	return "APU"
}
