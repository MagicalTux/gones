package cpu2a03

func ora(cpu *Cpu2A03, am AddressMode) {
	// Affects Flags: N Z
	v := am.Read(cpu)

	cpu.A |= v
	cpu.flagsNZ(cpu.A)
}

func bit(cpu *Cpu2A03, am AddressMode) {
	v := am.Read(cpu)
	// bits 7 and 6 of operand are transfered to bit 7 and 6 of SR (N,V);
	// the zero-flag is set to the result of operand AND accumulator.

	// A AND M, M7 -> N, M6 -> V

	cpu.setFlag(FlagOverflow, (v>>6)&1 == 1) // V
	cpu.flagsZ(v & cpu.A)
	cpu.flagsN(v)
}

func and(cpu *Cpu2A03, am AddressMode) {
	v := am.Read(cpu)

	cpu.A &= v
	cpu.flagsNZ(cpu.A)
}

func eor(cpu *Cpu2A03, am AddressMode) {
	v := am.Read(cpu)

	cpu.A ^= v
	cpu.flagsNZ(cpu.A)
}
