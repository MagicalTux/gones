package apu

import (
	"log"
	"unsafe"
)

func (apu *APU) MemRead(offset uint16) byte {
	offset &= 0x1fff
	switch offset {
	case 0x15: // status
		return apu.readStatus()
	case 0x16: // read from input 0
		if dev := apu.Input[0]; dev != nil {
			return dev.Read()
		}
		return 0
	case 0x17: // read from input 1
		if dev := apu.Input[1]; dev != nil {
			return dev.Read()
		}
		return 0
	default:
		if offset < 0x15 {
			// https://github.com/quackenbush/nestalgia/blob/master/docs/apu/nessound.txt
			// Note that $4015 is the only R/W register. All others are write only (attempt
			// to read them will most likely result in a returned 040H, due to heavy
			// capacitance on the NES's data bus). Reading a "write only" register, will
			// have no effect on the specific register, or channel.
			return 0x40
		}
		log.Printf("Unhandled APU read: $%04x", offset)
	}
	return 0
}

func (apu *APU) MemWrite(offset uint16, val byte) byte {
	offset &= 0x1fff

	if offset != 0x16 {
		// 16= controllers. Ignore it
		//log.Printf("APU: WRITE $%04x = $%02x", offset, val)
	}
	switch offset {
	case 0x00, 0x01, 0x02, 0x03: // pulse1
		return apu.pulse1.MemWrite(offset, val)
	case 0x04, 0x05, 0x06, 0x07: // pulse2
		return apu.pulse2.MemWrite(offset, val)
	case 0x10, 0x11, 0x12, 0x13: // DMC
		return apu.dmc.MemWrite(offset, val)

		// triangle
	case 0x08:
		apu.triangle.writeControl(val)
	case 0x0A:
		apu.triangle.writeTimerLow(val)
	case 0x0B:
		apu.triangle.writeTimerHigh(val)

		// noise
	case 0x0C:
		apu.noise.writeControl(val)
	case 0x0E:
		apu.noise.writePeriod(val)
	case 0x0F:
		apu.noise.writeLength(val)

		// control
	case 0x15:
		apu.writeControl(val)
	case 0x17:
		apu.writeFrameCounter(val)

	case 0x14: // OAM DMA
		// when writing to this port, send data of memory at HH=val to PPU's OAMDATA port
		// https://www.nesdev.org/wiki/PPU_registers#OAM_DMA_($4014)_%3E_write
		addr := uint16(val) << 8
		for i := uint16(0); i < 256; i++ {
			val = apu.Memory.MemRead(addr | i)
			apu.Memory.MemWrite(0x2004, val)
			apu.cpuDelay(2)
		}
		if apu.cpuDelay(0)&1 == 1 {
			// odd cpu cycle
			apu.cpuDelay(2)
		} else {
			apu.cpuDelay(1)
		}
	case 0x16: // controller polling mode
		val = val & 7
		// this sets OUT0 OUT1 and OUT2 depending on the bits values
		if dev := apu.Input[0]; dev != nil {
			dev.Write(val)
		}
		if dev := apu.Input[1]; dev != nil {
			dev.Write(val)
		}
	default:
		log.Printf("Unhandled APU write: $%04x = $%02x", offset, val)
	}
	return 0
}

func (p *APU) Ptr() uintptr {
	return uintptr(unsafe.Pointer(p))
}

func (p *APU) String() string {
	return "APU"
}

func (p *APU) Length() uint16 {
	return 0x2000
}
