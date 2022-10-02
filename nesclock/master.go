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

	StdMul = 25
)

type ClockInput func(uint64) uint64

type Master struct {
	freq uint64        // number of clocks per second
	mul  uint64        // hz per cycle
	intv time.Duration // duration of a cycle
	pos  uint64        // clocks so far
	now  time.Time
	next *Listener
	mu   sync.Mutex
}

// New returns a new clock running at the given frequency, using the specified
// mul value as increment value. Higher mul values means more calls will happen
// in batches, but improves how close this will run to real time.
func New(freq, mul uint64) *Master {
	// compute nanoseconds for `mul` Hz
	intv := time.Second * time.Duration(mul) / time.Duration(freq)
	realFreq := uint64(time.Second * time.Duration(mul) / intv)

	log.Printf("Clock: requested %d Hz clock, computed clock will be %d Hz (%d pulses/%s interval, a %01.6f%% diff)", freq, realFreq, mul, intv, float64(realFreq-freq)/float64(freq)*100)

	res := &Master{
		freq: freq,
		mul:  mul,
		intv: intv,
		now:  time.Now(),
	}
	go res.thread()

	// Sample usage:
	//res.Listen(freq/10, func(uint64) uint64 { log.Printf("Clock: test @1/10th of a sec"); return 1 })

	return res
}

// Listen will cause cb() to be called every `divider` tick of the master
// clock, allowing synchronization between various elements of the NES.
func (m *Master) Listen(divider uint64, cb ClockInput) *Listener {
	l := &Listener{
		cb:      cb,
		divider: divider,
	}

	l.nextRun = ((atomic.LoadUint64(&m.pos) / l.divider) + 1) * l.divider
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
			sleepCycles := sleepHz / m.mul
			// if not round, add 1
			if sleepHz%m.mul != 0 {
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
				pos += timCycles * m.mul
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
