package gsync

import (
	"context"
	"sync"
	"time"

	"github.com/drshriveer/gcommon/pkg/errors"
)

// ErrWGTimeout indicates a wait group timeout.
var ErrWGTimeout errors.Factory = &errors.GError{
	Name:    "ErrWGTimeout",
	Message: "timed out waiting for SelectableWaitGroup",
}

var closedChan chan struct{}

func init() {
	closedChan = make(chan struct{})
	close(closedChan)
}

// TODO: benchmark with various impl & against sync.WG
type SelectableWaitGroup struct {
	mu    sync.RWMutex
	count int
	wChan chan struct{}
}

func NewSelectableWaitGroup() *SelectableWaitGroup {
	return &SelectableWaitGroup{
		wChan: closedChan,
	}
}

// Inc adds 1 to the wait group.
func (wg *SelectableWaitGroup) Inc() int {
	return wg.Add(1)
}

// Dec adds -1 to the wait group.
func (wg *SelectableWaitGroup) Dec() int {
	return wg.Add(-1)
}

func (wg *SelectableWaitGroup) Add(in int) int {
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

func (wg *SelectableWaitGroup) Count() int {
	wg.mu.RLock()
	defer wg.mu.RUnlock()
	return wg.count
}

func (wg *SelectableWaitGroup) Wait() <-chan struct{} {
	wg.mu.RLock()
	defer wg.mu.RUnlock()
	return wg.wChan
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
		return ErrWGTimeout.Raw()
	case <-wg.Wait():
		return nil
	}
}
