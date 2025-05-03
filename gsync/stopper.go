package gsync

import (
	"sync"
)

type Stopper struct {
	stopChan         chan struct{}
	stopOnce         sync.Once
	stopCompleteChan chan struct{}
	stopCompleteOnce sync.Once
	cleanups         []func()
}

func NewStopper() Stopper {
	return Stopper{
		stopChan:         make(chan struct{}),
		stopCompleteChan: make(chan struct{}),
	}
}

func (s *Stopper) AddCleanup(fn func()) {
	s.cleanups = append(s.cleanups, fn)
}

func (s *Stopper) Stop() {
	s.stopOnce.Do(func() {
		for _, fn := range s.cleanups {
			fn()
		}
		close(s.stopChan)
	})
}

func (s *Stopper) StopSig() <-chan struct{} {
	return s.stopChan
}

func (s *Stopper) StopCompleteSig() <-chan struct{} {
	return s.stopCompleteChan
}

func (s *Stopper) SignalStopComplete() {
	s.stopCompleteOnce.Do(func() {
		close(s.stopCompleteChan)
	})
}
