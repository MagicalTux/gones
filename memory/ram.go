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
	return r[offset&uint16(len(r)-1)]
}

func (r RAM) MemWrite(offset uint16, val byte) byte {
	r[offset&uint16(len(r)-1)] = val
	return val
}

// Resize will return a copy of RAM with the specified size. If RAM is large
// enough to accomodate the new size (less or equal to capacity), then the
// exact same memory buffer will be returned. If however it is required to grow
// the memory buffer, a new buffer will be allocated, and data will be copied
//
// Reducing the size of a memory will not free memory, and re-growing it will
// cause data that was written there to be still there.
func (r RAM) Resize(newsize int) RAM {
	if cap(r) >= newsize {
		return r[:newsize]
	}
	newr := make([]byte, newsize)
	copy(newr, r[:cap(r)])
	return RAM(newr)
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
