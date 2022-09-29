package memory

type Handler interface {
	MemRead(offset uint16) byte
	MemWrite(offset uint16, val byte) byte
}

type Master interface {
	Handler
	MapHandler(offset uint16, length uint16, h Handler)
}
