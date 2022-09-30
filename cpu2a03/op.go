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
	&op{"ASL", nil, amZpg},
	&op{"SLO", nil, amZpg},

	// 0x08
	&op{"PHP", php, amImpl},
	&op{"ORA", ora, amImmed},
	&op{"ASL", nil, amAcc},
	&op{"ANC", nil, amImmed},
	&op{"NOP", nop, amAbs},
	&op{"ORA", ora, amAbs},
	&op{"ASL", nil, amAbs},
	&op{"SLO", nil, amAbs},

	// 0x10
	&op{"BPL", bpl, amRel},
	&op{"ORA", ora, amIndY},
	&op{"STP", nil, amImpl},
	&op{"SLO", nil, amIndY},
	&op{"NOP", nop, amZpgX},
	&op{"ORA", ora, amZpgX},
	&op{"ASL", nil, amZpgX},
	&op{"SLO", nil, amZpgX},

	// 0x18
	&op{"CLC", clc, amImpl},
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x20
	&op{"JSR", jsr, amAbs},
	nil,
	nil,
	nil,
	&op{"BIT", bit, amZpg},
	nil,
	nil,
	nil,

	// 0x28
	&op{"PLP", plp, amImpl},
	&op{"AND", and, amImmed},
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x30
	&op{"BMI", bmi, amRel},
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x38
	&op{"SEC", sec, amImpl},
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x40
	nil,
	&op{"EOR", eor, amIndX},
	nil,
	nil,
	nil,
	&op{"EOR", eor, amZpg},
	nil,
	nil,

	// 0x48
	&op{"PHA", pha, amImpl},
	&op{"EOR", eor, amImmed},
	nil,
	nil,
	&op{"JMP", jmp, amAbs},
	&op{"EOR", eor, amAbs},
	nil,
	nil,

	// 0x50
	&op{"BVC", bvc, amRel},
	&op{"EOR", eor, amIndY},
	nil,
	nil,
	nil,
	&op{"EOR", eor, amZpgX},
	nil,
	nil,

	// 0x58
	nil,
	&op{"EOR", eor, amAbsY},
	nil,
	nil,
	nil,
	&op{"EOR", eor, amAbsX},
	nil,
	nil,

	// 0x60
	&op{"RTS", rts, amImpl},
	&op{"ADC", adc, amIndX},
	nil,
	nil,
	nil,
	&op{"ADC", adc, amZpg},
	nil,
	nil,

	// 0x68
	&op{"PLA", pla, amImpl},
	&op{"ADC", adc, amImmed},
	nil,
	nil,
	nil,
	&op{"ADC", adc, amAbs},
	nil,
	nil,

	// 0x70
	&op{"BVS", bvs, amRel},
	&op{"ADC", adc, amIndY},
	nil,
	nil,
	nil,
	&op{"ADC", adc, amZpgX},
	nil,
	nil,

	// 0x78
	&op{"SEI", sei, amImpl},
	&op{"ADC", adc, amAbsY},
	nil,
	nil,
	nil,
	&op{"ADC", adc, amAbsX},
	nil,
	nil,

	// 0x80
	nil,
	nil,
	nil,
	nil,
	nil,
	&op{"STA", sta, amZpg},
	&op{"STX", stx, amZpg},
	nil,

	// 0x88
	&op{"DEY", dey, amImpl},
	nil,
	nil,
	nil,
	nil,
	&op{"STA", sta, amAbs},
	nil,
	nil,

	// 0x90
	&op{"BCC", bcc, amRel},
	&op{"STA", sta, amIndY},
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x98
	nil,
	nil,
	&op{"TXS", txs, amImpl},
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0xa0
	&op{"LDY", ldy, amImmed},
	nil,
	&op{"LDX", ldx, amImmed},
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0xa8
	nil,
	&op{"LDA", lda, amImmed},
	nil,
	nil,
	nil,
	&op{"LDA", lda, amAbs},
	nil,
	nil,

	// 0xb0
	&op{"BCS", bcs, amRel},
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0xb8
	&op{"CLV", clv, amImpl},
	nil,
	nil,
	nil,
	nil,
	&op{"LDA", lda, amAbsX},
	nil,
	nil,

	// 0xc0
	&op{"CPY", cpy, amImmed},
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0xc8
	&op{"INY", iny, amImpl},
	&op{"CMP", cmp, amImmed},
	&op{"DEX", dex, amImpl},
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0xd0
	&op{"BNE", bne, amRel},
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0xd8
	&op{"CLD", cld, amImpl},
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
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
	nil,
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
