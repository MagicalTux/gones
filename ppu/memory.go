package ppu

import (
	"log"
	"unsafe"
)

func (p *PPU) MemRead(offset uint16) byte {
	// only care about first 3 bits (&0x7)

	switch offset & 7 {
	case PPUSTATUS: // PPU status
		if p.scanline == 241 && p.cycle == 1 {
			// Race Condition Warning: Reading PPUSTATUS within two cycles of the start of vertical blank will return 0 in bit 7 but clear the latch anyway, causing NMI to not occur that frame.
			p.stat &= ^VBlankStarted // clear it now
			p.vblankDoNMI = false
		}
		stat := p.stat
		p.stat &= ^VBlankStarted // always clear VBlankStarted when reading PPU STATUS
		p.W = false              // reading PPUSTATUS resets the PPUADDR latch
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
	case OAMDATA:
		// read OAM data
		if p.oamAddr&0x03 == 0x02 {
			// see: https://www.nesdev.org/wiki/PPU_OAM#Byte_2
			// bits 2, 3, 4 of byte 2 always return zero
			return p.OAM[p.oamAddr] & 0xe3
		}
		return p.OAM[p.oamAddr]
	case PPUDATA:
		// read from memory at address p.ppuAddr
		// See: https://www.nesdev.org/wiki/PPU_registers#The_PPUDATA_read_buffer_(post-fetch)
		res := p.readBuf
		p.readBuf = p.Memory.MemRead(p.V & 0x3fff)
		if p.V > 0x3f00 {
			// return palette data instead
			res = p.Palette[p.V%0x20]
		}
		// increment p.V
		if !p.getFlag(LargeIncrements) {
			p.V += 1
		} else {
			p.V += 32
		}
		return res
	default:
		log.Printf("Unhandled PPU read: $%04x", offset)
	}
	return 0
}

func (p *PPU) MemWrite(offset uint16, val byte) byte {
	// only care about first 3 bits (&0x7)
	switch offset & 7 {
	case PPUCTRL:
		if !p.getFlag(VBlankStarted) && val&VBlankStarted == VBlankStarted {
			// enabling VBlankStarted will reset p.vblankNMI to p.vblankFlag and may trigger an extra NMI
			p.vblankNMI = p.vblankFlag
		}
		p.ctrl = val
		// also affects T
		// t: ....BA.. ........ = d: ......BA
		p.T = (p.T & 0xF3FF) | ((uint16(val) & 0x03) << 10)
	case PPUMASK:
		p.mask = val
	case OAMADDR:
		p.oamAddr = val
		return 0
	case OAMDATA:
		//log.Printf("got OAM data $%02x = $%02x", p.oamAddr, val)
		// write value & increment addr
		p.OAM[p.oamAddr] = val
		p.oamAddr += 1
		return 0
	case PPUSCROLL:
		// PPUSCROLL and PPUADDR share registers, see https://www.nesdev.org/wiki/PPU_scrolling#Register_controls
		if !p.W {
			p.T = (p.T & 0xffe0) | (uint16(val) >> 3)
			p.X = val & 0x07
			p.W = true
		} else {
			p.T = (p.T & 0x8fff) | ((uint16(val) & 0x07) << 12)
			p.T = (p.T & 0xfc1f) | ((uint16(val) & 0xf8) << 2)
			p.W = false
		}
		return 0
	case PPUADDR:
		if !p.W {
			p.T = (p.T & 0x80ff) | ((uint16(val) & 0x3f) << 8)
			p.W = true
		} else {
			p.T = (p.T & 0xFF00) | uint16(val)
			p.V = p.T
			p.W = false
		}
		return 0
	case PPUDATA:
		if p.V >= 0x3f00 {
			// write to palette
			p.Palette[p.V%0x20] = val
		} else {
			p.Memory.MemWrite(p.V&0x3fff, val)
		}
		// increment p.V
		if !p.getFlag(LargeIncrements) {
			p.V += 1
		} else {
			p.V += 32
		}
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

func (p *PPU) Length() uint16 {
	return 0x2000
}
