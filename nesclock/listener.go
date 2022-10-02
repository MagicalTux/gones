package nesclock

type Listener struct {
	next    *Listener
	divider uint64 // how many clocks between runs
	nextRun uint64 // time of next run
	cb      ClockInput
}
