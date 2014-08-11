package start_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/pat/start"
	"github.com/stretchr/pat/stop"
	"github.com/stretchr/testify/require"
)

type TestStarter struct {
	running bool
}

var _ start.StartStopper = (*TestStarter)(nil)

func (t *TestStarter) Start() error {
	time.Sleep(1 * time.Second)
	t.running = true
	return nil
}
func (t *TestStarter) Stop(time.Duration) {
	t.running = false
}
func (t *TestStarter) StopChan() <-chan stop.Signal {
	return stop.Stopped()
}

type ErrorStarter struct {
	running bool
}

var _ start.StartStopper = (*ErrorStarter)(nil)

func (t *ErrorStarter) Start() error {
	time.Sleep(1 * time.Second)
	return errors.New("something went wrong")
}
func (t *ErrorStarter) Stop(time.Duration) {
	t.running = false
}
func (t *ErrorStarter) StopChan() <-chan stop.Signal {
	return stop.Stopped()
}

func TestAll(t *testing.T) {

	s1 := &TestStarter{}
	s2 := &TestStarter{}
	s3 := &TestStarter{}

	errs := start.All(s1, s2, s3)
	require.Equal(t, 0, len(errs))

	require.True(t, s1.running)
	require.True(t, s2.running)
	require.True(t, s3.running)

}

func TestAllErr(t *testing.T) {

	s1 := &TestStarter{}
	s2 := &TestStarter{}
	s3 := &ErrorStarter{}

	errs := start.All(s1, s2, s3)
	require.Equal(t, 1, len(errs))

	require.Equal(t, errs[s3].Error(), "something went wrong")

	require.True(t, s1.running)
	require.True(t, s2.running)
	require.False(t, s3.running)

}

func TestMustAll(t *testing.T) {

	s1 := &TestStarter{}
	s2 := &TestStarter{}
	s3 := &ErrorStarter{}

	errs := start.MustAll(500*time.Millisecond, s1, s2, s3)
	require.Equal(t, 1, len(errs))

	require.Equal(t, errs[s3].Error(), "something went wrong")

	require.False(t, s1.running)
	require.False(t, s2.running)
	require.False(t, s3.running)

}
