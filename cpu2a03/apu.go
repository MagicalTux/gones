package cpu2a03

import (
	"log"
	"unsafe"
)

type APU struct {
	cpu *Cpu2A03
}

func (p *APU) MemRead(offset uint16) byte {
	offset &= 0x1fff
	switch offset {
	case 0x16: // read from input 0
		if dev := p.cpu.Input[0]; dev != nil {
			return dev.Read()
		}
		return 0
	case 0x17: // read from input 1
		if dev := p.cpu.Input[1]; dev != nil {
			return dev.Read()
		}
		return 0
	default:
		log.Printf("Unhandled APU read: $%04x", offset)
	}
	return 0
}

func (p *APU) MemWrite(offset uint16, val byte) byte {
	offset &= 0x1fff
	switch offset {
	case 0x16: // controller polling mode
		val = val & 7
		// this sets OUT0 OUT1 and OUT2 depending on the bits values
		if dev := p.cpu.Input[0]; dev != nil {
			dev.Write(val)
		}
		if dev := p.cpu.Input[1]; dev != nil {
			dev.Write(val)
		}
	default:
		log.Printf("Unhandled APU write: $%04x = %d", offset, val)
	}
	return 0
}

func (p *APU) Ptr() uintptr {
	return uintptr(unsafe.Pointer(p))
}

func (p *APU) String() string {
	return "APU"
}
