package cpu2a03

/*
MNEMONIC                       HEX
BPL (Branch on PLus)           $10
BMI (Branch on MInus)          $30
BVC (Branch on oVerflow Clear) $50
BVS (Branch on oVerflow Set)   $70
BCC (Branch on Carry Clear)    $90
BCS (Branch on Carry Set)      $B0
BNE (Branch on Not Equal)      $D0
BEQ (Branch on EQual)          $F0
*/

func bpl(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagNegative == 0 {
		cpu.PC = addr
	}
}

func bmi(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagNegative == FlagNegative {
		cpu.PC = addr
	}
}

func bvc(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagOverflow == 0 {
		cpu.PC = addr
	}
}

func bvs(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagOverflow == FlagOverflow {
		cpu.PC = addr
	}
}

func bcc(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagCarry == 0 {
		cpu.PC = addr
	}
}

func bcs(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagCarry == FlagCarry {
		cpu.PC = addr
	}
}

func bne(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagZero == 0 {
		cpu.PC = addr
	}
}

func beq(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	if cpu.P&FlagZero == FlagZero {
		cpu.PC = addr
	}
}

func jmp(cpu *Cpu2A03, am AddressMode) {
	cpu.PC = am.Addr(cpu)
}

func jsr(cpu *Cpu2A03, am AddressMode) {
	addr := am.Addr(cpu)
	cpu.msg("JSR push $%04x", cpu.PC)
	cpu.Push16(cpu.PC)
	cpu.PC = addr
}

func rts(cpu *Cpu2A03, am AddressMode) {
	cpu.PC = cpu.Pull16()
}
