package cpu2a03

import (
	"fmt"
	"log"

	"github.com/MagicalTux/gones/memory"
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
	e := cpu.ReadPC()
	o := cpu2a03op[e]
	if o == nil || o.f == nil {
		log.Printf("FATAL CPU ERROR - unsupported op $%02x at $%04x", e, cpu.PC-1)
		cpu.fault = true
		return
	}
	//log.Printf("CPU Step: $%02x o=%v", e, o)
	log.Printf("CPU Step: %s %s", o.i, o.am.Debug(cpu))

	o.f(cpu, o.am)
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

func (cpu *Cpu2A03) ReadPC() uint8 {
	v := cpu.Memory.MemRead(cpu.PC)
	cpu.PC += 1
	return v
}

func (cpu *Cpu2A03) ReadPC16() uint16 {
	v := cpu.Read16(cpu.PC)
	cpu.PC += 2
	return v
}

func (cpu *Cpu2A03) PeekPC() uint8 {
	return cpu.Memory.MemRead(cpu.PC)
}

func (cpu *Cpu2A03) PeekPC16() uint16 {
	return cpu.Read16(cpu.PC)
}

func (cpu *Cpu2A03) Push(v byte) {
	cpu.Memory.MemWrite(0x100+uint16(cpu.S), v)
	cpu.S -= 1
}

func (cpu *Cpu2A03) Pull() byte {
	cpu.S += 1
	v := cpu.Memory.MemRead(0x100 + uint16(cpu.S))
	return v
}

func (cpu *Cpu2A03) Push16(v uint16) {
	// TODO is this the right order?
	cpu.Push(uint8(v & 0xff))
	cpu.Push(uint8((v >> 8) & 0xff))
}

func (cpu *Cpu2A03) Pull16() uint16 {
	var v uint16
	v = uint16(cpu.Pull())
	v |= uint16(cpu.Pull()) << 8
	return v
}

func (cpu *Cpu2A03) flagsNZ(v byte) {
	// set flags N & Z based on value v
	if v == 0 {
		cpu.P |= FlagZero
	} else {
		cpu.P &= ^FlagZero
	}
	if v&0x80 == 0x80 {
		cpu.P |= FlagNegative
	} else {
		cpu.P &= ^FlagNegative
	}
}

func (cpu *Cpu2A03) Read16(offt uint16) uint16 {
	// little endian read
	a := cpu.Memory.MemRead(offt)
	b := cpu.Memory.MemRead(offt + 1)

	return uint16(a) | uint16(b)<<8
}

func (cpu *Cpu2A03) String() string {
	return fmt.Sprintf("CPU:2A03 [A=%02x X=%02x Y=%02x PC=%04x S=%02x P=%02x]", cpu.A, cpu.X, cpu.Y, cpu.PC, cpu.S, cpu.P)
}
