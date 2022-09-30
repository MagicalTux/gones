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
	&op{"SLO", nil, amIndX},  // invalid
	&op{"NOP", nop, amZpg},   // invalid
	&op{"ORA", ora, amZpg},
	&op{"ASL", asl, amZpg},
	&op{"SLO", nil, amZpg},

	// 0x08
	&op{"PHP", php, amImpl},
	&op{"ORA", ora, amImmed},
	&op{"ASL", asl, amAcc},
	&op{"ANC", nil, amImmed},
	&op{"NOP", nop, amAbs},
	&op{"ORA", ora, amAbs},
	&op{"ASL", asl, amAbs},
	&op{"SLO", nil, amAbs},

	// 0x10
	&op{"BPL", bpl, amRel},
	&op{"ORA", ora, amIndY},
	&op{"STP", nil, amImpl},
	&op{"SLO", nil, amIndY},
	&op{"NOP", nop, amZpgX},
	&op{"ORA", ora, amZpgX},
	&op{"ASL", asl, amZpgX},
	&op{"SLO", nil, amZpgX},

	// 0x18
	&op{"CLC", clc, amImpl},
	nil,
	nil,
	nil,
	nil,
	nil,
	&op{"ASL", asl, amAbsX},
	nil,

	// 0x20
	&op{"JSR", jsr, amAbs},
	&op{"AND", and, amIndX},
	nil,
	nil,
	&op{"BIT", bit, amZpg},
	nil,
	&op{"ROL", rol, amZpg},
	nil,

	// 0x28
	&op{"PLP", plp, amImpl},
	&op{"AND", and, amImmed},
	&op{"ROL", rol, amAcc},
	nil,
	nil,
	nil,
	&op{"ROL", rol, amAbs},
	nil,

	// 0x30
	&op{"BMI", bmi, amRel},
	nil,
	nil,
	nil,
	nil,
	nil,
	&op{"ROL", rol, amZpgX},
	nil,

	// 0x38
	&op{"SEC", sec, amImpl},
	nil,
	nil,
	nil,
	nil,
	nil,
	&op{"ROL", rol, amAbsX},
	nil,

	// 0x40
	&op{"RTI", rti, amImpl},
	&op{"EOR", eor, amIndX},
	nil,
	nil,
	&op{"NOP", nop, amZpg},
	&op{"EOR", eor, amZpg},
	&op{"LSR", lsr, amZpg},
	nil,

	// 0x48
	&op{"PHA", pha, amImpl},
	&op{"EOR", eor, amImmed},
	&op{"LSR", lsr, amAcc},
	nil,
	&op{"JMP", jmp, amAbs},
	&op{"EOR", eor, amAbs},
	&op{"LSR", lsr, amAbs},
	nil,

	// 0x50
	&op{"BVC", bvc, amRel},
	&op{"EOR", eor, amIndY},
	nil,
	nil,
	nil,
	&op{"EOR", eor, amZpgX},
	&op{"LSR", lsr, amZpgX},
	nil,

	// 0x58
	nil,
	&op{"EOR", eor, amAbsY},
	nil,
	nil,
	nil,
	&op{"EOR", eor, amAbsX},
	&op{"LSR", lsr, amAbsX},
	nil,

	// 0x60
	&op{"RTS", rts, amImpl},
	&op{"ADC", adc, amIndX},
	nil,
	nil,
	nil,
	&op{"ADC", adc, amZpg},
	&op{"ROR", ror, amZpg},
	nil,

	// 0x68
	&op{"PLA", pla, amImpl},
	&op{"ADC", adc, amImmed},
	&op{"ROR", ror, amAcc},
	nil,
	nil,
	&op{"ADC", adc, amAbs},
	&op{"ROR", ror, amAbs},
	nil,

	// 0x70
	&op{"BVS", bvs, amRel},
	&op{"ADC", adc, amIndY},
	nil,
	nil,
	nil,
	&op{"ADC", adc, amZpgX},
	&op{"ROR", ror, amZpgX},
	nil,

	// 0x78
	&op{"SEI", sei, amImpl},
	&op{"ADC", adc, amAbsY},
	nil,
	nil,
	nil,
	&op{"ADC", adc, amAbsX},
	&op{"ROR", ror, amAbsX},
	nil,

	// 0x80
	nil,
	&op{"STA", sta, amIndX},
	nil,
	nil,
	&op{"STY", sty, amZpg},
	&op{"STA", sta, amZpg},
	&op{"STX", stx, amZpg},
	nil,

	// 0x88
	&op{"DEY", dey, amImpl},
	nil,
	&op{"TXA", txa, amImpl},
	nil,
	&op{"STY", sty, amAbs},
	&op{"STA", sta, amAbs},
	&op{"STX", stx, amAbs},
	nil,

	// 0x90
	&op{"BCC", bcc, amRel},
	&op{"STA", sta, amIndY},
	nil,
	nil,
	&op{"STY", sty, amZpgX},
	nil,
	&op{"STX", stx, amZpgY},
	nil,

	// 0x98
	&op{"TYA", tya, amImpl},
	nil,
	&op{"TXS", txs, amImpl},
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0xa0
	&op{"LDY", ldy, amImmed},
	&op{"LDA", lda, amIndX},
	&op{"LDX", ldx, amImmed},
	nil,
	&op{"LDY", ldy, amZpg},
	&op{"LDA", lda, amZpg},
	nil,
	nil,

	// 0xa8
	&op{"TAY", tay, amImpl},
	&op{"LDA", lda, amImmed},
	&op{"TAX", tax, amImpl},
	nil,
	&op{"LDY", ldy, amAbs},
	&op{"LDA", lda, amAbs},
	&op{"LDX", ldx, amAbs},
	nil,

	// 0xb0
	&op{"BCS", bcs, amRel},
	&op{"LDA", lda, amIndY},
	nil,
	nil,
	&op{"LDY", ldy, amZpgX},
	&op{"LDA", lda, amZpgX},
	nil,
	nil,

	// 0xb8
	&op{"CLV", clv, amImpl},
	nil,
	&op{"TSX", tsx, amImpl},
	nil,
	&op{"LDY", ldy, amAbsX},
	&op{"LDA", lda, amAbsX},
	nil,
	nil,

	// 0xc0
	&op{"CPY", cpy, amImmed},
	&op{"CMP", cmp, amIndX},
	nil,
	nil,
	nil,
	&op{"CMP", cmp, amZpg},
	&op{"DEC", dec, amZpg},
	nil,

	// 0xc8
	&op{"INY", iny, amImpl},
	&op{"CMP", cmp, amImmed},
	&op{"DEX", dex, amImpl},
	nil,
	nil,
	&op{"CMP", cmp, amAbs},
	&op{"DEC", dec, amAbs},
	nil,

	// 0xd0
	&op{"BNE", bne, amRel},
	&op{"CMP", cmp, amIndY},
	nil,
	nil,
	nil,
	&op{"CMP", cmp, amZpgX},
	&op{"DEC", dec, amZpgX},
	nil,

	// 0xd8
	&op{"CLD", cld, amImpl},
	&op{"CMP", cmp, amAbsY},
	nil,
	nil,
	&op{"CMP", cmp, amAbsX},
	&op{"CMP", cmp, amAbsX},
	&op{"DEC", dec, amAbsX},
	nil,

	// 0xe0
	&op{"CPX", cpx, amImmed},
	&op{"SBC", sbc, amIndX},
	nil,
	nil,
	nil,
	&op{"SBC", sbc, amZpg},
	nil,
	nil,

	// 0xe8
	&op{"INX", inx, amImpl},
	&op{"SBC", sbc, amImmed},
	&op{"NOP", nop, amImpl},
	nil,
	nil,
	&op{"SBC", sbc, amAbs},
	nil,
	nil,

	// 0xf0
	&op{"BEQ", beq, amRel},
	&op{"SBC", sbc, amIndY},
	nil,
	nil,
	nil,
	&op{"SBC", sbc, amZpgX},
	nil,
	nil,

	// 0xf8
	&op{"SED", sed, amImpl},
	&op{"SBC", sbc, amAbsY},
	nil,
	nil,
	nil,
	&op{"SBC", sbc, amAbsX},
	nil,
	nil,
}
