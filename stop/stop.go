package stop

import "time"

// Signal is the type that gets sent down the stop channel.
type Signal struct{}

// Stopper represents types that implement
// the stop channel pattern.
type Stopper interface {
	// Stop instructs the type to halt operations and close
	// the stop channel when it is finished.
	Stop(wait time.Duration)
	// StopChan gets the stop channel which will block until
	// stopping has completed, at which point it is closed.
	// Callers should never close the stop channel.
	StopChan() <-chan Signal
}

// Stopped returns a channel that signals immediately. Useful for
// cases when no tear-down work is required and stopping is
// immediate.
func Stopped() <-chan Signal {
	c := Make()
	close(c)
	return c
}

// Make makes a new channel used to indicate when
// stopping has finished. Sends to channel will not block.
func Make() chan Signal {
	return make(chan Signal, 0)
}

// All stops all Stopper types and returns another channel
// which will close once all things have finished stopping.
func All(wait time.Duration, stoppers ...Stopper) <-chan Signal {
	all := Make()
	go func() {
		var allChans []<-chan Signal
		for _, stopper := range stoppers {
			stopper.Stop(wait)
			allChans = append(allChans, stopper.StopChan())
		}
		for _, ch := range allChans {
			<-ch
		}
		close(all)
	}()
	return all
}
