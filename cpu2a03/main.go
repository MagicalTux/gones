package cpu2a03

import (
	"fmt"
	"log"

	"github.com/MagicalTux/gones/memory"
)

type op2a03 func(cpu *Cpu2A03)

var (
	cpu2a03op [256]op2a03
)

type Cpu2A03 struct {
	A    byte   // accumulator
	X, Y byte   // registers
	PC   uint16 // program counter
	S    byte   // stack pointer
	P    byte   // status register

	Memory memory.Master
	PPU    *PPU
	APU    *APU
	fault  bool
}

func New() *Cpu2A03 {
	cpu := &Cpu2A03{
		Memory: memory.NewBus(),
		PPU:    &PPU{},
		APU:    &APU{},
	}
	cpu.APU.cpu = cpu

	// setup RAM (2kB=0x800 bytes) with its mirrors
	cpu.Memory.MapHandler(0x0000, 0x2000, memory.NewRAM(0x800))
	cpu.Memory.MapHandler(0x2000, 0x2000, cpu.PPU)
	cpu.Memory.MapHandler(0x4000, 0x2000, cpu.APU)

	return cpu
}

func (cpu *Cpu2A03) Step() {
	if cpu.fault {
		return
	}
	// read value at PC
	e := cpu.Memory.MemRead(cpu.PC)
	f := cpu2a03op[e]
	log.Printf("CPU Step: $%02x f=%v", e, f)
	if f == nil {
		log.Printf("FATAL CPU ERROR - unsupported operand")
		cpu.fault = true
		return
	}
}

func (cpu *Cpu2A03) Reset() {
	// reset
	cpu.A = 0
	cpu.X = 0
	cpu.Y = 0
	cpu.S = 0
	cpu.P = 0

	// $FFFC-$FFFD = Reset vector
	cpu.PC = cpu.Read16(0xfffc)

	log.Printf("CPU reset, new state = %s", cpu)
}

func (cpu *Cpu2A03) Read16(offt uint16) uint16 {
	// little endian read
	a := cpu.Memory.MemRead(offt)
	b := cpu.Memory.MemRead(offt + 1)

	return uint16(a) | uint16(b)<<8
}

func (cpu *Cpu2A03) String() string {
	return fmt.Sprintf("2A03 [A=%02x X=%02x Y=%02x PC=%04x S=%02x P=%02x]", cpu.A, cpu.X, cpu.Y, cpu.PC, cpu.S, cpu.P)
}
