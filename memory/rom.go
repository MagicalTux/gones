package memory

type ROM []byte

func (r ROM) MemRead(offset uint16) byte {
	return r[offset%uint16(len(r)-1)]
}

func (r ROM) MemWrite(offset uint16, val byte) byte {
	return r.MemRead(offset)
}
