package cpu2a03

func cmp(cpu *Cpu2A03, am AddressMode) {
	// compare A with value
	// A - M

	cpu.compare(cpu.A, am.Read(cpu))
}

func cpx(cpu *Cpu2A03, am AddressMode) {
	// compare X with value
	// X - M

	cpu.compare(cpu.X, am.Read(cpu))
}

func cpy(cpu *Cpu2A03, am AddressMode) {
	// compare X with value
	// Y - M

	cpu.compare(cpu.Y, am.Read(cpu))
}

func adc(cpu *Cpu2A03, am AddressMode) {
	// Add Memory to Accumulator with Carry
	// A + M + C -> A, C
	v := am.Read(cpu)
	c := byte(0)
	if cpu.getFlag(FlagCarry) {
		c = 1
	}

	a := cpu.A
	cpu.A = a + v + c

	cpu.flagsNZ(cpu.A)

	cpu.setFlag(FlagCarry, int(a)+int(v)+int(c) > 0xff)
	cpu.setFlag(FlagOverflow, (a^v)&0x80 == 0 && (a^cpu.A)&0x80 != 0)
}

func sbc(cpu *Cpu2A03, am AddressMode) {
	// Subtract Memory from Accumulator with Borrow
	// A - M - C -> A
	v := am.Read(cpu)
	c := byte(0)
	if cpu.getFlag(FlagCarry) {
		c = 1
	}

	a := cpu.A
	cpu.A = a - v - (1 - c)

	cpu.flagsNZ(cpu.A)

	cpu.setFlag(FlagCarry, int(a)-int(v)-int(1-c) >= 0)
	cpu.setFlag(FlagOverflow, (a^v)&0x80 != 0 && (a^cpu.A)&0x80 != 0)
}

func inc(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	v := cpu.Memory.MemRead(addr)
	v += 1
	cpu.Memory.MemWrite(addr, v)
	cpu.flagsNZ(v)
}

func dec(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	v := cpu.Memory.MemRead(addr)
	v -= 1
	cpu.Memory.MemWrite(addr, v)
	cpu.flagsNZ(v)
}

func dcp(cpu *Cpu2A03, am AddressMode) {
	// M - 1 -> M, A - M
	// Flags: N Z C
	addr := am.AddrFast(cpu)
	v := cpu.Memory.MemRead(addr)
	v -= 1
	cpu.Memory.MemWrite(addr, v)

	cpu.flagsNZ(cpu.A - v)
	cpu.setFlag(FlagCarry, int(cpu.A)-int(v) >= 0)
}

func isc(cpu *Cpu2A03, am AddressMode) {
	// INC oper + SBC oper
	// M + 1 -> M, A - M - C -> A
	// Flags: N Z C V

	addr := am.AddrFast(cpu)
	v := cpu.Memory.MemRead(addr)

	v += 1
	cpu.Memory.MemWrite(addr, v)

	// code after that is same as sbc

	c := byte(0)
	if cpu.getFlag(FlagCarry) {
		c = 1
	}

	a := cpu.A
	cpu.A = a - v - (1 - c)

	cpu.flagsNZ(cpu.A)

	cpu.setFlag(FlagCarry, int(a)-int(v)-int(1-c) >= 0)
	cpu.setFlag(FlagOverflow, (a^v)&0x80 != 0 && (a^cpu.A)&0x80 != 0)
}
