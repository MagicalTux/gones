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

func init() {
	cpu2a03op[0x10] = bpl
	cpu2a03op[0x30] = bmi
	cpu2a03op[0x50] = bvc
	cpu2a03op[0x70] = bvs
	cpu2a03op[0x90] = bcc
	cpu2a03op[0xb0] = bcs
	cpu2a03op[0xd0] = bne
	cpu2a03op[0xf0] = beq
}

func bpl(cpu *Cpu2A03) {
	offt := uint16(cpu.ReadPC())
	if cpu.P&FlagNegative == 0 {
		if offt&0x80 == 0x80 {
			offt |= 0xff00
		}
		cpu.PC += offt
	}
}

func bmi(cpu *Cpu2A03) {
	offt := uint16(cpu.ReadPC())
	if cpu.P&FlagNegative == FlagNegative {
		if offt&0x80 == 0x80 {
			offt |= 0xff00
		}
		cpu.PC += offt
	}
}

func bvc(cpu *Cpu2A03) {
	offt := uint16(cpu.ReadPC())
	if cpu.P&FlagOverflow == 0 {
		if offt&0x80 == 0x80 {
			offt |= 0xff00
		}
		cpu.PC += offt
	}
}

func bvs(cpu *Cpu2A03) {
	offt := uint16(cpu.ReadPC())
	if cpu.P&FlagOverflow == FlagOverflow {
		if offt&0x80 == 0x80 {
			offt |= 0xff00
		}
		cpu.PC += offt
	}
}

func bcc(cpu *Cpu2A03) {
	offt := uint16(cpu.ReadPC())
	if cpu.P&FlagCarry == 0 {
		if offt&0x80 == 0x80 {
			offt |= 0xff00
		}
		cpu.PC += offt
	}
}

func bcs(cpu *Cpu2A03) {
	offt := uint16(cpu.ReadPC())
	if cpu.P&FlagCarry == FlagCarry {
		if offt&0x80 == 0x80 {
			offt |= 0xff00
		}
		cpu.PC += offt
	}
}

func bne(cpu *Cpu2A03) {
	offt := uint16(cpu.ReadPC())
	if cpu.P&FlagZero == 0 {
		if offt&0x80 == 0x80 {
			offt |= 0xff00
		}
		cpu.PC += offt
	}
}

func beq(cpu *Cpu2A03) {
	offt := uint16(cpu.ReadPC())
	if cpu.P&FlagZero == FlagZero {
		if offt&0x80 == 0x80 {
			offt |= 0xff00
		}
		cpu.PC += offt
	}
}
