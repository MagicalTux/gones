package cartridge

import (
	"os"

	"github.com/MagicalTux/gones/memory"
	"github.com/MagicalTux/gones/nesppu"
	"github.com/MagicalTux/gones/pkgnes"
)

type Data struct {
	f      *os.File
	m      []byte // map+len
	Mapper Mapper

	numPRG          byte // Size of PRG ROM in 16 KB units
	numCHR          byte
	numPRGram       int
	mapperType      MapperType
	hasTrainer      bool
	hasMirroring    bool
	ignoreMirroring bool
}

func (d *Data) Close() error {
	if d.m != nil {
		d.unload()
	}
	if f := d.f; f != nil {
		d.f = nil
		return f.Close()
	}

	return nil
}

func (d *Data) PRG() []byte {
	// get PRG data
	// see: https://www.nesdev.org/wiki/INES#iNES_file_format
	offt := 16
	if d.hasTrainer {
		offt += 512
	}

	romSize := int(d.numPRG) << 14 // 1=16kB 2=32kB etc...

	return d.m[offt : offt+romSize]
}

func (d *Data) CHR() memory.Handler {
	if d.numCHR == 0 {
		// 5: Size of CHR ROM in 8 KB units (Value 0 means the board uses CHR RAM)
		return memory.NewRAM(0x2000)
	}

	// get CHR data
	offt := 16
	if d.hasTrainer {
		offt += 512
	}

	romSize := int(d.numPRG) << 14 // 1=16kB 2=32kB etc...
	offt += romSize

	romSize = int(d.numCHR) << 13 // 1=8kB, ...

	return memory.ROM(d.m[offt : offt+romSize])
}

func (d *Data) Setup(nes *pkgnes.NES) error {
	err := d.Mapper.setup(nes)
	if err != nil {
		return err
	}
	// see https://www.nesdev.org/wiki/Mirroring#Nametable_Mirroring
	if d.ignoreMirroring {
		// Ignore mirroring control or above mirroring bit; instead provide four-screen VRAM
		nes.PPU.SetMirroring(nesppu.FourScreenMirroring)
	} else if d.hasMirroring {
		// 1: vertical (horizontal arrangement) (CIRAM A10 = PPU A10)
		nes.PPU.SetMirroring(nesppu.VerticalMirroring)
	} else {
		// 0: horizontal (vertical arrangement) (CIRAM A10 = PPU A11)
		nes.PPU.SetMirroring(nesppu.HorizontalMirroring)
	}
	return nil
}
