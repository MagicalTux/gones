package memory

import (
	"fmt"
	"log"
	"strings"
	"unsafe"
)

type Bus [256][]Handler

func NewBus() Master {
	return &Bus{}
}

func (b *Bus) MapHandler(offset uint16, length uint16, h Handler) {
	offt := offset >> 8
	cnt := length >> 8
	if length%0x100 != 0 {
		cnt += 1
	}

	for i := uint16(0); i < cnt; i++ {
		b[offt+i] = append(b[offt+i], h)
	}
}

func (b *Bus) ClearMapping(offset, length uint16) {
	offt := offset >> 8
	cnt := length >> 8
	if length%0x100 != 0 {
		cnt += 1
	}

	for i := uint16(0); i < cnt; i++ {
		b[offt+i] = nil
	}
}

func (b Bus) MemRead(offset uint16) byte {
	offt := offset >> 8
	var res byte

	for _, h := range b[offt] {
		res |= h.MemRead(offset)
	}

	return res
}

func (b Bus) MemWrite(offset uint16, val byte) byte {
	offt := offset >> 8
	var res byte

	for _, h := range b[offt] {
		res = h.MemWrite(offset, val)
		if res != val && res != 0 {
			// see: https://www.nesdev.org/wiki/Bus_conflict
			// We don't overwrite val with this yet, but we might if it is required
			log.Printf("Bus conflict at address $%04x! Write=$%02x but got bits $%02x", offset, val, res)
		}
	}
	return res
}

func (b Bus) Ptr() uintptr {
	return uintptr(unsafe.Pointer(&b))
}

type debugInfo struct {
	h     Handler
	start uint16
	end   uint16
}

func (b Bus) String() string {
	var m []*debugInfo

	for n, v := range b {

	loop1:
		for _, h := range v {
			// TODO handle when same object is mapped at multiple places
			for _, i := range m {
				if i.h.Ptr() == h.Ptr() && i.end == uint16(n) {
					i.end = uint16(n) + 1
					continue loop1
				}
			}
			m = append(m, &debugInfo{h: h, start: uint16(n), end: uint16(n) + 1})
		}
	}

	var r []string

	for _, i := range m {
		r = append(r, fmt.Sprintf("$%02x00~$%02xff: %s", i.start, i.end-1, i.h))
	}

	return fmt.Sprintf("Memory Bus containing:\n%s", strings.Join(r, "\n"))
}
