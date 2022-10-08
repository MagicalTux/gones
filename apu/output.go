package apu

import (
	"encoding/binary"
	"io"
	"log"
	"time"
)

const bufferedSamples = 768

func BufferLength() time.Duration {
	// 1 audio sample at 44100Hz = 22.675µs or more exactly 22675.73ns
	return time.Duration(bufferedSamples * 22676 * time.Nanosecond)
}

func (apu *APU) fillBuffer() {
	for {
		select {
		case apu.channel <- 0:
			// ok
		default:
			return
		}
	}
}

func (apu *APU) sendSample() {
	output := apu.filterChain.Step(apu.output())
	select {
	case apu.channel <- output:
		apu.overrunWarning = false
	default:
		if !apu.overrunWarning {
			log.Printf("WARNING: buffer overrun, emptying half...")
			apu.overrunWarning = true
			ln := len(apu.channel) / 2
			for i := 0; i < ln; i += 1 {
				select {
				case <-apu.channel:
				default:
				}
			}
		}
	}
}

func (apu *APU) output() float32 {
	p1 := apu.pulse1.output()
	p2 := apu.pulse2.output()
	t := apu.triangle.output()
	n := apu.noise.output()
	d := apu.dmc.output()
	pulseOut := pulseTable[p1+p2]
	tndOut := tndTable[3*t+2*n+d]
	return pulseOut + tndOut
}

func (apu *APU) Read(b []byte) (int, error) {
	if len(b) < 4 {
		return 0, io.ErrShortBuffer
	}

	n := 0
	// signed 16bits little endian, 2 channel stereo
	var v float32

	//log.Printf("APU READ, len(channel) = %d", len(apu.channel))

	for len(b) > 0 {
		v = <-apu.channel // -1 ~ 1

		var i int16

		if v > 0 {
			i = int16(v * 32768)
		} else {
			i = int16(v * 32767)
		}

		//log.Printf("APU %f = %d", v, i)

		// make L & R channels the same by outputting this twice
		binary.LittleEndian.PutUint16(b[:2], uint16(i))
		binary.LittleEndian.PutUint16(b[2:4], uint16(i))

		n += 4
		b = b[4:]
	}

	return n, nil
}
