package cpu2a03

/*
MNEMONIC                       HEX
CLC (CLear Carry)              $18
SEC (SEt Carry)                $38
CLI (CLear Interrupt)          $58
SEI (SEt Interrupt)            $78
CLV (CLear oVerflow)           $B8
CLD (CLear Decimal)            $D8
SED (SEt Decimal)              $F8
*/

func init() {
	cpu2a03op[0x18] = clc
	cpu2a03op[0x38] = sec
	cpu2a03op[0x58] = cli
	cpu2a03op[0x78] = sei
	cpu2a03op[0xb8] = clv
	cpu2a03op[0xd8] = cld
	cpu2a03op[0xf8] = sed
}

func clc(cpu *Cpu2A03) {
	cpu.P &= ^FlagCarry
}

func sec(cpu *Cpu2A03) {
	cpu.P |= FlagCarry
}

func cli(cpu *Cpu2A03) {
	cpu.P &= ^FlagInterruptDisable
}

func sei(cpu *Cpu2A03) {
	cpu.P |= FlagInterruptDisable
}

func clv(cpu *Cpu2A03) {
	cpu.P &= ^FlagOverflow
}

func cld(cpu *Cpu2A03) {
	cpu.P &= ^FlagDecimal
}

func sed(cpu *Cpu2A03) {
	cpu.P |= FlagDecimal
}
