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

	var sleepDurs []time.Duration
	s.sleep = func(d time.Duration) Action {
		sleepDurs = append(sleepDurs, d)
		return carryon
	}

	require.Equal(t, 1*time.Second, s.Duration())

	var durations []time.Duration
	for s.Sleep() != Abort {
		durations = append(durations, s.Duration())
	}

	require.Equal(t, len(sleepDurs), 76)
	require.Equal(t, sleepDurs[0], 1*time.Second)   // 1st call
	require.Equal(t, sleepDurs[65], 10*time.Second) // 66th call
	require.Equal(t, sleepDurs[75], 1*time.Minute)  // 76th call

	require.Equal(t, durations[1], 1*time.Second)   // 1st call
	require.Equal(t, durations[64], 10*time.Second) // 66th call
	require.Equal(t, durations[74], 1*time.Minute)  // 76th call

	sleepDurs = make([]time.Duration, 0)
	require.Equal(t, 1*time.Second, s.Duration())

	for s.Sleep() != Abort {
	}

	require.Equal(t, len(sleepDurs), 76)
	require.Equal(t, sleepDurs[0], 1*time.Second)   // 1st call
	require.Equal(t, sleepDurs[65], 10*time.Second) // 66th call
	require.Equal(t, sleepDurs[75], 1*time.Minute)  // 76th call

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

	require.True(t, s.Reset())
	require.False(t, s.Reset())

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
