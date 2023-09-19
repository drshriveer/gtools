//nolint:all // this code exists for benchmark comparisons only.
package internal

import (
	"context"
	"sync"
	"time"

	"github.com/drshriveer/gtools/gsync"
)

var closedChan chan struct{}

func init() {
	closedChan = make(chan struct{})
	close(closedChan)
}

// TODO: benchmark with various impl & against sync.WG.
type SelectableWaitGroup1 struct {
	mu    sync.RWMutex
	count int
	wChan chan struct{}
}

func NewSelectableWaitGroup1() *SelectableWaitGroup1 {
	return &SelectableWaitGroup1{
		wChan: closedChan,
	}
}

// Inc adds 1 to the wait group.
func (wg *SelectableWaitGroup1) Inc() int {
	return wg.Add(1)
}

// Dec adds -1 to the wait group.
func (wg *SelectableWaitGroup1) Dec() int {
	return wg.Add(-1)
}

func (wg *SelectableWaitGroup1) Add(in int) int {
	wg.mu.Lock()
	defer wg.mu.Unlock()
	wg.count += in
	if wg.count < 0 {
		panic("wait group has gone negative")
	} else if wg.count == 0 {
		if wg.wChan != closedChan {
			close(wg.wChan)
			wg.wChan = closedChan
		}
	} else if wg.wChan == closedChan {
		wg.wChan = make(chan struct{})
	}
	return wg.count
}

func (wg *SelectableWaitGroup1) Count() int {
	wg.mu.RLock()
	defer wg.mu.RUnlock()
	return wg.count
}

func (wg *SelectableWaitGroup1) Wait() <-chan struct{} {
	wg.mu.RLock()
	defer wg.mu.RUnlock()
	return wg.wChan
}

// WaitCTX waits for the group to complete or for a context to be done.
// If the context ends first, this method will return the context error.
func (wg *SelectableWaitGroup1) WaitCTX(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-wg.Wait():
		return nil
	}
}

// WaitTimeout will wait for the group to be complete a specified time and will return
// ErrWGTimeout if the timeout passes.
func (wg *SelectableWaitGroup1) WaitTimeout(timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-timer.C:
		return gsync.ErrWGTimeout.Base()
	case <-wg.Wait():
		return nil
	}
}
