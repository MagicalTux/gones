package cpu2a03

import "github.com/MagicalTux/gones/nesclock"

type Model byte

const (
	NTSC Model = iota
	PAL
)

func (m Model) clock() *nesclock.Master {
	switch m {
	case NTSC:
		return nesclock.New(nesclock.NTSC, nesclock.StdMul)
	case PAL:
		return nesclock.New(nesclock.PAL, nesclock.StdMul)
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
