package cartridge

import (
	"io"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

func Load(fn string) (*Data, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	ln, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}
	f.Seek(0, io.SeekStart)

	// map it
	sc, err := f.SyscallConn()
	if err != nil {
		return nil, err
	}

	res := &Data{
		f: f,
	}

	err2 := sc.Control(func(fd uintptr) {
		res.m, err = unix.Mmap(int(fd), 0, int(ln), unix.PROT_READ, unix.MAP_SHARED|unix.MAP_POPULATE)
	})
	if err2 != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	log.Printf("Mapped %d bytes in memory", ln)

	if err = res.parse(); err != nil {
		res.Close()
		return nil, err
	}

	return res, nil
}
