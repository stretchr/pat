package sleep

import "time"

// Action is the return from Sleep() to indicate
// whether the process should be aborted or not.
// If Sleep() == sleep.Abort all added intervals have
// been exhausted.
type Action bool

const (
	// Abort is returned from Sleep if all intervals
	// have been exhausted.
	Abort   Action = true
	carryon Action = false
)

// Sleeper provides a mechanism by which sleep intervals
// change over time to reduce interactions with resources that may
// be unavailable now, but may come back in the future.
type Sleeper interface {
	// Add adds an interval to the Sleeper where it will sleep
	// for "sleep" duration, until "duration" has elapsed.
	Add(duration, sleep time.Duration)
	// Sleep sleeps for the next duration and returns sleep.Abort
	// when all intervals have been exhausted.
	// Sleeper automatically resets after all intervals have been exhausted.
	Sleep() Action
	// Reset manually resets the Sleeper.
	Reset()
	// Abort manually aborts the Sleeper, cancelling the current call
	// to Sleep.
	Abort()
	// Duration gets the next time.Duration that Sleep() will sleep for.
	Duration() time.Duration
}

type sleeper struct {
	sleep       func(time.Duration) Action
	ints        []*interval
	current     int
	abort       chan struct{}
	shouldReset bool
}

type interval struct {
	sleep   time.Duration
	count   int
	current int
}

var _ Sleeper = (*sleeper)(nil)

// New creates a new Sleeper.
func New() Sleeper {
	abort := make(chan struct{})
	return &sleeper{
		sleep: func(d time.Duration) Action {
			select {
			case <-abort:
				return Abort
			case <-time.After(d):
				return carryon
			}
		},
		abort: abort,
	}
}

func (s *sleeper) Add(duration, sleep time.Duration) {
	if sleep == 0 || duration == 0 {
		panic("sleep: duration and sleep must be > 0")
	}
	if sleep > duration {
		panic("sleep: sleep cannot exceed duration")
	}
	if duration%sleep != 0 {
		panic("sleep: duration must be divisible by sleep")
	}
	s.ints = append(s.ints, &interval{sleep: sleep, count: int(duration / sleep)})
}

func (s *sleeper) Duration() time.Duration {
	if s.current == len(s.ints) {
		return 0
	}
	return s.ints[s.current].sleep
}

func (s *sleeper) Sleep() Action {
	s.shouldReset = true
	if s.current == len(s.ints) {
		s.Reset()
		return Abort
	}
	next := s.ints[s.current]
	ret := s.sleep(next.sleep)
	next.current++
	if next.current == next.count {
		s.current++
	}
	return ret
}

func (s *sleeper) Reset() {
	if !s.shouldReset {
		return
	}
	s.current = 0
	for _, i := range s.ints {
		i.current = 0
	}
	s.shouldReset = false
}

func (s *sleeper) Abort() {
	s.abort <- struct{}{}
}
