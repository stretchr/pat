// Package stop represents a pattern for types that need to do some work
// when stopping. The Stop method returns a stop.Chan on which
// stop.Done is passed when the operation has completed.
// Stopper types can be stopped in the following ways:
//     // stop and forget
//     t.Stop()
//
//     // stop and wait
//     <-t.Stop()
//
//     // stop, do more work, then wait
//     stopped := t.Stop();
//     // do more work
//     <-stopped
//
//     // stop and timeout after 1 second
//     select {
//     case <-t.Stop():
//     case <-time.After(1 * time.Second):
//     }
//
//     // stop.All is the same as calling Stop() so
//     // all above patterns also work on many Stopper types,
//     // for example; stop and wait for many things:
//     <-stop.All(t1, t2, t3)
package stop
