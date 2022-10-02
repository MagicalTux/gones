package cpu2a03

import "time"

const (
	NTSC  = 559 * time.Nanosecond // 1.79 MHz
	PAL   = 601 * time.Nanosecond // 1.66 MHz
	Dendy = 564 * time.Nanosecond

	NTSCFreq = 1789773 // number of steps per second
)
