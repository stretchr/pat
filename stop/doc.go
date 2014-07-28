// Package stop represents a pattern for types that need to do some work
// when stopping. The Stop method returns a <-chan stop.Signal on which
// stop.Done is passed when the operation has completed.
//
// Stopper types when implementing the Stop method should use Make
// to create and return a stop channel, and pass stop.Done on the
// channel once stopping has completed:
//     func (t Type) Stop() <-chan stop.Signal {
//       c := stop.Make()
//       go func(){
//         // TODO: tear stuff down
//         c <- stop.Done
//       }()
//       return c
//     }
//
// If Stopper types do not need to do any background work, stop.Stopped() can
// be returned, for example:
//     func (t Type) Stop() <-chan stop.Signal {
//       t.stopped = true
//       return stop.Stopped()
//     }
//
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
