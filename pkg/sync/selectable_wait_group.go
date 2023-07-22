package sync

import "context"

// TODO: benchmark with various impl & against sync.WG
type SelectableWaitGroup struct {
}

func (wg *SelectableWaitGroup) Inc() int {

}

func (wg *SelectableWaitGroup) Dec() int {
}

func (wg *SelectableWaitGroup) Add(in int) int {

}

func (wg *SelectableWaitGroup) Wait() <-chan struct {
}

func (wg *SelectableWaitGroup) WaitCTX(ctx context.Context) error {

}

func (wg *SelectableWaitGroup) WaitTimeout(ctx context.Context) error {

}
