package main

type Cpu2A03 struct {
	A    byte   // accumulator
	X, Y byte   // registers
	PC   uint16 // program counter
	S    byte   // stack pointer
	P    byte   // status register

	Memory *MMU
}

func New2A03() *Cpu2A03 {
	cpu := &Cpu2A03{}

	cpu.Memory = NewMMU()

	// setup RAM (2kB=0x800 bytes) with its mirrors
	ram := make([]byte, 0x0800)
	cpu.Memory.Map(0x0000, ram)
	cpu.Memory.Map(0x0800, ram)
	cpu.Memory.Map(0x1000, ram)
	cpu.Memory.Map(0x1800, ram)

	return cpu
}
