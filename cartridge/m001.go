package cartridge

import (
	"log"
	"os"
	"unsafe"

	"github.com/MagicalTux/gones/cpu2a03"
	"github.com/MagicalTux/gones/memory"
	"github.com/MagicalTux/gones/ppu"
)

func init() {
	RegisterMapper(001, func(data *Data) Mapper {
		return &MMC1{
			data:    data,
			in:      0x10,
			prgMode: 3,
		}
	})
}

// https://www.nesdev.org/wiki/MMC1
type MMC1 struct {
	data *Data
	ppu  *ppu.PPU

	in byte

	prg memory.ROM
	chr memory.Handler

	prgRAM  memory.Handler
	prgMode byte
	chrMode byte

	chrBank0 memory.Handler
	chrBank1 memory.Handler
	prgBank0 memory.Handler
	prgBank1 memory.Handler

	chrBank0sel byte
	chrBank1sel byte
	prgBankSel  byte
}

type debugWrite struct {
	memory.RAM
}

func (d *debugWrite) MemWrite(offset uint16, val byte) byte {
	if val != 0 {
		os.Stdout.Write([]byte{val})
	}
	//log.Printf("MMC1: debug write $%04x = $%02x", offset, val)
	return d.RAM.MemWrite(offset, val)
}

func (m *MMC1) setup(cpu *cpu2a03.Cpu2A03) error {
	cpu.Memory.MapHandler(0x6000, 0x2000, m)
	cpu.Memory.MapHandler(0x8000, 0x8000, m)
	cpu.PPU.Memory.MapHandler(0x0000, 0x2000, m)

	m.ppu = cpu.PPU

	// CPU $6000-$7FFF: Family Basic only: PRG RAM, mirrored as necessary to fill entire 8 KiB window, write protectable with an external switch
	// we ignore numPRGram since value 0 means a 8kB RAM, value 1 means a 8kB ram, and higher values can't be addressed
	m.prgRAM = memory.NewRAM(0x2000)
	//m.prgRAM = &debugWrite{memory.NewRAM(0x2000)}

	// 2022/10/02 16:02:59 Parsed iNes1 file, 8*16kB PRG, 0*8kB CHR, 1*8kB PRG RAM, mapper=1, trainer=false mirroring=true/false

	m.prg = memory.ROM(m.data.PRG())
	m.chr = m.data.CHR()

	m.updateBanks()

	return nil
}

func (m *MMC1) MemRead(offset uint16) byte {
	bank := offset >> 12

	switch bank {
	case 0, 1:
		// CHR bank 1
		if m.chrBank0 != nil {
			return m.chrBank0.MemRead(offset)
		}
		return 0
	case 2, 3:
		// CHR bank 2
		if m.chrBank1 != nil {
			return m.chrBank1.MemRead(offset)
		}
		return 0
	case 6, 7:
		// PRG RAM bank
		if m.prgRAM == nil {
			return 0
		}
		return m.prgRAM.MemRead(offset)
	case 8, 9, 0xa, 0xb:
		// 16 KB PRG ROM bank, either switchable or fixed to the first bank
		if m.prgBank0 != nil {
			return m.prgBank0.MemRead(offset)
		}
		return 0
	case 0xc, 0xd, 0xe, 0xf:
		// 16 KB PRG ROM bank, either fixed to the last bank or switchable
		if m.prgBank1 != nil {
			return m.prgBank1.MemRead(offset)
		}
		return 0
	default:
		return 0
	}
}

func (m *MMC1) MemWrite(offset uint16, v byte) byte {
	bank := offset >> 12

	switch bank {
	case 0, 1:
		// CHR bank 1
		if m.chrBank0 != nil {
			return m.chrBank0.MemWrite(offset, v)
		}
		return 0
	case 2, 3:
		// CHR bank 2
		if m.chrBank1 != nil {
			return m.chrBank1.MemWrite(offset, v)
		}
		return 0
	case 4, 5:
		// NULL
		return 0
	case 6, 7:
		// PRG RAM bank
		if m.prgRAM == nil {
			return 0
		}
		return m.prgRAM.MemWrite(offset, v)
	}

	// we're writing to offset after 0x8000

	if v&0x80 == 0x80 {
		// clear
		m.in = 0x10
		return 0
	}

	v = (v & 1) << 4
	full := m.in&1 == 1
	m.in = (m.in >> 1) | v

	if !full {
		return 0
	}

	// we've read 5 bits
	v = m.in
	m.in = 0x10

	//log.Printf("MMC1: Set bank=$%02x v=%02x", bank, v)

	switch bank {
	case 0x8, 0x9:
		// Control (internal, $8000-$9FFF)
		mirrorMode := v & 3      // 0: one-screen, lower bank; 1: one-screen, upper bank; 2: vertical; 3: horizontal
		m.prgMode = (v >> 2) & 3 // 0, 1: switch 32 KB at $8000, ignoring low bit of bank number; 2: fix first bank at $8000 and switch 16 KB bank at $C000; 3: fix last bank at $C000 and switch 16 KB bank at $8000)
		m.chrMode = (v >> 4) & 1 // 0: switch 8 KB at a time; 1: switch two separate 4 KB banks

		log.Printf("MMC1: Set mirror mode=%d prgMode=%d chrMode=%d", mirrorMode, m.prgMode, m.chrMode)
		switch mirrorMode {
		case 0:
			m.ppu.SetMirroring(ppu.SingleScreenMirroring)
		case 1:
			m.ppu.SetMirroring(ppu.SingleScreen2Mirroring)
		case 2:
			m.ppu.SetMirroring(ppu.VerticalMirroring)
		case 3:
			m.ppu.SetMirroring(ppu.HorizontalMirroring)
		}
		m.updateBanks()
		return 0
	case 0xa, 0xb:
		// CHR bank 0 (internal, $A000-$BFFF)
		m.chrBank0sel = v
		log.Printf("MMC1: CHR Bank#0 set to $%02x", v)
		m.updateBanks()
		return 0
	case 0xc, 0xd:
		// CHR bank 1 (internal, $C000-$DFFF)
		m.chrBank1sel = v
		log.Printf("MMC1: CHR Bank#1 set to $%02x", v)
		m.updateBanks()
		return 0
	case 0xe, 0xf:
		// PRG bank (internal, $E000-$FFFF)
		if v != m.prgBankSel {
			m.prgBankSel = v
			log.Printf("MMC1: PRG bank set to $%02x", v)
			m.updateBanks()
		}
		return 0
	}
	return 0
}

func (m *MMC1) updateBanks() {
	switch m.prgMode {
	case 0, 1:
		// 0, 1: switch 32 KB at $8000, ignoring low bit of bank number
		n := int(m.prgBankSel) & 0x1e
		slice := memory.Slice(m.prg, 0x4000*n, 0x4000*n+0x8000)
		m.prgBank0 = slice
		m.prgBank1 = slice
	case 2:
		// 2: fix first bank at $8000 and switch 16 KB bank at $C000
		n := int(m.prgBankSel)
		m.prgBank0 = memory.Slice(m.prg, 0, 0x4000)
		m.prgBank1 = memory.Slice(m.prg, 0x4000*n, 0x4000*n+0x4000)
	case 3:
		// 3: fix last bank at $C000 and switch 16 KB bank at $8000
		n := int(m.prgBankSel)
		size := len(m.prg)
		m.prgBank0 = memory.Slice(m.prg, 0x4000*n, 0x4000*n+0x4000)
		m.prgBank1 = memory.Slice(m.prg, size-0x4000, size)
	}

	switch m.chrMode {
	case 0:
		// 0: switch 8 KB at a time
		n := int(m.chrBank0sel) & 0x1e
		slice := memory.Slice(m.chr, 0x1000*n, 0x1000*n+0x2000)
		m.chrBank0 = slice
		m.chrBank1 = slice
	case 1:
		// 1: switch two separate 4 KB banks
		n := int(m.chrBank0sel)
		m.chrBank0 = memory.Slice(m.chr, 0x1000*n, 0x1000*n+0x1000)
		n = int(m.chrBank1sel)
		m.chrBank1 = memory.Slice(m.chr, 0x1000*n, 0x1000*n+0x1000)
	}
}

func (m *MMC1) Ptr() uintptr {
	return uintptr(unsafe.Pointer(m))
}

func (m *MMC1) String() string {
	return "MMC1 Mapper"
}

func (m *MMC1) Length() uint16 {
	return 0
}
