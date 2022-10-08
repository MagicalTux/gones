package apu

import (
	"log"

	"github.com/MagicalTux/gones/memory"
)

type APU struct {
	Memory    memory.Master
	Input     [2]InputDevice // we put inputs here since the APU's buffer is used to talk to them
	Interrupt func()
	cpuDelay  func(uint64) uint64

	channel chan float32

	// instruments
	pulse1   *Pulse
	pulse2   *Pulse
	triangle *Triangle
	noise    *Noise
	dmc      *DMC

	cycle       uint64
	frameMode   byte // 0 or 1
	frameValue  byte
	frameIRQ    bool
	filterChain FilterChain

	overrunWarning bool
	interruptFlag  bool
}

func New(mem memory.Master, t func(uint64) uint64) *APU {
	res := &APU{
		Memory:   mem,
		cpuDelay: t,
		channel:  make(chan float32, bufferedSamples*4),
		pulse1:   &Pulse{channel: 1},
		pulse2:   &Pulse{channel: 2},
		triangle: &Triangle{},
		noise:    &Noise{},
		dmc:      &DMC{},
	}
	res.dmc.apu = res
	res.setSampleRate(44100) // standard NES sample rate
	res.fillBuffer()         // fill audio buffer with zeroes

	return res
}

func (apu *APU) setSampleRate(sampleRate float64) {
	// Initialize filters
	apu.filterChain = FilterChain{
		HighPassFilter(float32(sampleRate), 90),
		HighPassFilter(float32(sampleRate), 440),
		LowPassFilter(float32(sampleRate), 14000),
	}
}

func (apu *APU) Clock240(cnt uint64) uint64 {
	// 240Hz clock
	for i := uint64(0); i < cnt; i += 1 {
		apu.stepFrameCounter()
	}
	return cnt
}

func (apu *APU) Clock44100(cnt uint64) uint64 {
	// 44100Hz clock
	for i := uint64(0); i < cnt; i += 1 {
		apu.sendSample()
	}
	return cnt
}

func (apu *APU) ClockCPU(cnt uint64) uint64 {
	// Clock running at same speed as CPU
	for i := uint64(0); i < cnt; i += 1 {
		apu.cycle++
		apu.stepTimer()
	}
	return cnt
}

// https://www.nesdev.org/wiki/APU_Frame_Counter
func (apu *APU) stepFrameCounter() {
	switch apu.frameMode {
	case 0:
		apu.frameValue = (apu.frameValue + 1) & 3
		switch apu.frameValue {
		case 0, 2:
			apu.stepEnvelope()
		case 1:
			apu.stepEnvelope()
			apu.stepSweep()
			apu.stepLength()
		case 3:
			apu.stepEnvelope()
			apu.stepSweep()
			apu.stepLength()
			if apu.frameIRQ {
				apu.interruptFlag = true
				if i := apu.Interrupt; i != nil {
					i()
				}
			}
		}
	case 1:
		apu.frameValue = (apu.frameValue + 1) % 5
		switch apu.frameValue {
		case 0, 2:
			apu.stepEnvelope()
		case 1, 4:
			apu.stepEnvelope()
			apu.stepSweep()
			apu.stepLength()
		}
	}
}

func (apu *APU) stepTimer() {
	if apu.cycle%2 == 0 {
		apu.pulse1.stepTimer()
		apu.pulse2.stepTimer()
		apu.noise.stepTimer()
		apu.dmc.stepTimer()
	}
	apu.triangle.stepTimer()
}

func (apu *APU) stepEnvelope() {
	apu.pulse1.stepEnvelope()
	apu.pulse2.stepEnvelope()
	apu.triangle.stepCounter()
	apu.noise.stepEnvelope()
}

func (apu *APU) stepSweep() {
	apu.pulse1.stepSweep()
	apu.pulse2.stepSweep()
}

func (apu *APU) stepLength() {
	apu.pulse1.stepLength()
	apu.pulse2.stepLength()
	apu.triangle.stepLength()
	apu.noise.stepLength()
}

func (apu *APU) readStatus() byte {
	// Reading this register clears the frame interrupt flag (but not the DMC interrupt flag).
	// If an interrupt flag was set at the same moment of the read, it will read back as 1 but it will not be cleared.
	var res byte
	if apu.pulse1.lengthValue > 0 {
		res |= 0x01
	}
	if apu.pulse2.lengthValue > 0 {
		res |= 0x02
	}
	if apu.triangle.lengthValue > 0 {
		res |= 0x04
	}
	if apu.noise.lengthValue > 0 {
		res |= 0x08
	}
	if apu.dmc.currentLength > 0 {
		res |= 0x10
	}
	if apu.interruptFlag {
		res |= 0x40
		apu.interruptFlag = false
	}
	if apu.dmc.irqFlag {
		res |= 0x80
	}
	// $4015   if-d nt21   DMC IRQ, frame IRQ, length counter statuses
	log.Printf("APU: Read status = $%02x", res)
	return res
}

func (apu *APU) writeControl(value byte) {
	apu.pulse1.enabled = value&1 == 1
	apu.pulse2.enabled = value&2 == 2
	apu.triangle.enabled = value&4 == 4
	apu.noise.enabled = value&8 == 8
	apu.dmc.enabled = value&16 == 16

	if !apu.pulse1.enabled {
		apu.pulse1.lengthValue = 0
	}
	if !apu.pulse2.enabled {
		apu.pulse2.lengthValue = 0
	}
	if !apu.triangle.enabled {
		apu.triangle.lengthValue = 0
	}
	if !apu.noise.enabled {
		apu.noise.lengthValue = 0
	}
	apu.dmc.irqFlag = false
	if !apu.dmc.enabled {
		apu.dmc.currentLength = 0
	} else {
		if apu.dmc.currentLength == 0 {
			apu.dmc.restart()
		}
	}
}

func (apu *APU) writeFrameCounter(value byte) {
	apu.frameMode = (value >> 7) & 1
	apu.frameIRQ = (value>>6)&1 == 0
	if !apu.frameIRQ {
		apu.interruptFlag = false
	}
	// apu.frameValue = 0
	if apu.frameMode == 1 {
		apu.stepEnvelope()
		apu.stepSweep()
		apu.stepLength()
	}
}
