package cpu6502

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

func clc(cpu *CPU, am AddressMode) {
	am.Implied(cpu)
	cpu.P &= ^FlagCarry
}

func sec(cpu *CPU, am AddressMode) {
	am.Implied(cpu)
	cpu.P |= FlagCarry
}

func cli(cpu *CPU, am AddressMode) {
	am.Implied(cpu)
	cpu.P &= ^FlagInterruptDisable
}

func sei(cpu *CPU, am AddressMode) {
	am.Implied(cpu)
	cpu.P |= FlagInterruptDisable
}

func clv(cpu *CPU, am AddressMode) {
	am.Implied(cpu)
	cpu.P &= ^FlagOverflow
}

func cld(cpu *CPU, am AddressMode) {
	am.Implied(cpu)
	cpu.P &= ^FlagDecimal
}

func sed(cpu *CPU, am AddressMode) {
	am.Implied(cpu)
	cpu.P |= FlagDecimal
}
