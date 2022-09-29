package cpu2a03

/*
Affects Flags: none

MODE           SYNTAX       HEX LEN TIM
Zero Page     STA $44       $85  2   3
Zero Page,X   STA $44,X     $95  2   4
Absolute      STA $4400     $8D  3   4
Absolute,X    STA $4400,X   $9D  3   5
Absolute,Y    STA $4400,Y   $99  3   5
Indirect,X    STA ($44,X)   $81  2   6
Indirect,Y    STA ($44),Y   $91  2   6
*/

func init() {
	cpu2a03op[0x85] = staZero
	cpu2a03op[0x95] = staZeroX
	cpu2a03op[0x8d] = staAbs
	cpu2a03op[0x9d] = staAbsX
	cpu2a03op[0x99] = staAbsY
	cpu2a03op[0x81] = staIndX
	cpu2a03op[0x91] = staIndX
}

func staZero(cpu *Cpu2A03) {
	addr := uint16(cpu.ReadPC())
	cpu.Memory.MemWrite(addr, cpu.A)
}

func staZeroX(cpu *Cpu2A03) {
	addr := uint16(cpu.ReadPC() + cpu.X)
	cpu.Memory.MemWrite(addr, cpu.A)
}

func staAbs(cpu *Cpu2A03) {
	addr := cpu.ReadPC16()
	cpu.Memory.MemWrite(addr, cpu.A)
}

func staAbsX(cpu *Cpu2A03) {
	addr := cpu.ReadPC16() + uint16(cpu.X)
	cpu.Memory.MemWrite(addr, cpu.A)
}

func staAbsY(cpu *Cpu2A03) {
	addr := cpu.ReadPC16() + uint16(cpu.Y)
	cpu.Memory.MemWrite(addr, cpu.A)
}

func staIndX(cpu *Cpu2A03) {
	// operand is zeropage address; effective address is word in (LL + X, LL + X + 1), inc. without carry: C.w($00LL + X)
	addr := uint16(cpu.ReadPC()) + uint16(cpu.X)
	addr = cpu.Read16(addr)
	cpu.Memory.MemWrite(addr, cpu.A)
}

func staIndY(cpu *Cpu2A03) {
	// operand is zeropage address; effective address is word in (LL, LL + 1) incremented by Y with carry: C.w($00LL) + Y
	addr := uint16(cpu.ReadPC())
	addr = cpu.Read16(addr) + uint16(cpu.Y)
	cpu.Memory.MemWrite(addr, cpu.A)
}
