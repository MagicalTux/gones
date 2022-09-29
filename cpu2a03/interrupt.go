package cpu2a03

func brk(cpu *Cpu2A03, am AddressMode) {
	cpu.PC += 1
	cpu.Push16(cpu.PC)
	cpu.Push(cpu.P | FlagBreak)

	cpu.PC = cpu.Read16(0xfffe) // IRQ/BRK
}
