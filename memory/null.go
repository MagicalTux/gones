package memory

type Null struct{}

func (n Null) MemRead(offset uint16) byte {
	return 0
}

func (n Null) MemWrite(offset uint16, v byte) byte {
	return 0
}

func (n Null) Length() uint16 {
	return 0
}

func (n Null) Ptr() uintptr {
	return 0
}

func (n Null) String() string {
	return "Null (no device)"
}
