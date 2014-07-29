package stop

import "time"

// Signal is the type that gets sent down the stop channel.
type Signal struct{}

// Done is the Signal variable sent on the channel to indicate
// that something has stopped.
var Done = Signal{}

// Stopper represents types that implement
// the stop channel pattern.
type Stopper interface {
	// Stop instructs the type to halt operations and
	// returns a channel on which a stop.Done signal is sent
	// when stopping has completed.
	Stop(wait time.Duration)

	StopChan() <-chan Signal
}

// Stopped returns a channel that signals immediately. Useful for
// cases when no tear-down work is required and stopping is
// immediate.
func Stopped() <-chan Signal {
	c := Make()
	c <- Done
	return c
}

// Make makes a new channel used to indicate when
// stopping has finished. Sends to channel will not block.
func Make() chan Signal {
	return make(chan Signal, 1)
}

// All stops all Stopper types and returns another channel
// on which Done will be sent once all things have
// finished stopping.
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
		all <- Done
	}()
	return all
}
