package stop_test

import (
	"testing"
	"time"

	"github.com/stretchr/pat/stop"
)

type testStopper struct{}

func (t *testStopper) Stop() stop.Chan {
	s := stop.MakeChan()
	go func() {
		s <- stop.Done
	}()
	return s
}

func TestStop(t *testing.T) {

	s := new(testStopper)
	stopChan := s.Stop()
	select {
	case <-stopChan:
	case <-time.After(1 * time.Second):
		t.Error("Stop signal was never sent (timed out)")
	}

}

func TestAll(t *testing.T) {

	s1 := new(testStopper)
	s2 := new(testStopper)
	s3 := new(testStopper)

	select {
	case <-stop.All(s1, s2, s3):
	case <-time.After(1 * time.Second):
		t.Error("All signal was never sent (timed out)")
	}

}