package pkgnes

import (
	"github.com/MagicalTux/gones/apu"
	"github.com/MagicalTux/gones/cpu2a03"
	"github.com/MagicalTux/gones/memory"
	"github.com/MagicalTux/gones/nesclock"
	"github.com/MagicalTux/gones/ppu"
)

type NES struct {
	Memory memory.Master
	Clk    *nesclock.Master
	CPU    *cpu2a03.Cpu2A03
	PPU    *ppu.PPU
	APU    *apu.APU
	Input  []apu.InputDevice
	model  Model
}

func New(model Model) *NES {
	nes := &NES{
		Memory: memory.NewBus(),
		Clk:    model.newClock(),
		CPU:    cpu2a03.New(),
		PPU:    ppu.New(),
	}
	nes.CPU.Memory = nes.Memory              // connect main memory bus to CPU
	nes.PPU.VBlankInterrupt = nes.CPU.SetNMI // connect PPU's vblank to NMI

	nes.APU = apu.New(nes.Memory, nes.CPU.TimeFreeze) // APU has access to the cpu's memory & clock
	nes.Input = nes.APU.Input[:]
	nes.APU.Interrupt = nes.CPU.IRQ

	// setup RAM (2kB=0x800 bytes) with its mirrors
	nes.Memory.MapHandler(0x0000, 0x2000, memory.NewRAM(0x800))
	nes.Memory.MapHandler(0x2000, 0x2000, nes.PPU) // PPU at 0x2000
	nes.Memory.MapHandler(0x4000, 0x2000, nes.APU) // APU at 0x4000

	return nes
}

// Typically this runs into a goroutine
// go nes.Start(pkgnes.NTSC)
func (nes *NES) Start() {
	// trigger once every 12 clocks (if NTSC)
	nes.Clk.Listen(nes.model.cpuIntv(), 0, nes.CPU.Clock)

	// ppu & apu are run with a offset of 1 to ensure they run after the cpu
	nes.Clk.Listen(nes.model.ppuIntv(), 1, nes.PPU.Clock)

	// apu needs 3 clocks
	nes.Clk.Listen(nes.model.cpuIntv(), 1, nes.APU.ClockCPU)
	nes.Clk.Listen(nes.Clk.Frequency()/240, 1, nes.APU.Clock240)
	nes.Clk.Listen(nes.Clk.Frequency()/44100, 1, nes.APU.Clock44100)
}

func (nes *NES) Reset() {
	nes.CPU.Reset()
	nes.PPU.Reset()
}
