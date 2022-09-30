package cpu2a03

type op struct {
	i  string
	f  func(cpu *Cpu2A03, am AddressMode)
	am AddressMode
}

var cpu2a03op = [256]*op{
	// 0x00
	&op{"BRK", brk, amImmed}, // actually amImpl, but makes more sense like that
	&op{"ORA", ora, amIndX},
	&op{"STP", stop, amImpl}, // invalid
	&op{"SLO", slo, amIndX},  // invalid
	&op{"NOP", nop, amZpg},   // invalid
	&op{"ORA", ora, amZpg},
	&op{"ASL", asl, amZpg},
	&op{"SLO", slo, amZpg},

	// 0x08
	&op{"PHP", php, amImpl},
	&op{"ORA", ora, amImmed},
	&op{"ASL", asl, amAcc},
	&op{"ANC", nil, amImmed},
	&op{"NOP", nop, amAbs},
	&op{"ORA", ora, amAbs},
	&op{"ASL", asl, amAbs},
	&op{"SLO", slo, amAbs},

	// 0x10
	&op{"BPL", bpl, amRel},
	&op{"ORA", ora, amIndY},
	&op{"STP", nil, amImpl},
	&op{"SLO", slo, amIndY},
	&op{"NOP", nop, amZpgX},
	&op{"ORA", ora, amZpgX},
	&op{"ASL", asl, amZpgX},
	&op{"SLO", slo, amZpgX},

	// 0x18
	&op{"CLC", clc, amImpl},
	&op{"ORA", ora, amAbsY},
	&op{"NOP", nop, amImpl},
	&op{"SLO", slo, amAbsY},
	&op{"NOP", nop, amAbsX},
	&op{"ORA", ora, amAbsX},
	&op{"ASL", asl, amAbsX},
	&op{"SLO", slo, amAbsX},

	// 0x20
	&op{"JSR", jsr, amAbs},
	&op{"AND", and, amIndX},
	nil,
	&op{"RLA", rla, amIndX},
	&op{"BIT", bit, amZpg},
	&op{"AND", and, amZpg},
	&op{"ROL", rol, amZpg},
	&op{"RLA", rla, amZpg},

	// 0x28
	&op{"PLP", plp, amImpl},
	&op{"AND", and, amImmed},
	&op{"ROL", rol, amAcc},
	nil,
	&op{"BIT", bit, amAbs},
	&op{"AND", and, amAbs},
	&op{"ROL", rol, amAbs},
	&op{"RLA", rla, amAbs},

	// 0x30
	&op{"BMI", bmi, amRel},
	&op{"AND", and, amIndY},
	nil,
	&op{"RLA", rla, amIndY},
	&op{"NOP", nop, amZpgX},
	&op{"AND", and, amZpgX},
	&op{"ROL", rol, amZpgX},
	&op{"RLA", rla, amZpgX},

	// 0x38
	&op{"SEC", sec, amImpl},
	&op{"AND", and, amAbsY},
	&op{"NOP", nop, amImpl},
	&op{"RLA", rla, amAbsY},
	&op{"NOP", nop, amAbsX},
	&op{"AND", and, amAbsX},
	&op{"ROL", rol, amAbsX},
	&op{"RLA", rla, amAbsX},

	// 0x40
	&op{"RTI", rti, amImpl},
	&op{"EOR", eor, amIndX},
	nil,
	&op{"SRE", sre, amIndX},
	&op{"NOP", nop, amZpg},
	&op{"EOR", eor, amZpg},
	&op{"LSR", lsr, amZpg},
	&op{"SRE", sre, amZpg},

	// 0x48
	&op{"PHA", pha, amImpl},
	&op{"EOR", eor, amImmed},
	&op{"LSR", lsr, amAcc},
	nil,
	&op{"JMP", jmp, amAbs},
	&op{"EOR", eor, amAbs},
	&op{"LSR", lsr, amAbs},
	&op{"SRE", sre, amAbs},

	// 0x50
	&op{"BVC", bvc, amRel},
	&op{"EOR", eor, amIndY},
	nil,
	&op{"SRE", sre, amIndY},
	&op{"NOP", nop, amZpgX},
	&op{"EOR", eor, amZpgX},
	&op{"LSR", lsr, amZpgX},
	&op{"SRE", sre, amZpgX},

	// 0x58
	nil,
	&op{"EOR", eor, amAbsY},
	&op{"NOP", nop, amImpl},
	&op{"SRE", sre, amAbsY},
	&op{"NOP", nop, amAbsX},
	&op{"EOR", eor, amAbsX},
	&op{"LSR", lsr, amAbsX},
	&op{"SRE", sre, amAbsX},

	// 0x60
	&op{"RTS", rts, amImpl},
	&op{"ADC", adc, amIndX},
	nil,
	&op{"RRA", rra, amIndX},
	&op{"NOP", nop, amZpg},
	&op{"ADC", adc, amZpg},
	&op{"ROR", ror, amZpg},
	&op{"RRA", rra, amZpg},

	// 0x68
	&op{"PLA", pla, amImpl},
	&op{"ADC", adc, amImmed},
	&op{"ROR", ror, amAcc},
	nil,
	&op{"JMP", jmp, amInd},
	&op{"ADC", adc, amAbs},
	&op{"ROR", ror, amAbs},
	&op{"RRA", rra, amAbs},

	// 0x70
	&op{"BVS", bvs, amRel},
	&op{"ADC", adc, amIndY},
	nil,
	&op{"RRA", rra, amIndY},
	&op{"NOP", nop, amZpgX},
	&op{"ADC", adc, amZpgX},
	&op{"ROR", ror, amZpgX},
	&op{"RRA", rra, amZpgX},

	// 0x78
	&op{"SEI", sei, amImpl},
	&op{"ADC", adc, amAbsY},
	&op{"NOP", nop, amImpl},
	&op{"RRA", rra, amAbsY},
	&op{"NOP", nop, amAbsX},
	&op{"ADC", adc, amAbsX},
	&op{"ROR", ror, amAbsX},
	&op{"RRA", rra, amAbsX},

	// 0x80
	&op{"NOP", nop, amImmed},
	&op{"STA", sta, amIndX},
	nil,
	&op{"SAX", sax, amIndX},
	&op{"STY", sty, amZpg},
	&op{"STA", sta, amZpg},
	&op{"STX", stx, amZpg},
	&op{"SAX", sax, amZpg},

	// 0x88
	&op{"DEY", dey, amImpl},
	nil,
	&op{"TXA", txa, amImpl},
	nil,
	&op{"STY", sty, amAbs},
	&op{"STA", sta, amAbs},
	&op{"STX", stx, amAbs},
	&op{"SAX", sax, amAbs},

	// 0x90
	&op{"BCC", bcc, amRel},
	&op{"STA", sta, amIndY},
	nil,
	nil,
	&op{"STY", sty, amZpgX},
	&op{"STA", sta, amZpgX},
	&op{"STX", stx, amZpgY},
	&op{"SAX", sax, amZpgY},

	// 0x98
	&op{"TYA", tya, amImpl},
	&op{"STA", sta, amAbsY},
	&op{"TXS", txs, amImpl},
	nil,
	nil,
	&op{"STA", sta, amAbsX},
	nil,
	nil,

	// 0xa0
	&op{"LDY", ldy, amImmed},
	&op{"LDA", lda, amIndX},
	&op{"LDX", ldx, amImmed},
	&op{"LAX", lax, amIndX},
	&op{"LDY", ldy, amZpg},
	&op{"LDA", lda, amZpg},
	&op{"LDX", ldx, amZpg},
	&op{"LAX", lax, amZpg},

	// 0xa8
	&op{"TAY", tay, amImpl},
	&op{"LDA", lda, amImmed},
	&op{"TAX", tax, amImpl},
	nil,
	&op{"LDY", ldy, amAbs},
	&op{"LDA", lda, amAbs},
	&op{"LDX", ldx, amAbs},
	&op{"LAX", lax, amAbs},

	// 0xb0
	&op{"BCS", bcs, amRel},
	&op{"LDA", lda, amIndY},
	nil,
	&op{"LAX", lax, amIndY},
	&op{"LDY", ldy, amZpgX},
	&op{"LDA", lda, amZpgX},
	&op{"LDX", ldx, amZpgY},
	&op{"LAX", lax, amZpgY},

	// 0xb8
	&op{"CLV", clv, amImpl},
	&op{"LDA", lda, amAbsY},
	&op{"TSX", tsx, amImpl},
	nil,
	&op{"LDY", ldy, amAbsX},
	&op{"LDA", lda, amAbsX},
	&op{"LDX", ldx, amAbsY},
	&op{"LAX", lax, amAbsY},

	// 0xc0
	&op{"CPY", cpy, amImmed},
	&op{"CMP", cmp, amIndX},
	nil,
	&op{"DCP", dcp, amIndX},
	&op{"CPY", cpy, amZpg},
	&op{"CMP", cmp, amZpg},
	&op{"DEC", dec, amZpg},
	&op{"DCP", dcp, amZpg},

	// 0xc8
	&op{"INY", iny, amImpl},
	&op{"CMP", cmp, amImmed},
	&op{"DEX", dex, amImpl},
	nil,
	&op{"CPY", cpy, amAbs},
	&op{"CMP", cmp, amAbs},
	&op{"DEC", dec, amAbs},
	&op{"DCP", dcp, amAbs},

	// 0xd0
	&op{"BNE", bne, amRel},
	&op{"CMP", cmp, amIndY},
	nil,
	&op{"DCP", dcp, amIndY},
	&op{"NOP", nop, amZpgX},
	&op{"CMP", cmp, amZpgX},
	&op{"DEC", dec, amZpgX},
	&op{"DCP", dcp, amZpgX},

	// 0xd8
	&op{"CLD", cld, amImpl},
	&op{"CMP", cmp, amAbsY},
	&op{"NOP", nop, amImpl},
	&op{"DCP", dcp, amAbsY},
	&op{"CMP", cmp, amAbsX},
	&op{"CMP", cmp, amAbsX},
	&op{"DEC", dec, amAbsX},
	&op{"DCP", dcp, amAbsX},

	// 0xe0
	&op{"CPX", cpx, amImmed},
	&op{"SBC", sbc, amIndX},
	nil,
	&op{"ISC", isc, amIndX},
	&op{"CPX", cpx, amZpg},
	&op{"SBC", sbc, amZpg},
	&op{"INC", inc, amZpg},
	&op{"ISC", isc, amZpg},

	// 0xe8
	&op{"INX", inx, amImpl},
	&op{"SBC", sbc, amImmed},
	&op{"NOP", nop, amImpl},
	&op{"USBC", sbc, amImmed},
	&op{"CPX", cpx, amAbs},
	&op{"SBC", sbc, amAbs},
	&op{"INC", inc, amAbs},
	&op{"ISC", isc, amAbs},

	// 0xf0
	&op{"BEQ", beq, amRel},
	&op{"SBC", sbc, amIndY},
	nil,
	&op{"ISC", isc, amIndY},
	&op{"NOP", nop, amZpgX},
	&op{"SBC", sbc, amZpgX},
	&op{"INC", inc, amZpgX},
	&op{"ISC", isc, amZpgX},

	// 0xf8
	&op{"SED", sed, amImpl},
	&op{"SBC", sbc, amAbsY},
	&op{"NOP", nop, amImpl},
	&op{"ISC", isc, amAbsY},
	&op{"NOP", nop, amAbsX},
	&op{"SBC", sbc, amAbsX},
	&op{"INC", inc, amAbsX},
	&op{"ISC", isc, amAbsX},
}
