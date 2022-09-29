package memory

import (
	"fmt"
	"unsafe"
)

type RAM []byte

func NewRAM(siz int) RAM {
	return RAM(make([]byte, siz))
}

func (r RAM) MemRead(offset uint16) byte {
	return r[offset%uint16(len(r)-1)]
}

func (r RAM) MemWrite(offset uint16, val byte) byte {
	r[offset%uint16(len(r)-1)] = val
	return val
}

func (r RAM) Length() uint16 {
	return uint16(len(r))
}

func (r RAM) String() string {
	return fmt.Sprintf("%d bytes RAM", len(r))
}

func (r RAM) Ptr() uintptr {
	return uintptr(unsafe.Pointer(&r[0]))
}
