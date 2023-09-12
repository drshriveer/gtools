package internal

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/drshriveer/gtools/gsync"
)

// TODO: benchmark with various impl & against sync.WG
type SelectableWaitGroup3 struct {
	v atomic.Pointer[wg3internal]
}
type wg3internal struct {
	count int
	wChan chan struct{}
}

func NewSelectableWaitGroup3() *SelectableWaitGroup3 {
	wg := &SelectableWaitGroup3{
		v: atomic.Pointer[wg3internal]{},
	}
	wg.v.Store(&wg3internal{
		count: 0,
		wChan: closedChan,
	})
	return wg
}

// Inc adds 1 to the wait group.
func (wg *SelectableWaitGroup3) Inc() int {
	return wg.Add(1)
}

// Dec adds -1 to the wait group.
func (wg *SelectableWaitGroup3) Dec() int {
	return wg.Add(-1)
}

func (wg *SelectableWaitGroup3) Add(delta int) int {
	for {
		old := wg.v.Load()
		newV := &wg3internal{
			count: old.count + delta,
			wChan: old.wChan,
		}

		if newV.count == 0 && delta < 0 {
			newV.wChan = closedChan
			if wg.v.CompareAndSwap(old, newV) {
				if old.wChan != closedChan {
					close(old.wChan)
				}
				return newV.count
			}
		} else if newV.count > 0 {
			if newV.wChan == closedChan {
				newV.wChan = make(chan struct{})
				if wg.v.CompareAndSwap(old, newV) {
					return newV.count
				}
				close(newV.wChan)
			} else {
				if wg.v.CompareAndSwap(old, newV) {
					return newV.count
				}
			}
		}
	}
}

func (wg *SelectableWaitGroup3) Count() int {
	return wg.v.Load().count
}

func (wg *SelectableWaitGroup3) Wait() <-chan struct{} {
	return wg.v.Load().wChan
}

// WaitCTX waits for the group to complete or for a context to be done.
// If the context ends first, this method will return the context error.
func (wg *SelectableWaitGroup3) WaitCTX(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-wg.Wait():
		return nil
	}
}

// WaitTimeout will wait for the group to be complete a specified time and will return
// ErrWGTimeout if the timeout passes.
func (wg *SelectableWaitGroup3) WaitTimeout(timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-timer.C:
		return gsync.ErrWGTimeout.Raw()
	case <-wg.Wait():
		return nil
	}
}
