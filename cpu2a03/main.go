package cpu2a03

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MagicalTux/gones/memory"
)

const (
	NMIVector   = 0xfffa
	ResetVector = 0xfffc
	IRQVector   = 0xfffe
)

type Cpu2A03 struct {
	A    byte   // accumulator
	X, Y byte   // registers
	PC   uint16 // program counter
	S    byte   // stack pointer
	P    byte   // status register

	Memory    memory.Master
	PPU       *PPU
	APU       *APU
	fault     bool
	interrupt byte

	trace *os.File
	cyc   uint64
	clk   uint64
}

func New() *Cpu2A03 {
	cpu := &Cpu2A03{
		Memory: memory.NewBus(),
		APU:    &APU{},
	}
	cpu.PPU = NewPPU(cpu)
	cpu.APU.cpu = cpu

	trace, err := os.Create("trace_2a03.txt")
	if err == nil {
		cpu.trace = trace
	}

	// setup RAM (2kB=0x800 bytes) with its mirrors
	cpu.Memory.MapHandler(0x0000, 0x2000, memory.NewRAM(0x800))
	cpu.Memory.MapHandler(0x2000, 0x2000, cpu.PPU) // PPU at 0x2000
	cpu.Memory.MapHandler(0x4000, 0x2000, cpu.APU) // APU at 0x4000

	return cpu
}

// Typically this runs into a goroutine
// go cpu.Start(cpu2a03.NTSC)
func (cpu *Cpu2A03) Start(clockLn time.Duration) {
	t := time.NewTicker(clockLn)
	defer t.Stop()

	for !cpu.fault {
		cpu.Clock()

		select {
		case <-t.C:
		}
	}

	log.Printf("CPU stopped due to fault: %s", cpu)

	for i := uint16(2); i <= 3; i++ {
		log.Printf("Address value at %d: $%02x", i, cpu.Memory.MemRead(i))
	}
}

func (cpu *Cpu2A03) Clock() {
	if cpu.fault {
		return
	}
	cpu.clk += 1
	if cpu.clk < cpu.cyc {
		// wait for clock to catch up
		return
	}

	if cpu.interrupt != InterruptNone {
		cpu.handleInterrupt()
		cpu.interrupt = InterruptNone
	}

	pos := cpu.PC
	// read value at PC
	e := cpu.ReadPC()
	o := cpu2a03op[e]
	if o == nil || o.f == nil {
		cpu.fatal("FATAL CPU ERROR - unsupported op $%02x @ $%04x / %s", e, pos, cpu)
		return
	}
	//log.Printf("CPU Step: $%02x o=%v", e, o)
	//log.Printf("CPU Step: [$%04x] %s %s", pos, o.i, o.am.Debug(cpu))
	if cpu.trace != nil {
		fmt.Fprintf(cpu.trace, "CPU Step cyc=%d: [$%04x] %s % -32s %s\n", cpu.cyc, pos, o.i, o.am.Debug(cpu), cpu)
	}

	o.f(cpu, o.am)

	cyc := int(o.cyc)
	cpu.cyc += uint64(cyc)

	// move PPU forward
	cpu.PPU.Clock(cyc * 3)
}

func (cpu *Cpu2A03) Reset() {
	// reset
	cpu.A = 0
	cpu.X = 0
	cpu.Y = 0
	cpu.S = 0xfd
	cpu.P = FlagIgnored | FlagInterruptDisable

	cpu.cyc = 7 // cpu init typically takes 7 cycles
	cpu.clk = 7 // which is already done

	// $FFFC-$FFFD = Reset vector
	cpu.PC = cpu.Read16(ResetVector)

	cpu.PPU.Reset(cpu.cyc)

	log.Printf("CPU reset, new state = %s", cpu)
}

func (cpu *Cpu2A03) fatal(v string, a ...any) {
	cpu.fault = true
	log.Printf("CPU FAULT: "+v, a...)
	cpu.msg(v, a...)
}

func (cpu *Cpu2A03) msg(v string, a ...any) {
	if cpu.trace != nil {
		fmt.Fprintf(cpu.trace, "Debug: "+v+"\n", a...)
	}
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
	//cpu.msg("CPU Stack push $%02x", v)
	cpu.Memory.MemWrite(0x100+uint16(cpu.S), v)
	cpu.S -= 1
}

func (cpu *Cpu2A03) Pull() byte {
	cpu.S += 1
	v := cpu.Memory.MemRead(0x100 + uint16(cpu.S))
	//cpu.msg("CPU Stack pull $%02x S=%02x", v, cpu.S)
	return v
}

func (cpu *Cpu2A03) Push16(v uint16) {
	// TODO is this the right order?
	cpu.Push(uint8((v >> 8) & 0xff))
	cpu.Push(uint8(v & 0xff))
}

func (cpu *Cpu2A03) Pull16() uint16 {
	var v uint16
	v = uint16(cpu.Pull())
	v |= uint16(cpu.Pull()) << 8
	return v
}

func (cpu *Cpu2A03) Read16(offt uint16) uint16 {
	// little endian read
	a := cpu.Memory.MemRead(offt)
	b := cpu.Memory.MemRead(offt + 1)

	return uint16(a) | uint16(b)<<8
}

func (cpu *Cpu2A03) Read16W(offt uint16) uint16 {
	// little endian read wrapping around current page
	a := cpu.Memory.MemRead(offt)
	if offt&0xff == 0xff {
		// would wrap
		offt &= 0xff00
	} else {
		offt += 1
	}
	b := cpu.Memory.MemRead(offt)

	return uint16(a) | uint16(b)<<8
}

func (cpu *Cpu2A03) String() string {
	return fmt.Sprintf("CPU:2A03 [A=%02x X=%02x Y=%02x PC=%04x S=%02x P=%02x]", cpu.A, cpu.X, cpu.Y, cpu.PC, cpu.S, cpu.P)
}
