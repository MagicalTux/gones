package cartridge

import (
	"bytes"
	"fmt"
	"log"
)

const (
	iNesHeader = "NES\x1a"

	// flg6
	flgMirroring  = 1 // vertical (horizontal arrangement) (CIRAM A10 = PPU A10)
	flgBatteryRAM = 2 // Cartridge contains battery-backed PRG RAM ($6000-7FFF) or other persistent memory
	flgTrainer    = 4 // 512-byte trainer at $7000-$71FF (stored before PRG data)
	flgIgnoreMirr = 8 // Ignore mirroring control or above mirroring bit; instead provide four-screen VRAM

	// flg7
	flgUnisys     = 1 // VS Unisystem
	flgPlayChoice = 2 // PlayChoice-10 (8KB of Hint Screen data stored after CHR data)

	// flg9
	flgPAL = 1 // TV system (0: NTSC; 1: PAL)
)

// https://www.nesdev.org/wiki/INES
func (d *Data) parse() error {
	// parse data
	if len(d.m) < 16 {
		return fmt.Errorf("file is too small")
	}
	if !bytes.Equal(d.m[:4], []byte(iNesHeader)) {
		return fmt.Errorf("bad file header")
	}

	d.numPRG = d.m[4] // Size of PRG ROM in 16 KB units
	d.numCHR = d.m[5] // Size of CHR ROM in 8 KB units (Value 0 means the board uses CHR RAM)
	flg6 := d.m[6]
	flg7 := d.m[7]
	d.numPRGram = int(d.m[8]) + 1 // Size of PRG RAM in 8 KB units (Value 0 infers 8 KB for compatibility)
	//flg9 := d.m[9]
	//flg10 := d.m[10]
	// 11-15: Unused padding (should be filled with zero, but some rippers put their name across bytes 7-15)

	iNes2Flag := (flg7 >> 2) & 3
	if iNes2Flag == 2 {
		log.Printf("Image is in ines2 format, not supported yet")
		panic("ines2 format not supported")
	}

	d.mapperType = MapperType(flg6>>4 | flg7&0xf0)

	d.hasTrainer = flg6&flgTrainer == flgTrainer

	log.Printf("Parsed iNes1 file, %d*16kB PRG, %d*8kB CHR, %d*8kB PRG RAM, mapper=%d, trainer=%v", d.numPRG, d.numCHR, d.numPRGram, d.mapperType, d.hasTrainer)

	if f, ok := mappers[d.mapperType]; ok {
		d.Mapper = f(d)
	} else {
		return fmt.Errorf("unsupported mapper %d", d.mapperType)
	}

	return nil
}
