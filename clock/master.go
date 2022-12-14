package clock

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type ClockInput func(uint64) uint64

type ClockState int

const (
	Stopped ClockState = iota
	Running
	OneStep
)

type Master struct {
	freq  uint64        // number of clocks per second
	step  uint64        // hz per cycle
	intv  time.Duration // duration of a cycle
	pos   uint64        // clocks so far
	state ClockState
	now   time.Time
	next  *Listener
	mu    sync.Mutex
	cd    *sync.Cond
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
	res.cd = sync.NewCond(&res.mu)

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

				// This would be in a perfect world:
				//m.now = now.Add(sleep - eslap)
				//pos += sleepHz

				// As it turns out, Go's time.Sleep() can be very random, and while sometimes short sleeps are properly handled (50??s end sleeping 55??s),
				// other sleeps are not handled properly at all, and a 200??s may end taking 2ms, or a 1ms sleep could end taking 16ms, and someone even
				// reported a 1 hour sleep taking 1 hour and 3 minutes(!).
				// So we'll try to sleep the time we need to sleep using time.Sleep(), but re-read time.Now() after sleeping and update timers based
				// on the new value so we don't end with any surprises.
				// It looks like go1.20 may fix that (it's the target set on the bug report I found as of this writing), but it seems that even if the
				// large time difference is made a bit better, time.Sleep() is not a good option (and Go isn't meant to do real time stuff, the garbage
				// collector being a big proof of it).

				now = time.Now()
				eslap = now.Sub(m.now)
				timCycles := uint64(eslap / m.intv)
				pos += timCycles * m.step
				m.now = m.now.Add(time.Duration(timCycles) * m.intv)
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

stateLoop:
	for {
		switch m.state {
		case Running:
			break stateLoop
		case Stopped:
			m.cd.Wait()
		case OneStep:
			m.state = Stopped
			break stateLoop
		}
	}

	next := m.next
	if next != nil {
		m.next = next.next
	}
	return next
}

func (m *Master) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.state = Running
	m.now = time.Now() // reset wallclock time
	m.cd.Broadcast()   // wake thread if needed
}

func (m *Master) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.state = Stopped
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
