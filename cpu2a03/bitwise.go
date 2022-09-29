package cpu2a03

import "log"

func oraAbs(cpu *Cpu2A03) {
	// Affects Flags: N Z
	addr := cpu.ReadPC16()
	v := cpu.Memory.MemRead(addr)

	cpu.A |= v
	log.Printf("ORA $%04x", v)
}
