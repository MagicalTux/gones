package cpu6502

import (
	"fmt"
)

/*
A	Accumulator	OPC A	operand is AC (implied single byte instruction)
abs	absolute	OPC $LLHH	operand is address $HHLL *
abs,X	absolute, X-indexed	OPC $LLHH,X	operand is address; effective address is address incremented by X with carry **
abs,Y	absolute, Y-indexed	OPC $LLHH,Y	operand is address; effective address is address incremented by Y with carry **
#	immediate	OPC #$BB	operand is byte BB
impl	implied	OPC	operand implied
ind	indirect	OPC ($LLHH)	operand is address; effective address is contents of word at address: C.w($HHLL)
X,ind	X-indexed, indirect	OPC ($LL,X)	operand is zeropage address; effective address is word in (LL + X, LL + X + 1), inc. without carry: C.w($00LL + X)
ind,Y	indirect, Y-indexed	OPC ($LL),Y	operand is zeropage address; effective address is word in (LL, LL + 1) incremented by Y with carry: C.w($00LL) + Y
rel	relative	OPC $BB	branch target is PC + signed offset BB ***
zpg	zeropage	OPC $LL	operand is zeropage address (hi-byte is zero, address = $00LL)
zpg,X	zeropage, X-indexed	OPC $LL,X	operand is zeropage address; effective address is address incremented by X without carry **
zpg,Y	zeropage, Y-indexed	OPC $LL,Y	operand is zeropage address; effective address is address incremented by Y without carry **
*/

type AddressMode byte

const (
	amInvalid AddressMode = iota
	amAcc
	amAbs
	amAbsX
	amAbsY
	amImmed
	amImpl
	amInd
	amIndX
	amIndY
	amRel
	amZpg
	amZpgX
	amZpgY
)

func (am AddressMode) Addr(cpu *CPU) uint16 {
	switch am {
	case amAcc:
		panic("amAcc.Addr()")
	case amAbs:
		return cpu.ReadPC16()
	case amAbsX:
		// if page crossed add 1 cycle
		addr := cpu.ReadPC16()
		addr2 := addr + uint16(cpu.X)
		if addr&0xff00 != addr2&0xff00 {
			// different page, lose one cycle due to dummy read
			cpu.Memory.MemRead((addr & 0xff00) | (addr2 & 0xff)) // dummy read
			cpu.cyc += 1
		}
		return addr2
	case amAbsY:
		// if page crossed add 1 cycle
		addr := cpu.ReadPC16()
		addr2 := addr + uint16(cpu.Y)
		if addr&0xff00 != addr2&0xff00 {
			// different page, lose one cycle due to dummy read
			cpu.Memory.MemRead((addr & 0xff00) | (addr2 & 0xff)) // dummy read
			cpu.cyc += 1
		}
		return addr2
	case amImmed:
		panic("amImmed.Addr()")
	case amInd:
		addr := cpu.ReadPC16()
		return cpu.Read16W(addr)
	case amIndX:
		addr := uint16(cpu.ReadPC() + cpu.X)
		return cpu.Read16W(addr)
	case amIndY:
		addr := uint16(cpu.ReadPC())
		addr = cpu.Read16W(addr)
		addr2 := addr + uint16(cpu.Y)
		if addr&0xff00 != addr2&0xff00 {
			// different page, lose one cycle due to dummy read
			cpu.Memory.MemRead((addr & 0xff00) | (addr2 & 0xff)) // dummy read
			cpu.cyc += 1
		}
		return addr2
	case amRel:
		offt := uint16(cpu.ReadPC())
		if offt&0x80 == 0x80 {
			offt |= 0xff00
		}
		return cpu.PC + offt
	case amZpg:
		return uint16(cpu.ReadPC())
	case amZpgX:
		return uint16(cpu.ReadPC() + cpu.X)
	case amZpgY:
		return uint16(cpu.ReadPC() + cpu.Y)
	default:
		panic("unhandled address mode")
	}
}

func (am AddressMode) AddrFast(cpu *CPU) uint16 {
	switch am {
	case amAbsX:
		addr := cpu.ReadPC16()
		addr2 := addr + uint16(cpu.X)
		cpu.Memory.MemRead((addr & 0xff00) | (addr2 & 0xff)) // pre-write read
		return addr2
	case amAbsY:
		addr := cpu.ReadPC16()
		addr2 := addr + uint16(cpu.Y)
		cpu.Memory.MemRead((addr & 0xff00) | (addr2 & 0xff)) // pre-write read
		return addr2
	case amIndY:
		addr := cpu.Read16W(uint16(cpu.ReadPC()))
		addr2 := addr + uint16(cpu.Y)
		cpu.Memory.MemRead((addr & 0xff00) | (addr2 & 0xff)) // pre-write read
		return addr2
	default:
		return am.Addr(cpu)
	}
}

func (am AddressMode) Read(cpu *CPU) byte {
	switch am {
	case amAcc:
		return cpu.A
	case amImmed:
		return cpu.ReadPC()
	case amImpl:
		panic("amImpl.Read()")
	case amRel:
		// can only Addr() this
		panic("amRel.Read()")
	default:
		return cpu.Memory.MemRead(am.Addr(cpu))
	}
}

func (am AddressMode) Write(cpu *CPU, v byte) {
	switch am {
	case amAcc:
		cpu.A = v
	case amImmed:
		panic("amImmed.Write()")
	case amImpl:
		panic("amImpl.Write()")
	case amRel:
		// can only Addr() this
		panic("amRel.Write()")
	default:
		cpu.Memory.MemWrite(am.Addr(cpu), v)
	}
}

// WriteFast writes without losing cycles
func (am AddressMode) WriteFast(cpu *CPU, v byte) {
	switch am {
	case amAbsX, amAbsY, amIndY:
		cpu.Memory.MemWrite(am.AddrFast(cpu), v)
	default:
		am.Write(cpu, v)
	}
}

func (am AddressMode) Debug(cpu *CPU) string {
	switch am {
	case amAcc:
		return fmt.Sprintf("A = $%02x", cpu.A)
	case amAbs:
		return fmt.Sprintf("abs = $%04x", cpu.PeekPC16())
	case amAbsX:
		return fmt.Sprintf("abs,X = $%04x,$%02x", cpu.PeekPC16(), cpu.X)
	case amAbsY:
		return fmt.Sprintf("abs,Y = $%04x,$%02x", cpu.PeekPC16(), cpu.Y)
	case amImmed:
		return fmt.Sprintf("#$%02x", cpu.PeekPC())
	case amImpl:
		return "impl"
	case amInd:
		addr := cpu.PeekPC16()
		return fmt.Sprintf("ind = ($%04x)", addr)
	case amIndX:
		addr := cpu.PeekPC()
		return fmt.Sprintf("ind,X = ($%02x,$%02x)", addr, cpu.X)
	case amIndY:
		addr := uint16(cpu.PeekPC())
		return fmt.Sprintf("ind,Y = ($%02x),$%02x", addr, cpu.Y)
	case amRel:
		offt := uint16(cpu.PeekPC())
		if offt&0x80 == 0x80 {
			offt |= 0xff00
		}
		// add 1 because we used PeekPC instead of ReadPC
		return fmt.Sprintf("rel = %d = $%04x", int16(offt), cpu.PC+offt+1)
	case amZpg:
		return fmt.Sprintf("zpg = $%02x", cpu.PeekPC())
	case amZpgX:
		return fmt.Sprintf("zpg,X = $%02x,$%02x", cpu.PeekPC(), cpu.X)
	case amZpgY:
		return fmt.Sprintf("zpg,Y = $%02x,$%02x", cpu.PeekPC(), cpu.Y)
	default:
		return fmt.Sprintf("unknown $%02x", am)
	}
}

func (am AddressMode) Length() int {
	switch am {
	case amAcc, amImpl:
		return 0
	case amImmed, amIndX, amIndY, amRel, amZpg, amZpgX, amZpgY:
		return 1
	case amAbs, amAbsX, amAbsY, amInd:
		return 2
	default:
		return 0
	}
}

func (am AddressMode) Implied(cpu *CPU) {
	if am == amImpl {
		return
	}
	panic("expected amImpl")
}
