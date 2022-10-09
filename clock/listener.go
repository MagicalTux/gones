package clock

type Listener struct {
	next    *Listener
	divider uint64 // how many clocks between runs
	delta   uint64
	nextRun uint64 // time of next run
	cb      ClockInput

	// Stats
	//cnt uint64
	//tc  int64
}

func (l *Listener) run(v uint64) uint64 {
	/*
		now := time.Now().Unix()
		if l.tc != now {
			log.Printf("ran %d times", l.cnt)
			l.tc = now
			l.cnt = 0
		}

		res := l.cb(v)
		l.cnt += res
		return res
	*/
	return l.cb(v)
}
