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
		return nil
	}
}

func (m Model) cpuIntv() uint64 {
	switch m {
	case NTSC:
		return 12
	case PAL:
		return 16
	default:
		return 12 // ??
	}
}
