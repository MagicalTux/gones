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
	nil,
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
	nil,
	nil,
	nil,
	nil,

	// 0x28
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x30
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x38
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x40
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x48
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x50
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x58
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x60
	&op{"RTS", rts, amImpl},
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x68
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x70
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x78
	&op{"SEI", sei, amImpl},
	nil,
	nil,
	nil,
	nil,
	nil,
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
	nil,
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
	nil,
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
	nil,
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
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0xe8
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0xf0
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0xf8
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
}
