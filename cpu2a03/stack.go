package cpu2a03

/*
These instructions are implied mode, have a length of one byte and require machine cycles as indicated. The "PuLl" operations are known as "POP" on most other microprocessors. With the 6502, the stack is always on page one ($100-$1FF) and works top down.

MNEMONIC                        HEX TIM
TXS (Transfer X to Stack ptr)   $9A  2
TSX (Transfer Stack ptr to X)   $BA  2
PHA (PusH Accumulator)          $48  3
PLA (PuLl Accumulator)          $68  4
PHP (PusH Processor status)     $08  3
PLP (PuLl Processor status)     $28  4
*/

func txs(cpu *Cpu2A03, am AddressMode) {
	// typically, programs will start with "TXS $FF" to reset the stack
	am.Implied(cpu)
	cpu.S = cpu.X
}

func tsx(cpu *Cpu2A03, am AddressMode) {
	am.Implied(cpu)
	cpu.X = cpu.S
}

func pha(cpu *Cpu2A03, am AddressMode) {
	am.Implied(cpu)
	cpu.Push(cpu.A)
}

func pla(cpu *Cpu2A03, am AddressMode) {
	am.Implied(cpu)
	cpu.A = cpu.Pull()
}

func php(cpu *Cpu2A03, am AddressMode) {
	am.Implied(cpu)
	cpu.Push(cpu.P)
}

func plp(cpu *Cpu2A03, am AddressMode) {
	am.Implied(cpu)
	cpu.P = cpu.Pull()
}
