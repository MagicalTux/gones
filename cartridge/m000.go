package cartridge

import (
	"github.com/MagicalTux/gones/cpu2a03"
	"github.com/MagicalTux/gones/memory"
)

const (
	NROM MapperType = 0x00 // Nintendo cartridge boards NES-NROM-128, NES-NROM-256, their HVC counterparts, and clone boards
)

func init() {
	RegisterMapper(NROM, func(data *Data) Mapper {
		return &MapperNROM{data: data}
	})
}

// https://www.nesdev.org/wiki/NROM
type MapperNROM struct {
	data *Data
}

func (m *MapperNROM) setup(cpu *cpu2a03.Cpu2A03) error {
	// CPU $6000-$7FFF: Family Basic only: PRG RAM, mirrored as necessary to fill entire 8 KiB window, write protectable with an external switch
	// we ignore numPRGram since value 0 means a 8kB RAM, value 1 means a 8kB ram, and higher values can't be addressed
	cpu.Memory.MapHandler(0x6000, 0x2000, memory.NewRAM(0x2000))

	// CPU $8000-$BFFF: First 16 KB of ROM.
	// CPU $C000-$FFFF: Last 16 KB of ROM (NROM-256) or mirror of $8000-$BFFF (NROM-128).
	rom := memory.ROM(m.data.PRG())
	cpu.Memory.MapHandler(0x8000, 0x8000, rom)

	if chr := m.data.CHR(); chr != nil {
		cpu.PPU.Memory.MapHandler(0x0000, 0x2000, chr)
	}

	return nil
}
