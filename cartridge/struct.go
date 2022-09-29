package cartridge

import (
	"os"

	"golang.org/x/sys/unix"
)

type Data struct {
	f      *os.File
	m      []byte // map+len
	mapper Mapper

	numPRG     byte
	numCHR     byte
	numPRGram  int
	mapperType MapperType
	hasTrainer bool
}

func (d *Data) Close() error {
	if d.m != nil {
		// need unmap
		unix.Munmap(d.m)
		d.m = nil
	}

	return d.f.Close()
}
