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
