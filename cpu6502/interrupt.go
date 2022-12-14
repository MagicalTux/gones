package cpu6502

import "log"

const (
	InterruptNone = iota // 0
	InterruptIRQ
	InterruptNMI
)

func (cpu *CPU) NMI() {
	cpu.interrupt = InterruptNMI
}

func (cpu *CPU) IRQ() {
	if cpu.interrupt < InterruptIRQ {
		cpu.interrupt = InterruptIRQ
	}
}

func (cpu *CPU) handleInterrupt(i byte) {
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

func (cpu *CPU) SetNMI(v byte) {
	cpu.nmiSig = v
}

func (cpu *CPU) checkPendingNMI() {
	if cpu.nmiSig > 0 {
		cpu.nmiSig -= 1
		if cpu.nmiSig == 0 {
			cpu.NMI()
		}
	}
}

func brk(cpu *CPU, am AddressMode) {
	cpu.PC += 1
	cpu.Push16(cpu.PC)
	cpu.Push(cpu.P | FlagBreak)

	cpu.PC = cpu.Read16(IRQVector) // IRQ/BRK
	cpu.setFlag(FlagInterruptDisable, true)
}

func rti(cpu *CPU, am AddressMode) {
	p := cpu.Pull()
	p &= ^byte(0x30)  // ignore B and bit5
	p |= cpu.P & 0x30 // load B and bit5 from P
	cpu.P = p

	cpu.PC = cpu.Pull16()
}

func stop(cpu *CPU, am AddressMode) {
	log.Printf("CPU: STOP instruction at $%04x", cpu.PC-1)
	cpu.fault = true
}

func nop(cpu *CPU, am AddressMode) {
	switch am {
	case amImpl:
		// implicit: do nothing
	default:
		// read value so we advance PC if needed
		am.Read(cpu)
	}
	// do nothing
}
