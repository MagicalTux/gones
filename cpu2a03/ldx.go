package cpu2a03

func lda(cpu *Cpu2A03, am AddressMode) {
	cpu.A = am.Read(cpu)
	cpu.flagsNZ(cpu.X)
}

func ldx(cpu *Cpu2A03, am AddressMode) {
	cpu.X = am.Read(cpu)
	cpu.flagsNZ(cpu.X)
}

func ldy(cpu *Cpu2A03, am AddressMode) {
	cpu.Y = am.Read(cpu)
	cpu.flagsNZ(cpu.Y)
}
