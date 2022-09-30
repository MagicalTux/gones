package cpu2a03

func cmp(cpu *Cpu2A03, am AddressMode) {
	// compare A with value
	// A - M

	v := cpu.A - am.Read(cpu)

	cpu.flagsNZ(v) // TODO also set Carry
}

func cpx(cpu *Cpu2A03, am AddressMode) {
	// compare X with value
	// X - M

	v := cpu.X - am.Read(cpu)

	cpu.flagsNZ(v) // TODO also set Carry
}

func cpy(cpu *Cpu2A03, am AddressMode) {
	// compare X with value
	// Y - M

	v := cpu.Y - am.Read(cpu)

	cpu.flagsNZ(v) // TODO also set Carry
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

func dec(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	v := cpu.Memory.MemRead(addr)
	v -= 1
	cpu.Memory.MemWrite(addr, v)
	cpu.flagsNZ(v)
}
