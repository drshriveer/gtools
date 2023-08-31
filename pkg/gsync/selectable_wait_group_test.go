package gsync_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drshriveer/gcommon/pkg/gerrors"
	"github.com/drshriveer/gcommon/pkg/gsync"
)

func TestSelectableWaitGroup_Wait(t *testing.T) {
	wg := gsync.NewSelectableWaitGroup()
	wg.Add(1)
	err := wg.WaitTimeout(100 * time.Millisecond)
	assert.Equal(t, gsync.ErrWGTimeout, gerrors.Unwrap(err))

	ctx, done := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer done()
	err = wg.WaitCTX(ctx)
	assert.Equal(t, context.DeadlineExceeded, gerrors.Unwrap(err))

	ready := make(chan struct{})
	wg2 := gsync.NewSelectableWaitGroup()
	wg2.Inc()
	go func() {
		select {
		case <-wg.Wait():
			assert.Fail(t, "unexpected")
		case ready <- struct{}{}:
			// proceed
		}
		require.NoError(t, wg.WaitTimeout(time.Second))
		wg2.Dec()
	}()
	<-ready
	wg.Dec()
	require.NoError(t, wg2.WaitTimeout(time.Second))

	wg = gsync.NewSelectableWaitGroup()
	assert.NoError(t, wg.WaitTimeout(100*time.Millisecond))
	wg.Add(3)
	assert.Error(t, wg.WaitTimeout(100*time.Millisecond))
	wg.Dec()
	assert.Error(t, wg.WaitTimeout(100*time.Millisecond))
	wg.Dec()
	assert.Error(t, wg.WaitTimeout(100*time.Millisecond))
	wg.Dec()
	assert.NoError(t, wg.WaitTimeout(100*time.Millisecond))
}

func testWaitGroup(t *testing.T, wg1 *gsync.SelectableWaitGroup, wg2 *gsync.SelectableWaitGroup) {
	n := 16
	wg1.Add(n)
	wg2.Add(n)
	exited := make(chan bool, n)
	for i := 0; i != n; i++ {
		go func() {
			wg1.Dec()
			require.NoError(t, wg2.WaitTimeout(time.Second))
			exited <- true
		}()
	}
	require.NoError(t, wg1.WaitTimeout(time.Second))
	for i := 0; i != n; i++ {
		select {
		case <-exited:
			require.Fail(t, "WaitGroup released group too soon")
		default:
		}
		wg2.Dec()
	}
	for i := 0; i != n; i++ {
		<-exited // Will block if barrier fails to unlock someone.
	}
}

func TestWaitGroup(t *testing.T) {
	wg1 := gsync.NewSelectableWaitGroup()
	wg2 := gsync.NewSelectableWaitGroup()

	// Run the same test a few times to ensure barrier is in a proper state.
	for i := 0; i != 8; i++ {
		testWaitGroup(t, wg1, wg2)
	}
}

func TestWaitGroupRace(t *testing.T) {
	// Run this test for about 1ms.
	for i := 0; i < 1000; i++ {
		wg := gsync.NewSelectableWaitGroup()
		n := new(int32)
		// spawn goroutine 1
		wg.Add(1)
		go func() {
			atomic.AddInt32(n, 1)
			wg.Dec()
		}()
		// spawn goroutine 2
		wg.Add(1)
		go func() {
			atomic.AddInt32(n, 1)
			wg.Dec()
		}()
		// Wait for goroutine 1 and 2
		require.NoError(t, wg.WaitTimeout(time.Second))
		if atomic.LoadInt32(n) != 2 {
			assert.Fail(t, "Spurious wakeup from Wait")
		}
	}
}
