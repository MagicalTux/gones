package cpu2a03

func ora(cpu *Cpu2A03, am AddressMode) {
	// Affects Flags: N Z
	v := am.Read(cpu)

	cpu.A |= v
	cpu.flagsNZ(cpu.A)
}

func bit(cpu *Cpu2A03, am AddressMode) {
	v := am.Read(cpu)
	// bits 7 and 6 of operand are transfered to bit 7 and 6 of SR (N,V);
	// the zero-flag is set to the result of operand AND accumulator.

	// A AND M, M7 -> N, M6 -> V

	cpu.setFlag(FlagOverflow, (v>>6)&1 == 1) // V
	cpu.flagsZ(v & cpu.A)
	cpu.flagsN(v)
}

func and(cpu *Cpu2A03, am AddressMode) {
	v := am.Read(cpu)

	cpu.A &= v
	cpu.flagsNZ(cpu.A)
}

func eor(cpu *Cpu2A03, am AddressMode) {
	v := am.Read(cpu)

	cpu.A ^= v
	cpu.flagsNZ(cpu.A)
}

func ror(cpu *Cpu2A03, am AddressMode) {
	c := byte(0)
	if cpu.getFlag(FlagCarry) {
		c = 1
	}

	if am == amAcc {
		// act on cpu.A
		cpu.setFlag(FlagCarry, cpu.A&1 == 1)
		cpu.A = (cpu.A >> 1) | (c << 7)
		cpu.flagsNZ(cpu.A)
	} else {
		// act on mem
		addr := am.Addr(cpu)
		v := cpu.Memory.MemRead(addr)

		cpu.setFlag(FlagCarry, v&1 == 1)
		v = (v >> 1) | (c << 7)
		cpu.Memory.MemWrite(addr, v)
		cpu.flagsNZ(v)
	}
}

func lsr(cpu *Cpu2A03, am AddressMode) {
	if am == amAcc {
		cpu.setFlag(FlagCarry, cpu.A&1 == 1)
		cpu.A >>= 1
		cpu.flagsNZ(cpu.A)
	} else {
		addr := am.Addr(cpu)
		v := cpu.Memory.MemRead(addr)

		cpu.setFlag(FlagCarry, v&1 == 1)
		v >>= 1
		cpu.Memory.MemWrite(addr, v)
		cpu.flagsNZ(v)
	}
}

func asl(cpu *Cpu2A03, am AddressMode) {
	if am == amAcc {
		cpu.setFlag(FlagCarry, cpu.A&0x80 == 0x80)
		cpu.A <<= 1
		cpu.flagsNZ(cpu.A)
	} else {
		addr := am.Addr(cpu)
		v := cpu.Memory.MemRead(addr)

		cpu.setFlag(FlagCarry, v&0x80 == 0x80)
		v <<= 1
		cpu.Memory.MemWrite(addr, v)
		cpu.flagsNZ(v)
	}
}

func rol(cpu *Cpu2A03, am AddressMode) {
	c := byte(0)
	if cpu.getFlag(FlagCarry) {
		c = 1
	}

	if am == amAcc {
		// act on cpu.A
		cpu.setFlag(FlagCarry, cpu.A&0x80 == 0x80)
		cpu.A = cpu.A<<1 | c
		cpu.flagsNZ(cpu.A)
	} else {
		// act on mem
		addr := am.Addr(cpu)
		v := cpu.Memory.MemRead(addr)

		cpu.setFlag(FlagCarry, v&0x80 == 0x80)
		v = (v << 1) | c
		cpu.Memory.MemWrite(addr, v)
		cpu.flagsNZ(v)
	}
}

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
	v := am.Read(cpu) & cpu.A
	cpu.setFlag(FlagCarry, v&0x80 == 0x80)
	cpu.flagsNZ(v)
}

func alr(cpu *Cpu2A03, am AddressMode) {
	// AND oper + LSR
	// A AND oper, 0 -> [76543210] -> C
	// Flags: N Z C

	// AND
	addr := am.Addr(cpu)
	v := cpu.Memory.MemRead(addr)

	cpu.A &= v

	// LSR

	cpu.setFlag(FlagCarry, v&1 == 1)
	v >>= 1
	cpu.Memory.MemWrite(addr, v)
	cpu.flagsNZ(v)
}

func arr(cpu *Cpu2A03, am AddressMode) {
	// AND oper + ROR
	// A AND oper, C -> [76543210] -> C
	// Flags: N Z C V

	// TODO Fix this
	// This operation involves the adder:
	// V-flag is set according to (A AND oper) + oper
	// The carry is not set, but bit 7 (sign) is exchanged with the carry

	addr := am.Addr(cpu)
	v := cpu.Memory.MemRead(addr)

	cpu.A &= v

	c := byte(0)
	if cpu.getFlag(FlagCarry) {
		c = 1
	}

	// act on mem

	cpu.setFlag(FlagCarry, v&1 == 1)
	v = (v >> 1) | (c << 7)
	cpu.Memory.MemWrite(addr, v)
	cpu.flagsNZ(v)
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

func sha(cpu *Cpu2A03, am AddressMode) {
	// Stores A AND X AND (high-byte of addr. + 1) at addr.
	// unstable: sometimes 'AND (H+1)' is dropped, page boundary crossings may not work (with the high-byte of the value used as the high-byte of the address)
	// A AND X AND (H+1) -> M
	// Flags: none

	addr := am.Addr(cpu)
	v := cpu.A & cpu.X & uint8(addr>>8)
	cpu.Memory.MemWrite(addr, v)
}
