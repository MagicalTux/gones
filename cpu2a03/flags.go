package cpu2a03

// Flags set in cpu.P
const (
	FlagCarry            byte = 0x01
	FlagZero             byte = 0x02
	FlagInterruptDisable byte = 0x04
	FlagDecimal          byte = 0x08
	FlagBreak            byte = 0x10
	FlagIgnored          byte = 0x20
	FlagOverflow         byte = 0x40
	FlagNegative         byte = 0x80
)

/*
Note: The break flag is not an actual flag implemented in a register, and rather
appears only, when the status register is pushed onto or pulled from the stack.
When pushed, it will be 1 when transfered by a BRK or PHP instruction, and
zero otherwise (i.e., when pushed by a hardware interrupt).
When pulled into the status register (by PLP or on RTI), it will be ignored.

In other words, the break flag will be inserted, whenever the status register
is transferred to the stack by software (BRK or PHP), and will be zero, when
transferred by hardware. Since there is no actual slot for the break flag, it
will be always ignored, when retrieved (PLP or RTI).
The break flag is not accessed by the CPU at anytime and there is no internal
representation. Its purpose is more for patching, to discern an interrupt caused
by a BRK instruction from a normal interrupt initiated by hardware.

Note on the overflow flag: The overflow flag indicates overflow with signed
binary arithmetcis. As a signed byte represents a range of -128 to +127, an
overflow can never occure when the operands are of opposite sign, since the
result will never exceed this range. Thus, overflow may only occure, if both
operands are of the same sign. Then, the result must be also of the same sign.
Otherwise, overflow is detected and the overflow flag is set.
(I.e., both operands have a zero in the sign position at bit 7, but bit 7 of the
result is 1, or, both operands have the sign-bit set, but the result is positive.)
*/