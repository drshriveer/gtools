package internal

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/drshriveer/gcommon/pkg/gsync"
)

// TODO: benchmark with various impl & against sync.WG
type SelectableWaitGroup2 struct {
	mu    sync.RWMutex
	count atomic.Int32
	wChan chan struct{}
}

func NewSelectableWaitGroup2() *SelectableWaitGroup2 {
	return &SelectableWaitGroup2{
		wChan: closedChan,
		count: atomic.Int32{},
	}
}

// Inc adds 1 to the wait group.
func (wg *SelectableWaitGroup2) Inc() int {
	return wg.Add(1)
}

// Dec adds -1 to the wait group.
func (wg *SelectableWaitGroup2) Dec() int {
	return wg.Add(-1)
}

func (wg *SelectableWaitGroup2) Add(in int32) int {
	// This impl will try to modify count with a read lock only.
	// and will block all operations only if the old version is
	// 1. non-zero and new version is zero
	// 2. zero to non-zero
	wg.mu.RLock()
	if in < 0 {
		newCount := wg.count.Add(in)
		if newCount == 0 {
			wg.mu.RUnlock()
			wg.mu.Lock()
			close(wg.wChan)
			wg.wChan = closedChan
			wg.mu.Unlock()
		} else if newCount < 0 {
			panic("less than zero")
		} else {
			wg.mu.RUnlock()
			return int(newCount)
		}
	}
	newCount := wg.count.Add(in)
	if wg.wChan == closedChan {
		wg.mu.RUnlock()
		wg.mu.Lock()
		wg.wChan = make(chan struct{})
		wg.mu.Unlock()
	} else {
		wg.mu.RUnlock()
	}

	return int(newCount)
}

func (wg *SelectableWaitGroup2) Count() int {
	return int(wg.count.Load())
}

func (wg *SelectableWaitGroup2) Wait() <-chan struct{} {
	wg.mu.RLock()
	defer wg.mu.RUnlock()
	return wg.wChan
}

// WaitCTX waits for the group to complete or for a context to be done.
// If the context ends first, this method will return the context error.
func (wg *SelectableWaitGroup2) WaitCTX(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-wg.Wait():
		return nil
	}
}

// WaitTimeout will wait for the group to be complete a specified time and will return
// ErrWGTimeout if the timeout passes.
func (wg *SelectableWaitGroup2) WaitTimeout(timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-timer.C:
		return gsync.ErrWGTimeout.Raw()
	case <-wg.Wait():
		return nil
	}
}
