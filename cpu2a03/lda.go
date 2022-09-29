package cpu2a03

/*
Affects Flags: N Z

MODE           SYNTAX       HEX LEN TIM
Immediate     LDA #$44      $A9  2   2
Zero Page     LDA $44       $A5  2   3
Zero Page,X   LDA $44,X     $B5  2   4
Absolute      LDA $4400     $AD  3   4
Absolute,X    LDA $4400,X   $BD  3   4+
Absolute,Y    LDA $4400,Y   $B9  3   4+
Indirect,X    LDA ($44,X)   $A1  2   6
Indirect,Y    LDA ($44),Y   $B1  2   5+

+ add 1 cycle if page boundary crossed
*/

func init() {
	cpu2a03op[0xa9] = ldaImmed
	cpu2a03op[0xa5] = ldaZero
	cpu2a03op[0xb5] = ldaZeroX
	cpu2a03op[0xad] = ldaAbs
	cpu2a03op[0xbd] = ldaAbsX
	cpu2a03op[0xb9] = ldaAbsY
	cpu2a03op[0xa1] = ldaIndX
	cpu2a03op[0xb1] = ldaIndX
}

func ldaImmed(cpu *Cpu2A03) {
	cpu.A = cpu.ReadPC()
	cpu.flagsNZ(cpu.A)
}

func ldaZero(cpu *Cpu2A03) {
	addr := uint16(cpu.ReadPC())
	cpu.A = cpu.Memory.MemRead(addr)
	cpu.flagsNZ(cpu.A)
}

func ldaZeroX(cpu *Cpu2A03) {
	addr := uint16(cpu.ReadPC() + cpu.X)
	cpu.A = cpu.Memory.MemRead(addr)
	cpu.flagsNZ(cpu.A)
}

func ldaAbs(cpu *Cpu2A03) {
	addr := cpu.ReadPC16()
	cpu.A = cpu.Memory.MemRead(addr)
	cpu.flagsNZ(cpu.A)
}

func ldaAbsX(cpu *Cpu2A03) {
	addr := cpu.ReadPC16() + uint16(cpu.X)
	cpu.A = cpu.Memory.MemRead(addr)
	cpu.flagsNZ(cpu.A)
}

func ldaAbsY(cpu *Cpu2A03) {
	addr := cpu.ReadPC16() + uint16(cpu.Y)
	cpu.A = cpu.Memory.MemRead(addr)
	cpu.flagsNZ(cpu.A)
}

func ldaIndX(cpu *Cpu2A03) {
	// operand is zeropage address; effective address is word in (LL + X, LL + X + 1), inc. without carry: C.w($00LL + X)
	addr := uint16(cpu.ReadPC()) + uint16(cpu.X)
	addr = cpu.Read16(addr)
	cpu.A = cpu.Memory.MemRead(addr)
	cpu.flagsNZ(cpu.A)
}

func ldaIndY(cpu *Cpu2A03) {
	// operand is zeropage address; effective address is word in (LL, LL + 1) incremented by Y with carry: C.w($00LL) + Y
	addr := uint16(cpu.ReadPC())
	addr = cpu.Read16(addr) + uint16(cpu.Y)
	cpu.A = cpu.Memory.MemRead(addr)
	cpu.flagsNZ(cpu.A)
}
