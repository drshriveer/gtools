package gsync

import (
	"context"
	"sync/atomic"
	"time"
)

// TODO: benchmark with various impl & against sync.WG
type SelectableWaitGroup4 struct {
	count atomic.Int64
	wChan atomic.Pointer[chan struct{}]
}

func NewSelectableWaitGroup4() *SelectableWaitGroup4 {
	wg := &SelectableWaitGroup4{
		count: atomic.Int64{},
		wChan: atomic.Pointer[chan struct{}]{},
	}
	wg.wChan.Store(&closedChan)
	return wg
}

// Inc adds 1 to the wait group.
func (wg *SelectableWaitGroup4) Inc() int {
	return wg.Add(1)
}

// Dec adds -1 to the wait group.
func (wg *SelectableWaitGroup4) Dec() int {
	return wg.Add(-1)
}

func (wg *SelectableWaitGroup4) Add(delta int) int {
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

func (wg *SelectableWaitGroup4) Count() int {
	return int(wg.count.Load())
}

func (wg *SelectableWaitGroup4) Wait() <-chan struct{} {
	// there is a race between updating the counter and updating the channel
	// .. so to make sure we have a consistent state check that we don't have a > zero count
	// but a closed channel. Do this repeatedly until the result makes sense.
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
func (wg *SelectableWaitGroup4) WaitCTX(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-wg.Wait():
		return nil
	}
}

// WaitTimeout will wait for the group to be complete a specified time and will return
// ErrWGTimeout if the timeout passes.
func (wg *SelectableWaitGroup4) WaitTimeout(timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-timer.C:
		return ErrWGTimeout.Raw()
	case <-wg.Wait():
		return nil
	}
}
