package apu

type Noise struct {
	enabled         bool
	mode            bool
	shiftRegister   uint16
	lengthEnabled   bool
	lengthValue     byte
	timerPeriod     uint16
	timerValue      uint16
	envelopeEnabled bool
	envelopeLoop    bool
	envelopeStart   bool
	envelopePeriod  byte
	envelopeValue   byte
	envelopeVolume  byte
	constantVolume  byte
}

func (n *Noise) writeControl(value byte) {
	n.lengthEnabled = (value>>5)&1 == 0
	n.envelopeLoop = (value>>5)&1 == 1
	n.envelopeEnabled = (value>>4)&1 == 0
	n.envelopePeriod = value & 15
	n.constantVolume = value & 15
	n.envelopeStart = true
}

func (n *Noise) writePeriod(value byte) {
	n.mode = value&0x80 == 0x80
	n.timerPeriod = noiseTable[value&0x0F]
}

func (n *Noise) writeLength(value byte) {
	if !n.enabled {
		return
	}
	n.lengthValue = lengthTable[value>>3]
	n.envelopeStart = true
}

func (n *Noise) stepTimer() {
	if n.timerValue == 0 {
		n.timerValue = n.timerPeriod
		var shift byte
		if n.mode {
			shift = 6
		} else {
			shift = 1
		}
		b1 := n.shiftRegister & 1
		b2 := (n.shiftRegister >> shift) & 1
		n.shiftRegister >>= 1
		n.shiftRegister |= (b1 ^ b2) << 14
	} else {
		n.timerValue--
	}
}

func (n *Noise) stepEnvelope() {
	if n.envelopeStart {
		n.envelopeVolume = 15
		n.envelopeValue = n.envelopePeriod
		n.envelopeStart = false
	} else if n.envelopeValue > 0 {
		n.envelopeValue--
	} else {
		if n.envelopeVolume > 0 {
			n.envelopeVolume--
		} else if n.envelopeLoop {
			n.envelopeVolume = 15
		}
		n.envelopeValue = n.envelopePeriod
	}
}

func (n *Noise) stepLength() {
	if n.lengthEnabled && n.lengthValue > 0 {
		n.lengthValue--
	}
}

func (n *Noise) output() byte {
	if !n.enabled {
		return 0
	}
	if n.lengthValue == 0 {
		return 0
	}
	if n.shiftRegister&1 == 1 {
		return 0
	}
	if n.envelopeEnabled {
		return n.envelopeVolume
	} else {
		return n.constantVolume
	}
}
