# goNES

NES emulator using ebitengine.

PPU & APU are based on [Michael Fogleman's emulator](https://github.com/fogleman/nes), with added/removed bugs and more rewrites planned/needed.

MOS 6502 CPU is mostly implemented from scratch.

Memory mapping system is also implemented from scratch, as well as cartridge loading/mapping.

## Structure

* `cpu2a03` contains the CPU emulation
* `memory` contains memory primitives such as the bus, RAM and ROM
* `cartridge` has code to load a cartridge and map it on the CPU's bus
* `ppu` contains video rendering related code
* `apu` contains audio code

## References

* http://www.6502.org/tutorials/6502opcodes.html
* https://www.masswerk.at/6502/6502_instruction_set.html
* http://hp.vector.co.jp/authors/VA042397/nes/6502.html
* https://ersanio.gitbook.io/assembly-for-the-snes/mathemathics-and-logic/logic
