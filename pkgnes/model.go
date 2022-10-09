package pkgnes

import "github.com/MagicalTux/gones/clock"

type Model byte

const (
	NTSC Model = iota
	PAL
)

const (
	// See: https://www.nesdev.org/wiki/Cycle_reference_chart
	FreqNTSC = 21477470 // 21.47727 MHz (NTSC) 21.477272 MHz ± 40 Hz
	FreqPAL  = 26601700 // 26.6017 MHz (PAL) 26.601712 MHz ± 50 Hz

// Clock: requested 21477470 Hz clock, computed clock will be 21477663 Hz (25 steps/1.164µs interval, a 193Hz diff)
// Clock: requested 21477470 Hz clock, computed clock will be 21477484 Hz (207 steps/9.638µs interval, a 14Hz diff)
// Clock: requested 26601700 Hz clock, computed clock will be 26601723 Hz (71 steps/2.669µs interval, a 23Hz diff)
)

func (m Model) newClock() *clock.Master {
	switch m {
	case NTSC:
		return clock.New(FreqNTSC)
	case PAL:
		return clock.New(FreqPAL)
	default:
		panic("invalid model")
	}
}

func (m Model) cpuIntv() uint64 {
	switch m {
	case NTSC:
		return 12
	case PAL:
		return 16
	default:
		panic("invalid model")
	}
}

func (m Model) ppuIntv() uint64 {
	switch m {
	case NTSC:
		return 4
	case PAL:
		return 5
	default:
		panic("invalid model")
	}
}
