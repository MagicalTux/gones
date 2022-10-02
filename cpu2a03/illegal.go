package cpu2a03

func slo(cpu *Cpu2A03, am AddressMode) {
	// ASL oper + ORA oper
	// M = C <- [76543210] <- 0, A OR M -> A
	//Flags: N Z C

	addr := am.AddrFast(cpu)
	v := cpu.Memory.MemRead(addr) // input M or 0 ?

	cpu.setFlag(FlagCarry, v&0x80 == 0x80)
	v <<= 1
	cpu.Memory.MemWrite(addr, v)
	cpu.flagsNZ(v)

	// ORA
	//v := am.Read(cpu)

	cpu.A |= v
	cpu.flagsNZ(cpu.A)
}

func rla(cpu *Cpu2A03, am AddressMode) {
	// ROL oper + AND oper
	// M = C <- [76543210] <- C, A AND M -> A
	// Flags: N Z C

	// ROL
	c := byte(0)
	if cpu.getFlag(FlagCarry) {
		c = 1
	}

	addr := am.AddrFast(cpu)
	v := cpu.Memory.MemRead(addr)

	cpu.setFlag(FlagCarry, v&0x80 == 0x80)
	v = (v << 1) | c
	cpu.Memory.MemWrite(addr, v)

	cpu.A &= v
	cpu.flagsNZ(cpu.A)
}

func sre(cpu *Cpu2A03, am AddressMode) {
	// LSR oper + EOR oper
	// M = 0 -> [76543210] -> C, A EOR M -> A
	// Flags: N Z C

	addr := am.AddrFast(cpu)
	v := cpu.Memory.MemRead(addr)

	cpu.setFlag(FlagCarry, v&1 == 1)
	v >>= 1
	cpu.Memory.MemWrite(addr, v)

	cpu.A ^= v
	cpu.flagsNZ(cpu.A)
}

func rra(cpu *Cpu2A03, am AddressMode) {
	// ROR oper + ADC oper
	// M = C -> [76543210] -> C, A + M + C -> A, C
	// Flags: N Z C V
	c := byte(0)
	if cpu.getFlag(FlagCarry) {
		c = 1
	}

	// act on mem
	addr := am.AddrFast(cpu)
	v := cpu.Memory.MemRead(addr)

	c2 := v & 1
	cpu.setFlag(FlagCarry, v&1 == 1)
	v = (v >> 1) | (c << 7)
	cpu.Memory.MemWrite(addr, v)

	// ADC
	c = c2
	a := cpu.A
	cpu.A = a + v + c

	cpu.flagsNZ(cpu.A)

	cpu.setFlag(FlagCarry, int(a)+int(v)+int(c) > 0xff)
	cpu.setFlag(FlagOverflow, (a^v)&0x80 == 0 && (a^cpu.A)&0x80 != 0)
}

func anc(cpu *Cpu2A03, am AddressMode) {
	// AND oper + set C as ASL
	// A AND oper, bit(7) -> C
	// Flags: N Z C

	// ANC    ***
	// ANC ANDs the contents of the A register with an immediate value and then
	// moves bit 7 of A into the Carry flag.  This opcode works basically
	// identically to AND #immed. except that the Carry flag is set to the same
	// state that the Negative flag is set to.

	cpu.A &= am.Read(cpu)
	cpu.setFlag(FlagCarry, cpu.A&0x80 == 0x80)
	cpu.flagsNZ(cpu.A)
}

func alr(cpu *Cpu2A03, am AddressMode) {
	// AND oper + LSR
	// A AND oper, 0 -> [76543210] -> C
	// Flags: N Z C

	// This opcode ANDs the contents of the A register with an immediate value and
	// then LSRs the result.

	// Equivalent instructions:
	// AND #$FE
	// LSR A

	// AND
	cpu.A &= am.Read(cpu)

	// LSR

	cpu.setFlag(FlagCarry, cpu.A&1 == 1)
	cpu.A >>= 1
	//cpu.Memory.MemWrite(addr, v) // can't set back value as it was an immed value
	cpu.flagsNZ(cpu.A)
}

func arr(cpu *Cpu2A03, am AddressMode) {
	// AND oper + ROR
	// A AND oper, C -> [76543210] -> C
	// Flags: N Z C V

	// This operation involves the adder:
	// V-flag is set according to (A AND oper) + oper
	// The carry is not set, but bit 7 (sign) is exchanged with the carry

	// The opcode ARR operates more complexily than actually described in the list
	// above.  Here is a brief rundown on this.  The following assumes the decimal
	// flag is clear.  You see, the sub-instruction for ARR ($6B) is in fact ADC
	// ($69), not AND.  While ADC is not performed, some of the ADC mechanics are
	// evident.  Like ADC, ARR affects the overflow flag.  The following effects
	// occur after ANDing but before RORing.  The V flag is set to the result of
	// exclusive ORing bit 7 with bit 6.  Unlike ROR, bit 0 does not go into the
	// carry flag.  The state of bit 7 is exchanged with the carry flag.  Bit 0 is
	// lost.  All of this may appear strange, but it makes sense if you consider
	// the probable internal operations of ADC itself.

	v := am.Read(cpu)
	cpu.A &= v

	// something like that?
	cpu.setFlag(FlagOverflow, ((cpu.A>>7&1)^(cpu.A>>6&1)) == 1)

	c := byte(0)
	if cpu.getFlag(FlagCarry) {
		c = 1
	}

	// Update A
	cpu.setFlag(FlagCarry, cpu.A&0x80 == 0x80) // "The state of bit 7 is exchanged with the carry flag." ?
	cpu.A = (cpu.A >> 1) | (c << 7)
	cpu.flagsNZ(cpu.A)
}

func ane(cpu *Cpu2A03, am AddressMode) {
	// * AND X + AND oper
	// Highly unstable, do not use.
	// A base value in A is determined based on the contets of A and a constant, which may be typically $00, $ff, $ee, etc. The value of this constant depends on temerature, the chip series, and maybe other factors, as well.
	// In order to eliminate these uncertaincies from the equation, use either 0 as the operand or a value of $FF in the accumulator.
	// (A OR CONST) AND X AND oper -> A

	// Flags: N Z

	// I chose $ee as the constant because it feels right
	cpu.A = cpu.A & 0xee & am.Read(cpu)
	cpu.flagsNZ(cpu.A)
}

func lxa(cpu *Cpu2A03, am AddressMode) {
	// Opcode AB
	// also known as OAL or ATX
	// Store * AND oper in A and X
	// Highly unstable, involves a 'magic' constant, see ANE
	// (A OR CONST) AND oper -> A -> X

	// This opcode ORs the A register with #xx, ANDs the result with an immediate
	// value, and then stores the result in both A and X.

	// ORA #$EE
	// AND #$AA
	// TAX

	// Flags: N Z

	// http://visual6502.org/JSSim/expert.html?graphics=f&steps=12&a=5555&d=44&a=0&d=a9ffab88&loglevel=2&logmore=dpc3_SBX,dpc23_SBAC,plaOutputs,DPControl

	// Let's use the same constant as for ANE
	cpu.A = cpu.A & am.Read(cpu)
	cpu.X = cpu.A
	cpu.flagsNZ(cpu.A)
}

func las(cpu *Cpu2A03, am AddressMode) {
	// LDA/TSX oper
	// M AND SP -> A, X, SP
	// Flags: N Z

	// This opcode ANDs the contents of a memory location with the contents of the
	// stack pointer register and stores the result in the accumulator, the X
	// register, and the stack pointer.  Affected flags: N Z.

	v := am.Read(cpu)
	v &= cpu.S

	cpu.A = v
	cpu.X = v
	cpu.S = v
	cpu.Memory.MemWrite(0x100+uint16(cpu.S), v)
	cpu.flagsNZ(v)
}

func sha(cpu *Cpu2A03, am AddressMode) {
	// Stores A AND X AND (high-byte of addr. + 1) at addr.
	// unstable: sometimes 'AND (H+1)' is dropped, page boundary crossings may not work (with the high-byte of the value used as the high-byte of the address)
	// A AND X AND (H+1) -> M
	// Flags: none

	addr := am.AddrFast(cpu)
	v := cpu.A & cpu.X & uint8(addr>>8)
	cpu.Memory.MemWrite(addr, v)
}

func sbx(cpu *Cpu2A03, am AddressMode) {
	// CMP and DEX at once, sets flags like CMP
	// (A AND X) - oper -> X

	// SAX ANDs the contents of the A and X registers (leaving the contents of A
	// intact), subtracts an immediate value, and then stores the result in X.
	// ... A few points might be made about the action of subtracting an immediate
	// value.  It actually works just like the CMP instruction, except that CMP
	// does not store the result of the subtraction it performs in any register.
	// This subtract operation is not affected by the state of the Carry flag,
	// though it does affect the Carry flag.  It does not affect the Overflow
	// flag.

	a := cpu.A & cpu.X
	v := am.Read(cpu)
	cpu.X = a - v

	cpu.compare(a, v)
}

func shy(cpu *Cpu2A03, am AddressMode) {
	// Opcode 9C
	// Also known as: A11, SYA, SAY
	// Stores Y AND (high-byte of addr. + 1) at addr.
	// unstable: sometimes 'AND (H+1)' is dropped, page boundary crossings may not work (with the high-byte of the value used as the high-byte of the address)
	// Y AND (H+1) -> M

	addr := am.AddrFast(cpu)
	cpu.Memory.MemWrite(addr, cpu.Y&uint8((addr>>8)+1))
}

func tas(cpu *Cpu2A03, am AddressMode) {
	// Opcode 9B
	// Also known as XAS, SHS
	// Puts A AND X in SP and stores A AND X AND (high-byte of addr. + 1) at addr.
	// unstable: sometimes 'AND (H+1)' is dropped, page boundary crossings may not work (with the high-byte of the value used as the high-byte of the address)
	// A AND X -> SP, A AND X AND (H+1) -> M

	addr := am.AddrFast(cpu)
	cpu.S = cpu.A & cpu.X
	cpu.Memory.MemWrite(addr, cpu.A&cpu.X&uint8((addr>>8)+1))
}

func shx(cpu *Cpu2A03, am AddressMode) {
	// Opcode 9E
	// Also known as A11, SXA, XAS
	// Stores X AND (high-byte of addr. + 1) at addr.
	// unstable: sometimes 'AND (H+1)' is dropped, page boundary crossings may not work (with the high-byte of the value used as the high-byte of the address)
	// X AND (H+1) -> M

	addr := am.AddrFast(cpu)
	cpu.Memory.MemWrite(addr, cpu.X&uint8((addr>>8)+1))
}
