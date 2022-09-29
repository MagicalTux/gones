package cartridge

import "os"

type Data struct {
	f *os.File
	m []byte // map+len
}
