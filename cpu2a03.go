package main

import (
	"github.com/MagicalTux/gones/memory"
)

type Cpu2A03 struct {
	A    byte   // accumulator
	X, Y byte   // registers
	PC   uint16 // program counter
	S    byte   // stack pointer
	P    byte   // status register

	Memory memory.Master
}

func New2A03() *Cpu2A03 {
	cpu := &Cpu2A03{}

	cpu.Memory = memory.NewBus()

	// setup RAM (2kB=0x800 bytes) with its mirrors
	cpu.Memory.MapHandler(0x0000, 0x2000, memory.NewRAM(0x800))

	return cpu
}
