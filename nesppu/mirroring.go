package nesppu

import (
	"fmt"
	"unsafe"

	"github.com/MagicalTux/gones/memory"
)

type MirroringOption byte

// See: https://www.nesdev.org/wiki/Mirroring

const (
	InvalidMirroring MirroringOption = iota
	HorizontalMirroring
	VerticalMirroring
	SingleScreenMirroring
	SingleScreen2Mirroring
	FourScreenMirroring
	DiagonalMirroring
	LShapedMirroring
	ThreeScreenMirroring
	ThreeScreenHorizontalMirroring
	ThreeScreenDiagonalMirroring
)

func (ppu *PPU) SetMirroring(mopt MirroringOption) {
	ppu.Memory.ClearMapping(0x2000, 0x2000)

	switch mopt {
	case InvalidMirroring:
		// keep memory not mapped, it'll cause all writes to be lost (as if the RAM chip was absent)
		return
	case HorizontalMirroring:
		ppu.nameTableMemory = ppu.nameTableMemory.Resize(0x800) // set memory size to 0x800
		ppu.Memory.MapHandler(0x2000, 0x2000, &ppuMirrorRouter{keys: [4]byte{0, 0, 1, 1}, parent: ppu.nameTableMemory})
	case VerticalMirroring:
		ppu.nameTableMemory = ppu.nameTableMemory.Resize(0x800) // set memory size to 0x800
		ppu.Memory.MapHandler(0x2000, 0x2000, ppu.nameTableMemory)
	case SingleScreenMirroring:
		ppu.nameTableMemory = ppu.nameTableMemory.Resize(0x800)
		ppu.Memory.MapHandler(0x2000, 0x2000, ppu.nameTableMemory[:0x400])
	case SingleScreen2Mirroring:
		ppu.nameTableMemory = ppu.nameTableMemory.Resize(0x800)
		ppu.Memory.MapHandler(0x2000, 0x2000, ppu.nameTableMemory[0x400:])
	case FourScreenMirroring:
		ppu.nameTableMemory = ppu.nameTableMemory.Resize(0x1000)
		ppu.Memory.MapHandler(0x2000, 0x2000, ppu.nameTableMemory)
	case DiagonalMirroring:
		ppu.nameTableMemory = ppu.nameTableMemory.Resize(0x800)
		ppu.Memory.MapHandler(0x2000, 0x2000, &ppuMirrorRouter{keys: [4]byte{0, 1, 1, 0}, parent: ppu.nameTableMemory})
	case LShapedMirroring:
		ppu.nameTableMemory = ppu.nameTableMemory.Resize(0x800)
		ppu.Memory.MapHandler(0x2000, 0x2000, &ppuMirrorRouter{keys: [4]byte{0, 1, 1, 1}, parent: ppu.nameTableMemory})
	case ThreeScreenMirroring:
		ppu.nameTableMemory = ppu.nameTableMemory.Resize(0x1000)
		ppu.Memory.MapHandler(0x2000, 0x2000, &ppuMirrorRouter{keys: [4]byte{0, 2, 1, 2}, parent: ppu.nameTableMemory})
	case ThreeScreenHorizontalMirroring:
		ppu.nameTableMemory = ppu.nameTableMemory.Resize(0x1000)
		ppu.Memory.MapHandler(0x2000, 0x2000, &ppuMirrorRouter{keys: [4]byte{0, 1, 2, 2}, parent: ppu.nameTableMemory})
	case ThreeScreenDiagonalMirroring:
		ppu.nameTableMemory = ppu.nameTableMemory.Resize(0x1000)
		ppu.Memory.MapHandler(0x2000, 0x2000, &ppuMirrorRouter{keys: [4]byte{0, 1, 1, 2}, parent: ppu.nameTableMemory})
	}
}

// SetCustomNametables configures nametables memory address (0x2000~0x2fff) to point to the
// given memory device. keys can be used to configure mirroring, allowing bits 10 and 11 of
// the address to be altered as needed. giving keys=[4]byte{0,1,2,3} means no rewriting
// happens, while other values will give various behaviors.
// If device is null, the PPU's internal WRAM will be used after being resized depending on
// the largest key provided (for example 0,1,2,3 will allocate 4kB of WRAM).
func (ppu *PPU) SetCustomNametables(device memory.Handler, keys [4]byte) {
	ppu.Memory.ClearMapping(0x2000, 0x2000)

	if device == nil {
		maxKey := byte(0)
		for _, v := range keys {
			if v > maxKey {
				maxKey = v
			}
		}
		newSize := 0x400
		switch maxKey {
		case 0:
			newSize = 0x400
		case 1:
			newSize = 0x800
		case 2, 3:
			newSize = 0x1000
		}
		ppu.nameTableMemory = ppu.nameTableMemory.Resize(newSize)
		device = ppu.nameTableMemory
	}

	ppu.Memory.MapHandler(0x2000, 0x2000, &ppuMirrorRouter{keys: keys, parent: device})
}

type ppuMirrorRouter struct {
	keys   [4]byte
	parent memory.Handler
}

func (m *ppuMirrorRouter) MemRead(offset uint16) byte {
	opt := (offset >> 10) & 3                 // selected screen
	offset &= 0x3ff | uint16(m.keys[opt])<<10 // redirected screen

	return m.parent.MemRead(offset)
}

func (m *ppuMirrorRouter) MemWrite(offset uint16, v byte) byte {
	opt := (offset >> 10) & 3                 // selected screen
	offset &= 0x3ff | uint16(m.keys[opt])<<10 // redirected screen

	return m.parent.MemWrite(offset, v)
}

func (m *ppuMirrorRouter) Ptr() uintptr {
	return uintptr(unsafe.Pointer(m))
}

func (m *ppuMirrorRouter) Length() uint16 {
	return m.parent.Length()
}

func (m *ppuMirrorRouter) String() string {
	return fmt.Sprintf("PPU Mirror proxy (set=%v) of %s", m.keys, m.parent)
}
