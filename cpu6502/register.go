package cpu6502

func lda(cpu *CPU, am AddressMode) {
	cpu.A = am.Read(cpu)
	cpu.flagsNZ(cpu.A)
}

func ldx(cpu *CPU, am AddressMode) {
	cpu.X = am.Read(cpu)
	cpu.flagsNZ(cpu.X)
}

func ldy(cpu *CPU, am AddressMode) {
	cpu.Y = am.Read(cpu)
	cpu.flagsNZ(cpu.Y)
}

func sta(cpu *CPU, am AddressMode) {
	am.WriteFast(cpu, cpu.A)
}

func stx(cpu *CPU, am AddressMode) {
	am.WriteFast(cpu, cpu.X)
}

func sty(cpu *CPU, am AddressMode) {
	am.WriteFast(cpu, cpu.Y)
}

func tax(cpu *CPU, am AddressMode) {
	cpu.X = cpu.A
	cpu.flagsNZ(cpu.X)
}

func txa(cpu *CPU, am AddressMode) {
	cpu.A = cpu.X
	cpu.flagsNZ(cpu.A)
}

func tay(cpu *CPU, am AddressMode) {
	cpu.Y = cpu.A
	cpu.flagsNZ(cpu.Y)
}

func tya(cpu *CPU, am AddressMode) {
	cpu.A = cpu.Y
	cpu.flagsNZ(cpu.A)
}

func dex(cpu *CPU, am AddressMode) {
	// X - 1 -> X
	am.Implied(cpu)
	cpu.X -= 1
	cpu.flagsNZ(cpu.X)
}

func inx(cpu *CPU, am AddressMode) {
	am.Implied(cpu)
	cpu.X += 1
	cpu.flagsNZ(cpu.X)
}

func dey(cpu *CPU, am AddressMode) {
	// Y - 1 -> Y
	am.Implied(cpu)
	cpu.Y -= 1
	cpu.flagsNZ(cpu.Y)
}

func iny(cpu *CPU, am AddressMode) {
	am.Implied(cpu)
	cpu.Y += 1
	cpu.flagsNZ(cpu.Y)
}

func lax(cpu *CPU, am AddressMode) {
	v := am.Read(cpu)
	cpu.A = v
	cpu.X = v
	cpu.flagsNZ(v)
}

func sax(cpu *CPU, am AddressMode) {
	// A AND X -> M
	// no flags
	v := cpu.A & cpu.X
	am.Write(cpu, v)
}
