package cpu2a03

import "log"

const (
	InterruptNone = iota // 0
	InterruptIRQ
	InterruptNMI
)

func (cpu *Cpu2A03) NMI() {
	cpu.interrupt = InterruptNMI
}

func (cpu *Cpu2A03) IRQ() {
	if cpu.interrupt < InterruptIRQ {
		cpu.interrupt = InterruptIRQ
	}
}

func (cpu *Cpu2A03) handleInterrupt(i byte) {
	// handle interrupt
	cpu.Push16(cpu.PC)
	cpu.Push(cpu.P) // | FlagBreak)

	switch i {
	case InterruptNMI:
		cpu.PC = cpu.Read16(NMIVector) // NMI interrupt vector
	case InterruptIRQ:
		cpu.PC = cpu.Read16(IRQVector) // IRQ interrupt vector
	}
	cpu.setFlag(FlagInterruptDisable, true)
	cpu.cyc += 7
}

func (cpu *Cpu2A03) setNMI(v byte) {
	cpu.nmiSig = v
}

func (cpu *Cpu2A03) checkPendingNMI() {
	if cpu.nmiSig > 0 {
		cpu.nmiSig -= 1
		if cpu.nmiSig == 0 {
			cpu.NMI()
		}
	}
}

func brk(cpu *Cpu2A03, am AddressMode) {
	cpu.PC += 1
	cpu.Push16(cpu.PC)
	cpu.Push(cpu.P | FlagBreak)

	cpu.PC = cpu.Read16(IRQVector) // IRQ/BRK
	cpu.setFlag(FlagInterruptDisable, true)
}

func rti(cpu *Cpu2A03, am AddressMode) {
	p := cpu.Pull()
	p &= ^byte(0x30)  // ignore B and bit5
	p |= cpu.P & 0x30 // load B and bit5 from P
	cpu.P = p

	cpu.PC = cpu.Pull16()
}

func stop(cpu *Cpu2A03, am AddressMode) {
	log.Printf("CPU: STOP instruction at $%04x", cpu.PC-1)
	cpu.fault = true
}

func nop(cpu *Cpu2A03, am AddressMode) {
	switch am {
	case amImpl:
		// implicit: do nothing
	default:
		// read value so we advance PC if needed
		am.Read(cpu)
	}
	// do nothing
}
