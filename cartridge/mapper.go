package cartridge

import "github.com/MagicalTux/gones/memory"

type MapperType byte

type Mapper interface {
	init(d *Data) error
	setup(m *memory.MMU) error
}

var mappers = make(map[MapperType]func() Mapper)

const (
	NROM MapperType = 0x00 // Nintendo cartridge boards NES-NROM-128, NES-NROM-256, their HVC counterparts, and clone boards
)
