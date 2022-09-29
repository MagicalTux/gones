package cpu2a03

/*
Affects Flags: N Z

MODE           SYNTAX       HEX LEN TIM
Immediate     LDX #$44      $A2  2   2
Zero Page     LDX $44       $A6  2   3
Zero Page,Y   LDX $44,Y     $B6  2   4
Absolute      LDX $4400     $AE  3   4
Absolute,Y    LDX $4400,Y   $BE  3   4+

+ add 1 cycle if page boundary crossed
*/

func init() {
	cpu2a03op[0xa2] = ldxImmed
	cpu2a03op[0xa6] = ldxZero
	cpu2a03op[0xb6] = ldxZeroY
	cpu2a03op[0xae] = ldxAbs
	cpu2a03op[0xbe] = ldxAbsY
}

func ldxImmed(cpu *Cpu2A03) {
	cpu.X = cpu.ReadPC()
}

func ldxZero(cpu *Cpu2A03) {
	addr := uint16(cpu.ReadPC())
	cpu.X = cpu.Memory.MemRead(addr)
}

func ldxZeroY(cpu *Cpu2A03) {
	addr := uint16(cpu.ReadPC() + cpu.Y)
	cpu.X = cpu.Memory.MemRead(addr)
}

func ldxAbs(cpu *Cpu2A03) {
	addr := cpu.ReadPC16()
	cpu.X = cpu.Memory.MemRead(addr)
}

func ldxAbsY(cpu *Cpu2A03) {
	addr := cpu.ReadPC16() + uint16(cpu.Y)
	cpu.X = cpu.Memory.MemRead(addr)
}
