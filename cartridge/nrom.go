package cartridge

import "github.com/MagicalTux/gones/memory"

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

func (m *MapperNROM) Setup(mem memory.Master) error {
	// CPU $6000-$7FFF: Family Basic only: PRG RAM, mirrored as necessary to fill entire 8 KiB window, write protectable with an external switch
	mem.MapHandler(0x6000, 0x2000, memory.NewRAM(0x2000))

	// CPU $8000-$BFFF: First 16 KB of ROM.
	// CPU $C000-$FFFF: Last 16 KB of ROM (NROM-256) or mirror of $8000-$BFFF (NROM-128).
	mem.MapHandler(0x8000, 0x8000, memory.ROM(m.data.PRG()))

	return nil
}
