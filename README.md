# goNES

NES emulator using [ebitengine](https://ebitengine.org/).

PPU & APU are based on [Michael Fogleman's emulator](https://github.com/fogleman/nes), with added/removed bugs and more rewrites planned/needed. Because the CPU and the memory code isn't the same, the PPU/APU code isn't the same either, and I'm hoping to fully rewrite it eventually.

MOS 6502 CPU is mostly implemented from scratch.

Memory mapping system is also implemented from scratch, as well as cartridge loading/mapping.

Clock uses a master clock at the frequency specified for NES, and divisors to feed cpu/ppu/apu/etc synchronized signals.

## Structure

* `pkgnes` is the base NES package that will instanciate the various required elements
* `cpu6502` contains the CPU emulation
* `clock` generate clock signals for the other parts of the system
* `memory` contains memory primitives such as the bus, RAM and ROM
* `nescartridge` has code to load a cartridge and map it on the CPU's bus
* `nesppu` contains video rendering related code
* `nesapu` contains audio code
* `nesinput` manages input devices (keyboard only for now)

## References

### CPU

* http://www.6502.org/tutorials/6502opcodes.html
* https://www.masswerk.at/6502/6502_instruction_set.html
* http://hp.vector.co.jp/authors/VA042397/nes/6502.html
* https://ersanio.gitbook.io/assembly-for-the-snes/mathemathics-and-logic/logic
