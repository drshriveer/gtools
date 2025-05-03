package gsync

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

type AccumulatorFn[TVal any, Tout any] func(add TVal, current Tout) Tout

type taskResult[TVal any] struct {
	value TVal
	err   error
}
type Executor[TVal any, Tout any] struct {
	ctx context.Context

	runErrors atomic.Pointer[[]error]
	mu        sync.Mutex
	result    Tout
	accum     AccumulatorFn[TVal, Tout]

	stopper    Stopper
	resultChan chan taskResult[TVal]
	inflight   SelectableWaitGroup
}

func NewSliceExecutor[TVal any](ctx context.Context) (*Executor[TVal, []TVal], func()) {
	return NewExecutor[TVal, []TVal](ctx, func(add TVal, current []TVal) []TVal {
		return append(current, add)
	}, nil)
}

func NewExecutor[TVal any, Tout any](
	ctx context.Context,
	accum AccumulatorFn[TVal, Tout],
	initial Tout,
) (*Executor[TVal, Tout], func()) {
	result := &Executor[TVal, Tout]{
		result:     initial,
		accum:      accum,
		runErrors:  atomic.Pointer[[]error]{},
		resultChan: make(chan taskResult[TVal], 1),
		stopper:    NewStopper(),
		inflight: SelectableWaitGroup{
			count: atomic.Int64{},
			wChan: atomic.Pointer[chan struct{}]{},
		},
	}
	var cancel context.CancelFunc
	result.ctx, cancel = context.WithCancel(ctx)
	result.stopper.AddCleanup(cancel)

	result.inflight.wChan.Store(&closedChan)
	result.runErrors.Store(nil)
	go result.backgroundWorker()
	return result, result.stop
}

func (e *Executor[TVal, Tout]) AddTask(fn func(ctx context.Context) (TVal, error)) (err error) {
	err = e.checkErrors()
	if err != nil {
		return err
	}
	e.inflight.Inc()
	go func() {
		result, err := fn(e.ctx)
		select {
		case e.resultChan <- taskResult[TVal]{result, err}:
		case <-e.ctx.Done():
		case <-e.stopper.StopSig():
		}
	}()

	return nil
}

func (e *Executor[TVal, Tout]) WaitForCompletion() error {
	if e.ctx.Err() != nil {
		return e.ctx.Err()
	}
	select {
	case <-e.inflight.Wait():
		return e.checkErrors()
	case <-e.ctx.Done():
		return e.ctx.Err()
	case <-e.stopper.StopSig():
		return e.checkErrors()
	}
}

func (e *Executor[TVal, Tout]) WaitAndResult() (Tout, error) {
	err := e.WaitForCompletion()
	if err != nil {
		return *new(Tout), err
	}
	return e.Result(), nil
}

func (e *Executor[TVal, Tout]) Result() Tout {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.result
}

func (e *Executor[TCtx, TVal]) stop() {
	e.stopper.Stop()
	<-e.stopper.StopCompleteSig()
}

func (e *Executor[TCtx, TVal]) backgroundWorker() {
	defer e.stopper.SignalStopComplete()
	for {
		select {
		case r := <-e.resultChan:
			if r.err != nil {
				e.addRunError(r.err)
			} else {
				e.mu.Lock()
				e.result = e.accum(r.value, e.result)
				e.mu.Unlock()
			}
			e.inflight.Dec()
		case <-e.ctx.Done():
			return
		case <-e.stopper.StopSig():
			return
		}
	}
}

func (e *Executor[TVal, Tout]) checkErrors() error {
	if e.ctx.Err() != nil {
		return e.ctx.Err()
	}
	errs := e.runErrors.Load()
	if errs != nil && len(*errs) > 0 {
		return errors.Join(*errs...)
	}
	return nil
}

func (e *Executor[TVal, Tout]) addRunError(err error) {
	if err == nil {
		return
	}

	for {
		old := e.runErrors.Load()
		errs := []error{err}
		if old == nil {
			if e.runErrors.CompareAndSwap(nil, &errs) {
				return
			}
		}
		errs = append(errs, *old...)
		if e.runErrors.CompareAndSwap(old, &errs) {
			return
		}
	}
}
