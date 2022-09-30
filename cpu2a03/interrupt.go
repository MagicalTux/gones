package cpu2a03

import "log"

func brk(cpu *Cpu2A03, am AddressMode) {
	cpu.PC += 1
	cpu.Push16(cpu.PC)
	cpu.Push(cpu.P | FlagBreak)

	cpu.PC = cpu.Read16(0xfffe) // IRQ/BRK
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
