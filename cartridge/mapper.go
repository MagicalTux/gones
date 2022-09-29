package cartridge

import "github.com/MagicalTux/gones/memory"

type MapperType byte

type Mapper interface {
	Setup(m memory.Master) error
}

var mappers = make(map[MapperType]func(*Data) Mapper)

func RegisterMapper(mt MapperType, f func(*Data) Mapper) {
	mappers[mt] = f
}
