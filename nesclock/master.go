package nesclock

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// See: https://www.nesdev.org/wiki/Cycle_reference_chart
	NTSC = 21477470 // 21.47727 MHz (NTSC) 21.477272 MHz ± 40 Hz
	PAL  = 26601700 // 26.6017 MHz (PAL) 26.601712 MHz ± 50 Hz

	// Clock: requested 21477470 Hz clock, computed clock will be 21477663 Hz (25 steps/1.164µs interval, a 193Hz diff)
	// Clock: requested 21477470 Hz clock, computed clock will be 21477484 Hz (207 steps/9.638µs interval, a 14Hz diff)
	// Clock: requested 26601700 Hz clock, computed clock will be 26601723 Hz (71 steps/2.669µs interval, a 23Hz diff)
)

type ClockInput func(uint64) uint64

type Master struct {
	freq uint64        // number of clocks per second
	step uint64        // hz per cycle
	intv time.Duration // duration of a cycle
	pos  uint64        // clocks so far
	now  time.Time
	next *Listener
	mu   sync.Mutex
}

// New returns a new clock running at the given frequency, using the specified
// step value as increment value. Higher step values means more calls will
// happen in batches, but improves how close this will run to real time.
func New(freq uint64) *Master {
	step := uint64(1)

	var intv time.Duration
	var realFreq, diff, bestStep, bestDiff uint64

	for ; step < 255; step++ {
		// compute nanoseconds for a step
		intv = time.Second * time.Duration(step) / time.Duration(freq)
		realFreq = uint64(time.Second * time.Duration(step) / intv)

		if realFreq > freq {
			diff = realFreq - freq
		} else {
			diff = freq - realFreq
		}

		if bestStep == 0 || diff < bestDiff {
			bestStep = step
			bestDiff = diff
		}
	}
	step = bestStep

	// compute nanoseconds for a step
	intv = time.Second * time.Duration(step) / time.Duration(freq)
	realFreq = uint64(time.Second * time.Duration(step) / intv)
	if realFreq > freq {
		diff = realFreq - freq
	} else {
		diff = freq - realFreq
	}

	log.Printf("Clock: requested %d Hz clock, computed clock will be %d Hz (%d steps/%s interval, a %dHz diff)", freq, realFreq, step, intv, diff)

	res := &Master{
		freq: freq,
		step: step,
		intv: intv,
		now:  time.Now(),
	}
	go res.thread()

	// Sample usage:
	//res.Listen(freq/10, func(uint64) uint64 { log.Printf("Clock: test @1/10th of a sec"); return 1 })

	return res
}

func (m *Master) Frequency() uint64 {
	return m.freq
}

// Listen will cause cb() to be called every `divider` tick of the master
// clock, allowing synchronization between various elements of the NES.
func (m *Master) Listen(divider, delta uint64, cb ClockInput) *Listener {
	l := &Listener{
		cb:      cb,
		divider: divider,
		delta:   delta,
	}

	l.nextRun = ((atomic.LoadUint64(&m.pos)/l.divider)+1)*l.divider + l.delta
	m.insert(l)
	return l
}

// The NES cpu runs with a divider of 12 or 16, the NES PPU runs with a divider of 4, and the NES APU runs with a divider of 12*240

func (m *Master) thread() {
	var pos uint64 // master value for pos (m.pos is a slave)

	for {
		cur := m.takeNext()
		if cur == nil {
			time.Sleep(50 * time.Millisecond)
			m.now = time.Now()
			continue
		}
		if cur.nextRun > pos {
			// this doesn't need to run yet?
			now := time.Now()
			eslap := now.Sub(m.now)

			// how long to sleep until cur > nextRun?
			sleepHz := (cur.nextRun - pos)
			sleepCycles := sleepHz / m.step
			// if not round, add 1
			if sleepHz%m.step != 0 {
				sleepCycles += 1
			}
			// convert to time duration
			sleep := time.Duration(sleepCycles) * m.intv
			if eslap < sleep {
				// can sleep more
				time.Sleep(sleep - eslap)
				m.now = now.Add(sleep - eslap)
				pos += sleepHz
			} else {
				// convert eslap back into time unit
				timCycles := uint64(eslap / m.intv)
				pos += timCycles * m.step
				m.now = m.now.Add(time.Duration(timCycles) * m.intv)
			}
			atomic.StoreUint64(&m.pos, pos)
		}

		// call it
		// TODO: compute how much time we have until cur.next & pass it as parameter

		cnt := cur.run(1)

		// add time
		cur.nextRun += cur.divider * cnt
		m.insert(cur)
	}
}

func (m *Master) takeNext() *Listener {
	m.mu.Lock()
	defer m.mu.Unlock()

	next := m.next
	if next != nil {
		m.next = next.next
	}
	return next
}

func (m *Master) insert(l *Listener) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.next == nil {
		l.next = nil
		m.next = l
		return
	}
	cur := m.next
	if cur.nextRun > l.nextRun {
		// this one happens first
		l.next = cur
		m.next = l
		return
	}
	for {
		if cur.next == nil {
			l.next = nil
			cur.next = l
			return
		}
		if cur.next.nextRun > l.nextRun {
			l.next = cur.next
			cur.next = l
			return
		}

		cur = cur.next
	}
}
