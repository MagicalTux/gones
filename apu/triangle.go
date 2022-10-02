package apu

type Triangle struct {
	enabled       bool
	lengthEnabled bool
	lengthValue   byte
	timerPeriod   uint16
	timerValue    uint16
	dutyValue     byte
	counterPeriod byte
	counterValue  byte
	counterReload bool
}

func (t *Triangle) writeControl(value byte) {
	t.lengthEnabled = (value>>7)&1 == 0
	t.counterPeriod = value & 0x7F
}

func (t *Triangle) writeTimerLow(value byte) {
	if !t.enabled {
		return
	}
	t.timerPeriod = (t.timerPeriod & 0xFF00) | uint16(value)
}

func (t *Triangle) writeTimerHigh(value byte) {
	if !t.enabled {
		return
	}
	t.lengthValue = lengthTable[value>>3]
	t.timerPeriod = (t.timerPeriod & 0x00FF) | (uint16(value&7) << 8)
	t.timerValue = t.timerPeriod
	t.counterReload = true
}

func (t *Triangle) stepTimer() {
	if t.timerValue == 0 {
		t.timerValue = t.timerPeriod
		if t.lengthValue > 0 && t.counterValue > 0 {
			t.dutyValue = (t.dutyValue + 1) % 32
		}
	} else {
		t.timerValue--
	}
}

func (t *Triangle) stepLength() {
	if t.lengthEnabled && t.lengthValue > 0 {
		t.lengthValue--
	}
}

func (t *Triangle) stepCounter() {
	if t.counterReload {
		t.counterValue = t.counterPeriod
	} else if t.counterValue > 0 {
		t.counterValue--
	}
	if t.lengthEnabled {
		t.counterReload = false
	}
}

func (t *Triangle) output() byte {
	if !t.enabled {
		return 0
	}
	if t.timerPeriod < 3 {
		return 0
	}
	if t.lengthValue == 0 {
		return 0
	}
	if t.counterValue == 0 {
		return 0
	}
	return triangleTable[t.dutyValue]
}
