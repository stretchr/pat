package sleep

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSleeper(t *testing.T) {

	s := New().(*sleeper)
	s.Add(1*time.Minute, 1*time.Second)
	s.Add(1*time.Minute, 10*time.Second)
	s.Add(10*time.Minute, 1*time.Minute)

	var durations []time.Duration
	s.sleep = func(d time.Duration) Action {
		durations = append(durations, d)
		return carryon
	}

	for s.Sleep() != Abort {
	}

	require.Equal(t, len(durations), 76)
	require.Equal(t, durations[0], 1*time.Second)   // 1st call
	require.Equal(t, durations[65], 10*time.Second) // 66th call
	require.Equal(t, durations[75], 1*time.Minute)  // 76th call

	durations = make([]time.Duration, 0)

	for s.Sleep() != Abort {
	}

	require.Equal(t, len(durations), 76)
	require.Equal(t, durations[0], 1*time.Second)   // 1st call
	require.Equal(t, durations[65], 10*time.Second) // 66th call
	require.Equal(t, durations[75], 1*time.Minute)  // 76th call

}

func TestSleeping(t *testing.T) {

	s := New()
	s.Add(10*time.Millisecond, 1*time.Millisecond)

	sleeps := 0
	for s.Sleep() != Abort {
		sleeps++
	}

	require.Equal(t, 10, sleeps)

}

func TestAbort(t *testing.T) {

	s := New()
	s.Add(1*time.Second, 10*time.Millisecond)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		time.Sleep(15 * time.Millisecond)
		s.Abort()
		wg.Done()
	}()
	sleeps := 0
	go func() {
		for s.Sleep() != Abort {
			sleeps++
		}
		wg.Done()
	}()

	wg.Wait()
	require.Equal(t, sleeps, 1)

}

func TestPanics(t *testing.T) {

	require.Panics(t, func() {
		New().Add(0, 0)
	}, "Zeros")
	require.Panics(t, func() {
		New().Add(0, 1)
	}, "Zero duration")
	require.Panics(t, func() {
		New().Add(1, 0)
	}, "Zero sleep")
	require.Panics(t, func() {
		New().Add(1, 2)
	}, "sleep > duration")
	require.Panics(t, func() {
		New().Add(10, 8)
	}, "duration % sleep")

}
