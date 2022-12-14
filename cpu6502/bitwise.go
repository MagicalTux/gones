package cpu6502

func ora(cpu *CPU, am AddressMode) {
	// Affects Flags: N Z
	v := am.Read(cpu)

	cpu.A |= v
	cpu.flagsNZ(cpu.A)
}

func bit(cpu *CPU, am AddressMode) {
	v := am.Read(cpu)
	// bits 7 and 6 of operand are transfered to bit 7 and 6 of SR (N,V);
	// the zero-flag is set to the result of operand AND accumulator.

	// A AND M, M7 -> N, M6 -> V

	cpu.setFlag(FlagOverflow, (v>>6)&1 == 1) // V
	cpu.flagsZ(v & cpu.A)
	cpu.flagsN(v)
}

func and(cpu *CPU, am AddressMode) {
	v := am.Read(cpu)

	cpu.A &= v
	cpu.flagsNZ(cpu.A)
}

func eor(cpu *CPU, am AddressMode) {
	v := am.Read(cpu)

	cpu.A ^= v
	cpu.flagsNZ(cpu.A)
}

func ror(cpu *CPU, am AddressMode) {
	c := byte(0)
	if cpu.getFlag(FlagCarry) {
		c = 1
	}

	if am == amAcc {
		// act on cpu.A
		cpu.setFlag(FlagCarry, cpu.A&1 == 1)
		cpu.A = (cpu.A >> 1) | (c << 7)
		cpu.flagsNZ(cpu.A)
	} else {
		// act on mem
		addr := am.AddrFast(cpu)
		v := cpu.Memory.MemRead(addr)

		cpu.setFlag(FlagCarry, v&1 == 1)
		v = (v >> 1) | (c << 7)
		cpu.Memory.MemWrite(addr, v)
		cpu.flagsNZ(v)
	}
}

func lsr(cpu *CPU, am AddressMode) {
	if am == amAcc {
		cpu.setFlag(FlagCarry, cpu.A&1 == 1)
		cpu.A >>= 1
		cpu.flagsNZ(cpu.A)
	} else {
		addr := am.AddrFast(cpu)
		v := cpu.Memory.MemRead(addr)

		cpu.setFlag(FlagCarry, v&1 == 1)
		v >>= 1
		cpu.Memory.MemWrite(addr, v)
		cpu.flagsNZ(v)
	}
}

func asl(cpu *CPU, am AddressMode) {
	if am == amAcc {
		cpu.setFlag(FlagCarry, cpu.A&0x80 == 0x80)
		cpu.A <<= 1
		cpu.flagsNZ(cpu.A)
	} else {
		addr := am.AddrFast(cpu)
		v := cpu.Memory.MemRead(addr)

		cpu.setFlag(FlagCarry, v&0x80 == 0x80)
		v <<= 1
		cpu.Memory.MemWrite(addr, v)
		cpu.flagsNZ(v)
	}
}

func rol(cpu *CPU, am AddressMode) {
	c := byte(0)
	if cpu.getFlag(FlagCarry) {
		c = 1
	}

	if am == amAcc {
		// act on cpu.A
		cpu.setFlag(FlagCarry, cpu.A&0x80 == 0x80)
		cpu.A = cpu.A<<1 | c
		cpu.flagsNZ(cpu.A)
	} else {
		// act on mem
		addr := am.AddrFast(cpu)
		v := cpu.Memory.MemRead(addr)

		// re-write the original value as per 6502
		// See: https://www.nesdev.org/6502_cpu.txt
		// ??? "Read-Modify-Write instructions"
		cpu.Memory.MemWrite(addr, v)

		cpu.setFlag(FlagCarry, v&0x80 == 0x80)
		v = (v << 1) | c
		cpu.Memory.MemWrite(addr, v)
		cpu.flagsNZ(v)
	}
}
