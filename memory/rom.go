package memory

import (
	"fmt"
	"unsafe"
)

type ROM []byte

func (r ROM) MemRead(offset uint16) byte {
	return r[offset%uint16(len(r)-1)]
}

func (r ROM) MemWrite(offset uint16, val byte) byte {
	return r.MemRead(offset)
}

func (r ROM) Length() uint16 {
	return uint16(len(r))
}

func (r ROM) String() string {
	return fmt.Sprintf("%d bytes ROM", len(r))
}

func (r ROM) Ptr() uintptr {
	return uintptr(unsafe.Pointer(&r[0]))
}
