package nesppu

import (
	"log"
	"unsafe"
)

func (p *PPU) MemRead(offset uint16) byte {
	// only care about first 3 bits (&0x7)
	if offset != 0x2002 {
		p.trace("Memory read $%04x", offset)
	}

	switch offset & 7 {
	case PPUSTATUS: // PPU status
		stat := p.stat
		if p.getStatus(VBlankStarted) {
			p.trace("Crear VBlankStarted because PPUSTATUS read")
		}
		p.stat &= ^VBlankStarted // always clear VBlankStarted when reading PPU STATUS
		p.W = false              // reading PPUSTATUS resets the PPUADDR latch
		if p.scanline == 241 && p.cycle == 1 {
			// special case, hide vblank flag and don't send the NMI
			p.vblankNMI = false
			p.vblankDoNMI = false
			stat &= ^VBlankStarted
		}
		if p.scanline == 241 && p.cycle <= 3 {
			// if we are within 2 PPU clocks of setting p.vblankFlag we should inhibit the NMI
			p.vblankNMI = false
			p.vblankDoNMI = false
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
		if p.V >= 0x3f00 {
			// return palette data instead
			res = p.Palette[palAddr(p.V)]
		}
		p.trace("PPUDATA read: $%04x = $%02x", p.V, res)
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
		p.trace("Write PPUCTRL → $%02x", val)
		if !p.getFlag(GenerateNMI) && val&GenerateNMI == GenerateNMI {
			if p.getStatus(VBlankStarted) {
				p.trace("PPUCTRL: enabled NMI with VBL set, vblDoNmi=%v", p.vblankDoNMI)
				// enabling VBlankStarted will reset p.vblankNMI to p.vblankFlag and may trigger an extra NMI
				p.vblankNMI = p.vblankDoNMI
			} else {
				p.trace("PPUCTRL: enabled NMI, no VBL flag")
				// we don't want to do NMI since it was enabled late
				p.vblankDoNMI = false
			}
		}
		if p.getFlag(GenerateNMI) && val&GenerateNMI == 0 {
			p.trace("PPUCTRL: disabled NMI")
		}
		p.ctrl = val
		// also affects T
		// t: ....BA.. ........ = d: ......BA
		p.T = (p.T & 0xF3FF) | ((uint16(val) & 0x03) << 10)
	case PPUMASK:
		p.trace("Write PPUMASK → $%02x", val)
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
		p.trace("PPUDATA write: $%04x = $%02x", p.V, val)
		if p.V >= 0x3f00 {
			// write to palette
			p.Palette[palAddr(p.V)] = val
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

// palAddr returns the offset within the palette (0~31) for a given memory access address
func palAddr(v uint16) uint16 {
	v %= 0x20

	// https://www.nesdev.org/wiki/PPU_palettes says:
	// Addresses $3F10/$3F14/$3F18/$3F1C are mirrors of $3F00/$3F04/$3F08/$3F0C.
	switch v {
	case 0x10, 0x14, 0x18, 0x1c:
		v &= 0x0c
	}
	return v
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
