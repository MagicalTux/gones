package cartridge

import "github.com/MagicalTux/gones/cpu2a03"

type MapperType byte

type Mapper interface {
	setup(cpu *cpu2a03.Cpu2A03) error
}

var mappers = make(map[MapperType]func(*Data) Mapper)

func RegisterMapper(mt MapperType, f func(*Data) Mapper) {
	mappers[mt] = f
}
