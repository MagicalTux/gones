package memory

import (
	"fmt"
	"unsafe"
)

// Slice can slice a RAM or ROM device
func Slice(h Handler, start, end int) Handler {
	switch v := h.(type) {
	case RAM:
		if start > cap(v) || end > cap(v) {
			return Null{}
		}
		return v[start:end]
	case ROM:
		if start > cap(v) || end > cap(v) {
			return Null{}
		}
		return v[start:end]
	default:
		return &memSlice{h, uint16(start), uint16(end - start - 1)}
	}
}

type memSlice struct {
	dev    Handler
	offset uint16
	mask   uint16
}

func (ms *memSlice) MemRead(offset uint16) byte {
	return ms.dev.MemRead((ms.offset & ms.mask) + ms.offset)
}

func (ms *memSlice) MemWrite(offset uint16, v byte) byte {
	return ms.dev.MemWrite((ms.offset&ms.mask)+ms.offset, v)
}

func (ms *memSlice) Length() uint16 {
	return ms.mask + 1
}

func (ms *memSlice) Ptr() uintptr {
	return uintptr(unsafe.Pointer(ms))
}

func (ms *memSlice) String() string {
	return fmt.Sprintf("Slice(offset=$%04x, mask=$%04x) of %s", ms.offset, ms.mask, ms.dev)
}
