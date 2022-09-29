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

func clc(cpu *Cpu2A03, am AddressMode) {
	cpu.P &= ^FlagCarry
}

func sec(cpu *Cpu2A03, am AddressMode) {
	cpu.P |= FlagCarry
}

func cli(cpu *Cpu2A03, am AddressMode) {
	cpu.P &= ^FlagInterruptDisable
}

func sei(cpu *Cpu2A03, am AddressMode) {
	cpu.P |= FlagInterruptDisable
}

func clv(cpu *Cpu2A03, am AddressMode) {
	cpu.P &= ^FlagOverflow
}

func cld(cpu *Cpu2A03, am AddressMode) {
	cpu.P &= ^FlagDecimal
}

func sed(cpu *Cpu2A03, am AddressMode) {
	cpu.P |= FlagDecimal
}
