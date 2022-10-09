package cartridge

import "github.com/MagicalTux/gones/pkgnes"

type MapperType byte

type Mapper interface {
	setup(nes *pkgnes.NES) error
}

var mappers = make(map[MapperType]func(*Data) Mapper)

func RegisterMapper(mt MapperType, f func(*Data) Mapper) {
	mappers[mt] = f
}
