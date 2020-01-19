package timerx

import "time"

type state int8

const (
	running state = iota
	paused  state = iota
	stopped state = iota
)

type Timer struct {
	timer    *time.Timer
	started  time.Time
	duration time.Duration
	state    state
	C        <-chan time.Time
}

func NewTimer(d time.Duration) *Timer {
	timer := time.NewTimer(d)
	return &Timer{
		timer:    timer,
		started:  time.Now(),
		duration: d,
		state:    running,
		C:        timer.C,
	}
}

func (t *Timer) Stop() bool {
	if t.state != running {
		return false
	} else if !t.timer.Stop() {
		return false
	}

	t.state = stopped
	return true
}

func (t *Timer) Reset(d time.Duration) bool {
	if !t.timer.Reset(d) {
		return false
	}

	t.started = time.Now()
	t.duration = d
	t.state = running
	return true
}

func (t *Timer) Pause() bool {
	if t.state != running {
		return false
	} else if !t.timer.Stop() {
		return false
	}

	t.duration = t.duration - time.Now().Sub(t.started)
	t.state = paused
	return true
}

func (t *Timer) Start() bool {
	if t.state != paused {
		return false
	} else if t.duration <= 0 {
		return false
	}

	return t.Reset(t.duration)
}
