package main

// A page size of 0x100 means 256 pages for the whole 16bit system
const (
	PageBits   = 8
	OffsetBits = 16 - PageBits
	PageSize   = 1 << PageBits
	OffsetSize = 1 << OffsetBits
)

// see: https://www.nesdev.org/wiki/CPU_memory_map

type mmuHandler interface {
	MemRead(offset uint16) byte
	MemWrite(offset uint16, val byte)
}

type MMU struct {
	direct   [OffsetSize][]byte
	indirect [OffsetSize]mmuHandler

	ram []byte // 2kB of ram at 0x0000~
}

func NewMMU() *MMU {
	res := &MMU{}
	res.ram = make([]byte, 2048)
	res.Map(0, res.ram)

	return res
}

func (m *MMU) Map(offset uint16, buf []byte) {
	// direct mapping
	offt := offset >> PageBits
	cnt := uint16(len(buf) >> PageBits)
	if len(buf)%PageBits != 0 {
		cnt += 1
	}

	for i := uint16(0); i < cnt; i++ {
		inoff := i << PageBits // offset in page
		m.direct[offt+i] = buf[inoff : inoff+PageSize]
	}
}

func (m *MMU) MemRead(offset uint16) byte {
	offt := offset >> PageBits
	if v := m.direct[offt]; v != nil {
		inoff := int(offset) & (PageSize - 1)
		if len(v) > inoff {
			return v[inoff]
		}
	}
	if v := m.indirect[offt]; v != nil {
		return v.MemRead(offset)
	}
	return 0 // page fault
}

func (m *MMU) MemWrite(offset uint16, val byte) {
	offt := offset >> PageBits
	if v := m.direct[offt]; v != nil {
		inoff := int(offset) & (PageSize - 1)
		if len(v) > inoff {
			v[inoff] = val
			return
		}
	}
	if v := m.indirect[offt]; v != nil {
		v.MemWrite(offset, val)
		return
	}
	// page fault
}
