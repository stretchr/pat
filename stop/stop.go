package stop

// Done is the variable sent on the channel to indicate
// that something has stopped.
var Done = struct{}{}

// Stopper represents types that implement
// the stop channel pattern.
type Stopper interface {
	// Stop instructs the type to halt operations and
	// returns a channel on which a stop.Done signal is sent
	// when stopping has completed.
	Stop() <-chan struct{}
}

// Stopped returns a channel that signals immediately. Useful for
// cases when no tear-down work is required and stopping is
// immediate.
func Stopped() <-chan struct{} {
	c := Make()
	c <- Done
	return c
}

// Make makes a new channel used to indicate when
// stopping has finished. Sends to channel will not block.
func Make() chan struct{} {
	return make(chan struct{}, 1)
}

// All stops all Stopper types and returns another channel
// on which Done will be sent once all things have
// finished stopping.
func All(stoppers ...Stopper) <-chan struct{} {
	all := Make()
	go func() {
		var allChans []<-chan struct{}
		for _, stopper := range stoppers {
			allChans = append(allChans, stopper.Stop())
		}
		for _, ch := range allChans {
			<-ch
		}
		all <- Done
	}()
	return all
}
