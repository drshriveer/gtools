package gsync

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

// ErrWGTimeout indicates a wait group timeout.
var ErrWGTimeout = errors.New("WaitGroupTimeout: timed out waiting for SelectableWaitGroup")

var closedChan chan struct{}

func init() {
	closedChan = make(chan struct{})
	close(closedChan)
}

// SelectableWaitGroup is a wait group that can be used in a select block!
// you _must_ use the `NewSelectableWaitGroup` function to construct one.
type SelectableWaitGroup struct {
	count atomic.Int64
	wChan atomic.Pointer[chan struct{}]
}

// NewSelectableWaitGroup creates a new SelectableWaitGroup.
func NewSelectableWaitGroup() *SelectableWaitGroup {
	wg := &SelectableWaitGroup{
		count: atomic.Int64{},
		wChan: atomic.Pointer[chan struct{}]{},
	}
	wg.wChan.Store(&closedChan)
	return wg
}

// Inc adds 1 to the wait group.
func (wg *SelectableWaitGroup) Inc() int {
	return wg.Add(1)
}

// Dec adds -1 to the wait group.
func (wg *SelectableWaitGroup) Dec() int {
	return wg.Add(-1)
}

// Add can be used to add or subtract a number from the wait group.
func (wg *SelectableWaitGroup) Add(delta int) int {
	newV := wg.count.Add(int64(delta))
	if newV == 0 {
		oldChan := wg.wChan.Swap(&closedChan)
		if oldChan != &closedChan {
			close(*oldChan)
		}
	} else if delta > 0 && newV == int64(delta) {
		newChan := make(chan struct{})
		if !wg.wChan.CompareAndSwap(&closedChan, &newChan) {
			close(newChan)
		}
	}

	return int(newV)
}

// Count returns current count of the wait group.
func (wg *SelectableWaitGroup) Count() int {
	return int(wg.count.Load())
}

// Wait returns a channel to use in a select block when the wait group reaches zero.
func (wg *SelectableWaitGroup) Wait() <-chan struct{} {
	// there is a race between updating the counter and updating the channel
	// .. so to make sure we have a consistent state check that we don't have a > zero count
	// but a closed channel. Do this repeatedly until the taskResult makes sense.
	for {
		count := wg.count.Load()
		wgChan := wg.wChan.Load()
		if count == 0 || (count > 0 && wgChan != &closedChan) {
			return *wgChan
		}
	}
}

// WaitCTX waits for the group to complete or for a context to be done.
// If the context ends first, this method will return the context error.
func (wg *SelectableWaitGroup) WaitCTX(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-wg.Wait():
		return nil
	}
}

// WaitTimeout will wait for the group to be complete a specified time and will return
// ErrWGTimeout if the timeout passes.
func (wg *SelectableWaitGroup) WaitTimeout(timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-timer.C:
		return ErrWGTimeout
	case <-wg.Wait():
		return nil
	}
}
