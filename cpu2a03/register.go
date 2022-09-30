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

func tax(cpu *Cpu2A03, am AddressMode) {
	cpu.X = cpu.A
	cpu.flagsNZ(cpu.X)
}

func txa(cpu *Cpu2A03, am AddressMode) {
	cpu.A = cpu.X
	cpu.flagsNZ(cpu.A)
}

func tay(cpu *Cpu2A03, am AddressMode) {
	cpu.Y = cpu.A
	cpu.flagsNZ(cpu.Y)
}

func tya(cpu *Cpu2A03, am AddressMode) {
	cpu.A = cpu.Y
	cpu.flagsNZ(cpu.A)
}

func dex(cpu *Cpu2A03, am AddressMode) {
	// X - 1 -> X
	am.Implied(cpu)
	cpu.X -= 1
	cpu.flagsNZ(cpu.X)
}

func inx(cpu *Cpu2A03, am AddressMode) {
	am.Implied(cpu)
	cpu.X += 1
	cpu.flagsNZ(cpu.X)
}

func dey(cpu *Cpu2A03, am AddressMode) {
	// Y - 1 -> Y
	am.Implied(cpu)
	cpu.Y -= 1
	cpu.flagsNZ(cpu.Y)
}

func iny(cpu *Cpu2A03, am AddressMode) {
	am.Implied(cpu)
	cpu.Y += 1
	cpu.flagsNZ(cpu.Y)
}

func lax(cpu *Cpu2A03, am AddressMode) {
	v := am.Read(cpu)
	cpu.A = v
	cpu.X = v
	cpu.flagsNZ(v)
}

func sax(cpu *Cpu2A03, am AddressMode) {
	// A AND X -> M
	// no flags
	v := cpu.A & cpu.X
	am.Write(cpu, v)
}
