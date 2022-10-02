package cpu2a03

type op struct {
	i   string
	f   func(cpu *Cpu2A03, am AddressMode)
	am  AddressMode
	cyc byte // number of cycles for operation itself (addressmode may have something too)
}

var cpu2a03op = [256]*op{
	// 0x00
	&op{"BRK", brk, amImmed, 7}, // actually amImpl, but makes more sense like that
	&op{"ORA", ora, amIndX, 6},
	&op{"JAM", stop, amImpl, 2}, // invalid
	&op{"SLO", slo, amIndX, 8},  // invalid
	&op{"NOP", nop, amZpg, 3},   // invalid
	&op{"ORA", ora, amZpg, 3},
	&op{"ASL", asl, amZpg, 5},
	&op{"SLO", slo, amZpg, 5},

	// 0x08
	&op{"PHP", php, amImpl, 3},
	&op{"ORA", ora, amImmed, 2},
	&op{"ASL", asl, amAcc, 2},
	&op{"ANC", anc, amImmed, 2},
	&op{"NOP", nop, amAbs, 4},
	&op{"ORA", ora, amAbs, 4},
	&op{"ASL", asl, amAbs, 6},
	&op{"SLO", slo, amAbs, 6},

	// 0x10
	&op{"BPL", bpl, amRel, 2},
	&op{"ORA", ora, amIndY, 5},
	&op{"JAM", stop, amImpl, 2},
	&op{"SLO", slo, amIndY, 8},
	&op{"NOP", nop, amZpgX, 4},
	&op{"ORA", ora, amZpgX, 4},
	&op{"ASL", asl, amZpgX, 6},
	&op{"SLO", slo, amZpgX, 6},

	// 0x18
	&op{"CLC", clc, amImpl, 2},
	&op{"ORA", ora, amAbsY, 4},
	&op{"NOP", nop, amImpl, 2},
	&op{"SLO", slo, amAbsY, 7},
	&op{"NOP", nop, amAbsX, 4},
	&op{"ORA", ora, amAbsX, 4},
	&op{"ASL", asl, amAbsX, 7},
	&op{"SLO", slo, amAbsX, 7},

	// 0x20
	&op{"JSR", jsr, amAbs, 6},
	&op{"AND", and, amIndX, 6},
	&op{"JAM", stop, amImpl, 2},
	&op{"RLA", rla, amIndX, 8},
	&op{"BIT", bit, amZpg, 3},
	&op{"AND", and, amZpg, 3},
	&op{"ROL", rol, amZpg, 5},
	&op{"RLA", rla, amZpg, 5},

	// 0x28
	&op{"PLP", plp, amImpl, 4},
	&op{"AND", and, amImmed, 2},
	&op{"ROL", rol, amAcc, 2},
	&op{"ANC", anc, amImmed, 2},
	&op{"BIT", bit, amAbs, 4},
	&op{"AND", and, amAbs, 4},
	&op{"ROL", rol, amAbs, 6},
	&op{"RLA", rla, amAbs, 6},

	// 0x30
	&op{"BMI", bmi, amRel, 2},
	&op{"AND", and, amIndY, 5},
	&op{"JAM", stop, amImpl, 2},
	&op{"RLA", rla, amIndY, 8},
	&op{"NOP", nop, amZpgX, 4},
	&op{"AND", and, amZpgX, 4},
	&op{"ROL", rol, amZpgX, 6},
	&op{"RLA", rla, amZpgX, 6},

	// 0x38
	&op{"SEC", sec, amImpl, 2},
	&op{"AND", and, amAbsY, 4},
	&op{"NOP", nop, amImpl, 2},
	&op{"RLA", rla, amAbsY, 7},
	&op{"NOP", nop, amAbsX, 4},
	&op{"AND", and, amAbsX, 4},
	&op{"ROL", rol, amAbsX, 7},
	&op{"RLA", rla, amAbsX, 7},

	// 0x40
	&op{"RTI", rti, amImpl, 6},
	&op{"EOR", eor, amIndX, 6},
	&op{"JAM", stop, amImpl, 2},
	&op{"SRE", sre, amIndX, 8},
	&op{"NOP", nop, amZpg, 3},
	&op{"EOR", eor, amZpg, 3},
	&op{"LSR", lsr, amZpg, 5},
	&op{"SRE", sre, amZpg, 5},

	// 0x48
	&op{"PHA", pha, amImpl, 3},
	&op{"EOR", eor, amImmed, 2},
	&op{"LSR", lsr, amAcc, 2},
	&op{"ALR", alr, amImmed, 2},
	&op{"JMP", jmp, amAbs, 3},
	&op{"EOR", eor, amAbs, 4},
	&op{"LSR", lsr, amAbs, 6},
	&op{"SRE", sre, amAbs, 6},

	// 0x50
	&op{"BVC", bvc, amRel, 2},
	&op{"EOR", eor, amIndY, 5},
	&op{"JAM", stop, amImpl, 2},
	&op{"SRE", sre, amIndY, 8},
	&op{"NOP", nop, amZpgX, 4},
	&op{"EOR", eor, amZpgX, 4},
	&op{"LSR", lsr, amZpgX, 6},
	&op{"SRE", sre, amZpgX, 6},

	// 0x58
	&op{"CLI", cli, amImpl, 2},
	&op{"EOR", eor, amAbsY, 4},
	&op{"NOP", nop, amImpl, 2},
	&op{"SRE", sre, amAbsY, 7},
	&op{"NOP", nop, amAbsX, 4},
	&op{"EOR", eor, amAbsX, 4},
	&op{"LSR", lsr, amAbsX, 7},
	&op{"SRE", sre, amAbsX, 7},

	// 0x60
	&op{"RTS", rts, amImpl, 6},
	&op{"ADC", adc, amIndX, 6},
	&op{"JAM", stop, amImpl, 2},
	&op{"RRA", rra, amIndX, 8},
	&op{"NOP", nop, amZpg, 3},
	&op{"ADC", adc, amZpg, 3},
	&op{"ROR", ror, amZpg, 5},
	&op{"RRA", rra, amZpg, 5},

	// 0x68
	&op{"PLA", pla, amImpl, 4},
	&op{"ADC", adc, amImmed, 2},
	&op{"ROR", ror, amAcc, 2},
	&op{"ARR", arr, amImmed, 2},
	&op{"JMP", jmp, amInd, 5},
	&op{"ADC", adc, amAbs, 4},
	&op{"ROR", ror, amAbs, 6},
	&op{"RRA", rra, amAbs, 6},

	// 0x70
	&op{"BVS", bvs, amRel, 2},
	&op{"ADC", adc, amIndY, 5},
	&op{"JAM", stop, amImpl, 2},
	&op{"RRA", rra, amIndY, 8},
	&op{"NOP", nop, amZpgX, 4},
	&op{"ADC", adc, amZpgX, 4},
	&op{"ROR", ror, amZpgX, 6},
	&op{"RRA", rra, amZpgX, 6},

	// 0x78
	&op{"SEI", sei, amImpl, 2},
	&op{"ADC", adc, amAbsY, 4},
	&op{"NOP", nop, amImpl, 2},
	&op{"RRA", rra, amAbsY, 7},
	&op{"NOP", nop, amAbsX, 4},
	&op{"ADC", adc, amAbsX, 4},
	&op{"ROR", ror, amAbsX, 7},
	&op{"RRA", rra, amAbsX, 7},

	// 0x80
	&op{"NOP", nop, amImmed, 2},
	&op{"STA", sta, amIndX, 6},
	&op{"NOP", nop, amImmed, 2},
	&op{"SAX", sax, amIndX, 6},
	&op{"STY", sty, amZpg, 3},
	&op{"STA", sta, amZpg, 3},
	&op{"STX", stx, amZpg, 3},
	&op{"SAX", sax, amZpg, 3},

	// 0x88
	&op{"DEY", dey, amImpl, 2},
	&op{"NOP", nop, amImmed, 2},
	&op{"TXA", txa, amImpl, 2},
	&op{"ANE", ane, amImmed, 2},
	&op{"STY", sty, amAbs, 4},
	&op{"STA", sta, amAbs, 4},
	&op{"STX", stx, amAbs, 4},
	&op{"SAX", sax, amAbs, 4},

	// 0x90
	&op{"BCC", bcc, amRel, 2},
	&op{"STA", sta, amIndY, 6},
	&op{"JAM", stop, amImpl, 2},
	&op{"SHA", sha, amIndY, 6},
	&op{"STY", sty, amZpgX, 4},
	&op{"STA", sta, amZpgX, 4},
	&op{"STX", stx, amZpgY, 4},
	&op{"SAX", sax, amZpgY, 4},

	// 0x98
	&op{"TYA", tya, amImpl, 2},
	&op{"STA", sta, amAbsY, 5},
	&op{"TXS", txs, amImpl, 2},
	nil, // 5
	&op{"SHY", shy, amAbsX, 5},
	&op{"STA", sta, amAbsX, 5},
	nil, // 5
	&op{"SHA", sha, amAbsY, 5},

	// 0xa0
	&op{"LDY", ldy, amImmed, 2},
	&op{"LDA", lda, amIndX, 6},
	&op{"LDX", ldx, amImmed, 2},
	&op{"LAX", lax, amIndX, 6},
	&op{"LDY", ldy, amZpg, 3},
	&op{"LDA", lda, amZpg, 3},
	&op{"LDX", ldx, amZpg, 3},
	&op{"LAX", lax, amZpg, 3},

	// 0xa8
	&op{"TAY", tay, amImpl, 2},
	&op{"LDA", lda, amImmed, 2},
	&op{"TAX", tax, amImpl, 2},
	&op{"LXA", lxa, amImmed, 2},
	&op{"LDY", ldy, amAbs, 4},
	&op{"LDA", lda, amAbs, 4},
	&op{"LDX", ldx, amAbs, 4},
	&op{"LAX", lax, amAbs, 4},

	// 0xb0
	&op{"BCS", bcs, amRel, 2},
	&op{"LDA", lda, amIndY, 5},
	&op{"JAM", stop, amImpl, 2},
	&op{"LAX", lax, amIndY, 5},
	&op{"LDY", ldy, amZpgX, 4},
	&op{"LDA", lda, amZpgX, 4},
	&op{"LDX", ldx, amZpgY, 4},
	&op{"LAX", lax, amZpgY, 4},

	// 0xb8
	&op{"CLV", clv, amImpl, 2},
	&op{"LDA", lda, amAbsY, 4},
	&op{"TSX", tsx, amImpl, 2},
	&op{"LAS", las, amAbsY, 2},
	&op{"LDY", ldy, amAbsX, 4},
	&op{"LDA", lda, amAbsX, 4},
	&op{"LDX", ldx, amAbsY, 4},
	&op{"LAX", lax, amAbsY, 4},

	// 0xc0
	&op{"CPY", cpy, amImmed, 2},
	&op{"CMP", cmp, amIndX, 6},
	&op{"NOP", nop, amImmed, 2},
	&op{"DCP", dcp, amIndX, 8},
	&op{"CPY", cpy, amZpg, 3},
	&op{"CMP", cmp, amZpg, 3},
	&op{"DEC", dec, amZpg, 5},
	&op{"DCP", dcp, amZpg, 5},

	// 0xc8
	&op{"INY", iny, amImpl, 2},
	&op{"CMP", cmp, amImmed, 2},
	&op{"DEX", dex, amImpl, 2},
	&op{"SBX", sbx, amImmed, 2},
	&op{"CPY", cpy, amAbs, 4},
	&op{"CMP", cmp, amAbs, 4},
	&op{"DEC", dec, amAbs, 6},
	&op{"DCP", dcp, amAbs, 6},

	// 0xd0
	&op{"BNE", bne, amRel, 2},
	&op{"CMP", cmp, amIndY, 5},
	&op{"JAM", stop, amImpl, 2},
	&op{"DCP", dcp, amIndY, 8},
	&op{"NOP", nop, amZpgX, 4},
	&op{"CMP", cmp, amZpgX, 4},
	&op{"DEC", dec, amZpgX, 6},
	&op{"DCP", dcp, amZpgX, 6},

	// 0xd8
	&op{"CLD", cld, amImpl, 2},
	&op{"CMP", cmp, amAbsY, 4},
	&op{"NOP", nop, amImpl, 2},
	&op{"DCP", dcp, amAbsY, 7},
	&op{"NOP", nop, amAbsX, 4},
	&op{"CMP", cmp, amAbsX, 4},
	&op{"DEC", dec, amAbsX, 7},
	&op{"DCP", dcp, amAbsX, 7},

	// 0xe0
	&op{"CPX", cpx, amImmed, 2},
	&op{"SBC", sbc, amIndX, 6},
	&op{"NOP", nop, amImmed, 2},
	&op{"ISC", isc, amIndX, 8},
	&op{"CPX", cpx, amZpg, 3},
	&op{"SBC", sbc, amZpg, 3},
	&op{"INC", inc, amZpg, 5},
	&op{"ISC", isc, amZpg, 5},

	// 0xe8
	&op{"INX", inx, amImpl, 2},
	&op{"SBC", sbc, amImmed, 2},
	&op{"NOP", nop, amImpl, 2},
	&op{"USBC", sbc, amImmed, 2},
	&op{"CPX", cpx, amAbs, 4},
	&op{"SBC", sbc, amAbs, 4},
	&op{"INC", inc, amAbs, 6},
	&op{"ISC", isc, amAbs, 6},

	// 0xf0
	&op{"BEQ", beq, amRel, 2},
	&op{"SBC", sbc, amIndY, 5},
	&op{"JAM", stop, amImpl, 2},
	&op{"ISC", isc, amIndY, 8},
	&op{"NOP", nop, amZpgX, 4},
	&op{"SBC", sbc, amZpgX, 4},
	&op{"INC", inc, amZpgX, 6},
	&op{"ISC", isc, amZpgX, 6},

	// 0xf8
	&op{"SED", sed, amImpl, 2},
	&op{"SBC", sbc, amAbsY, 4},
	&op{"NOP", nop, amImpl, 2},
	&op{"ISC", isc, amAbsY, 7},
	&op{"NOP", nop, amAbsX, 4},
	&op{"SBC", sbc, amAbsX, 4},
	&op{"INC", inc, amAbsX, 7},
	&op{"ISC", isc, amAbsX, 7},
}
