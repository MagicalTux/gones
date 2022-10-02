package apu

type DMC struct {
	apu            *APU
	enabled        bool
	value          byte
	sampleAddress  uint16
	sampleLength   uint16
	currentAddress uint16
	currentLength  uint16
	shiftRegister  byte
	bitCount       byte
	tickPeriod     byte
	tickValue      byte
	loop           bool
	irq            bool
	irqFlag        bool
}

func (d *DMC) MemWrite(addr uint16, val byte) byte {
	switch addr & 3 {
	case 0x00:
		d.writeControl(val)
		return val
	case 0x01:
		d.writeValue(val)
		return val
	case 0x02:
		d.writeAddress(val)
		return val
	case 0x03:
		d.writeLength(val)
		return val
	}
	return 0
}

func (d *DMC) writeControl(value byte) {
	d.irq = value&0x80 == 0x80
	if !d.irq {
		d.irqFlag = false
	}
	d.loop = value&0x40 == 0x40
	d.tickPeriod = dmcTable[value&0x0F]
}

func (d *DMC) writeValue(value byte) {
	d.value = value & 0x7F
}

func (d *DMC) writeAddress(value byte) {
	// Sample address = %11AAAAAA.AA000000
	d.sampleAddress = 0xC000 | (uint16(value) << 6)
}

func (d *DMC) writeLength(value byte) {
	// Sample length = %0000LLLL.LLLL0001
	d.sampleLength = (uint16(value) << 4) | 1
}

func (d *DMC) restart() {
	d.currentAddress = d.sampleAddress
	d.currentLength = d.sampleLength
}

func (d *DMC) stepTimer() {
	if !d.enabled {
		return
	}
	d.stepReader()
	if d.tickValue == 0 {
		d.tickValue = d.tickPeriod
		d.stepShifter()
	} else {
		d.tickValue--
	}
}

func (d *DMC) stepReader() {
	if d.currentLength > 0 && d.bitCount == 0 {
		d.apu.cpuDelay(4)
		d.shiftRegister = d.apu.Memory.MemRead(d.currentAddress)
		d.bitCount = 8
		d.currentAddress++
		if d.currentAddress == 0 {
			d.currentAddress = 0x8000
		}
		d.currentLength--
		if d.currentLength == 0 {
			if d.loop {
				d.restart()
			} else if d.irq {
				d.irqFlag = true
				if i := d.apu.Interrupt; i != nil {
					i()
				}
			}
		}
	}
}

func (d *DMC) stepShifter() {
	if d.bitCount == 0 {
		return
	}
	if d.shiftRegister&1 == 1 {
		if d.value <= 125 {
			d.value += 2
		}
	} else {
		if d.value >= 2 {
			d.value -= 2
		}
	}
	d.shiftRegister >>= 1
	d.bitCount--
}

func (d *DMC) output() byte {
	return d.value
}
