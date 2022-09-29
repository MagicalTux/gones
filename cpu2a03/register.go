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

func sta(cpu *Cpu2A03, am AddressMode) {
	am.Write(cpu, cpu.A)
}

func stx(cpu *Cpu2A03, am AddressMode) {
	am.Write(cpu, cpu.X)
}

func sty(cpu *Cpu2A03, am AddressMode) {
	am.Write(cpu, cpu.Y)
}

func dex(cpu *Cpu2A03, am AddressMode) {
	// X - 1 -> X
	am.Implied(cpu)
	cpu.X -= 1
	cpu.flagsNZ(cpu.X)
}
