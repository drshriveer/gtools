package gsync

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// AccumulatorFn is a function that takes a value and the current accumulated value
// and returns the new accumulated value. It is used to accumulate results from
// multiple tasks in the Executor.
// This should be implemented by the caller to define how the results are combined.
type AccumulatorFn[TVal any, Tout any] func(add TVal, current Tout) Tout

type taskResult[TVal any] struct {
	value TVal
	err   error
}

// Executor is a generic executor and accumulator that runs tasks concurrently and accumulates results.
// Into an arbitrary type. It is designed to be used with a context.Context to allow for cancellation and
// timeout handling.
type Executor[TVal any, Tout any] struct {
	ctx context.Context

	runErrors atomic.Pointer[[]error]
	mu        sync.Mutex
	result    Tout
	accum     AccumulatorFn[TVal, Tout]

	stopper    Shutdown
	resultChan chan taskResult[TVal]
	inflight   SelectableWaitGroup
}

// NewSliceExecutor is boilerplate code to create a new Executor that accumulates results into a slice.
func NewSliceExecutor[TVal any](ctx context.Context) (*Executor[TVal, []TVal], func()) {
	return NewExecutor[TVal, []TVal](ctx, func(add TVal, current []TVal) []TVal {
		return append(current, add)
	}, nil)
}

// NewExecutor creates a new Executor with the given context, accumulator function, and initial value.
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
		stopper:    NewShutdown(),
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

// AddTask adds a task to the executor which will be immediately executed in the background.
// results are joined into the result of the executor using the provided accumulator function.
// This function will return an error if the context has already been canceled or if
// any task has already failed.
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
		case <-e.stopper.ShutdownSignal():
		}
	}()

	return nil
}

// WaitForCompletion waits for all tasks to complete and returns an error if any task failed.
// It also checks for context cancellation and handles it appropriately.
func (e *Executor[TVal, Tout]) WaitForCompletion() error {
	if e.ctx.Err() != nil {
		return e.ctx.Err()
	}
	select {
	case <-e.inflight.Wait():
		return e.checkErrors()
	case <-e.ctx.Done():
		return e.ctx.Err()
	case <-e.stopper.ShutdownSignal():
		return e.checkErrors()
	}
}

// WaitAndResult is boilerplate to wait for the executor to finish and return the result.
func (e *Executor[TVal, Tout]) WaitAndResult() (Tout, error) {
	err := e.WaitForCompletion()
	if err != nil {
		return *new(Tout), err
	}
	return e.Result(), nil
}

// Result returns the accumulated result of all tasks executed by the executor.
// It is safe to call this function even if the executor has not completed yet.
func (e *Executor[TVal, Tout]) Result() Tout {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.result
}

func (e *Executor[TCtx, TVal]) stop() {
	e.stopper.Shutdown()
	<-e.stopper.ShutdownCompleteSignal()
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
		case <-e.stopper.ShutdownSignal():
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
