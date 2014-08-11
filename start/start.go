package start

import (
	"sync"
	"time"

	"github.com/stretchr/pat/stop"
)

// StartStopper represents types that need starting and stopping.
type StartStopper interface {
	stop.Stopper
	// Start is called to start the operations of the
	// type, an error is returned if starting fails, otherwise
	// nil.
	Start() error
}

// All starts all StartStoppers and returns a map of any
// errors that occurred keyed by the StartStopper object.
// len(All(s...)) == 0 when everything started successfully.
func All(startStoppers ...StartStopper) map[StartStopper]error {

	var wg sync.WaitGroup
	errs := make(map[StartStopper]error)
	var errsL sync.Mutex
	wg.Add(len(startStoppers))

	// start everything
	for _, s := range startStoppers {
		go func(s StartStopper) {
			if err := s.Start(); err != nil {
				errsL.Lock()
				errs[s] = err
				errsL.Unlock()
			}
			wg.Done()
		}(s)
	}

	// wait for all things to start
	wg.Wait()

	return errs
}

// MustAll starts all StartStoppers and if any one fails to start,
// stops all of the others and returns a map of errors that occurred, keyed
// by the StartStopper object.
// len(All(s...)) == 0 when everything started successfully.
// By the time MustAll returns, everything will have been given the chance to
// start, and in the event of an error, everything will be properly stopped.
func MustAll(stopGrace time.Duration, startStoppers ...StartStopper) map[StartStopper]error {
	errs := All(startStoppers...)
	if len(errs) > 0 {
		// at least one thing failed to start - stop everything that
		// started
		var wg sync.WaitGroup
		for _, s := range startStoppers {
			if _, failed := errs[s]; failed == false {
				s.Stop(stopGrace)
				wg.Add(1)
				go func() {
					<-s.StopChan()
					wg.Done()
				}()
			}
		}
		wg.Wait() // wait for everything to stop
	}
	return errs
}
