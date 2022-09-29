package cpu2a03

func cmp(cpu *Cpu2A03, am AddressMode) {
	// compare A with value
	// A - M

	v := cpu.A - am.Read(cpu)

	cpu.flagsNZ(v) // also set Carry
}
