package memory

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
