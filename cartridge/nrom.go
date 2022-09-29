package cartridge

import "github.com/MagicalTux/gones/memory"

func init() {
	mappers[NROM] = func() Mapper {
		return &MapperNROM{}
	}
}

// https://www.nesdev.org/wiki/NROM
type MapperNROM struct {
	data *Data
}

func (m *MapperNROM) init(data *Data) error {
	m.data = data
	return nil
}

func (m *MapperNROM) setup(mem *memory.MMU) error {
	// CPU $6000-$7FFF: Family Basic only: PRG RAM, mirrored as necessary to fill entire 8 KiB window, write protectable with an external switch
	mem.MapAnonymous(0x6000, 0x2000)
	// CPU $8000-$BFFF: First 16 KB of ROM.
	// CPU $C000-$FFFF: Last 16 KB of ROM (NROM-256) or mirror of $8000-$BFFF (NROM-128).
	return nil
}
