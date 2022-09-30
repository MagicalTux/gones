package cpu2a03

// branchTo branches execution to the given address while accounting for cycles
func (cpu *Cpu2A03) branchTo(addr uint16) {
	if cpu.PC&0xff00 != addr&0xff00 {
		// different page
		cpu.cyc += 2
	} else {
		cpu.cyc += 1
	}
	cpu.PC = addr
}

func bpl(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagNegative == 0 {
		cpu.branchTo(addr)
	}
}

func bmi(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagNegative == FlagNegative {
		cpu.branchTo(addr)
	}
}

func bvc(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagOverflow == 0 {
		cpu.branchTo(addr)
	}
}

func bvs(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagOverflow == FlagOverflow {
		cpu.branchTo(addr)
	}
}

func bcc(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagCarry == 0 {
		cpu.branchTo(addr)
	}
}

func bcs(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagCarry == FlagCarry {
		cpu.branchTo(addr)
	}
}

func bne(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagZero == 0 {
		cpu.branchTo(addr)
	}
}

func beq(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagZero == FlagZero {
		cpu.branchTo(addr)
	}
}

func jmp(cpu *Cpu2A03, am AddressMode) {
	cpu.PC = am.Addr(cpu)
}

func jsr(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	cpu.Push16(cpu.PC - 1) // push PC+2 (we are at PC+3 now)
	cpu.PC = addr
}

func rts(cpu *Cpu2A03, am AddressMode) {
	cpu.PC = cpu.Pull16() + 1
}
